package osmimporter

import (
	"context"
	"fmt"
	"net/http"
	"runtime"

	"github.com/amuttsch/bahnglei.se/pkg/config"
	"github.com/amuttsch/bahnglei.se/pkg/repository"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/paulmach/osm"
	"github.com/paulmach/osm/osmpbf"

	log "github.com/sirupsen/logrus"
)

type parser interface {
	parse(object osm.Object)
}

type osmImporter struct {
	config *config.Config
	repo   *repository.Queries
	db     *pgxpool.Pool
}

func New(config *config.Config, repo *repository.Queries, db *pgxpool.Pool) *osmImporter {
	return &osmImporter{
		repo:   repo,
		config: config,
		db:     db,
	}
}

func (i *osmImporter) Import(ctx context.Context) {
	for _, c := range i.config.Countries {
		country, err := i.repo.SaveCountry(ctx, repository.SaveCountryParams{
			IsoCode: c.Iso,
			OsmUrl:  c.Url,
			Name:    c.Name,
		})
		if err != nil {
			log.Errorf("Failed to create country %s: %+v", c.Name, err)
			return
		}

		state, err := i.repo.CreateImportState(ctx, country.IsoCode)
		if err != nil {
			log.Errorf("Failed to import country %s: %+v", country.Name, err)
			return
		}
		err = i.importFirstPass(ctx, country, state.ID)
		if err != nil {
			log.Errorf("Failed to import country %s: %+v", country.Name, err)
			return
		}
		err = i.importSecondPass(ctx, country, state.ID)
		if err != nil {
			log.Errorf("Failed to import country %s: %+v", country.Name, err)
			return
		}
	}
}

func (i *osmImporter) importFirstPass(ctx context.Context, country repository.Country, stateId int32) error {
	log.Infof("Importing country %s from %s", country.IsoCode, country.OsmUrl)

	response, err := http.Get(country.OsmUrl)
	if err != nil {
		fmt.Println(err)
		i.repo.UpdateImportState(ctx, repository.UpdateImportStateParams{
			ID:              stateId,
			NumberStations:  0,
			NumberPlatforms: 0,
			State:           "failed: Get OSM data",
		})

		return err
	}

	defer response.Body.Close()

	platformParser := newPlatformParser(i.db, ctx, i.repo)
	stationParser := newStationParser(i.db, ctx, i.repo, country.IsoCode)
	stopPositionParser := newStopPositionParser(i.db, ctx, i.repo)

	scanner := osmpbf.New(ctx, response.Body, runtime.GOMAXPROCS(-1))
	defer scanner.Close()

	for scanner.Scan() {
		osmObject := scanner.Object()
		platformParser.parse(country.IsoCode, osmObject)
		stationParser.parse(osmObject)
		stopPositionParser.parse(country.IsoCode, osmObject)
	}

	log.Infof("Got %d stations", stationParser.numElements)
	log.Infof("Got %d platforms", platformParser.numElements)
	log.Infof("Got %d stop positions", stopPositionParser.numElements)

	if err := scanner.Err(); err != nil {
		i.repo.UpdateImportState(ctx, repository.UpdateImportStateParams{
			ID:              stateId,
			NumberStations:  0,
			NumberPlatforms: 0,
			State:           "failed: " + err.Error(),
		})
		log.Errorf("Failed to import country %s: %+v", country.Name, err)
		return err
	}

	log.Info("Calculating stations for stop positions")
	i.repo.SetStopPositionStationIdToNearestStation(ctx, country.IsoCode)

	i.repo.UpdateImportState(ctx, repository.UpdateImportStateParams{
		ID:              stateId,
		NumberStations:  int32(stationParser.numElements),
		NumberPlatforms: int32(platformParser.numElements),
		State:           "1st pass done",
	})

	return nil

}

func (i *osmImporter) importSecondPass(ctx context.Context, country repository.Country, stateId int32) error {

	response, err := http.Get(country.OsmUrl)
	if err != nil {
		i.repo.UpdateImportState(ctx, repository.UpdateImportStateParams{
			ID:              stateId,
			NumberStations:  0,
			NumberPlatforms: 0,
			State:           "failed: Get OSM data",
		})

		return err
	}

	defer response.Body.Close()

	platformNodeParser := newPlatformNodeParser(i.db, ctx, i.repo)

	scanner := osmpbf.New(ctx, response.Body, runtime.GOMAXPROCS(-1))
	defer scanner.Close()

	for scanner.Scan() {
		osmObject := scanner.Object()
		platformNodeParser.parse(osmObject, country.IsoCode)
	}

	if err := scanner.Err(); err != nil {
		i.repo.UpdateImportState(ctx, repository.UpdateImportStateParams{
			ID:              stateId,
			NumberStations:  0,
			NumberPlatforms: 0,
			State:           "failed: " + err.Error(),
		})
		return err
	}

	log.Info("Calculating center coordinate for platform")
	err = i.repo.SetPlatformCoordinates(ctx, country.IsoCode)
	if err != nil {
		log.Errorf("Failed to set center for platforms: %+v", err)
	}

	log.Info("Calculating stations for platforms")
	err = i.repo.SetPlatformToNearestStation(ctx, country.IsoCode)
	if err != nil {
		log.Errorf("Failed to set stations for platforms: %+v", err)
	}

	log.Info("Setting stop position neighbors")
	err = i.repo.SetStopPositionNeighbors(ctx, country.IsoCode)
	if err != nil {
		log.Errorf("Failed to set center for platforms: %+v", err)
	}

	log.Info("Setting number of tracks for stations")
	err = i.repo.SetStationNumberOfTracks(ctx, country.IsoCode)
	if err != nil {
		log.Errorf("Failed to set number of tracks: %+v", err)
	}

	i.repo.UpdateImportState(ctx, repository.UpdateImportStateParams{
		ID:    stateId,
		State: "Done",
	})

	return nil

}

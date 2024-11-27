package osmimporter

import (
	"context"

	"github.com/amuttsch/bahnglei.se/pkg/config"
	"github.com/amuttsch/bahnglei.se/pkg/repository"
	"github.com/jackc/pgx/v5/pgxpool"

	log "github.com/sirupsen/logrus"
)

type osmImporter struct {
	config *config.Config
	repo   *repository.Queries
	db     *pgxpool.Pool
}

func Run(ctx context.Context, config *config.Config, repo *repository.Queries, db *pgxpool.Pool) error {
	overpass := NewOverpassApi(ctx, db, repo, config.OverpassUrl)

	for _, c := range config.Countries {
		log.Infof("Start importing country %s", c.Name)
		country, err := repo.SaveCountry(ctx, repository.SaveCountryParams{
			IsoCode: c.Iso,
			Name:    c.Name,
		})
		if err != nil {
			log.Errorf("Failed to create country %s: %+v", c.Name, err)
			return err
		}

		//overpass.fetchStations(c.Area, c.Iso)
		// overpass.fetchStopPositions(c.Area, c.Iso)
		//overpass.fetchPlatforms(c.Area, c.Iso)
		overpass.fetchStopAreas(c.Area, c.Iso)

		calculateDistances(ctx, repo, country)
	}

	return nil
}

func calculateDistances(ctx context.Context, repo *repository.Queries, country repository.Country) error {
	log.Info("Calculating stations for stop positions")
	err := repo.SetStopPositionStationIdToNearestStation(ctx, country.IsoCode)
	if err != nil {
		log.Errorf("Failed to set stations for platforms: %+v", err)
	}

	log.Info("Calculating stations for platforms")
	err = repo.SetPlatformToNearestStation(ctx, country.IsoCode)
	if err != nil {
		log.Errorf("Failed to set stations for platforms: %+v", err)
	}

	log.Info("Setting stop position neighbors")
	err = repo.SetStopPositionNeighbors(ctx, country.IsoCode)
	if err != nil {
		log.Errorf("Failed to set center for platforms: %+v", err)
	}

	log.Info("Setting number of tracks for stations")
	err = repo.SetStationNumberOfTracks(ctx, country.IsoCode)
	if err != nil {
		log.Errorf("Failed to set number of tracks: %+v", err)
	}

	return nil

}

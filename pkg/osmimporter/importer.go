package osmimporter

import (
	"context"
	"fmt"
	"net/http"
	"runtime"
	"strings"

	"github.com/amuttsch/bahnglei.se/pkg/config"
	"github.com/amuttsch/bahnglei.se/pkg/repository"
	"github.com/paulmach/osm"
	"github.com/paulmach/osm/osmpbf"
	"github.com/samber/lo"

	"github.com/sirupsen/logrus"
	log "github.com/sirupsen/logrus"
)

type osmImporter struct {
	config *config.Config
	repo   *repository.Queries
}

type osmStation struct {
	id        int64
	name      string
	lat       float64
	lng       float64
	operator  string
	wikidata  string
	wikipedia string
}

type osmPlatform struct {
	id        int64
	positions string
}

type osmStopPosition struct {
	id       int64
	position string
	lat      float64
	lng      float64
}

type osmStopArea struct {
	id      int
	members []int64
}

func New(config *config.Config, repo *repository.Queries) *osmImporter {
	return &osmImporter{
		repo:   repo,
		config: config,
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

		ci := countryImporter{
			country:       country,
			osmImporter:   i,
			stations:      make(map[int64]osmStation),
			platforms:     make(map[int64]osmPlatform),
			stopPositions: make(map[int64]osmStopPosition),
			stopAreas:     []osmStopArea{},
		}
		ci.importCountry(ctx)
	}
}

type countryImporter struct {
	country       repository.Country
	osmImporter   *osmImporter
	stations      map[int64]osmStation
	platforms     map[int64]osmPlatform
	stopPositions map[int64]osmStopPosition
	stopAreas     []osmStopArea
}

func (i *countryImporter) importCountry(ctx context.Context) {
	state, err := i.osmImporter.repo.CreateImportState(ctx, i.country.IsoCode)
	if err != nil {
		fmt.Println(err)
	}

	log.Infof("Importing country %s from %s", i.country.IsoCode, i.country.OsmUrl)

	response, err := http.Get(i.country.OsmUrl)
	if err != nil {
		fmt.Println(err)
		i.osmImporter.repo.UpdateImportState(ctx, repository.UpdateImportStateParams{
			ID:              state.ID,
			NumberStations:  0,
			NumberPlatforms: 0,
			State:           "failed: Get OSM data",
		})

		return
	}

	defer response.Body.Close()

	scanner := osmpbf.New(ctx, response.Body, runtime.GOMAXPROCS(-1))
	defer scanner.Close()

	for scanner.Scan() {
		i.parseOsmData(scanner)
	}

	log.Infof("Got %d stations", len(i.stations))
	log.Infof("Got %d platforms", len(i.platforms))
	log.Infof("Got %d stop areas", len(i.stopAreas))
	i.saveStations(ctx)

	if err := scanner.Err(); err != nil {
		i.osmImporter.repo.UpdateImportState(ctx, repository.UpdateImportStateParams{
			ID:              state.ID,
			NumberStations:  0,
			NumberPlatforms: 0,
			State:           "failed: " + err.Error(),
		})
		log.Errorf("Failed to import country %s: %+v", i.country.Name, err)
		return
	}

	i.osmImporter.repo.UpdateImportState(ctx, repository.UpdateImportStateParams{
		ID:              state.ID,
		NumberStations:  int32(len(i.stations)),
		NumberPlatforms: int32(len(i.platforms)),
		State:           "finished",
	})

}

func (i *countryImporter) parseOsmData(scanner *osmpbf.Scanner) {
	switch o := scanner.Object().(type) {
	case *osm.Node:
		isStation := o.Tags.Find("public_transport") == "station"
		isStopPosition := o.Tags.Find("public_transport") == "stop_position"
		railwayTag := o.Tags.Find("railway")
		if railwayTag == "" {
			break
		}

		isRailwayStation := railwayTag == "halt" || railwayTag == "station"
		isRailwayStop := railwayTag == "stop"
		isTrain := o.Tags.Find("train") == "yes"

		if isStation && isRailwayStation {
			elementID := o.ID
			i.stations[int64(elementID)] = osmStation{
				id:        int64(elementID),
				name:      o.Tags.Find("name"),
				lat:       o.Lat,
				lng:       o.Lon,
				operator:  o.Tags.Find("operator"),
				wikidata:  o.Tags.Find("wikidata"),
				wikipedia: o.Tags.Find("wikipedia"),
			}
		}

		if isTrain && isStopPosition && isRailwayStop {
			ref := o.Tags.Find("ref")
			localRef := o.Tags.Find("local_ref")
			position := ref
			if localRef != "" {
				position = localRef
			}
			i.stopPositions[int64(o.ID)] = osmStopPosition{
				id:       int64(o.ID),
				position: position,
				lat:      o.Lat,
				lng:      o.Lon,
			}
		}
	case *osm.Way:
		isTrain := o.Tags.Find("train") == "yes"
		ref := o.Tags.Find("ref")
		isPlatform := o.Tags.Find("public_transport") == "platform" || o.Tags.Find("railway") == "platform"

		if isTrain && isPlatform {
			i.platforms[int64(o.ID)] = osmPlatform{
				id:        int64(o.ID),
				positions: ref,
			}
		}

	case *osm.Relation:
		isStopArea := o.Tags.Find("public_transport") == "stop_area"
		isPublicTransport := o.Tags.Find("type") == "public_transport"

		isTrain := o.Tags.Find("train") == "yes"
		ref := o.Tags.Find("ref")
		isPlatform := o.Tags.Find("public_transport") == "platform" || o.Tags.Find("railway") == "platform"

		if isTrain && isPlatform {
			i.platforms[int64(o.ID)] = osmPlatform{
				id:        int64(o.ID),
				positions: ref,
			}
		}

		if isStopArea && isPublicTransport {
			var members []int64
			for _, member := range o.Members {
				members = append(members, member.Ref)
			}
			i.stopAreas = append(i.stopAreas, osmStopArea{
				id:      int(o.ID),
				members: members,
			})
		}
	}
}

func (i *countryImporter) saveStations(ctx context.Context) {
	log.Info("Start saving stations")
	for _, s := range i.stations {
		i.osmImporter.repo.DeletePlatformsForStation(ctx, s.id)
		i.osmImporter.repo.DeleteStopPositionsForStation(ctx, s.id)
		i.osmImporter.repo.DeleteStation(ctx, s.id)

		_, err := i.osmImporter.repo.CreateStation(ctx, repository.CreateStationParams{
			ID:             s.id,
			CountryIsoCode: i.country.IsoCode,
			Name:           s.name,
			Lat:            s.lat,
			Lng:            s.lng,
			Operator:       s.operator,
			Wikidata:       s.wikidata,
			Wikipedia:      s.wikipedia,
		})
		if err != nil {
			log.Errorf("Failed to save station: %+v", err)
		}
	}

	for _, sa := range i.stopAreas {
		var stopAreaStation osmStation
		var stopAreaPlatforms []osmPlatform
		var stopAreaStopPositions []osmStopPosition
		for _, m := range sa.members {
			s := i.stations[m]
			if s != (osmStation{}) {
				stopAreaStation = s
			}
			p := i.platforms[m]
			if p != (osmPlatform{}) {
				stopAreaPlatforms = append(stopAreaPlatforms, p)
			}
			sp := i.stopPositions[m]
			if sp != (osmStopPosition{}) {
				stopAreaStopPositions = append(stopAreaStopPositions, sp)
			}
		}

		if stopAreaStation == (osmStation{}) {
			continue
		}

		i.osmImporter.repo.UpdateStationNumberOfTracks(ctx, repository.UpdateStationNumberOfTracksParams{
			ID:     stopAreaStation.id,
			Tracks: int64(len(stopAreaStopPositions)),
		})

		positions := make([][]string, 3)
		for _, sap := range stopAreaPlatforms {
			i.osmImporter.repo.CreatePlatform(ctx, repository.CreatePlatformParams{
				ID:        sap.id,
				StationID: stopAreaStation.id,
				Positions: sap.positions,
			})
			positions = append(positions, strings.Split(sap.positions, ";"))
		}
		for _, sp := range stopAreaStopPositions {
			neighbors, _ := lo.Find(positions, func(p []string) bool {
				return lo.Contains(p, sp.position)
			})
			i.osmImporter.repo.CreateStopPosition(ctx, repository.CreateStopPositionParams{
				ID:        sp.id,
				StationID: stopAreaStation.id,
				Platform:  sp.position,
				Lat:       sp.lat,
				Lng:       sp.lng,
				Neighbors: strings.Join(neighbors, ";"),
			})
		}
	}
}

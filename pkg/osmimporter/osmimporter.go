package osmimporter

import (
	"context"
	"fmt"
	"net/http"
	"runtime"
	"strings"

	"github.com/amuttsch/bahnglei.se/pkg/config"
	"github.com/amuttsch/bahnglei.se/pkg/repo/country"
	importerRepo "github.com/amuttsch/bahnglei.se/pkg/repo/importer"
	stationRepo "github.com/amuttsch/bahnglei.se/pkg/repo/station"
	"github.com/paulmach/osm"
	"github.com/paulmach/osm/osmpbf"
	"gorm.io/gorm"
      "github.com/samber/lo"

	log "github.com/sirupsen/logrus"
)

type osmImporter struct {
	config       *config.Config
	countryRepo  country.Repo
	importerRepo importerRepo.Repo
	stationRepo  stationRepo.Repo
}

type station struct {
	id        int64
	name      string
	lat       float64
	lng       float64
	operator  string
	wikidata  string
	wikipedia string
}

type platform struct {
	id        int64
	positions string
}

type stopPosition struct {
	id       int64
	position string
	lat      float64
	lng      float64
}

type stopArea struct {
	id      int
	members []int64
}

func New(config *config.Config, countryRepo country.Repo, importerRepo importerRepo.Repo, ststationRepo stationRepo.Repo) *osmImporter {
	return &osmImporter{
		countryRepo:  countryRepo,
		importerRepo: importerRepo,
		stationRepo:  ststationRepo,
		config:       config,
	}
}

func (i *osmImporter) Import() {
	for _, c := range i.config.Countries {
		country := country.Country{
			IsoCode: c.Iso,
			OsmUrl:  c.Url,
			Name:    c.Name,
		}
		i.countryRepo.Save(country)

		ci := countryImporter{
			country:       country,
			osmImporter:   i,
			stations:      make(map[int64]station),
			platforms:     make(map[int64]platform),
			stopPositions: make(map[int64]stopPosition),
			stopAreas:     []stopArea{},
		}
		ci.importCountry()
	}
}

type countryImporter struct {
	country       country.Country
	osmImporter   *osmImporter
	stations      map[int64]station
	platforms     map[int64]platform
	stopPositions map[int64]stopPosition
	stopAreas     []stopArea
}

func (i *countryImporter) importCountry() {
	importerState := importerRepo.ImporterModel{
		Country:         i.country,
		State:           "started",
		NumberStations:  0,
		NumberPlatforms: 0,
	}
	i.osmImporter.importerRepo.Save(&importerState)

	log.Infof("Importing country %s from %s", i.country.IsoCode, i.country.OsmUrl)

	response, err := http.Get(i.country.OsmUrl)
	if err != nil {
		fmt.Println(err)
		return
	}

	defer response.Body.Close()

	scanner := osmpbf.New(context.Background(), response.Body, runtime.GOMAXPROCS(-1))
	defer scanner.Close()

	for scanner.Scan() {
		i.parseOsmData(scanner)
	}

	log.Infof("Got %d stations", len(i.stations))
	log.Infof("Got %d platforms", len(i.platforms))
	log.Infof("Got %d stop areas", len(i.stopAreas))

	i.saveStations()

	importerState.State = "finished"
	importerState.NumberStations = int64(len(i.stations))
	importerState.NumberPlatforms = int64(len(i.platforms))
	i.osmImporter.importerRepo.Save(&importerState)

	if err := scanner.Err(); err != nil {
		panic(err)
	}
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
			i.stations[int64(elementID)] = station{
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
			i.stopPositions[int64(o.ID)] = stopPosition{
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
			i.platforms[int64(o.ID)] = platform{
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
			i.platforms[int64(o.ID)] = platform{
				id:        int64(o.ID),
				positions: ref,
			}
		}

		if isStopArea && isPublicTransport {
			var members []int64
			for _, member := range o.Members {
				members = append(members, member.Ref)
			}
			i.stopAreas = append(i.stopAreas, stopArea{
				id:      int(o.ID),
				members: members,
			})
		}
	}
}

func (i *countryImporter) saveStations() {
	for _, s := range i.stations {
		bahnStation := stationRepo.Station{
			Model: gorm.Model{
				ID: uint(s.id),
			},
			Country:   i.country,
			Name:      s.name,
			Lat:       s.lat,
			Lng:       s.lng,
			Operator:  s.operator,
			Wikidata:  s.wikidata,
			Wikipedia: s.wikipedia,
		}

		i.osmImporter.stationRepo.Save(&bahnStation)
	}

	for _, sa := range i.stopAreas {
		var stopAreaStation station
		var stopAreaPlatforms []platform
		var stopAreaStopPositions []stopPosition
		for _, m := range sa.members {
			s := i.stations[m]
			if s != (station{}) {
				stopAreaStation = s
			}
			p := i.platforms[m]
			if p != (platform{}) {
				stopAreaPlatforms = append(stopAreaPlatforms, p)
			}
			sp := i.stopPositions[m]
			if sp != (stopPosition{}) {
				stopAreaStopPositions = append(stopAreaStopPositions, sp)
			}
		}

		if stopAreaStation == (station{}) {
			continue
		}

		bahnStation := i.osmImporter.stationRepo.Get(uint(stopAreaStation.id))
		bahnStation.Tracks = len(stopAreaStopPositions)
    positions := make([][]string, 3)
		for _, sap := range stopAreaPlatforms {
			bahnStation.Platforms = append(bahnStation.Platforms, stationRepo.Platform{
				Model:     gorm.Model{ID: uint(sap.id)},
				Positions: sap.positions,
			})
      positions = append(positions, strings.Split(sap.positions, ";"))
		}
		for _, sp := range stopAreaStopPositions {
      neighbors, _ := lo.Find(positions, func(p []string) bool {
          return lo.Contains(p, sp.position)
      })
			bahnStation.StopPosition = append(bahnStation.StopPosition, stationRepo.StopPosition{
				Model:    gorm.Model{ID: uint(sp.id)},
				Platform: sp.position,
				Lat:      sp.lat,
				Lng:      sp.lng,
        Neighbors: strings.Join(neighbors, ";"),
			})
		}

		i.osmImporter.stationRepo.Save(bahnStation)
	}
}

package osmimporter

import (
	"context"

	"github.com/amuttsch/bahnglei.se/pkg/repository"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/paulmach/osm"
	"github.com/sirupsen/logrus"
)

type stationParser struct {
	db          *pgxpool.Pool
	ctx         context.Context
	repo        *repository.Queries
	numElements int64
	countryIso  string
}

func newStationParser(db *pgxpool.Pool, ctx context.Context, repo *repository.Queries, countryIso string) *stationParser {
	return &stationParser{
		db,
		ctx,
		repo,
		0,
		countryIso,
	}
}

func (p *stationParser) parse(object osm.Object) {
	switch o := object.(type) {
	case *osm.Node:
		isStation := o.Tags.Find("public_transport") == "station"
		railwayTag := o.Tags.Find("railway")
		if railwayTag == "" {
			break
		}

		isRailwayStation := railwayTag == "halt" || railwayTag == "station"

		if isStation || isRailwayStation {
			_, err := p.repo.CreateStation(p.ctx, repository.CreateStationParams{
				CountryIsoCode: p.countryIso,
				ID:             int64(o.ID),
				Name:           o.Tags.Find("name"),
				Coordinate: pgtype.Point{
					P: pgtype.Vec2{
						X: o.Lon,
						Y: o.Lat,
					},
					Valid: true,
				},
				Operator:  o.Tags.Find("operator"),
				Wikidata:  o.Tags.Find("wikidata"),
				Wikipedia: o.Tags.Find("wikipedia"),
			})

			if err != nil {
				logrus.Errorf("Failed to save station: %+v\n", err)
				break
			}

			p.numElements = p.numElements + 1
		}

	case *osm.Way:

	case *osm.Relation:
	}
}

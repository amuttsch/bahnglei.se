package osmimporter

import (
	"context"

	"github.com/amuttsch/bahnglei.se/pkg/repository"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/paulmach/osm"
	"github.com/sirupsen/logrus"
)

type stopPositionParser struct {
	db          *pgxpool.Pool
	ctx         context.Context
	repo        *repository.Queries
	numElements int64
}

func newStopPositionParser(db *pgxpool.Pool, ctx context.Context, repo *repository.Queries) *stopPositionParser {
	return &stopPositionParser{
		db,
		ctx,
		repo,
		0,
	}
}

func (p *stopPositionParser) parse(countryIso string, object osm.Object) {
	switch o := object.(type) {
	case *osm.Node:
		isStopPosition := o.Tags.Find("public_transport") == "stop_position"
		railwayTag := o.Tags.Find("railway")
		if railwayTag == "" {
			break
		}

		isRailwayStop := railwayTag == "stop"
		isTrain := o.Tags.Find("train") == "yes"

		if isTrain && isStopPosition && isRailwayStop {
			ref := o.Tags.Find("ref")
			localRef := o.Tags.Find("local_ref")
			platform := ref
			if localRef != "" {
				platform = localRef
			}
			_, err := p.repo.CreateStopPosition(p.ctx, repository.CreateStopPositionParams{
				ID:       int64(o.ID),
				Platform: platform,
				Coordinate: pgtype.Point{
					P: pgtype.Vec2{
						X: o.Lon,
						Y: o.Lat,
					},
					Valid: true,
				},
				CountryIsoCode: countryIso,
			})

			if err != nil {
				logrus.Errorf("Failed to save stop position: %+v\n", err)
				break
			}

			p.numElements = p.numElements + 1
		}
	case *osm.Way:

	case *osm.Relation:
	}
}

package osmimporter

import (
	"context"

	"github.com/amuttsch/bahnglei.se/pkg/repository"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/paulmach/osm"
	"github.com/sirupsen/logrus"
)

type platformParser struct {
	db          *pgxpool.Pool
	ctx         context.Context
	repo        *repository.Queries
	numElements int64
}

func newPlatformParser(db *pgxpool.Pool, ctx context.Context, repo *repository.Queries) *platformParser {
	return &platformParser{
		db,
		ctx,
		repo,
		0,
	}
}

func (p *platformParser) parse(countryIso string, object osm.Object) {
	switch o := object.(type) {
	case *osm.Node:

	case *osm.Way:
		isTrain := o.Tags.Find("train") == "yes"
		ref := o.Tags.Find("ref")
		isPlatform := o.Tags.Find("public_transport") == "platform" || o.Tags.Find("railway") == "platform"

		if isTrain && isPlatform {
			_, err := p.repo.CreatePlatform(p.ctx, repository.CreatePlatformParams{
				ID:             int64(o.ID),
				Positions:      ref,
				CountryIsoCode: countryIso,
			})
			if err != nil {
				logrus.Errorf("Failed to save platform: %+v\n", err)
				break
			}

			p.numElements = p.numElements + 1

			for _, node := range o.Nodes.NodeIDs() {
				err = p.repo.CreatePlatformNode(p.ctx, repository.CreatePlatformNodeParams{
					ID:             int64(node),
					PlatformID:     int64(o.ID),
					CountryIsoCode: countryIso,
				})
				if err != nil {
					logrus.Errorf("Failed to save platform node: %+v\n", err)
					break
				}
			}
		}

	case *osm.Relation:
		isTrain := o.Tags.Find("train") == "yes"
		ref := o.Tags.Find("ref")
		isPlatform := o.Tags.Find("public_transport") == "platform" || o.Tags.Find("railway") == "platform"

		if isTrain && isPlatform {
			_, err := p.repo.CreatePlatform(p.ctx, repository.CreatePlatformParams{
				ID:             int64(o.ID),
				Positions:      ref,
				CountryIsoCode: countryIso,
			})
			if err != nil {
				logrus.Errorf("Failed to save platform: %+v\n", err)
				break
			}

			p.numElements = p.numElements + 1

			for _, member := range o.Members {
				switch member.Type {
				case osm.TypeNode:
					err = p.repo.CreatePlatformNode(p.ctx, repository.CreatePlatformNodeParams{
						ID:             int64(member.ElementID()),
						PlatformID:     int64(o.ID),
						CountryIsoCode: countryIso,
					})
					if err != nil {
						logrus.Errorf("Failed to save platform node: %+v\n", err)
					}
				case osm.TypeWay:
					err = p.repo.CreatePlatformWay(p.ctx, repository.CreatePlatformWayParams{
						ID:             member.Ref,
						PlatformID:     int64(o.ID),
						CountryIsoCode: countryIso,
					})
					if err != nil {
						logrus.Errorf("Failed to save platform way: %+v\n", err)
					}
				}
			}
		}
	}
}

package osmimporter

import (
	"context"

	"github.com/amuttsch/bahnglei.se/pkg/repository"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/paulmach/osm"
	"github.com/sirupsen/logrus"
	"golang.org/x/text/language"
	"golang.org/x/text/message"
)

type platformWayParser struct {
	db            *pgxpool.Pool
	ctx           context.Context
	repo          *repository.Queries
	nodeBuffer    []repository.InsertTemporaryWaysParams
	processedWays int64
}

func newPlatformWayParser(db *pgxpool.Pool, ctx context.Context, repo *repository.Queries) *platformWayParser {
	err := repo.CreateTemporaryWaysTable(ctx)
	if err != nil {
		logrus.Errorf("Failed to create temporary nodes table: %+v", err)
	}

	return &platformWayParser{
		db,
		ctx,
		repo,
		make([]repository.InsertTemporaryWaysParams, 0, NODE_BUFFER_SIZE),
		0,
	}
}

func (p *platformWayParser) parse(object osm.Object, countryIso string) {
	switch o := object.(type) {
	case *osm.Node:

	case *osm.Way:
		for _, node := range o.Nodes {
			if len(p.nodeBuffer) == cap(p.nodeBuffer) {
				p.saveNodeBuffer(countryIso)
			}

			p.nodeBuffer = append(p.nodeBuffer, repository.InsertTemporaryWaysParams{
				ID:   int64(o.ID),
				Node: int64(node.ID),
			})
		}
		p.processedWays += 1

	case *osm.Relation:
	}
}

func (p *platformWayParser) saveNodeBuffer(countryIso string) {
	tx, err := p.db.BeginTx(p.ctx, pgx.TxOptions{})
	txRepo := p.repo.WithTx(tx)
	if err != nil {
		logrus.Errorf("Failed to open transaction: %+v", err)
	}
	_, err = txRepo.InsertTemporaryWays(p.ctx, p.nodeBuffer)
	if err != nil {
		logrus.Errorf("Failed to save temporary ways: %+v", err)
	}

	ids, err := txRepo.InsertPlatformNodesFromPlatformWays(p.ctx)
	if err != nil {
		logrus.Errorf("Failed to insert temporary way nodes: %+v", err)
	}

	printer := message.NewPrinter(language.English)
	logrus.Infof(
		"Saving and merging %d / %d ways for country %s. Total ways processed: %s",
		len(ids),
		cap(p.nodeBuffer),
		countryIso,
		printer.Sprint(p.processedWays),
	)
	p.nodeBuffer = p.nodeBuffer[:0]

	tx.Commit(p.ctx)
}

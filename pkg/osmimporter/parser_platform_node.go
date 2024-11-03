package osmimporter

import (
	"context"

	"github.com/amuttsch/bahnglei.se/pkg/repository"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/paulmach/osm"
	"github.com/sirupsen/logrus"
	"golang.org/x/text/language"
	"golang.org/x/text/message"
)

const NODE_BUFFER_SIZE = 10_000

type platformNodeParser struct {
	db             *pgxpool.Pool
	ctx            context.Context
	repo           *repository.Queries
	nodeBuffer     []repository.InsertTemporaryNodesParams
	processedNodes int64
}

func newPlatformNodeParser(db *pgxpool.Pool, ctx context.Context, repo *repository.Queries) *platformNodeParser {
	err := repo.CreateTemporaryNodeTable(ctx)
	if err != nil {
		logrus.Errorf("Failed to create temporary nodes table: %+v", err)
	}

	return &platformNodeParser{
		db,
		ctx,
		repo,
		make([]repository.InsertTemporaryNodesParams, 0, NODE_BUFFER_SIZE),
		0,
	}
}

func (p *platformNodeParser) parse(object osm.Object, countryIso string) {
	switch o := object.(type) {
	case *osm.Node:
		if len(p.nodeBuffer) == cap(p.nodeBuffer) {
			p.saveNodeBuffer(countryIso)
		}

		p.nodeBuffer = append(p.nodeBuffer, repository.InsertTemporaryNodesParams{
			ID: int64(o.ID),
			Coordinate: pgtype.Point{
				P: pgtype.Vec2{
					X: o.Lon,
					Y: o.Lat,
				},
				Valid: true,
			},
		})
		p.processedNodes += 1

	case *osm.Way:

	case *osm.Relation:
	}
}

func (p *platformNodeParser) saveNodeBuffer(countryIso string) {
	tx, err := p.db.BeginTx(p.ctx, pgx.TxOptions{})
	txRepo := p.repo.WithTx(tx)
	if err != nil {
		logrus.Errorf("Failed to open transaction: %+v", err)
	}

	_, err = txRepo.InsertTemporaryNodes(p.ctx, p.nodeBuffer)
	if err != nil {
		logrus.Errorf("Failed to save temporary nodes: %+v", err)
	}

	ids, err := txRepo.MergeNodesIntoPlatformNodes(p.ctx)
	if err != nil {
		logrus.Errorf("Failed to save temporary nodes: %+v", err)
	}

	printer := message.NewPrinter(language.English)
	logrus.Infof(
		"Saving and merging %d / %d nodes for country %s. Total nodes processed: %s",
		len(ids),
		cap(p.nodeBuffer),
		countryIso,
		printer.Sprint(p.processedNodes),
	)
	p.nodeBuffer = p.nodeBuffer[:0]

	tx.Commit(p.ctx)
}

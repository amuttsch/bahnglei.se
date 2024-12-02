package integrationtest

import (
	"context"
	"os"
	"testing"

	"github.com/amuttsch/bahnglei.se/pkg/config"
	"github.com/amuttsch/bahnglei.se/pkg/repository"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/labstack/gommon/log"
)

type StationData struct {
	Station       repository.Station
	StopPositions []repository.GetStopPositionsForStationRow
	Platforms     []repository.Platform
}

func RepoTestCase(t *testing.T) (*repository.Queries, context.Context) {
	t.Helper()

	conf := config.Read()
	context := context.Background()

	// Do Stuff Here
	dbPool, err := pgxpool.New(context, conf.DatabaseUrl)
	if err != nil {
		log.Errorf("Unable to connect to database: %v\n", err)
		os.Exit(1)
	}

	t.Cleanup(func() {
		dbPool.Close()
	})

	return repository.New(dbPool), context
}

func getStationData(t *testing.T, ctx context.Context, repo *repository.Queries, id int64) *StationData {
	station, err := repo.GetStation(ctx, id)
	if err != nil {
		t.Fatalf("Could not find station with if %d", id)
	}

	stopPositions, err := repo.GetStopPositionsForStation(ctx, pgtype.Int8{
		Int64: id,
		Valid: true,
	})
	if err != nil {
		t.Fatalf("Could not find stop positions for station_id %d", id)
	}

	platforms, err := repo.GetPlatformsForStation(ctx, pgtype.Int8{
		Int64: id,
		Valid: true,
	})
	if err != nil {
		t.Fatalf("Could not find stop positions for station_id %d", id)
	}

	return &StationData{
		Station:       station,
		StopPositions: stopPositions,
		Platforms:     platforms,
	}
}

func GetStopPosition(platform string, stopPositions []repository.GetStopPositionsForStationRow) *repository.GetStopPositionsForStationRow {
	for _, sp := range stopPositions {
		if sp.Platform == platform {
			return &sp
		}
	}
	return nil
}

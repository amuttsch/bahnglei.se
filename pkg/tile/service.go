package tile

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/amuttsch/bahnglei.se/pkg/coordinates"
	"github.com/amuttsch/bahnglei.se/pkg/repository"
	"github.com/sirupsen/logrus"
)

type TileService interface {
	Tile(context context.Context, x, y, z int64, targetLat, targetLng float64) ([]byte, error)
}

type tileService struct {
	repository      *repository.Queries
	tileForstApiKey string
}

func NewTileService(repository *repository.Queries, tileForestApiKey string) *tileService {
	return &tileService{
		repository:      repository,
		tileForstApiKey: tileForestApiKey,
	}
}

func (ts *tileService) Tile(ctx context.Context, x, y, z int64, targetLat, targetLng float64) ([]byte, error) {
	lat := coordinates.Y2lat(y, z)
	lon := coordinates.X2lon(x, z)

	distance := coordinates.Distance(lat, lon, targetLat, targetLng)

	if distance > 10000 {
		return nil, errors.New("tile not in range")
	}

	osmTile, err := ts.repository.GetTile(ctx, repository.GetTileParams{
		X: x,
		Y: y,
		Z: z,
	})
	if err == nil {
		return osmTile.Data, nil
	}

	osmLink := fmt.Sprintf("https://tile.thunderforest.com/transport/%d/%d/%d@2x.png?apikey=%s", z, x, y, ts.tileForstApiKey)
	resp, err := http.Get(osmLink)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	image, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	_, err = ts.repository.CreateTile(ctx, repository.CreateTileParams{
		Z:    z,
		X:    x,
		Y:    y,
		Data: image,
	})
	if err != nil {
		logrus.Errorf("Could not create tile: %+v", err)
		return nil, err
	}

	return image, nil
}

package tile

import (
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/amuttsch/bahnglei.se/pkg/coordinates"
)

type TileService interface {
	Tile(x, y, z int64, targetLat, targetLng float64) ([]byte, error)
}

type tileService struct {
	tileRepo        Repo
	tileForstApiKey string
}

func NewTileService(tileRepo Repo, tileForestApiKey string) *tileService {
	return &tileService{
		tileRepo:        tileRepo,
		tileForstApiKey: tileForestApiKey,
	}
}

func (ts *tileService) Tile(x, y, z int64, targetLat, targetLng float64) ([]byte, error) {
	lat := coordinates.Y2lat(y, z)
	lon := coordinates.X2lon(x, z)

	distance := coordinates.Distance(lat, lon, targetLat, targetLng)

	if distance > 10000 {
		return nil, errors.New("tile not in range")
	}

	osmTile := ts.tileRepo.Get(uint(x), uint(y), uint(z))
	if osmTile != nil {
		return osmTile.Data, nil
	}

	osmLink := fmt.Sprintf("https://tile.thunderforest.com/transport/%d/%d/%d.png?apikey=%s", z, x, y, ts.tileForstApiKey)
	resp, err := http.Get(osmLink)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	image, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	osmTile = &OsmTile{
		Z:    uint(z),
		X:    uint(x),
		Y:    uint(y),
		Data: image,
	}
	ts.tileRepo.Save(osmTile)

	return image, nil
}

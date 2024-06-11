package station

import (
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"

	"github.com/amuttsch/bahnglei.se/pkg/config"
	"github.com/amuttsch/bahnglei.se/pkg/coordinates"
	"github.com/amuttsch/bahnglei.se/pkg/tile"
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
)

type controller struct {
	e      *echo.Echo
	config *config.Config
}

type StationListData struct {
	Stations []Station
}

type StationData struct {
	StationListData
	Station *Station
}

func Http(e *echo.Echo, config *config.Config, stationRepo Repo, tileRepo tile.Repo) *controller {
	e.POST("/station", func(c echo.Context) error {
		data := StationListData{
			Stations: stationRepo.Search(c.FormValue("station")),
		}
		return c.Render(200, "stationlist", data)
	})

	e.GET("/station/:id", func(c echo.Context) error {
		id, _ := strconv.Atoi(c.Param("id"))
		station := stationRepo.Get(uint(id))
		data := StationData{
			StationListData: StationListData{},
			Station:         station,
		}
		return c.Render(200, "station.html", data)
	})

	e.GET("/station/:id/tile/:z/:x/:y", func(c echo.Context) error {
		z, _ := strconv.Atoi(c.Param("z"))
		y, _ := strconv.Atoi(strings.TrimRight(c.Param("y"), ".png"))
		x, _ := strconv.Atoi(c.Param("x"))

		id, _ := strconv.Atoi(c.Param("id"))
		station := stationRepo.Get(uint(id))
		if station == nil {
			return c.NoContent(404)
		}

		lat := coordinates.Y2lat(y, z)
		lon := coordinates.X2lon(x, z)

		distance := coordinates.Distance(lat, lon, station.Lat, station.Lng)

		logrus.Infof("Distance: %fm\n", distance)

		if distance > 10000 {
			return c.NoContent(404)
		}

		osmTile := tileRepo.Get(uint(x), uint(y), uint(z))
		if osmTile != nil {
			return c.Blob(200, "image/png", osmTile.Data)
		}

		osmLink := fmt.Sprintf("https://tile.thunderforest.com/transport/%d/%d/%d.png?apikey=%s", z, x, y, config.ThunderforestConfig.ApiKey)
		resp, err := http.Get(osmLink)
		if err != nil {
			logrus.Error(err)
			return c.NoContent(502)
		}

		defer resp.Body.Close()

		image, err := io.ReadAll(resp.Body)
		if err != nil {
			logrus.Error(err)
			return c.NoContent(502)
		}

		osmTile = &tile.OsmTile{
			Z:    uint(z),
			X:    uint(x),
			Y:    uint(y),
			Data: image,
		}
		tileRepo.Save(osmTile)

		return c.Blob(200, "image/png", image)
	})

	return &controller{
		e:      e,
		config: config,
	}
}

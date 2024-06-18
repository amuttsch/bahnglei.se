package station

import (
	"strconv"
	"strings"

	"github.com/amuttsch/bahnglei.se/pkg/config"
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
	Station *Station
	StationListData
}

func Http(e *echo.Echo, config *config.Config, stationRepo Repo, tileService tile.TileService) *controller {
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

		image, err := tileService.Tile(int64(x), int64(y), int64(z), station.Lat, station.Lng)
		if err != nil {
			logrus.Error(err)
			return c.NoContent(404)
		}

		c.Response().Header().Set("Cache-Control", "max-age=2592000")
		return c.Blob(200, "image/png", image)
	})

	return &controller{
		e:      e,
		config: config,
	}
}

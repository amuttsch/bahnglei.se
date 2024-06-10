package http

import (
	"fmt"
	"io"
	"net/http"
	"strconv"

	"github.com/amuttsch/bahnglei.se/pkg/config"
	stationRepo "github.com/amuttsch/bahnglei.se/pkg/repo/station"
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
)

type controller struct {
	e      *echo.Echo
	config *config.Config
}

type StationListData struct {
	Stations []stationRepo.Station
}

type StationData struct {
    StationListData
    Station *stationRepo.Station
}

func Setup(e *echo.Echo, config *config.Config, stationRepo stationRepo.Repo) *controller {
	e.GET("/", func(c echo.Context) error {
		return c.Render(200, "index.html", nil)
	})

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
            Station: station,
		}
		return c.Render(200, "station.html", data)
	})

    e.GET("/tile/:z/:x/:y", func(c echo.Context) error {
		z := c.Param("z")
		x := c.Param("x")
		y := c.Param("y")
		osmLink := fmt.Sprintf("https://tile.thunderforest.com/transport/%s/%s/%s?apikey=%s", z, x, y, config.ThunderforestConfig.ApiKey)
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
		return c.Blob(200, "image/png", image)
	})

	return &controller{
		e:      e,
		config: config,
	}
}

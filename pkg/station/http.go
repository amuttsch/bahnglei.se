package station

import (
	"strconv"
	"strings"

	"github.com/amuttsch/bahnglei.se/pkg/config"
	"github.com/amuttsch/bahnglei.se/pkg/repository"
	"github.com/amuttsch/bahnglei.se/pkg/tile"
	"github.com/amuttsch/bahnglei.se/templates/components"
	"github.com/amuttsch/bahnglei.se/templates/pages"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/sirupsen/logrus"
)

type controller struct {
	e      *echo.Echo
	config *config.Config
}

type StationData struct {
	components.StationSearchProps
}

func Http(e *echo.Echo, config *config.Config, repo *repository.Queries, tileService tile.TileService) *controller {
	e.POST("/station", func(c echo.Context) error {
		searchString := "%" + c.FormValue("station") + "%"
		stations, _ := repo.SearchStations(c.Request().Context(), searchString)
		stationData := make([]components.StationSearchElement, len(stations))

		for i, station := range stations {
			stationData[i] = components.StationSearchElement{
				ID:   strconv.Itoa(int(station.ID)),
				Name: station.Name,
			}
		}
		data := components.StationSearchProps{
			Stations: stationData,
		}
		stationComponent := components.StationSearchResultList(data)
		return stationComponent.Render(c.Request().Context(), c.Response().Writer)
	})

	e.GET("/station/:id", func(c echo.Context) error {
		id, _ := strconv.Atoi(c.Param("id"))
		station, _ := repo.GetStation(c.Request().Context(), int64(id))
		stationStopPositions, _ := repo.GetStopPositionsForStation(c.Request().Context(), int64(id))
		stationPlatforms, _ := repo.GetPlatformsForStation(c.Request().Context(), int64(id))

		data := pages.StationPageProps{
			StationSearchProps: components.StationSearchProps{
				CSRFToken: c.Get(middleware.DefaultCSRFConfig.ContextKey).(string),
			},
			Station:      station,
			StopPosition: stationStopPositions,
			Platforms:    stationPlatforms,
		}

		stationPage := pages.StationPage(data)
		return stationPage.Render(c.Request().Context(), c.Response().Writer)
	})

	e.POST("/station/:id/report", func(c echo.Context) error {
		id, _ := strconv.Atoi(c.Param("id"))
		report := c.FormValue("report")
		station, err := repo.GetStation(c.Request().Context(), int64(id))
		if err != nil {
			return c.NoContent(404)
		}
		logrus.Warnf("Reported station %s (%d): %s", station.Name, station.ID, report)
		return c.String(200, "Reported station")
	})

	e.GET("/station/:id/tile/:z/:x/:y", func(c echo.Context) error {
		z, _ := strconv.Atoi(c.Param("z"))
		y, _ := strconv.Atoi(strings.TrimRight(c.Param("y"), ".png"))
		x, _ := strconv.Atoi(c.Param("x"))

		id, _ := strconv.Atoi(c.Param("id"))
		station, err := repo.GetStation(c.Request().Context(), int64(id))
		if err != nil {
			return c.NoContent(404)
		}

		image, err := tileService.Tile(c.Request().Context(), int64(x), int64(y), int64(z), station.Lat, station.Lng)
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

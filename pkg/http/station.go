package http

import (
	"slices"
	"strconv"
	"strings"

	"github.com/amuttsch/bahnglei.se/pkg/config"
	"github.com/amuttsch/bahnglei.se/pkg/cookies"
	"github.com/amuttsch/bahnglei.se/pkg/repository"
	"github.com/amuttsch/bahnglei.se/pkg/tile"
	"github.com/amuttsch/bahnglei.se/templates/components"
	"github.com/amuttsch/bahnglei.se/templates/pages"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/sirupsen/logrus"
)

type controllerStation struct {
	e      *echo.Echo
	config *config.Config
}

type StationData struct {
	components.StationSearchProps
}

func Station(e *echo.Echo, config *config.Config, repo *repository.Queries, tileService tile.TileService) *controllerStation {
	e.POST("/station", func(c echo.Context) error {
		searchString := "%" + c.FormValue("station") + "%"
		stations, _ := repo.SearchStations(c.Request().Context(), searchString)
		stationData := make([]components.StationSearchElement, len(stations))

		for i, station := range stations {
			stationData[i] = components.StationSearchElement{
				ID:         strconv.Itoa(int(station.ID)),
				Name:       station.Name,
				CountryIso: station.CountryIsoCode,
				NumTracks:  strconv.FormatInt(station.Tracks, 10),
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
		stationStopPositions, _ := repo.GetStopPositionsForStation(c.Request().Context(),
			pgtype.Int8{
				Int64: int64(id),
				Valid: true,
			},
		)

		stationId := pgtype.Int8{
			Int64: int64(id),
			Valid: true,
		}
		stationPlatforms, _ := repo.GetPlatformsForStation(c.Request().Context(), stationId)

		data := pages.StationPageProps{
			StationSearchProps: components.StationSearchProps{
				CSRFToken: c.Get(middleware.DefaultCSRFConfig.ContextKey).(string),
			},
			Station:      station,
			StopPosition: stationStopPositions,
			Platforms:    stationPlatforms,
		}

		stationPage := pages.StationPage(data)
		err := cookies.SetStationCookie(c, station)
		if err != nil {
			logrus.Error(err)
		}
		return stationPage.Render(c.Request().Context(), c.Response().Writer)
	})

	e.GET("/station/:id/details/:platform", func(c echo.Context) error {
		id, _ := strconv.Atoi(c.Param("id"))
		platform := c.Param("platform")
		stationStopPositions, _ := repo.GetStopPositionsForStationAndPlatform(
			c.Request().Context(), repository.GetStopPositionsForStationAndPlatformParams{
				StationID: pgtype.Int8{
					Int64: int64(id),
					Valid: true,
				},
				Platform: platform,
			})

		neighbors := strings.Split(stationStopPositions.Neighbors, ";")
		neighbors = slices.DeleteFunc(neighbors, func(n string) bool {
			return n == platform || n == ""
		})

		trackDetails := pages.TrackDetails(platform, neighbors)
		return trackDetails.Render(c.Request().Context(), c.Response().Writer)
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

		image, err := tileService.Tile(c.Request().Context(), int64(x), int64(y), int64(z), station.Coordinate.P.Y, station.Coordinate.P.X)
		if err != nil {
			logrus.Error(err)
			return c.NoContent(404)
		}

		c.Response().Header().Set("Cache-Control", "max-age=2592000")
		return c.Blob(200, "image/png", image)
	})

	return &controllerStation{
		e:      e,
		config: config,
	}
}

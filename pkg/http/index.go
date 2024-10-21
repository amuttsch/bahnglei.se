package http

import (
	"strconv"

	"github.com/amuttsch/bahnglei.se/pkg/config"
	"github.com/amuttsch/bahnglei.se/pkg/cookies"
	"github.com/amuttsch/bahnglei.se/pkg/repository"
	"github.com/amuttsch/bahnglei.se/templates/components"
	"github.com/amuttsch/bahnglei.se/templates/pages"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/sirupsen/logrus"
)

type controllerIndex struct {
	e      *echo.Echo
	config *config.Config
}

func Index(e *echo.Echo, config *config.Config, repo *repository.Queries) *controllerIndex {
	e.HEAD("/", func(c echo.Context) error {
		return c.NoContent(204)
	})

	e.GET("/", func(c echo.Context) error {
		stationCount, _ := repo.CountStations(c.Request().Context())
		countryCount, _ := repo.CountCountries(c.Request().Context())
		recentStations, err := cookies.GetRecentStations(c)
		if err != nil {
			logrus.Error(err)
		}

		data := pages.IndexProps{
			CountryCount: strconv.Itoa(int(countryCount)),
			StationCount: strconv.Itoa(int(stationCount)),
			StationSearchProps: components.StationSearchProps{
				CSRFToken: c.Get(middleware.DefaultCSRFConfig.ContextKey).(string),
			},
			RecentStations: recentStations,
		}
		index := pages.IndexPage(data)
		return index.Render(c.Request().Context(), c.Response().Writer)
	})

	e.GET("/about", func(c echo.Context) error {
		countryCount, _ := repo.CountCountries(c.Request().Context())
		countries, _ := repo.GetCountries(c.Request().Context())

		data := pages.AboutPageProps{
			CountryCount: countryCount,
			Countries:    countries,
		}
		index := pages.AboutPage(data)
		return index.Render(c.Request().Context(), c.Response().Writer)
	})

	return &controllerIndex{
		e:      e,
		config: config,
	}
}

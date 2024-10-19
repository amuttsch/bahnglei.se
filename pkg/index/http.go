package index

import (
	"strconv"

	"github.com/amuttsch/bahnglei.se/pkg/config"
	"github.com/amuttsch/bahnglei.se/pkg/repository"
	"github.com/amuttsch/bahnglei.se/templates/components"
	"github.com/amuttsch/bahnglei.se/templates/pages"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type controller struct {
	e      *echo.Echo
	config *config.Config
}

func Http(e *echo.Echo, config *config.Config, repo *repository.Queries) *controller {
	e.HEAD("/", func(c echo.Context) error {
		return c.NoContent(204)
	})

	e.GET("/", func(c echo.Context) error {
		stationCount, _ := repo.CountStations(c.Request().Context())
		countryCount, _ := repo.CountCountries(c.Request().Context())

		data := pages.IndexProps{
			CountryCount: strconv.Itoa(int(countryCount)),
			StationCount: strconv.Itoa(int(stationCount)),
			StationSearchProps: components.StationSearchProps{
				CSRFToken: c.Get(middleware.DefaultCSRFConfig.ContextKey).(string),
			},
		}
		index := pages.IndexPage(data)
		return index.Render(c.Request().Context(), c.Response().Writer)
	})

	return &controller{
		e:      e,
		config: config,
	}
}

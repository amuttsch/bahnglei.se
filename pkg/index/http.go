package index

import (
	"github.com/amuttsch/bahnglei.se/pkg/config"
	"github.com/amuttsch/bahnglei.se/pkg/country"
	"github.com/amuttsch/bahnglei.se/pkg/station"
	"github.com/labstack/echo/v4"
)

type controller struct {
	e      *echo.Echo
	config *config.Config
}

type IndexData struct {
	station.StationListData
	CountryCount int64
	StationCount int64
}

func Http(e *echo.Echo, config *config.Config, countryRepo country.Repo, stationRepo station.Repo) *controller {
	e.GET("/", func(c echo.Context) error {
		stationCount := stationRepo.Count()
		countryCount := countryRepo.Count()

		data := IndexData{
			CountryCount: countryCount,

			StationCount: stationCount,
		}
		return c.Render(200, "index.html", data)
	})

	return &controller{
		e:      e,
		config: config,
	}
}

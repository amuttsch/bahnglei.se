package http

import (
	"net/http"
	"strconv"

	"github.com/amuttsch/bahnglei.se/pkg/config"
	"github.com/amuttsch/bahnglei.se/pkg/cookies"
	"github.com/amuttsch/bahnglei.se/pkg/repository"
	"github.com/amuttsch/bahnglei.se/templates/pages"
	"github.com/labstack/echo/v4"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"github.com/sirupsen/logrus"
)

type controllerIndex struct {
	e      *echo.Echo
	config *config.Config
}

func Index(e *echo.Echo, config *config.Config, repo *repository.Queries, bundle *i18n.Bundle) *controllerIndex {
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
			CountryCount:   strconv.Itoa(int(countryCount)),
			StationCount:   strconv.Itoa(int(stationCount)),
			RecentStations: recentStations,
		}

		lang, _ := cookies.GetLanguage(c)
		accept := c.Request().Header.Get("Accept-Language")
		localizer := i18n.NewLocalizer(bundle, lang, accept)

		index := pages.IndexPage(data, localizer)
		return index.Render(c.Request().Context(), c.Response().Writer)
	})

	e.GET("/about", func(c echo.Context) error {
		countryCount, _ := repo.CountCountries(c.Request().Context())
		countries, _ := repo.GetCountries(c.Request().Context())

		data := pages.AboutPageProps{
			CountryCount: countryCount,
			Countries:    countries,
		}

		lang, _ := cookies.GetLanguage(c)
		accept := c.Request().Header.Get("Accept-Language")
		localizer := i18n.NewLocalizer(bundle, lang, accept)

		about := pages.AboutPage(data, localizer)
		return about.Render(c.Request().Context(), c.Response().Writer)
	})

	e.POST("/lang", func(c echo.Context) error {
		lang := c.FormValue("lang")
		to := c.FormValue("to")
		if to == "" {
			to = "/"
		}

		cookies.SetLanguageCookie(c, lang)

		return c.Redirect(http.StatusFound, to)
	})

	return &controllerIndex{
		e:      e,
		config: config,
	}
}

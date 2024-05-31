package http

import (
	"strconv"

	stationRepo "github.com/amuttsch/bahnglei.se/pkg/repo/station"
	"github.com/labstack/echo/v4"
)

type controller struct {
    e *echo.Echo
}

type StationListData struct {
    Stations []stationRepo.Station
}

func Setup(e *echo.Echo, stationRepo stationRepo.Repo) *controller {
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
        return c.Render(200, "station.html", station)
    })

    return &controller{
        e: e,
    }
}

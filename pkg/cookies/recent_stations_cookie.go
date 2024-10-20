package cookies

import (
	"encoding/base64"
	"encoding/json"
	"net/http"
	"slices"

	"github.com/amuttsch/bahnglei.se/pkg/repository"
	"github.com/labstack/echo/v4"
)

const STATION_COOKIE_NAME = "recentStations"
const MAX_RECENT_STATIONS = 5

type RecentStationCookieItem struct {
	ID   int64  `json:"ID"`
	Name string `json:"Name"`
}

func SetStationCookie(c echo.Context, station repository.Station) error {
	recentStations, err := GetRecentStations(c)
	if err != nil {
		return err
	}

	if len(recentStations) > 0 && recentStations[0].ID == station.ID {
		return nil
	}

	recentStations = slices.Insert(recentStations, 0, RecentStationCookieItem{
		ID:   station.ID,
		Name: station.Name,
	})

	if len(recentStations) > MAX_RECENT_STATIONS {
		recentStations = slices.Delete(recentStations, MAX_RECENT_STATIONS, len(recentStations))
	}

	stationJson, err := json.Marshal(recentStations)
	if err != nil {
		return err
	}

	stationCookie := &http.Cookie{
		Name:  STATION_COOKIE_NAME,
		Value: base64.StdEncoding.EncodeToString(stationJson),
		Path:  "/",
	}
	c.SetCookie(stationCookie)

	return nil
}

func GetRecentStations(c echo.Context) ([]RecentStationCookieItem, error) {
	recentStations := make([]RecentStationCookieItem, 0, MAX_RECENT_STATIONS)
	stationCookie, err := c.Cookie(STATION_COOKIE_NAME)
	if err != nil {
		return recentStations, nil
	}

	stationJson, err := base64.StdEncoding.DecodeString(stationCookie.Value)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(stationJson, &recentStations)
	if err != nil {
		return nil, err
	}

	return recentStations, nil
}

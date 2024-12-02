package osmimporter

import (
	"fmt"

	"github.com/amuttsch/bahnglei.se/pkg/repository"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/sirupsen/logrus"
)

const fetchStopAreasQuery = `
[out:json];
area[name="%s"];
(
  //relation["public_transport"="stop_area"][train=yes](area);
  relation["public_transport"="stop_area"][railway=facility](area);
);

out;
`

const (
	RolePlatform = "platform"
	RoleStop     = "stop"
)

func (o *Overpass) fetchStopAreas(area string, countryIso string) error {
	logrus.Infof("Fetching stop areas")
	resp, err := o.fetch(fmt.Sprintf(fetchStopAreasQuery, area))
	if err != nil {
		return err
	}

	logrus.Infof("Processing stop areas")
	for _, stopArea := range resp.Elements {
		stationId, err := o.findStationInStopArea(stopArea)
		if err != nil {
			continue
		}

		o.updateStationIdForPlatforms(stationId, stopArea)
		o.updateStationIdForStopPositions(stationId, stopArea)
	}

	logrus.Infof("Found %d stop areas", len(resp.Elements))

	return nil
}

func (o *Overpass) findStationInStopArea(element overpassResponseElement) (int64, error) {
	for _, member := range element.Members {
		if member.Role == RolePlatform || member.Role == RoleStop {
			continue
		}

		station, err := o.repo.GetStation(o.ctx, member.Ref)
		if err != nil {
			continue
		}

		return station.ID, nil
	}

	station, err := o.repo.GetStationByName(o.ctx, element.Tags["name"])
	if err != nil {
		return 0, fmt.Errorf("No station found in members")
	}

	return station.ID, nil
}

func (o *Overpass) updateStationIdForStopPositions(stationId int64, element overpassResponseElement) {
	for _, member := range element.Members {
		if member.Role != RoleStop {
			continue
		}

		o.repo.UpdateStopPositionSetStationId(o.ctx, repository.UpdateStopPositionSetStationIdParams{
			ID: member.Ref,
			StationID: pgtype.Int8{
				Int64: stationId,
				Valid: true,
			},
		})
	}
}

func (o *Overpass) updateStationIdForPlatforms(stationId int64, element overpassResponseElement) {
	for _, member := range element.Members {
		if member.Role != RolePlatform {
			continue
		}

		o.repo.UpdatePlatformSetStationId(o.ctx, repository.UpdatePlatformSetStationIdParams{
			ID: member.Ref,
			StationID: pgtype.Int8{
				Int64: stationId,
				Valid: true,
			},
		})
	}
}

package osmimporter

import (
	"fmt"

	"github.com/amuttsch/bahnglei.se/pkg/repository"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/sirupsen/logrus"
)

const fetchStopPositionsQuery = `
[out:json];
area[name="%s"];
(
	node["public_transport"="stop_position"]["railway"~"stop"][train=yes](area);
);

out;
`

func (o *Overpass) fetchStopPositions(area string, countryIso string) error {
	logrus.Infof("Fetching stop positions ")
	resp, err := o.fetch(fmt.Sprintf(fetchStopPositionsQuery, area))
	if err != nil {
		return err
	}

	logrus.Infof("Saving stop positions ")
	for _, stopPosition := range resp.Elements {
		ref := stopPosition.Tags["ref"]
		localRef := stopPosition.Tags["local_ref"]
		platform := ref
		if localRef != "" {
			platform = localRef
		}
		_, err := o.repo.CreateStopPosition(o.ctx, repository.CreateStopPositionParams{
			ID:       stopPosition.ID,
			Platform: platform,
			Coordinate: pgtype.Point{
				P: pgtype.Vec2{
					X: stopPosition.Lon,
					Y: stopPosition.Lat,
				},
				Valid: true,
			},
			CountryIsoCode: countryIso,
		})

		if err != nil {
			logrus.Errorf("Failed to save stop position: %+v\n", err)
			break
		}
	}

	logrus.Infof("Found %d stop positions", len(resp.Elements))

	return nil
}

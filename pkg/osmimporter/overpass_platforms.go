package osmimporter

import (
	"fmt"

	"github.com/amuttsch/bahnglei.se/pkg/repository"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/sirupsen/logrus"
)

const fetchPlatformsQuery = `
[out:json];
area[name="%s"];
(
	wr["public_transport"="platform"][train=yes](area);
	wr["railwaiy"="platform"][train=yes](area);
);

out tags center;
`

func (o *Overpass) fetchPlatforms(area string, countryIso string) error {
	logrus.Infof("Fetching platforms")
	resp, err := o.fetch(fmt.Sprintf(fetchPlatformsQuery, area))
	if err != nil {
		return err
	}

	logrus.Infof("Saving platforms")
	for _, platform := range resp.Elements {

		_, err := o.repo.CreatePlatform(o.ctx, repository.CreatePlatformParams{
			ID:             platform.ID,
			Positions:      platform.Tags["ref"],
			CountryIsoCode: countryIso,
			Coordinate: pgtype.Point{
				P: pgtype.Vec2{
					X: platform.Center.Lon,
					Y: platform.Center.Lat,
				},
				Valid: true,
			},
		})
		if err != nil {
			logrus.Errorf("Failed to save platform: %+v\n", err)
			break
		}
	}

	logrus.Infof("Found %d platforms", len(resp.Elements))

	return nil
}

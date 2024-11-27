package osmimporter

import (
	"fmt"

	"github.com/amuttsch/bahnglei.se/pkg/repository"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/sirupsen/logrus"
)

const fetchStationsQuery = `
[out:json];
area[name="%s"];
(
	node["public_transport"="station"]["railway"~"station|halt"](area);
);

out;
`

func (o *Overpass) fetchStations(area string, countryIso string) error {
	logrus.Infof("Fetching stations")
	resp, err := o.fetch(fmt.Sprintf(fetchStationsQuery, area))
	if err != nil {
		return err
	}

	logrus.Infof("Saving stations")
	for _, station := range resp.Elements {
		_, err := o.repo.CreateStation(o.ctx, repository.CreateStationParams{
			CountryIsoCode: countryIso,
			ID:             station.ID,
			Name:           station.Tags["name"],
			Coordinate: pgtype.Point{
				P: pgtype.Vec2{
					X: station.Lon,
					Y: station.Lat,
				},
				Valid: true,
			},
			Operator:  station.Tags["operator"],
			Wikidata:  station.Tags["wikidata"],
			Wikipedia: station.Tags["wikipedia"],
		})

		if err != nil {
			logrus.Errorf("Failed to save station: %+v\n", err)
			break
		}
	}

	logrus.Infof("Found %d stations", len(resp.Elements))

	return nil
}

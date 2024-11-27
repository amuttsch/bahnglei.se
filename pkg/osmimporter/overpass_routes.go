package osmimporter

import (
	"fmt"

	"github.com/amuttsch/bahnglei.se/pkg/repository"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/sirupsen/logrus"
)

const fetchRoutesQuery = `
[out:json];
area[name="%s"];
(
  node
      ["public_transport"="stop_position"]
      ["railway"="stop"]
      ["train"="yes"]
      (area);

  foreach->.stopPosition(
    relation[type=route][route=train](bn.stopPosition) -> .rels;
    
    (.stopPosition;.rels;);
    out tags;
  );
);

out;
`

func (o *Overpass) fetchRoutes(area string) error {
	logrus.Infof("Fetching routes")
	resp, err := o.fetch(fmt.Sprintf(fetchRoutesQuery, area))
	if err != nil {
		return err
	}

	logrus.Infof("Saving routes")

	countStopPositions := 0
	countRoutes := 0

	currentStopPosition := int64(0)
	for _, route := range resp.Elements {
		if currentStopPosition == 0 && route.Type != ElementTypeNode {
			return fmt.Errorf("First element should be a node, was: %v", route)
		}

		if route.Type == ElementTypeNode {
			currentStopPosition = route.ID
			countStopPositions += 1
			continue
		}

		if route.Type != ElementTypeRelation {
			return fmt.Errorf("Expected a relation, got: %v", route)
		}

		_, err := o.repo.CreateRoute(o.ctx, repository.CreateRouteParams{
			RouteID:        route.ID,
			StopPositionID: currentStopPosition,
			FromStation: pgtype.Text{
				String: route.Tags["from"],
				Valid:  true,
			},
			ToStation: pgtype.Text{
				String: route.Tags["to"],
				Valid:  true,
			},
			Via: pgtype.Text{
				String: route.Tags["via"],
				Valid:  true,
			},
			Ref: pgtype.Text{
				String: route.Tags["ref"],
				Valid:  true,
			},
			Name: pgtype.Text{
				String: route.Tags["name"],
				Valid:  true,
			},
			Service: pgtype.Text{
				String: route.Tags["service"],
				Valid:  true,
			},
			Network: pgtype.Text{
				String: route.Tags["network"],
				Valid:  true,
			},
			Operator: pgtype.Text{
				String: route.Tags["operator"],
				Valid:  true,
			},
		})
		if err != nil {
			logrus.Errorf("Failed to save route: %+v\n", err)
			break
		}

		countRoutes += 1
	}

	logrus.Infof("Found %d routes for %d stop positions", countRoutes, countStopPositions)

	return nil
}

package osmimporter

import (
	"context"
	"time"

	"github.com/amuttsch/bahnglei.se/pkg/config"
	"github.com/amuttsch/bahnglei.se/pkg/repository"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"

	log "github.com/sirupsen/logrus"
)

type OsmImportFlags struct {
	Stations      bool
	StopPositions bool
	Platforms     bool
	StopAreas     bool
	Routes        bool
	ComputeData   bool
}

type osmImporter struct {
	config *config.Config
	repo   *repository.Queries
	db     *pgxpool.Pool
}

func Run(ctx context.Context, config *config.Config, repo *repository.Queries, db *pgxpool.Pool, importflags *OsmImportFlags) error {
	overpass := NewOverpassApi(ctx, db, repo, config.OverpassUrl)

	for _, c := range config.Countries {
		startTime := time.Now()
		log.Infof("Start importing country %s", c.Name)
		country, err := repo.SaveCountry(ctx, repository.SaveCountryParams{
			IsoCode: c.Iso,
			Name:    c.Name,
		})
		if err != nil {
			log.Errorf("Failed to create country %s: %+v", c.Name, err)
			return err
		}

		if importflags.Stations {
			err = overpass.fetchStations(c.Area, c.Iso)
			if err != nil {
				return err
			}
			err = repo.DeleteStationsUpdatedBefore(ctx, repository.DeleteStationsUpdatedBeforeParams{
				CountryIsoCode: c.Iso,
				UpdatedAt: pgtype.Timestamptz{
					Time:  startTime,
					Valid: true,
				},
			})
			if err != nil {
				log.Errorf("Failed to cleanup stations: %+v", err)
				return err
			}
		}

		if importflags.StopPositions {
			err = overpass.fetchStopPositions(c.Area, c.Iso)
			if err != nil {
				return err
			}
			err = repo.DeleteStopPositionsUpdatedBefore(ctx, repository.DeleteStopPositionsUpdatedBeforeParams{
				CountryIsoCode: c.Iso,
				UpdatedAt: pgtype.Timestamptz{
					Time:  startTime,
					Valid: true,
				},
			})
			if err != nil {
				log.Errorf("Failed to cleanup stop positions: %+v", err)
				return err
			}
		}

		if importflags.Platforms {
			err = overpass.fetchPlatforms(c.Area, c.Iso)
			if err != nil {
				return err
			}
			err = repo.DeletePlatformsUpdatedBefore(ctx, repository.DeletePlatformsUpdatedBeforeParams{
				CountryIsoCode: c.Iso,
				UpdatedAt: pgtype.Timestamptz{
					Time:  startTime,
					Valid: true,
				},
			})
			if err != nil {
				log.Errorf("Failed to cleanup platforms: %+v", err)
				return err
			}
		}

		if importflags.StopAreas {
			err = overpass.fetchStopAreas(c.Area, c.Iso)
			if err != nil {
				return err
			}
		}

		if importflags.Routes {
			err = overpass.fetchRoutes(c.Area, c.Iso)
			if err != nil {
				return err
			}
			err = repo.DeleteRoutesUpdatedBefore(ctx, repository.DeleteRoutesUpdatedBeforeParams{
				CountryIsoCode: c.Iso,
				UpdatedAt: pgtype.Timestamptz{
					Time:  startTime,
					Valid: true,
				},
			})
			if err != nil {
				log.Errorf("Failed to cleanup routes: %+v", err)
				return err
			}
		}

		if importflags.ComputeData {
			err = calculateDistances(ctx, repo, country)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func calculateDistances(ctx context.Context, repo *repository.Queries, country repository.Country) error {
	log.Info("Calculating stations for stop positions")
	err := repo.SetStopPositionStationIdToNearestStation(ctx, country.IsoCode)
	if err != nil {
		log.Errorf("Failed to set stations for platforms: %+v", err)
	}

	log.Info("Calculating stations for platforms")
	err = repo.SetPlatformToNearestStation(ctx, country.IsoCode)
	if err != nil {
		log.Errorf("Failed to set stations for platforms: %+v", err)
	}

	log.Info("Setting stop position neighbors")
	err = repo.SetStopPositionNeighbors(ctx, country.IsoCode)
	if err != nil {
		log.Errorf("Failed to set center for platforms: %+v", err)
	}

	log.Info("Setting number of tracks for stations")
	err = repo.SetStationNumberOfTracks(ctx, country.IsoCode)
	if err != nil {
		log.Errorf("Failed to set number of tracks: %+v", err)
	}

	return nil

}

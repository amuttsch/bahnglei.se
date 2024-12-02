package cmd

import (
	"os"

	migrate "github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/spf13/cobra"

	"github.com/amuttsch/bahnglei.se/pkg/config"
	"github.com/amuttsch/bahnglei.se/pkg/osmimporter"
	"github.com/amuttsch/bahnglei.se/pkg/repository"
	log "github.com/sirupsen/logrus"
)

var importAll *bool
var importStations *bool
var importStopPositions *bool
var importPlatforms *bool
var importStopAreas *bool
var importRoutes *bool
var computeData *bool

var importCmd = &cobra.Command{
	Use:   "import",
	Short: "Import OSM railway data",
	Long:  `Load OSM data given from the config file and parse the railway station data.`,
	Run: func(cmd *cobra.Command, args []string) {
		log.Infoln("Starting OSM importer")
		conf := config.Read()
		context := cmd.Context()

		// Do Stuff Here
		dbPool, err := pgxpool.New(context, conf.DatabaseUrl)
		if err != nil {
			log.Errorf("Unable to connect to database: %v\n", err)
			os.Exit(1)
		}
		defer dbPool.Close()

		m, err := migrate.New(
			"file://db/migrations",
			conf.DatabaseUrl)
		if err != nil {
			log.Errorf("Unable to create migration: %v\n", err)
			os.Exit(1)
		}
		if err := m.Up(); err != nil && err != migrate.ErrNoChange {
			log.Errorf("Unable to run migrations: %v\n", err)
			os.Exit(1)
		}

		repo := repository.New(dbPool)

		*importAll = *importAll || cmd.Flags().NFlag() == 0

		importFlags := &osmimporter.OsmImportFlags{
			Stations:      *importAll || *importStations,
			StopPositions: *importAll || *importStopPositions,
			Platforms:     *importAll || *importPlatforms,
			StopAreas:     *importAll || *importStopAreas,
			Routes:        *importAll || *importRoutes,
			ComputeData:   *importAll || *computeData,
		}
		err = osmimporter.Run(context, conf, repo, dbPool, importFlags)
		if err != nil {
			log.Errorf("Import failed: %v\n", err)
			os.Exit(1)
		}

		log.Infoln("Finished importing OSM data")
	},
}

func init() {
	importAll = importCmd.Flags().Bool("all", false, "run all imports")
	importStations = importCmd.Flags().Bool("stations", false, "import stations")
	importStopPositions = importCmd.Flags().Bool("stopPositions", false, "import stop positions")
	importPlatforms = importCmd.Flags().Bool("platforms", false, "import platforms")
	importStopAreas = importCmd.Flags().Bool("stopAreas", false, "import stop areas")
	importRoutes = importCmd.Flags().Bool("routes", false, "import routes")
	computeData = importCmd.Flags().Bool("compute", true, "compute distances, number of tracks, etc. from data")
	rootCmd.AddCommand(importCmd)
}

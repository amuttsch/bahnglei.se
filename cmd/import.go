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
			log.Errorf("Unable to run migrations: %v\n", err)
			os.Exit(1)
		}
		m.Up()

		repo := repository.New(dbPool)

		osmImporter := osmimporter.New(conf, repo)
		osmImporter.Import(context)

		log.Infoln("Finished importing OSM data")
	},
}

func init() {
	rootCmd.AddCommand(importCmd)
}

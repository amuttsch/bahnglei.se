package cmd

import (
	"os"

	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/spf13/cobra"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/amuttsch/bahnglei.se/pkg/config"
	"github.com/amuttsch/bahnglei.se/pkg/osmimporter"
	"github.com/amuttsch/bahnglei.se/pkg/repo/country"
	importerRepo "github.com/amuttsch/bahnglei.se/pkg/repo/importer"
	stationRepo "github.com/amuttsch/bahnglei.se/pkg/repo/station"
	log "github.com/sirupsen/logrus"
)

var importCmd = &cobra.Command{
	Use:   "import",
	Short: "Import OSM railway data",
	Long:  `Load OSM data given from the config file and parse the railway station data.`,
	Run: func(cmd *cobra.Command, args []string) {
        log.Infoln("Starting OSM importer")
        conf := config.Read()

		// Do Stuff Here
		db, err := gorm.Open(postgres.Open(conf.DatabaseUrl))
		if err != nil {
			log.Errorf("Unable to connect to database: %v\n", err)
			os.Exit(1)
		}

        context := cmd.Context()
        countryRepo := country.New(db, context)
        importRepo := importerRepo.New(db, context)
        stationRepo := stationRepo.New(db, context)

        osmImporter := osmimporter.New(conf, countryRepo, importRepo, stationRepo)
        osmImporter.Import()

        log.Infoln("Finished importing OSM data")
	},
}

func init() {
	rootCmd.AddCommand(importCmd)
}

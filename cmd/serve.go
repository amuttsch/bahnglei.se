package cmd

import (
	"errors"
	"net/http"
	"os"

	"github.com/amuttsch/bahnglei.se/pkg/config"
	bahnHttp "github.com/amuttsch/bahnglei.se/pkg/http"
	"github.com/amuttsch/bahnglei.se/pkg/repository"
	"github.com/amuttsch/bahnglei.se/pkg/tile"
	"github.com/benbjohnson/hashfs"
	migrate "github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/labstack/echo-contrib/echoprometheus"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var AssetFS *hashfs.FS

var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Run the server",
	Run: func(cmd *cobra.Command, args []string) {
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
		tileService := tile.NewTileService(repo, conf.ThunderforestConfig.ApiKey)

		e := echo.New()
		e.Use(middleware.Logger())
		e.Use(echoprometheus.NewMiddleware("bahngleise"))
		e.Use(middleware.Gzip())
		// e.Use(middleware.RateLimiter(middleware.NewRateLimiterMemoryStore(20)))
		e.Use(middleware.CSRFWithConfig(middleware.CSRFConfig{
			TokenLookup:    "form:_csrf",
			CookieSameSite: http.SameSiteStrictMode,
		}))

		// Add default cache for non hashfs files
		cacheControlHeaderMiddleware := func(next echo.HandlerFunc) echo.HandlerFunc {
			return func(c echo.Context) error {
				c.Response().Header().Add("Cache-Control", "max-age=604800")
				return next(c)
			}
		}
		e.Add(http.MethodGet,
			"/assets/*",
			echo.WrapHandler(hashfs.FileServer(AssetFS)),
			cacheControlHeaderMiddleware,
		)

		bahnHttp.Index(e, conf, repo)
		bahnHttp.Station(e, conf, repo, tileService)

		go func() {
			metrics := echo.New()                                // this Echo will run on separate port 8081
			metrics.GET("/metrics", echoprometheus.NewHandler()) // adds route to serve gathered metrics
			if err := metrics.Start(":9091"); err != nil && !errors.Is(err, http.ErrServerClosed) {
				log.Fatal(err)
			}
		}()

		e.Logger.Fatal(e.Start(":1323"))
	},
}

func init() {
	rootCmd.AddCommand(serveCmd)
}

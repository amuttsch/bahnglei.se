package cmd

import (
	"crypto/subtle"
	"errors"
	"io"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"strings"
	"text/template"

	"github.com/amuttsch/bahnglei.se/pkg/config"
	"github.com/amuttsch/bahnglei.se/pkg/country"
	"github.com/amuttsch/bahnglei.se/pkg/index"
	"github.com/amuttsch/bahnglei.se/pkg/station"
	"github.com/amuttsch/bahnglei.se/pkg/tile"
	"github.com/benbjohnson/hashfs"
	"github.com/labstack/echo-contrib/echoprometheus"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"golang.org/x/exp/slices"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var AssetFS *hashfs.FS

type Template struct {
	tmpl *template.Template
}

func newTemplate() *Template {
	funcMap := template.FuncMap{
		"sortStopPositions": func(slice []station.StopPosition) []station.StopPosition {
			slices.SortFunc(slice, func(i station.StopPosition, j station.StopPosition) int {
				r := regexp.MustCompile("[^0-9]")
				iPlatform, _ := strconv.Atoi(r.ReplaceAllString(i.Platform, ""))
				jPlatform, _ := strconv.Atoi(r.ReplaceAllString(j.Platform, ""))
				return iPlatform - jPlatform
			})
			return slice
		},
		"splitString": func(s string, sep string) []string {
			return strings.Split(s, sep)
		},
		"replaceSpace": func(s string) string {
			return strings.ReplaceAll(s, " ", "-")
		},
		"asset": func(name string) string {
			log.Info(AssetFS.HashName(name))

			return "/assets/" + AssetFS.HashName(name)
		},
	}
	return &Template{
		tmpl: template.Must(template.New("").Funcs(funcMap).ParseFS(AssetFS, "views/*.html")).Funcs(funcMap),
	}
}

func (t *Template) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	if strings.HasSuffix(name, ".html") {
		tmpl := template.Must(t.tmpl.Clone())
		tmpl = template.Must(tmpl.ParseFS(AssetFS, "views/"+name))
		return tmpl.ExecuteTemplate(w, name, data)
	}
	return t.tmpl.ExecuteTemplate(w, name, data)
}

var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Run the server",
	Run: func(cmd *cobra.Command, args []string) {
		conf := config.Read()

		// Do Stuff Here
		db, err := gorm.Open(postgres.Open(conf.DatabaseUrl))
		if err != nil {
			log.Errorf("Unable to connect to database: %v\n", err)
			os.Exit(1)
		}

		context := cmd.Context()
		countryRepo := country.NewRepo(db, context)
		stationRepo := station.NewRepo(db, context)
		tileRepo := tile.NewRepo(db, context)
		tileService := tile.NewTileService(tileRepo, conf.ThunderforestConfig.ApiKey)

		e := echo.New()
		e.Renderer = newTemplate()
		e.Use(middleware.Logger())
		e.Use(echoprometheus.NewMiddleware("bahngleise"))
		e.Use(middleware.Gzip())
		e.Use(middleware.RateLimiter(middleware.NewRateLimiterMemoryStore(20)))

		e.Use(middleware.BasicAuth(func(username, password string, c echo.Context) (bool, error) {
			// Be careful to use constant time comparison to prevent timing attacks
			if subtle.ConstantTimeCompare([]byte(username), []byte("joe")) == 1 &&
				subtle.ConstantTimeCompare([]byte(password), []byte("secret")) == 1 {
				return true, nil
			}
			return false, nil
		}))

		e.GET("/assets/*", echo.WrapHandler(http.StripPrefix("/assets", hashfs.FileServer(AssetFS))))

		index.Http(e, conf, countryRepo, stationRepo)
		station.Http(e, conf, stationRepo, tileService)

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

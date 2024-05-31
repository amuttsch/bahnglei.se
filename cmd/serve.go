package cmd

import (
	"crypto/subtle"
	"io"
	"os"
	"strings"
	"text/template"

	"github.com/amuttsch/bahnglei.se/pkg/config"
	"github.com/amuttsch/bahnglei.se/pkg/http"
	stationRepo "github.com/amuttsch/bahnglei.se/pkg/repo/station"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Template struct {
	tmpl *template.Template
}

func newTemplate() *Template {
	return &Template{
		tmpl: template.Must(template.ParseGlob("views/*.html")),
	}
}

func (t *Template) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	if strings.HasSuffix(name, ".html") {
		tmpl := template.Must(t.tmpl.Clone())
		tmpl = template.Must(tmpl.ParseGlob("views/"+name))
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
		stationRepo := stationRepo.New(db, context)

		e := echo.New()
		e.Renderer = newTemplate()
		e.Use(middleware.Logger())
		e.Static("/images", "images")
		e.Static("/css", "css")

		e.Use(middleware.BasicAuth(func(username, password string, c echo.Context) (bool, error) {
			// Be careful to use constant time comparison to prevent timing attacks
			if subtle.ConstantTimeCompare([]byte(username), []byte("joe")) == 1 &&
				subtle.ConstantTimeCompare([]byte(password), []byte("secret")) == 1 {
				return true, nil
			}
			return false, nil
		}))

		http.Setup(e, stationRepo)

		e.Logger.Fatal(e.Start(":1323"))
	},
}

func init() {
	rootCmd.AddCommand(serveCmd)
}

package http

import "github.com/labstack/echo/v4"

type controller struct {
    e *echo.Echo
}

func Setup(e *echo.Echo) *controller {
    e.GET("/", func(c echo.Context) error {
        return c.Render(200, "index.html", nil)
    })

    return &controller{
        e: e,
    }
}

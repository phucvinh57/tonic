package echoAdapter

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/TickLabVN/tonic/core"
	"github.com/TickLabVN/tonic/core/docs"
	"github.com/labstack/echo/v4"
)

func UIHandle(e *echo.Echo, spec *docs.OpenApi, path string, ui ...core.UI) {
	New(spec).UIHandle(e, path, ui...)
}

func (a *Adapter) UIHandle(e *echo.Echo, path string, ui ...core.UI) {
	spec := a.spec
	if strings.HasSuffix(path, "/") {
		path = path[:len(path)-1]
	} else {
		e.GET(path, func(c echo.Context) error {
			return c.Redirect(http.StatusMovedPermanently, fmt.Sprintf("%s/index.html", path))
		})
	}
	selected := core.SwaggerUI
	if len(ui) > 0 {
		selected = ui[0]
	}
	swaggerPath := fmt.Sprintf("%s.json", path)
	e.GET(swaggerPath, echo.WrapHandler(core.JsonHttpHandler(spec)))
	e.GET(fmt.Sprintf("%s/*", path), echo.WrapHandler(core.UIHandler(selected, swaggerPath)))
}

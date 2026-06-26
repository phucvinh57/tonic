package fiberAdapter

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/TickLabVN/tonic/core"
	"github.com/TickLabVN/tonic/core/docs"
	"github.com/gofiber/fiber/v3"
)

func UIHandle(f *fiber.App, spec *docs.OpenApi, path string, ui ...core.UI) {
	New(spec).UIHandle(f, path, ui...)
}

func (a *Adapter) UIHandle(f *fiber.App, path string, ui ...core.UI) {
	if strings.HasSuffix(path, "/") {
		path = path[:len(path)-1]
	} else {
		f.Get(path, func(c fiber.Ctx) error {
			return c.Redirect().Status(http.StatusMovedPermanently).To(fmt.Sprintf("%s/index.html", path))
		})
	}
	selected := core.SwaggerUI
	if len(ui) > 0 {
		selected = ui[0]
	}
	swaggerPath := fmt.Sprintf("%s.json", path)
	spec := a.spec
	f.Get(swaggerPath, core.JsonHttpHandler(spec))
	f.Get(fmt.Sprintf("%s/*", path), core.UIHandler(selected, swaggerPath))
}

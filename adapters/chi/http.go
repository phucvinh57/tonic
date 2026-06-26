package chiAdapter

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/TickLabVN/tonic/core"
	"github.com/TickLabVN/tonic/core/docs"
	"github.com/go-chi/chi/v5"
)

func UIHandle(r chi.Router, spec *docs.OpenApi, path string, ui ...core.UI) {
	New(spec).UIHandle(r, path, ui...)
}

func (a *Adapter) UIHandle(r chi.Router, path string, ui ...core.UI) {
	if strings.HasSuffix(path, "/") {
		path = path[:len(path)-1]
	} else {
		r.Get(path, func(w http.ResponseWriter, req *http.Request) {
			http.Redirect(w, req, fmt.Sprintf("%s/index.html", path), http.StatusMovedPermanently)
		})
	}
	selected := core.SwaggerUI
	if len(ui) > 0 {
		selected = ui[0]
	}
	swaggerPath := fmt.Sprintf("%s.json", path)
	spec := a.spec
	r.Method(http.MethodGet, swaggerPath, core.JsonHttpHandler(spec))
	r.Method(http.MethodGet, fmt.Sprintf("%s/*", path), core.UIHandler(selected, swaggerPath))
}

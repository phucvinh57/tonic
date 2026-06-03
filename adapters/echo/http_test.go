package echoAdapter

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/TickLabVN/tonic/core/docs"
	"github.com/labstack/echo/v4"
)

func TestUIHandleRedirectsBasePath(t *testing.T) {
	e := echo.New()
	adapter := New(&docs.OpenApi{
		OpenAPI: docs.VERSION,
		Info: docs.InfoObject{
			Title:   "test",
			Version: "1.0.0",
		},
	})
	adapter.UIHandle(e, "/docs")

	req := httptest.NewRequest(http.MethodGet, "/docs", nil)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	if rec.Code != http.StatusMovedPermanently {
		t.Fatalf("expected redirect status, got %d", rec.Code)
	}
	if location := rec.Header().Get("Location"); location != "/docs/index.html" {
		t.Fatalf("expected redirect to /docs/index.html, got %q", location)
	}
}

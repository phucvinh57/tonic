package chiAdapter

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/TickLabVN/tonic/core/docs"
	"github.com/go-chi/chi/v5"
)

func TestUIHandleRedirectsBasePath(t *testing.T) {
	router := chi.NewRouter()
	adapter := New(&docs.OpenApi{
		OpenAPI: docs.VERSION,
		Info: docs.InfoObject{
			Title:   "test",
			Version: "1.0.0",
		},
	})
	adapter.UIHandle(router, "/docs")

	req := httptest.NewRequest(http.MethodGet, "/docs", nil)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusMovedPermanently {
		t.Fatalf("expected redirect status, got %d", rec.Code)
	}
	if location := rec.Header().Get("Location"); location != "/docs/index.html" {
		t.Fatalf("expected redirect to /docs/index.html, got %q", location)
	}
}

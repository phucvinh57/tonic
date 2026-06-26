package fiberAdapter

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/TickLabVN/tonic/core/docs"
	"github.com/gofiber/fiber/v3"
)

func TestUIHandleRedirectsBasePath(t *testing.T) {
	app := fiber.New()
	adapter := New(&docs.OpenApi{
		OpenAPI: docs.VERSION,
		Info: docs.InfoObject{
			Title:   "test",
			Version: "1.0.0",
		},
	})
	adapter.UIHandle(app, "/docs")

	resp, err := app.Test(httptest.NewRequest(http.MethodGet, "/docs", nil))
	if err != nil {
		t.Fatalf("expected request to succeed: %v", err)
	}

	if resp.StatusCode != http.StatusMovedPermanently {
		t.Fatalf("expected redirect status, got %d", resp.StatusCode)
	}
	if location := resp.Header.Get("Location"); location != "/docs/index.html" {
		t.Fatalf("expected redirect to /docs/index.html, got %q", location)
	}
}

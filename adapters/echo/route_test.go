package echoAdapter

import (
	"net/http"
	"testing"

	"github.com/TickLabVN/tonic/core/docs"
	"github.com/labstack/echo/v4"
)

type echoRequest struct {
	ID string `param:"id" validate:"required"`
}

type echoResponse struct {
	ID string `json:"id"`
}

func echoHandler(c echo.Context) error {
	return c.NoContent(http.StatusOK)
}

func TestAddRouteNormalizesPathAndAddsDescriptions(t *testing.T) {
	e := echo.New()
	spec := &docs.OpenApi{
		OpenAPI: docs.VERSION,
		Info: docs.InfoObject{
			Title:   "Echo",
			Version: "1.0.0",
		},
	}

	adapter := New(spec)
	For[echoRequest, echoResponse](adapter).AddRoute(e.GET("/users/:id/", echoHandler))

	pathItem, ok := spec.Paths["/users/{id}"]
	if !ok {
		t.Fatalf("expected normalized path to be present")
	}
	if pathItem.Get == nil {
		t.Fatalf("expected GET operation")
	}
	if pathItem.Get.Responses["200"].ResponseObject.Description != "OK" {
		t.Fatalf("expected generated response description OK")
	}
}

package fiberAdapter

import (
	"bytes"
	"log"
	"net/http"
	"testing"

	"github.com/TickLabVN/tonic/core/docs"
	"github.com/gofiber/fiber/v3"
)

type fiberRequest struct {
	ID        string `uri:"id" validate:"required"`
	Q         string `query:"q"`
	RequestID string `header:"x-request-id"`
	Body      string `json:"body"`
}

type fiberResponse struct {
	ID string `json:"id"`
}

func fiberHandler(c fiber.Ctx) error {
	return c.SendStatus(http.StatusOK)
}

func newTestSpec() *docs.OpenApi {
	return &docs.OpenApi{
		OpenAPI: docs.VERSION,
		Info: docs.InfoObject{
			Title:   "Fiber",
			Version: "1.0.0",
		},
	}
}

func TestMethodHelperDocumentsRequestAndNormalizesPath(t *testing.T) {
	app := fiber.New()
	group := app.Group("/api/v1/users")
	spec := newTestSpec()

	adapter := New(spec)
	For[fiberRequest, fiberResponse](adapter).POST(group, "/:id/", fiberHandler)

	pathItem, ok := spec.Paths["/api/v1/users/{id}"]
	if !ok {
		t.Fatalf("expected normalized path to be present")
	}
	if pathItem.Post == nil {
		t.Fatalf("expected POST operation")
	}
	if pathItem.Post.RequestBody == nil {
		t.Fatalf("expected request body to be documented")
	}
	if len(pathItem.Post.Parameters) < 3 {
		t.Fatalf("expected path, query, and header parameters to be documented")
	}
	if pathItem.Post.Responses["200"].ResponseObject.Description != "OK" {
		t.Fatalf("expected generated response description OK")
	}
}

func TestMethodHelperAppliesOperationOption(t *testing.T) {
	app := fiber.New()
	spec := newTestSpec()

	adapter := New(spec)
	For[fiberRequest, fiberResponse](adapter).GET(app, "/users/:id", fiberHandler, WithOperation(docs.OperationObject{
		Summary: "Get user",
	}))

	pathItem := spec.Paths["/users/{id}"]
	if pathItem.Get == nil || pathItem.Get.Summary != "Get user" {
		t.Fatalf("expected operation option to override generated operation")
	}
	if pathItem.Get.RequestBody != nil {
		t.Fatalf("expected GET request body to be omitted")
	}
}

func TestMethodHelperRejectsInvalidMixedArgs(t *testing.T) {
	app := fiber.New()
	spec := newTestSpec()

	var logs bytes.Buffer
	previousWriter := log.Writer()
	log.SetOutput(&logs)
	defer log.SetOutput(previousWriter)

	adapter := New(spec)
	For[fiberRequest, fiberResponse](adapter).GET(app, "/users/:id", fiberHandler, "invalid")

	if _, ok := spec.Paths["/users/{id}"]; ok {
		t.Fatalf("expected docs to be skipped for invalid helper args")
	}
	if logs.Len() == 0 {
		t.Fatalf("expected warning log for invalid helper args")
	}
}

func TestMethodHelperRejectsMissingHandlers(t *testing.T) {
	app := fiber.New()
	spec := newTestSpec()

	var logs bytes.Buffer
	previousWriter := log.Writer()
	log.SetOutput(&logs)
	defer log.SetOutput(previousWriter)

	adapter := New(spec)
	For[fiberRequest, fiberResponse](adapter).GET(app, "/users/:id", WithOperation(docs.OperationObject{
		Summary: "No handler",
	}))

	if _, ok := spec.Paths["/users/{id}"]; ok {
		t.Fatalf("expected docs to be skipped when no handlers are provided")
	}
	if logs.Len() == 0 {
		t.Fatalf("expected warning log for missing handlers")
	}
}

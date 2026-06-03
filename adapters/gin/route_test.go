package ginAdapter

import (
	"bytes"
	"log"
	"net/http"
	"testing"

	"github.com/TickLabVN/tonic/core/docs"
	"github.com/gin-gonic/gin"
)

type ginRequest struct {
	ID   string `uri:"id" binding:"required"`
	Q    string `form:"q"`
	Body string `json:"body"`
}

type ginResponse struct {
	ID string `json:"id"`
}

func ginHandler(c *gin.Context) {
	c.Status(http.StatusOK)
}

func TestAddRouteDocumentsRequestAndNormalizesPath(t *testing.T) {
	gin.SetMode(gin.ReleaseMode)
	engine := gin.New()
	group := engine.Group("/users")
	spec := &docs.OpenApi{
		OpenAPI: docs.VERSION,
		Info: docs.InfoObject{
			Title:   "Gin",
			Version: "1.0.0",
		},
	}

	adapter := New(spec)
	For[ginRequest, ginResponse](adapter).POST(group, "/:id/", ginHandler)

	pathItem, ok := spec.Paths["/users/{id}"]
	if !ok {
		t.Fatalf("expected normalized path to be present")
	}
	if pathItem.Post == nil {
		t.Fatalf("expected POST operation")
	}
	if pathItem.Post.RequestBody == nil {
		t.Fatalf("expected request body to be documented")
	}
	if len(pathItem.Post.Parameters) < 2 {
		t.Fatalf("expected path and query parameters to be documented")
	}
	if pathItem.Post.Responses["200"].ResponseObject.Description != "OK" {
		t.Fatalf("expected generated response description OK")
	}
}

func TestMethodHelperAppliesOperationOption(t *testing.T) {
	gin.SetMode(gin.ReleaseMode)
	engine := gin.New()
	spec := &docs.OpenApi{
		OpenAPI: docs.VERSION,
		Info: docs.InfoObject{
			Title:   "Gin",
			Version: "1.0.0",
		},
	}

	adapter := New(spec)
	For[ginRequest, ginResponse](adapter).GET(engine, "/users/:id", ginHandler, WithOperation(docs.OperationObject{
		Summary: "Get user",
	}))

	pathItem := spec.Paths["/users/{id}"]
	if pathItem.Get == nil || pathItem.Get.Summary != "Get user" {
		t.Fatalf("expected operation option to override generated operation")
	}
}

func TestMethodHelperRejectsInvalidMixedArgs(t *testing.T) {
	gin.SetMode(gin.ReleaseMode)
	engine := gin.New()
	spec := &docs.OpenApi{
		OpenAPI: docs.VERSION,
		Info: docs.InfoObject{
			Title:   "Gin",
			Version: "1.0.0",
		},
	}

	var logs bytes.Buffer
	previousWriter := log.Writer()
	log.SetOutput(&logs)
	defer log.SetOutput(previousWriter)

	adapter := New(spec)
	For[ginRequest, ginResponse](adapter).GET(engine, "/users/:id", ginHandler, "invalid")

	if _, ok := spec.Paths["/users/{id}"]; ok {
		t.Fatalf("expected docs to be skipped for invalid helper args")
	}
	if logs.Len() == 0 {
		t.Fatalf("expected warning log for invalid helper args")
	}
}

func TestMethodHelperRejectsMissingHandlers(t *testing.T) {
	gin.SetMode(gin.ReleaseMode)
	engine := gin.New()
	spec := &docs.OpenApi{
		OpenAPI: docs.VERSION,
		Info: docs.InfoObject{
			Title:   "Gin",
			Version: "1.0.0",
		},
	}

	var logs bytes.Buffer
	previousWriter := log.Writer()
	log.SetOutput(&logs)
	defer log.SetOutput(previousWriter)

	adapter := New(spec)
	For[ginRequest, ginResponse](adapter).GET(engine, "/users/:id", WithOperation(docs.OperationObject{
		Summary: "No handler",
	}))

	if _, ok := spec.Paths["/users/{id}"]; ok {
		t.Fatalf("expected docs to be skipped when no handlers are provided")
	}
	if logs.Len() == 0 {
		t.Fatalf("expected warning log for missing handlers")
	}
}

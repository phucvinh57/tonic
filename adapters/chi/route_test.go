package chiAdapter

import (
	"bytes"
	"log"
	"net/http"
	"testing"

	"github.com/TickLabVN/tonic/core/docs"
	"github.com/go-chi/chi/v5"
)

type chiRequest struct {
	ID        string `path:"id" validate:"required"`
	Q         string `query:"q"`
	RequestID string `header:"x-request-id"`
	Body      string `json:"body"`
}

type chiResponse struct {
	ID string `json:"id"`
}

func chiHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}

func chiMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		next.ServeHTTP(w, r)
	})
}

func newTestSpec() *docs.OpenApi {
	return &docs.OpenApi{
		OpenAPI: docs.VERSION,
		Info: docs.InfoObject{
			Title:   "Chi",
			Version: "1.0.0",
		},
	}
}

func TestMethodHelperDocumentsRequestAndNormalizesPath(t *testing.T) {
	router := chi.NewRouter()
	spec := newTestSpec()

	adapter := New(spec)
	wrapped := adapter.Wrap(router)
	wrapped.Route("/api/v1/users", func(r chi.Router) {
		For[chiRequest, chiResponse](adapter).POST(r, "/{id}/", chiMiddleware, chiHandler)
	})

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
	router := chi.NewRouter()
	spec := newTestSpec()

	adapter := New(spec)
	For[chiRequest, chiResponse](adapter).GET(router, "/users/{id}", http.HandlerFunc(chiHandler), WithOperation(docs.OperationObject{
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
	router := chi.NewRouter()
	spec := newTestSpec()

	var logs bytes.Buffer
	previousWriter := log.Writer()
	log.SetOutput(&logs)
	defer log.SetOutput(previousWriter)

	adapter := New(spec)
	For[chiRequest, chiResponse](adapter).GET(router, "/users/{id}", chiHandler, "invalid")

	if _, ok := spec.Paths["/users/{id}"]; ok {
		t.Fatalf("expected docs to be skipped for invalid helper args")
	}
	if logs.Len() == 0 {
		t.Fatalf("expected warning log for invalid helper args")
	}
}

func TestMethodHelperRejectsMissingHandlers(t *testing.T) {
	router := chi.NewRouter()
	spec := newTestSpec()

	var logs bytes.Buffer
	previousWriter := log.Writer()
	log.SetOutput(&logs)
	defer log.SetOutput(previousWriter)

	adapter := New(spec)
	For[chiRequest, chiResponse](adapter).GET(router, "/users/{id}", WithOperation(docs.OperationObject{
		Summary: "No handler",
	}))

	if _, ok := spec.Paths["/users/{id}"]; ok {
		t.Fatalf("expected docs to be skipped when no handlers are provided")
	}
	if logs.Len() == 0 {
		t.Fatalf("expected warning log for missing handler")
	}
}

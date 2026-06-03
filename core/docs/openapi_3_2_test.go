package docs

import (
	"encoding/json"
	"testing"

	"github.com/go-playground/validator/v10"
)

func TestOpenAPI32FieldsMarshal(t *testing.T) {
	doc := OpenApi{
		OpenAPI: VERSION,
		Self:    "https://example.com/openapi.json",
		Info: InfoObject{
			Title:   "Test",
			Version: "1.0.0",
		},
		Components: ComponentsObject{
			MediaTypes: map[string]MediaTypeOrReference{
				"problem": {
					ReferenceObject: &ReferenceObject{Ref: "#/components/mediaTypes/problem"},
				},
			},
		},
	}

	raw, err := json.Marshal(doc)
	if err != nil {
		t.Fatalf("marshal openapi: %v", err)
	}

	var got map[string]any
	if err := json.Unmarshal(raw, &got); err != nil {
		t.Fatalf("unmarshal openapi: %v", err)
	}

	if got["$self"] != "https://example.com/openapi.json" {
		t.Fatalf("expected $self field to be marshaled, got %v", got["$self"])
	}

	components := got["components"].(map[string]any)
	if _, ok := components["mediaTypes"]; !ok {
		t.Fatalf("expected components.mediaTypes to be marshaled")
	}
}

func TestExampleObjectNewFieldsMarshal(t *testing.T) {
	example := ExampleObject{
		Summary:         "tuple",
		DataValue:       []any{"first", 2},
		SerializedValue: "",
	}

	raw, err := json.Marshal(example)
	if err != nil {
		t.Fatalf("marshal example: %v", err)
	}

	var got map[string]any
	if err := json.Unmarshal(raw, &got); err != nil {
		t.Fatalf("unmarshal example: %v", err)
	}

	if _, ok := got["dataValue"]; !ok {
		t.Fatalf("expected dataValue field to be marshaled")
	}
}

func TestResponseObjectSummaryAndJSONResponse(t *testing.T) {
	resp := JSONResponse(200, "#/components/schemas/User_json")
	if resp.ResponseObject == nil {
		t.Fatalf("expected response object")
	}
	if resp.ResponseObject.Description != "OK" {
		t.Fatalf("expected default description OK, got %q", resp.ResponseObject.Description)
	}

	raw, err := json.Marshal(ResponseObject{
		Description: "ok",
		Summary:     "short",
	})
	if err != nil {
		t.Fatalf("marshal response: %v", err)
	}

	var got map[string]any
	if err := json.Unmarshal(raw, &got); err != nil {
		t.Fatalf("unmarshal response: %v", err)
	}

	if got["summary"] != "short" {
		t.Fatalf("expected summary to be marshaled")
	}
}

func TestCallbackObjectSerializesPathItems(t *testing.T) {
	callback := CallbackObject{
		"{$request.query.queryUrl}": {
			PathItemObject: &PathItemObject{
				Post: &OperationObject{
					OperationId: "notify",
				},
			},
		},
	}

	raw, err := json.Marshal(callback)
	if err != nil {
		t.Fatalf("marshal callback: %v", err)
	}

	var got map[string]any
	if err := json.Unmarshal(raw, &got); err != nil {
		t.Fatalf("unmarshal callback: %v", err)
	}

	entry := got["{$request.query.queryUrl}"].(map[string]any)
	if _, ok := entry["post"]; !ok {
		t.Fatalf("expected callback entry to contain a path item operation")
	}
}

func TestValidatorAcceptsOpenAPI32Fields(t *testing.T) {
	validate := validator.New()

	if err := validate.Struct(ServerObject{Url: "/"}); err != nil {
		t.Fatalf("expected relative server url to validate, got %v", err)
	}

	if err := validate.Struct(ParameterObject{
		Name:   "filter",
		In:     "querystring",
		Schema: &SchemaOrReference{SchemaObject: &SchemaObject{Type: "string"}},
	}); err != nil {
		t.Fatalf("expected querystring parameter to validate, got %v", err)
	}
}

func TestPathItemQueryMarshals(t *testing.T) {
	item := PathItemObject{
		Query: &OperationObject{
			OperationId: "query_users",
		},
	}

	raw, err := json.Marshal(item)
	if err != nil {
		t.Fatalf("marshal path item: %v", err)
	}

	var got map[string]any
	if err := json.Unmarshal(raw, &got); err != nil {
		t.Fatalf("unmarshal path item: %v", err)
	}

	if _, ok := got["query"]; !ok {
		t.Fatalf("expected query operation to be marshaled")
	}
}

func TestContentCanMarshalMediaTypeReference(t *testing.T) {
	resp := ResponseObject{
		Description: "problem",
		Content: map[string]MediaTypeOrReference{
			"application/problem+json": {
				ReferenceObject: &ReferenceObject{Ref: "#/components/mediaTypes/problem"},
			},
		},
	}

	raw, err := json.Marshal(resp)
	if err != nil {
		t.Fatalf("marshal response: %v", err)
	}

	var got map[string]any
	if err := json.Unmarshal(raw, &got); err != nil {
		t.Fatalf("unmarshal response: %v", err)
	}

	content := got["content"].(map[string]any)
	mediaType := content["application/problem+json"].(map[string]any)
	if mediaType["$ref"] != "#/components/mediaTypes/problem" {
		t.Fatalf("expected media type reference to be preserved, got %v", mediaType["$ref"])
	}
}

func TestResponseDescriptionFallback(t *testing.T) {
	if got := ResponseDescription(999); got != "Response" {
		t.Fatalf("expected fallback response description, got %q", got)
	}
}

func TestValidatorRejectsMutuallyExclusiveFields(t *testing.T) {
	validate := validator.New()

	tests := []struct {
		name string
		data any
	}{
		{
			name: "parameter schema and content",
			data: ParameterObject{
				Name:    "filter",
				In:      "query",
				Schema:  &SchemaOrReference{SchemaObject: &SchemaObject{Type: "string"}},
				Content: map[string]MediaTypeOrReference{"application/json": JSONSchemaRef("#/components/schemas/Filter_json")},
			},
		},
		{
			name: "media type schema and item schema",
			data: MediaTypeObject{
				Schema:     &SchemaOrReference{SchemaObject: &SchemaObject{Type: "array"}},
				ItemSchema: &SchemaOrReference{SchemaObject: &SchemaObject{Type: "string"}},
			},
		},
		{
			name: "link operation ref and id",
			data: LinkObject{
				OperationRef: "#/paths/~1users/get",
				OperationId:  "listUsers",
			},
		},
		{
			name: "example value and external value",
			data: ExampleObject{
				Value:         map[string]any{"id": "1"},
				ExternalValue: "https://example.com/example.json",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := validate.Struct(tt.data); err == nil {
				t.Fatalf("expected validation error for mutually exclusive fields")
			}
		})
	}
}

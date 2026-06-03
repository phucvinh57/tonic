package docs

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/TickLabVN/tonic/core/utils"
)

// https://swagger.io/specification/#components-object
type ComponentsObject struct {
	Schemas         map[string]SchemaOrReference         `json:"schemas,omitempty"`
	MediaTypes      map[string]MediaTypeOrReference      `json:"mediaTypes,omitempty"`
	Responses       map[string]ResponseOrReference       `json:"responses,omitempty"`
	Parameters      map[string]ParameterOrReference      `json:"parameters,omitempty"`
	Examples        map[string]ExampleOrReference        `json:"examples,omitempty"`
	RequestBodies   map[string]RequestBodyOrReference    `json:"requestBodies,omitempty"`
	SecuritySchemes map[string]SecuritySchemeOrReference `json:"securitySchemes,omitempty"`
	Links           map[string]LinkOrReference           `json:"links,omitempty"`
	Callbacks       map[string]CallbackOrReference       `json:"callbacks,omitempty"`
	PathItems       map[string]PathItemOrReference       `json:"pathItems,omitempty"`
}

func (c *ComponentsObject) AddSchema(t reflect.Type, parsingKey string, bindingKey string) (*SchemaObject, error) {
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	if c.Schemas == nil {
		c.Schemas = make(map[string]SchemaOrReference)
	}
	schemaName := utils.GetSchemaName(t)
	schemaName = strings.Join([]string{schemaName, parsingKey}, "_")
	if schema, exists := c.Schemas[schemaName]; exists {
		// Schema already exists, no need to add it again
		return schema.SchemaObject, nil
	}

	schema, err := SchemaFromType(t, parsingKey, bindingKey, nil)
	if err != nil {
		return nil, fmt.Errorf("create schema from type: %w", err)
	}
	c.Schemas[schemaName] = SchemaOrReference{SchemaObject: &schema}
	return &schema, nil
}

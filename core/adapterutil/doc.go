package adapterutil

import (
	"fmt"
	"log"
	"slices"

	"github.com/TickLabVN/tonic/core/docs"
	"github.com/TickLabVN/tonic/core/utils"
)

func Warn(format string, args ...any) {
	log.Printf("tonic docs warning: "+format, args...)
}

func ValidatePathParametersMatch(routePath string, schema *docs.SchemaObject) error {
	pathParamNames := utils.ExtractPathParamNames(routePath)
	schemaParamNames := make([]string, 0, len(schema.Properties))
	for name := range schema.Properties {
		schemaParamNames = append(schemaParamNames, name)
	}
	slices.Sort(pathParamNames)
	slices.Sort(schemaParamNames)
	if !slices.Equal(pathParamNames, schemaParamNames) {
		return fmt.Errorf("route path parameters %v do not match schema parameters %v for %s", pathParamNames, schemaParamNames, routePath)
	}
	return nil
}

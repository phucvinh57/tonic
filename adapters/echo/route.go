package echoAdapter

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/TickLabVN/tonic/core/adapterutil"
	"github.com/TickLabVN/tonic/core/docs"
	"github.com/TickLabVN/tonic/core/utils"
	"github.com/labstack/echo/v4"
)

type Adapter struct {
	spec *docs.OpenApi
}

type TypedAdapter[D any, R any] struct {
	base *Adapter
}

type BindingOptions struct {
	Param  bool
	Query  bool
	Header bool
	Body   bool
}

func New(spec *docs.OpenApi) *Adapter {
	return &Adapter{spec: spec}
}

func For[D any, R any](adapter *Adapter) TypedAdapter[D, R] {
	return TypedAdapter[D, R]{base: adapter}
}

func getParsingOptions(t reflect.Type) BindingOptions {
	opts := BindingOptions{}
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	if t.Kind() == reflect.Struct {
		for i := 0; i < t.NumField(); i++ {
			field := t.Field(i)
			opts.Param = opts.Param || field.Tag.Get("param") != ""
			opts.Query = opts.Query || field.Tag.Get("query") != ""
			opts.Header = opts.Header || field.Tag.Get("header") != ""
			opts.Body = opts.Body || field.Tag.Get("json") != ""
		}
	}
	return opts
}

func AddRoute[D any, R any](spec *docs.OpenApi, route *echo.Route, opts ...docs.OperationObject) {
	For[D, R](New(spec)).AddRoute(route, opts...)
}

func (a TypedAdapter[D, R]) AddRoute(route *echo.Route, opts ...docs.OperationObject) {
	spec := a.base.spec
	routeTypes := adapterutil.NewRouteTypes[D, R]()
	input, resp := routeTypes.Request, routeTypes.Response
	opId := fmt.Sprintf("%s_%s", route.Method, route.Name)
	opId = strings.ReplaceAll(opId, ".", "_")
	routePath := utils.NormalizeAPIPath(route.Path)
	op := docs.OperationObject{
		OperationId: opId,
	}
	schemaBasePath := utils.GetSchemaPath(input)

	parsingOpts := getParsingOptions(input)
	if parsingOpts.Param {
		schema, err := spec.Components.AddSchema(input, "param", "validate")
		if err != nil {
			adapterutil.Warn("skipping docs for %s %s: add param schema: %v", route.Method, route.Path, err)
			return
		}
		if err := adapterutil.ValidatePathParametersMatch(routePath, schema); err != nil {
			adapterutil.Warn("skipping docs for %s %s: %v", route.Method, route.Path, err)
			return
		}
		op.AddParameter("path", schema, schemaBasePath+"_param")
	}
	if parsingOpts.Query {
		schema, err := spec.Components.AddSchema(input, "query", "validate")
		if err != nil {
			adapterutil.Warn("skipping docs for %s %s: add query schema: %v", route.Method, route.Path, err)
			return
		}
		op.AddParameter("query", schema, schemaBasePath+"_query")
	}
	if parsingOpts.Header {
		schema, err := spec.Components.AddSchema(input, "header", "validate")
		if err != nil {
			adapterutil.Warn("skipping docs for %s %s: add header schema: %v", route.Method, route.Path, err)
			return
		}
		op.AddParameter("header", schema, schemaBasePath+"_header")
	}
	if parsingOpts.Body {
		_, err := spec.Components.AddSchema(input, "json", "validate")
		if err != nil {
			adapterutil.Warn("skipping docs for %s %s: add body schema: %v", route.Method, route.Path, err)
			return
		}
		if route.Method == echo.POST || route.Method == echo.PUT || route.Method == echo.PATCH {
			op.RequestBody = &docs.RequestBodyOrReference{
				RequestBodyObject: &docs.RequestBodyObject{
					Content: map[string]docs.MediaTypeOrReference{
						"application/json": docs.JSONSchemaRef(schemaBasePath + "_json"),
					},
				},
			}
		}
	}
	_, err := spec.Components.AddSchema(resp, "json", "validate")
	if err != nil {
		adapterutil.Warn("skipping docs for %s %s: add response schema: %v", route.Method, route.Path, err)
		return
	}

	op = utils.MergeStructs(op, docs.OperationObject{
		Responses: map[string]docs.ResponseOrReference{
			"200": docs.JSONResponse(200, utils.GetSchemaPath(resp)+"_json"),
		},
	})
	op = utils.MergeStructs(append([]docs.OperationObject{op}, opts...)...)

	if spec.Paths == nil {
		spec.Paths = make(docs.Paths)
	}
	pathItem := docs.PathItemObject{}
	switch route.Method {
	case echo.GET:
		pathItem.Get = &op
	case echo.POST:
		pathItem.Post = &op
	case echo.PUT:
		pathItem.Put = &op
	case echo.PATCH:
		pathItem.Patch = &op
	case echo.DELETE:
		pathItem.Delete = &op
	case echo.OPTIONS:
		pathItem.Options = &op
	case echo.HEAD:
		pathItem.Head = &op
	default:
		adapterutil.Warn("skipping docs for %s %s: unsupported HTTP method", route.Method, route.Path)
		return
	}

	spec.Paths.Update(routePath, pathItem)
}

package fiberAdapter

import (
	"fmt"
	"net/http"
	"reflect"
	"strings"

	"github.com/TickLabVN/tonic/core/adapterutil"
	"github.com/TickLabVN/tonic/core/docs"
	"github.com/TickLabVN/tonic/core/utils"
	"github.com/gofiber/fiber/v3"
)

type Adapter struct {
	spec *docs.OpenApi
}

type TypedAdapter[D any, R any] struct {
	base *Adapter
}

type Option func(*routeConfig)

type routeConfig struct {
	operations []docs.OperationObject
}

type BindingOptions struct {
	Path   bool
	Query  bool
	Header bool
	Body   bool
}

type methodAdapter struct {
	register   func(g fiber.Router, path string, handlers []fiber.Handler)
	assign     func(item *docs.PathItemObject, op *docs.OperationObject)
	allowsBody bool
}

type parameterBinding struct {
	enabled      bool
	parsingKey   string
	location     string
	suffix       string
	validatePath bool
}

var methodAdapters = map[string]methodAdapter{
	http.MethodGet: {
		register: func(g fiber.Router, path string, handlers []fiber.Handler) { register(g.Get, path, handlers) },
		assign:   func(item *docs.PathItemObject, op *docs.OperationObject) { item.Get = op },
	},
	http.MethodPost: {
		register:   func(g fiber.Router, path string, handlers []fiber.Handler) { register(g.Post, path, handlers) },
		assign:     func(item *docs.PathItemObject, op *docs.OperationObject) { item.Post = op },
		allowsBody: true,
	},
	http.MethodPut: {
		register:   func(g fiber.Router, path string, handlers []fiber.Handler) { register(g.Put, path, handlers) },
		assign:     func(item *docs.PathItemObject, op *docs.OperationObject) { item.Put = op },
		allowsBody: true,
	},
	http.MethodPatch: {
		register:   func(g fiber.Router, path string, handlers []fiber.Handler) { register(g.Patch, path, handlers) },
		assign:     func(item *docs.PathItemObject, op *docs.OperationObject) { item.Patch = op },
		allowsBody: true,
	},
	http.MethodDelete: {
		register: func(g fiber.Router, path string, handlers []fiber.Handler) { register(g.Delete, path, handlers) },
		assign:   func(item *docs.PathItemObject, op *docs.OperationObject) { item.Delete = op },
	},
	http.MethodOptions: {
		register: func(g fiber.Router, path string, handlers []fiber.Handler) { register(g.Options, path, handlers) },
		assign:   func(item *docs.PathItemObject, op *docs.OperationObject) { item.Options = op },
	},
	http.MethodHead: {
		register: func(g fiber.Router, path string, handlers []fiber.Handler) { register(g.Head, path, handlers) },
		assign:   func(item *docs.PathItemObject, op *docs.OperationObject) { item.Head = op },
	},
}

func New(spec *docs.OpenApi) *Adapter {
	return &Adapter{spec: spec}
}

func For[D any, R any](adapter *Adapter) TypedAdapter[D, R] {
	return TypedAdapter[D, R]{base: adapter}
}

func WithOperation(op docs.OperationObject) Option {
	return func(cfg *routeConfig) {
		cfg.operations = append(cfg.operations, op)
	}
}

func WithOperations(ops ...docs.OperationObject) Option {
	return func(cfg *routeConfig) {
		cfg.operations = append(cfg.operations, ops...)
	}
}

func GET[D any, R any](spec *docs.OpenApi, g fiber.Router, path string, args ...any) {
	For[D, R](New(spec)).GET(g, path, args...)
}

func POST[D any, R any](spec *docs.OpenApi, g fiber.Router, path string, args ...any) {
	For[D, R](New(spec)).POST(g, path, args...)
}

func PUT[D any, R any](spec *docs.OpenApi, g fiber.Router, path string, args ...any) {
	For[D, R](New(spec)).PUT(g, path, args...)
}

func PATCH[D any, R any](spec *docs.OpenApi, g fiber.Router, path string, args ...any) {
	For[D, R](New(spec)).PATCH(g, path, args...)
}

func DELETE[D any, R any](spec *docs.OpenApi, g fiber.Router, path string, args ...any) {
	For[D, R](New(spec)).DELETE(g, path, args...)
}

func OPTIONS[D any, R any](spec *docs.OpenApi, g fiber.Router, path string, args ...any) {
	For[D, R](New(spec)).OPTIONS(g, path, args...)
}

func HEAD[D any, R any](spec *docs.OpenApi, g fiber.Router, path string, args ...any) {
	For[D, R](New(spec)).HEAD(g, path, args...)
}

func (a TypedAdapter[D, R]) GET(g fiber.Router, path string, args ...any) {
	a.registerAndDocument(http.MethodGet, g, path, args...)
}

func (a TypedAdapter[D, R]) POST(g fiber.Router, path string, args ...any) {
	a.registerAndDocument(http.MethodPost, g, path, args...)
}

func (a TypedAdapter[D, R]) PUT(g fiber.Router, path string, args ...any) {
	a.registerAndDocument(http.MethodPut, g, path, args...)
}

func (a TypedAdapter[D, R]) PATCH(g fiber.Router, path string, args ...any) {
	a.registerAndDocument(http.MethodPatch, g, path, args...)
}

func (a TypedAdapter[D, R]) DELETE(g fiber.Router, path string, args ...any) {
	a.registerAndDocument(http.MethodDelete, g, path, args...)
}

func (a TypedAdapter[D, R]) OPTIONS(g fiber.Router, path string, args ...any) {
	a.registerAndDocument(http.MethodOptions, g, path, args...)
}

func (a TypedAdapter[D, R]) HEAD(g fiber.Router, path string, args ...any) {
	a.registerAndDocument(http.MethodHead, g, path, args...)
}

func getParsingOptions(t reflect.Type) BindingOptions {
	opts := BindingOptions{}
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	if t.Kind() != reflect.Struct {
		return opts
	}
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		opts.Path = opts.Path || field.Tag.Get("uri") != ""
		opts.Query = opts.Query || field.Tag.Get("query") != ""
		opts.Header = opts.Header || field.Tag.Get("header") != ""
		opts.Body = opts.Body || field.Tag.Get("json") != ""
	}
	return opts
}

func parseArgs(method string, path string, args []any) ([]fiber.Handler, []Option, bool) {
	handlers := make([]fiber.Handler, 0)
	options := make([]Option, 0)
	seenOption := false

	for _, arg := range args {
		if !seenOption {
			if handler, ok := arg.(fiber.Handler); ok {
				handlers = append(handlers, handler)
				continue
			}
			if handler, ok := arg.(func(fiber.Ctx) error); ok {
				handlers = append(handlers, fiber.Handler(handler))
				continue
			}
		}

		option, ok := arg.(Option)
		if !ok {
			adapterutil.Warn("skipping route %s %s: invalid helper argument of type %T", method, path, arg)
			return nil, nil, false
		}
		seenOption = true
		options = append(options, option)
	}

	if len(handlers) == 0 {
		adapterutil.Warn("skipping route %s %s: no handlers provided", method, path)
		return nil, nil, false
	}

	return handlers, options, true
}

func (a TypedAdapter[D, R]) registerAndDocument(method string, g fiber.Router, path string, args ...any) {
	spec := a.base.spec
	methodCfg, ok := methodAdapters[method]
	if !ok {
		adapterutil.Warn("skipping docs and route registration for %s %s: unsupported HTTP method", method, path)
		return
	}

	handlers, options, ok := parseArgs(method, path, args)
	if !ok {
		return
	}
	methodCfg.register(g, path, handlers)

	routeTypes := adapterutil.NewRouteTypes[D, R]()
	input, resp := routeTypes.Request, routeTypes.Response
	if _, err := spec.Components.AddSchema(resp, "json", "validate"); err != nil {
		adapterutil.Warn("skipping docs for %s %s: add response schema: %v", method, path, err)
		return
	}

	basePath, ok := getBasePath(g)
	if !ok {
		adapterutil.Warn("skipping docs for %s %s: unsupported fiber.Router type %T", method, path, g)
		return
	}

	normalizedPath := utils.NormalizeAPIPath(joinPaths(basePath, path))
	schemaBasePath := utils.GetSchemaPath(input)
	op := docs.OperationObject{
		OperationId: fmt.Sprintf("%s_%s", method, normalizedPath),
	}

	parsingOpts := getParsingOptions(input)
	bindings := []parameterBinding{
		{enabled: parsingOpts.Path, parsingKey: "uri", location: "path", suffix: "_uri", validatePath: true},
		{enabled: parsingOpts.Query, parsingKey: "query", location: "query", suffix: "_query"},
		{enabled: parsingOpts.Header, parsingKey: "header", location: "header", suffix: "_header"},
	}
	for _, binding := range bindings {
		if !binding.enabled {
			continue
		}
		schema, err := spec.Components.AddSchema(input, binding.parsingKey, "validate")
		if err != nil {
			adapterutil.Warn("skipping docs for %s %s: add %s schema: %v", method, path, binding.location, err)
			return
		}
		if binding.validatePath {
			if err := adapterutil.ValidatePathParametersMatch(normalizedPath, schema); err != nil {
				adapterutil.Warn("skipping docs for %s %s: %v", method, path, err)
				return
			}
		}
		op.AddParameter(binding.location, schema, schemaBasePath+binding.suffix)
	}

	if parsingOpts.Body && methodCfg.allowsBody {
		if _, err := spec.Components.AddSchema(input, "json", "validate"); err != nil {
			adapterutil.Warn("skipping docs for %s %s: add body schema: %v", method, path, err)
			return
		}
		op.RequestBody = &docs.RequestBodyOrReference{
			RequestBodyObject: &docs.RequestBodyObject{
				Content: map[string]docs.MediaTypeOrReference{
					"application/json": docs.JSONSchemaRef(schemaBasePath + "_json"),
				},
			},
		}
	}

	cfg := routeConfig{}
	for _, option := range options {
		option(&cfg)
	}

	op = utils.MergeStructs(op, docs.OperationObject{
		Responses: map[string]docs.ResponseOrReference{
			"200": docs.JSONResponse(200, utils.GetSchemaPath(resp)+"_json"),
		},
	})
	op = utils.MergeStructs(append([]docs.OperationObject{op}, cfg.operations...)...)

	if spec.Paths == nil {
		spec.Paths = make(docs.Paths)
	}
	pathItem := docs.PathItemObject{}
	methodCfg.assign(&pathItem, &op)
	spec.Paths.Update(normalizedPath, pathItem)
}

func register(fn func(string, any, ...any) fiber.Router, path string, handlers []fiber.Handler) {
	if len(handlers) == 1 {
		fn(path, handlers[0])
		return
	}

	rest := make([]any, 0, len(handlers)-1)
	for _, handler := range handlers[1:] {
		rest = append(rest, handler)
	}
	fn(path, handlers[0], rest...)
}

func getBasePath(g fiber.Router) (string, bool) {
	if group, ok := g.(*fiber.Group); ok {
		return group.Prefix, true
	}
	if _, ok := g.(*fiber.App); ok {
		return "", true
	}
	return "", false
}

func joinPaths(basePath string, path string) string {
	basePath = strings.TrimSuffix(basePath, "/")
	if basePath == "" {
		return path
	}
	if path == "" {
		return basePath
	}
	if strings.HasPrefix(path, "/") {
		return basePath + path
	}
	return basePath + "/" + path
}

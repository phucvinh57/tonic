package chiAdapter

import (
	"fmt"
	"net/http"
	"reflect"
	"strings"

	"github.com/TickLabVN/tonic/core/adapterutil"
	"github.com/TickLabVN/tonic/core/docs"
	"github.com/TickLabVN/tonic/core/utils"
	"github.com/go-chi/chi/v5"
)

type Adapter struct {
	spec *docs.OpenApi
}

type TypedAdapter[D any, R any] struct {
	base *Adapter
}

type Router struct {
	chi.Router
	basePath string
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
	register   func(r chi.Router, path string, middlewares []func(http.Handler) http.Handler, handler http.Handler)
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
		register: func(r chi.Router, path string, middlewares []func(http.Handler) http.Handler, handler http.Handler) {
			register(r, http.MethodGet, path, middlewares, handler)
		},
		assign: func(item *docs.PathItemObject, op *docs.OperationObject) { item.Get = op },
	},
	http.MethodPost: {
		register: func(r chi.Router, path string, middlewares []func(http.Handler) http.Handler, handler http.Handler) {
			register(r, http.MethodPost, path, middlewares, handler)
		},
		assign:     func(item *docs.PathItemObject, op *docs.OperationObject) { item.Post = op },
		allowsBody: true,
	},
	http.MethodPut: {
		register: func(r chi.Router, path string, middlewares []func(http.Handler) http.Handler, handler http.Handler) {
			register(r, http.MethodPut, path, middlewares, handler)
		},
		assign:     func(item *docs.PathItemObject, op *docs.OperationObject) { item.Put = op },
		allowsBody: true,
	},
	http.MethodPatch: {
		register: func(r chi.Router, path string, middlewares []func(http.Handler) http.Handler, handler http.Handler) {
			register(r, http.MethodPatch, path, middlewares, handler)
		},
		assign:     func(item *docs.PathItemObject, op *docs.OperationObject) { item.Patch = op },
		allowsBody: true,
	},
	http.MethodDelete: {
		register: func(r chi.Router, path string, middlewares []func(http.Handler) http.Handler, handler http.Handler) {
			register(r, http.MethodDelete, path, middlewares, handler)
		},
		assign: func(item *docs.PathItemObject, op *docs.OperationObject) { item.Delete = op },
	},
	http.MethodOptions: {
		register: func(r chi.Router, path string, middlewares []func(http.Handler) http.Handler, handler http.Handler) {
			register(r, http.MethodOptions, path, middlewares, handler)
		},
		assign: func(item *docs.PathItemObject, op *docs.OperationObject) { item.Options = op },
	},
	http.MethodHead: {
		register: func(r chi.Router, path string, middlewares []func(http.Handler) http.Handler, handler http.Handler) {
			register(r, http.MethodHead, path, middlewares, handler)
		},
		assign: func(item *docs.PathItemObject, op *docs.OperationObject) { item.Head = op },
	},
}

func New(spec *docs.OpenApi) *Adapter {
	return &Adapter{spec: spec}
}

func For[D any, R any](adapter *Adapter) TypedAdapter[D, R] {
	return TypedAdapter[D, R]{base: adapter}
}

func Wrap(r chi.Router) *Router {
	return &Router{Router: r}
}

func (a *Adapter) Wrap(r chi.Router) *Router {
	return Wrap(r)
}

func (r *Router) BasePath() string {
	return r.basePath
}

func (r *Router) With(middlewares ...func(http.Handler) http.Handler) chi.Router {
	return &Router{
		Router:   r.Router.With(middlewares...),
		basePath: r.basePath,
	}
}

func (r *Router) Group(fn func(r chi.Router)) chi.Router {
	grouped := r.Router.Group(func(cr chi.Router) {
		fn(&Router{Router: cr, basePath: r.basePath})
	})
	return &Router{Router: grouped, basePath: r.basePath}
}

func (r *Router) Route(pattern string, fn func(r chi.Router)) chi.Router {
	basePath := joinPaths(r.basePath, pattern)
	routed := r.Router.Route(pattern, func(cr chi.Router) {
		fn(&Router{Router: cr, basePath: basePath})
	})
	return &Router{Router: routed, basePath: basePath}
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

func GET[D any, R any](spec *docs.OpenApi, r chi.Router, path string, args ...any) {
	For[D, R](New(spec)).GET(r, path, args...)
}

func POST[D any, R any](spec *docs.OpenApi, r chi.Router, path string, args ...any) {
	For[D, R](New(spec)).POST(r, path, args...)
}

func PUT[D any, R any](spec *docs.OpenApi, r chi.Router, path string, args ...any) {
	For[D, R](New(spec)).PUT(r, path, args...)
}

func PATCH[D any, R any](spec *docs.OpenApi, r chi.Router, path string, args ...any) {
	For[D, R](New(spec)).PATCH(r, path, args...)
}

func DELETE[D any, R any](spec *docs.OpenApi, r chi.Router, path string, args ...any) {
	For[D, R](New(spec)).DELETE(r, path, args...)
}

func OPTIONS[D any, R any](spec *docs.OpenApi, r chi.Router, path string, args ...any) {
	For[D, R](New(spec)).OPTIONS(r, path, args...)
}

func HEAD[D any, R any](spec *docs.OpenApi, r chi.Router, path string, args ...any) {
	For[D, R](New(spec)).HEAD(r, path, args...)
}

func (a TypedAdapter[D, R]) GET(r chi.Router, path string, args ...any) {
	a.registerAndDocument(http.MethodGet, r, path, args...)
}

func (a TypedAdapter[D, R]) POST(r chi.Router, path string, args ...any) {
	a.registerAndDocument(http.MethodPost, r, path, args...)
}

func (a TypedAdapter[D, R]) PUT(r chi.Router, path string, args ...any) {
	a.registerAndDocument(http.MethodPut, r, path, args...)
}

func (a TypedAdapter[D, R]) PATCH(r chi.Router, path string, args ...any) {
	a.registerAndDocument(http.MethodPatch, r, path, args...)
}

func (a TypedAdapter[D, R]) DELETE(r chi.Router, path string, args ...any) {
	a.registerAndDocument(http.MethodDelete, r, path, args...)
}

func (a TypedAdapter[D, R]) OPTIONS(r chi.Router, path string, args ...any) {
	a.registerAndDocument(http.MethodOptions, r, path, args...)
}

func (a TypedAdapter[D, R]) HEAD(r chi.Router, path string, args ...any) {
	a.registerAndDocument(http.MethodHead, r, path, args...)
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
		opts.Path = opts.Path || field.Tag.Get("path") != ""
		opts.Query = opts.Query || field.Tag.Get("query") != ""
		opts.Header = opts.Header || field.Tag.Get("header") != ""
		opts.Body = opts.Body || field.Tag.Get("json") != ""
	}
	return opts
}

func parseArgs(method string, path string, args []any) ([]func(http.Handler) http.Handler, http.Handler, []Option, bool) {
	middlewares := make([]func(http.Handler) http.Handler, 0)
	options := make([]Option, 0)
	var endpoint http.Handler
	seenEndpoint := false
	seenOption := false

	for _, arg := range args {
		if !seenOption && !seenEndpoint {
			switch value := arg.(type) {
			case func(http.Handler) http.Handler:
				middlewares = append(middlewares, value)
				continue
			case chi.Middlewares:
				middlewares = append(middlewares, value...)
				continue
			case http.HandlerFunc:
				endpoint = value
				seenEndpoint = true
				continue
			case http.Handler:
				endpoint = value
				seenEndpoint = true
				continue
			case func(http.ResponseWriter, *http.Request):
				endpoint = http.HandlerFunc(value)
				seenEndpoint = true
				continue
			}
		}

		option, ok := arg.(Option)
		if !ok {
			adapterutil.Warn("skipping route %s %s: invalid helper argument of type %T", method, path, arg)
			return nil, nil, nil, false
		}
		seenOption = true
		options = append(options, option)
	}

	if endpoint == nil {
		adapterutil.Warn("skipping route %s %s: no handler provided", method, path)
		return nil, nil, nil, false
	}

	return middlewares, endpoint, options, true
}

func (a TypedAdapter[D, R]) registerAndDocument(method string, r chi.Router, path string, args ...any) {
	spec := a.base.spec
	methodCfg, ok := methodAdapters[method]
	if !ok {
		adapterutil.Warn("skipping docs and route registration for %s %s: unsupported HTTP method", method, path)
		return
	}

	middlewares, handler, options, ok := parseArgs(method, path, args)
	if !ok {
		return
	}
	methodCfg.register(r, path, middlewares, handler)

	routeTypes := adapterutil.NewRouteTypes[D, R]()
	input, resp := routeTypes.Request, routeTypes.Response
	if _, err := spec.Components.AddSchema(resp, "json", "validate"); err != nil {
		adapterutil.Warn("skipping docs for %s %s: add response schema: %v", method, path, err)
		return
	}

	basePath := getBasePath(r)
	normalizedPath := utils.NormalizeAPIPath(joinPaths(basePath, path))
	schemaBasePath := utils.GetSchemaPath(input)
	op := docs.OperationObject{
		OperationId: fmt.Sprintf("%s_%s", method, normalizedPath),
	}

	parsingOpts := getParsingOptions(input)
	bindings := []parameterBinding{
		{enabled: parsingOpts.Path, parsingKey: "path", location: "path", suffix: "_path", validatePath: true},
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

func register(r chi.Router, method string, path string, middlewares []func(http.Handler) http.Handler, handler http.Handler) {
	if len(middlewares) > 0 {
		handler = chi.Chain(middlewares...).Handler(handler)
	}
	r.Method(method, path, handler)
}

func getBasePath(r chi.Router) string {
	base, ok := r.(interface{ BasePath() string })
	if !ok {
		return ""
	}
	return base.BasePath()
}

func joinPaths(basePath string, path string) string {
	basePath = strings.TrimSuffix(basePath, "/")
	if basePath == "" {
		if path == "" {
			return "/"
		}
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

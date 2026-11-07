package core

import "github.com/TickLabVN/tonic/core/docs"

type Option func(*docs.OpenApi)

func WithTagConfig(tc docs.TagConfig) Option {
	return func(o *docs.OpenApi) {
		o.Components.TagConfig = tc
	}
}

func Init(opts ...Option) *docs.OpenApi {
	c := &docs.OpenApi{
		OpenAPI: docs.VERSION,
		Info: docs.InfoObject{
			Title:   "Tonic API",
			Version: "0.0.0",
		},
	}
	c.Components.TagConfig = docs.DefaultTagConfig()

	for _, opt := range opts {
		opt(c)
	}

	return c
}

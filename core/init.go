package core

import "github.com/TickLabVN/tonic/core/docs"

func Init() *docs.OpenApi {
	c := &docs.OpenApi{
		OpenAPI: docs.VERSION,
		Info: docs.InfoObject{
			Title:   "Tonic API",
			Version: "0.0.0",
		},
	}

	return c
}

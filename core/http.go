package core

import (
	"encoding/json"
	"net/http"

	"github.com/TickLabVN/tonic/core/docs"
)

func JsonHttpHandler(spec *docs.OpenApi) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(spec)
	})
}

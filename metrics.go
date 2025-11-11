package main

import (
	"fmt"
	"net/http"
)

func (cfg *apiConfig) metricsResponse(writer http.ResponseWriter, request *http.Request) {
	writer.WriteHeader(200)
	writer.Write([]byte(fmt.Sprintf("Hits: %d", cfg.fileserverHits.Load())))
}

func (cfg *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(
		func(writer http.ResponseWriter, request *http.Request) {
			cfg.fileserverHits.Add(1)
			next.ServeHTTP(writer, request)
		},
	)
}

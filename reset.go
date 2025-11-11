package main

import "net/http"

func (cfg *apiConfig) metricsReset(writer http.ResponseWriter, request *http.Request) {
	cfg.fileserverHits.Store(0)
	writer.WriteHeader(200)
}

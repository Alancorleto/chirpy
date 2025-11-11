package main

import (
	"fmt"
	"net/http"
	"path/filepath"
	"sync/atomic"
)

type apiConfig struct {
	fileserverHits atomic.Int32
}

func main() {
	serveMux := http.NewServeMux()
	server := http.Server{
		Handler: serveMux,
		Addr:    ":8080",
	}

	apiCfg := apiConfig{}

	basePath := filepath.Dir(".")
	handler := http.FileServer(http.Dir(basePath))
	serveMux.Handle("/app/", apiCfg.middlewareMetricsInc(http.StripPrefix("/app", handler)))

	serveMux.HandleFunc("/healthz/", healthzResponse)

	serveMux.HandleFunc("/metrics/", apiCfg.metricsResponse)

	serveMux.HandleFunc("/reset/", apiCfg.metricsReset)

	err := server.ListenAndServe()
	if err != nil {
		fmt.Printf("Error listening to server: %v\n", err)
	}
}

func healthzResponse(writer http.ResponseWriter, request *http.Request) {
	writer.Header().Add("Content-Type", "text/plain; charset=utf-8")
	writer.WriteHeader(200)
	writer.Write([]byte("OK"))
}

type middlewareMetricsHandler struct {
	ApiCfg  *apiConfig
	Handler http.Handler
}

func (mmh *middlewareMetricsHandler) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	mmh.ApiCfg.fileserverHits.Add(1)
	mmh.Handler.ServeHTTP(writer, request)
}

func (cfg *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(
		func(writer http.ResponseWriter, request *http.Request) {
			cfg.fileserverHits.Add(1)
			next.ServeHTTP(writer, request)
		},
	)
}

func (cfg *apiConfig) metricsResponse(writer http.ResponseWriter, request *http.Request) {
	writer.WriteHeader(200)
	writer.Write([]byte(fmt.Sprintf("Hits: %d", cfg.fileserverHits.Load())))
}

func (cfg *apiConfig) metricsReset(writer http.ResponseWriter, request *http.Request) {
	cfg.fileserverHits.Store(0)
	writer.WriteHeader(200)
}

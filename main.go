package main

import (
	"fmt"
	"log"
	"net/http"
	"sync/atomic"
)

type apiConfig struct {
	fileserverHits atomic.Int32
}

func main() {
	const filepathRoot = "."
	const port = "8080"

	serveMux := http.NewServeMux()
	server := http.Server{
		Handler: serveMux,
		Addr:    ":" + port,
	}

	apiCfg := apiConfig{}

	serveMux.Handle("/app/", apiCfg.middlewareMetricsInc(http.StripPrefix("/app", http.FileServer(http.Dir(filepathRoot)))))
	serveMux.HandleFunc("GET /api/healthz", healthzResponse)
	serveMux.HandleFunc("GET /admin/metrics", apiCfg.metricsResponse)
	serveMux.HandleFunc("POST /admin/reset", apiCfg.metricsReset)
	serveMux.HandleFunc("POST /api/validate_chirp", validateChirpHandler)

	log.Printf("Serving files from %s on port: %s\n", filepathRoot, port)
	err := server.ListenAndServe()
	if err != nil {
		fmt.Printf("Error listening to server: %v\n", err)
	}
}

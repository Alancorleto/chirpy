package main

import (
	"fmt"
	"net/http"
	"os"
)

func (cfg *apiConfig) metricsReset(writer http.ResponseWriter, request *http.Request) {
	if os.Getenv("PLATFORM") != "dev" {
		respondWithError(writer, 403, "forbidden action")
		return
	}
	err := cfg.db.ResetUsers(request.Context())
	if err != nil {
		fmt.Printf("Error reseting users table: %v", err)
		return
	}
	cfg.fileserverHits.Store(0)
	writer.WriteHeader(200)
}

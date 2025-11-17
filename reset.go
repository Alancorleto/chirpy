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
		respondWithError(writer, 500, fmt.Sprintf("Error reseting users table: %v", err))
		return
	}

	err = cfg.db.ResetChirps(request.Context())
	if err != nil {
		respondWithError(writer, 500, fmt.Sprintf("Error reseting chirps table: %v", err))
		return
	}

	cfg.fileserverHits.Store(0)
	writer.WriteHeader(200)
}

package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"sync/atomic"

	"github.com/alancorleto/chirpy/internal/database"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

type apiConfig struct {
	fileserverHits atomic.Int32
	dbQueries      *database.Queries
}

func main() {
	godotenv.Load()
	dbURL := os.Getenv("DB_URL")
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		fmt.Printf("Error creating SQL database: %v\n", err)
		return
	}
	dbQueries := database.New(db)

	const filepathRoot = "."
	const port = "8080"

	serveMux := http.NewServeMux()
	server := http.Server{
		Handler: serveMux,
		Addr:    ":" + port,
	}

	apiCfg := apiConfig{
		dbQueries: dbQueries,
	}

	serveMux.Handle("/app/", apiCfg.middlewareMetricsInc(http.StripPrefix("/app", http.FileServer(http.Dir(filepathRoot)))))
	serveMux.HandleFunc("GET /api/healthz", healthzResponse)
	serveMux.HandleFunc("GET /admin/metrics", apiCfg.metricsResponse)
	serveMux.HandleFunc("POST /admin/reset", apiCfg.metricsReset)
	serveMux.HandleFunc("POST /api/validate_chirp", validateChirpHandler)

	log.Printf("Serving files from %s on port: %s\n", filepathRoot, port)
	err = server.ListenAndServe()
	if err != nil {
		fmt.Printf("Error listening to server: %v\n", err)
	}
}

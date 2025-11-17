package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/alancorleto/chirpy/internal/database"
	"github.com/google/uuid"
)

type Chirp struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Body      string    `json:"body"`
	UserID    uuid.UUID `json:"user_id"`
}

func (cfg *apiConfig) createChirpHandler(writer http.ResponseWriter, request *http.Request) {
	type parameters struct {
		Body   string    `json:"body"`
		UserID uuid.UUID `json:"user_id"`
	}

	params, err := parseRequestParameters[parameters](request)
	if err != nil {
		respondWithError(writer, 500, fmt.Sprintf("Error decoding parameters: %s", err))
		return
	}

	chirp, err := cfg.db.CreateChirp(
		request.Context(),
		database.CreateChirpParams{
			Body:   params.Body,
			UserID: params.UserID,
		},
	)
	if err != nil {
		respondWithError(writer, 500, fmt.Sprintf("Error creating chirp: %s", err))
		return
	}

	respondWithJSON(writer, 201, fromDatabaseChirpToChirp(chirp))
}

func (cfg *apiConfig) getChirpsHandler(writer http.ResponseWriter, request *http.Request) {
	dbChirps, err := cfg.db.GetChirps(request.Context())
	if err != nil {
		respondWithError(writer, 500, fmt.Sprintf("Error getting chirps: %s", err))
		return
	}

	chirps := []Chirp{}
	for _, dbChirp := range dbChirps {
		chirps = append(chirps, fromDatabaseChirpToChirp(dbChirp))
	}
	respondWithJSON(writer, 200, chirps)
}

func fromDatabaseChirpToChirp(dbChirp database.Chirp) Chirp {
	return Chirp{
		ID:        dbChirp.ID,
		CreatedAt: dbChirp.CreatedAt,
		UpdatedAt: dbChirp.UpdatedAt,
		Body:      dbChirp.Body,
		UserID:    dbChirp.UserID,
	}
}

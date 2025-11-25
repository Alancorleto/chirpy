package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/alancorleto/chirpy/internal/auth"
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
		Body string `json:"body"`
	}

	params, err := parseRequestParameters[parameters](request)
	if err != nil {
		respondWithError(writer, 500, fmt.Sprintf("Error decoding parameters: %s", err))
		return
	}

	bearerToken, err := auth.GetBearerToken(request.Header)
	if err != nil {
		respondWithError(writer, 401, fmt.Sprintf("Error getting bearer token: %s", err))
		return
	}

	userID, err := auth.ValidateJWT(bearerToken, cfg.secret)
	if err != nil {
		respondWithError(writer, 401, "Unauthorized operation")
		return
	}

	chirp, err := cfg.db.CreateChirp(
		request.Context(),
		database.CreateChirpParams{
			Body:   params.Body,
			UserID: userID,
		},
	)
	if err != nil {
		respondWithError(writer, 500, fmt.Sprintf("Error creating chirp: %s", err))
		return
	}

	respondWithJSON(writer, 201, fromDatabaseChirpToChirp(chirp))
}

func (cfg *apiConfig) getChirpsHandler(writer http.ResponseWriter, request *http.Request) {
	authorIDString := request.URL.Query().Get("author_id")
	if authorIDString == "" {
		cfg.getAllChirps(writer, request)
	} else {
		authorID, err := uuid.Parse(authorIDString)
		if err != nil {
			respondWithError(writer, 401, fmt.Sprintf("Error parsing author ID: %v", err))
			return
		}
		cfg.getAllChirpsByAuthor(writer, request, authorID)
	}
}

func (cfg *apiConfig) getAllChirps(writer http.ResponseWriter, request *http.Request) {
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

func (cfg *apiConfig) getAllChirpsByAuthor(writer http.ResponseWriter, request *http.Request, authorID uuid.UUID) {
	dbChirps, err := cfg.db.GetChirpsByAuthor(request.Context(), authorID)
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

func (cfg *apiConfig) getChirpHandler(writer http.ResponseWriter, request *http.Request) {
	pathParameter := request.PathValue("chirpID")
	chirpID, err := uuid.Parse(pathParameter)
	if err != nil {
		respondWithError(writer, 500, fmt.Sprintf("Could not parse chirp id: %s", err))
		return
	}
	chirp, err := cfg.db.GetChirp(request.Context(), chirpID)
	if err != nil {
		respondWithError(writer, 404, "chirp not found")
		return
	}
	respondWithJSON(writer, 200, fromDatabaseChirpToChirp(chirp))
}

func (cfg *apiConfig) deleteChirpHandler(writer http.ResponseWriter, request *http.Request) {
	pathParameter := request.PathValue("chirpID")
	chirpID, err := uuid.Parse(pathParameter)
	if err != nil {
		respondWithError(writer, 500, fmt.Sprintf("Could not parse chirp id: %s", err))
		return
	}
	chirp, err := cfg.db.GetChirp(request.Context(), chirpID)
	if err != nil {
		respondWithError(writer, 404, "chirp not found")
		return
	}

	token, err := auth.GetBearerToken(request.Header)
	if err != nil {
		respondWithError(writer, 401, fmt.Sprintf("Error getting bearer token: %s", err))
		return
	}

	tokenUserID, err := auth.ValidateJWT(token, cfg.secret)
	if err != nil {
		respondWithError(writer, 401, fmt.Sprintf("Error validating access token: %s", err))
		return
	}

	if chirp.UserID != tokenUserID {
		respondWithError(writer, 403, "Unauthorized operation")
		return
	}

	err = cfg.db.DeleteChirp(request.Context(), chirpID)
	if err != nil {
		respondWithError(writer, 500, fmt.Sprintf("Error vdeleting chirp: %s", err))
		return
	}

	writer.WriteHeader(204)
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

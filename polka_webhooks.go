package main

import (
	"fmt"
	"net/http"

	"github.com/alancorleto/chirpy/internal/auth"
	"github.com/google/uuid"
)

func (cfg *apiConfig) polkaWebhooksHandler(writer http.ResponseWriter, request *http.Request) {
	apiKey, err := auth.GetAPIKey(request.Header)
	if err != nil {
		respondWithError(writer, 401, fmt.Sprintf("Error getting API key: %s", err))
		return
	}

	if apiKey != cfg.polkaKey {
		respondWithError(writer, 401, "Incorrect API key")
		return
	}

	type parameters struct {
		Event string `json:"event"`
		Data  struct {
			UserID string `json:"user_id"`
		} `json:"data"`
	}

	params, err := parseRequestParameters[parameters](request)
	if err != nil {
		respondWithError(writer, 500, fmt.Sprintf("Error decoding parameters: %s", err))
		return
	}

	if params.Event != "user.upgraded" {
		writer.WriteHeader(204)
		return
	}

	userID, err := uuid.Parse(params.Data.UserID)
	if err != nil {
		respondWithError(writer, 401, fmt.Sprintf("Error parsing user ID: %s", err))
		return
	}

	err = cfg.db.UpgradeUserToChirpyRed(request.Context(), userID)
	if err != nil {
		respondWithError(writer, 404, "User not found")
		return
	}

	writer.WriteHeader(204)
}

package main

import (
	"fmt"
	"net/http"

	"github.com/alancorleto/chirpy/internal/auth"
)

func (cfg *apiConfig) revokeHandler(writer http.ResponseWriter, request *http.Request) {
	type returnVals struct {
		Token string `json:"token"`
	}

	refreshToken, err := auth.GetBearerToken(request.Header)
	if err != nil {
		respondWithError(writer, 500, fmt.Sprintf("Error getting refresh bearer token: %s", err))
		return
	}

	err = cfg.db.RevokeToken(request.Context(), refreshToken)
	if err != nil {
		respondWithError(writer, 401, "invalid refresh token")
		return
	}

	writer.WriteHeader(204)
}

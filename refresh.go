package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/alancorleto/chirpy/internal/auth"
)

func (cfg *apiConfig) refreshHandler(writer http.ResponseWriter, request *http.Request) {
	type returnVals struct {
		Token string `json:"token"`
	}

	refreshToken, err := auth.GetBearerToken(request.Header)
	if err != nil {
		respondWithError(writer, 500, fmt.Sprintf("Error getting refresh bearer token: %s", err))
		return
	}

	token, err := cfg.db.GetRefreshToken(request.Context(), refreshToken)
	if err != nil {
		respondWithError(writer, 401, "invalid refresh token")
		return
	}
	if token.ExpiresAt.Before(time.Now()) {
		respondWithError(writer, 401, "token has expired")
		return
	}
	if token.RevokedAt.Valid && token.RevokedAt.Time.Before(time.Now()) {
		respondWithError(writer, 401, "token has been revoked")
		return
	}

	user, err := cfg.db.GetUserFromRefreshToken(request.Context(), refreshToken)
	if err != nil {
		respondWithError(writer, 401, "invalid refresh token")
		return
	}

	jwt, err := auth.MakeJWT(user.ID, cfg.secret, 1*time.Hour)
	if err != nil {
		respondWithError(writer, 500, "error creating access token")
		return
	}

	respondWithJSON(writer, 200, returnVals{Token: jwt})
}

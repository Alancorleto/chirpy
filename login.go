package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/alancorleto/chirpy/internal/auth"
	"github.com/alancorleto/chirpy/internal/database"
)

type loggedUser struct {
	User
	Token string `json:"token"`
}

func (cfg *apiConfig) loginHandler(writer http.ResponseWriter, request *http.Request) {
	type parameters struct {
		Email            string `json:"email"`
		Password         string `json:"password"`
		ExpiresInSeconds int    `json:"expires_in_seconds"`
	}

	params, err := parseRequestParameters[parameters](request)
	if err != nil {
		respondWithError(writer, 500, fmt.Sprintf("Error decoding parameters: %s", err))
		return
	}

	user, err := cfg.db.GetUserByEmail(request.Context(), params.Email)
	if err != nil {
		respondWithError(writer, 401, "Incorrect email or password")
		return
	}

	isPasswordCorrect, err := auth.CheckPasswordHash(params.Password, user.HashedPassword)
	if err != nil {
		respondWithError(writer, 500, fmt.Sprintf("Error verifying password: %s", err))
		return
	}
	if !isPasswordCorrect {
		respondWithError(writer, 401, "Incorrect email or password")
		return
	}

	expirationTime := time.Duration(params.ExpiresInSeconds)
	if expirationTime == 0 {
		expirationTime = time.Hour
	}
	if expirationTime > time.Hour {
		expirationTime = time.Hour
	}

	token, err := auth.MakeJWT(user.ID, cfg.secret, expirationTime)
	if err != nil {
		respondWithError(writer, 500, fmt.Sprintf("Error creating authentication token: %s", err))
		return
	}

	respondWithJSON(writer, 200, fromDatabaseUserToLoggedUser(user, token))
}

func fromDatabaseUserToLoggedUser(dbUser database.User, token string) loggedUser {
	return loggedUser{
		User:  fromDatabaseUserToUser(dbUser),
		Token: token,
	}
}

package main

import (
	"fmt"
	"net/http"

	"github.com/alancorleto/chirpy/internal/auth"
)

func (cfg *apiConfig) loginHandler(writer http.ResponseWriter, request *http.Request) {
	type parameters struct {
		Email    string `json:"email"`
		Password string `json:"password"`
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

	respondWithJSON(writer, 200, fromDatabaseUserToUser(user))
}

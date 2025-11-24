package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/alancorleto/chirpy/internal/database"
	"github.com/google/uuid"

	auth "github.com/alancorleto/chirpy/internal/auth"
)

type User struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Email     string    `json:"email"`
}

func (cfg *apiConfig) postUsersHandler(writer http.ResponseWriter, request *http.Request) {
	type parameters struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	params, err := parseRequestParameters[parameters](request)
	if err != nil {
		respondWithError(writer, 500, fmt.Sprintf("Error decoding parameters: %s", err))
		return
	}

	hashedPassword, err := auth.HashPassword(params.Password)
	if err != nil {
		respondWithError(writer, 500, fmt.Sprintf("Error hashing password: %s", err))
	}

	user, err := cfg.db.CreateUser(
		request.Context(),
		database.CreateUserParams{
			Email:          params.Email,
			HashedPassword: hashedPassword,
		},
	)
	if err != nil {
		respondWithError(writer, 500, fmt.Sprintf("Error creating user: %s", err))
		return
	}

	respondWithJSON(writer, 201, fromDatabaseUserToUser(user))
}

func (cfg *apiConfig) putUsersHandler(writer http.ResponseWriter, request *http.Request) {
	type parameters struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	params, err := parseRequestParameters[parameters](request)
	if err != nil {
		respondWithError(writer, 500, fmt.Sprintf("Error decoding parameters: %s", err))
		return
	}

	hashedPassword, err := auth.HashPassword(params.Password)
	if err != nil {
		respondWithError(writer, 500, fmt.Sprintf("Error hashing password: %s", err))
	}

	token, err := auth.GetBearerToken(request.Header)
	if err != nil {
		respondWithError(writer, 401, fmt.Sprintf("Error getting bearer token: %s", err))
		return
	}

	userID, err := auth.ValidateJWT(token, cfg.secret)
	if err != nil {
		respondWithError(writer, 401, fmt.Sprintf("Error validating access token: %s", err))
		return
	}

	user, err := cfg.db.UpdateUser(
		request.Context(),
		database.UpdateUserParams{
			ID:             userID,
			Email:          params.Email,
			HashedPassword: hashedPassword,
		},
	)
	if err != nil {
		respondWithError(writer, 500, fmt.Sprintf("Error updating user: %s", err))
		return
	}

	respondWithJSON(writer, 200, fromDatabaseUserToUser(user))
}

func fromDatabaseUserToUser(dbUser database.User) User {
	return User{
		ID:        dbUser.ID,
		CreatedAt: dbUser.CreatedAt,
		UpdatedAt: dbUser.UpdatedAt,
		Email:     dbUser.Email,
	}
}

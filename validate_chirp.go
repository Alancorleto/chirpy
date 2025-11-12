package main

import (
	"fmt"
	"net/http"
	"strings"
)

func validateChirpHandler(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Body string `json:"body"`
	}
	type returnVals struct {
		CleanedBody string `json:"cleaned_body"`
	}

	params, err := parseRequestParameters[parameters](r)
	if err != nil {
		respondWithError(w, 500, fmt.Sprintf("Error decoding parameters: %s", err))
		return
	}

	const maxChirpLength = 140
	if len(params.Body) > maxChirpLength {
		respondWithError(w, 400, "Chirp is too long")
	} else {
		cleanedBody := cleanProfanities(params.Body)
		respondWithJSON(w, 200, returnVals{CleanedBody: cleanedBody})
	}
}

func cleanProfanities(body string) string {
	profanities := map[string]struct{}{
		"kerfuffle": {},
		"sharbert":  {},
		"fornax":    {},
	}
	censoredWord := "****"
	words := strings.Split(body, " ")
	for i, word := range words {
		loweredWord := strings.ToLower(word)
		if _, ok := profanities[loweredWord]; ok {
			words[i] = censoredWord
		}
	}
	return strings.Join(words, " ")
}

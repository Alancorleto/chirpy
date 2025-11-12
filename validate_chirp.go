package main

import (
	"encoding/json"
	"log"
	"net/http"
)

func validateChirpHandler(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Body string `json:"body"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		log.Printf("Error decoding parameters: %s", err)
		w.WriteHeader(500)
		return
	}

	if len(params.Body) > 140 {
		type returnValsError struct {
			// the key will be the name of struct field unless you give it an explicit JSON tag
			Error string `json:"error"`
		}
		respBodyError := returnValsError{
			Error: "Chirp is too long",
		}
		dat, err := json.Marshal(respBodyError)
		if err != nil {
			log.Printf("Error marshalling JSON: %s", err)
			w.WriteHeader(500)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(400)
		w.Write(dat)
	} else {
		type returnVals struct {
			// the key will be the name of struct field unless you give it an explicit JSON tag
			Valid bool `json:"valid"`
		}
		respBody := returnVals{
			Valid: true,
		}
		dat, err := json.Marshal(respBody)
		if err != nil {
			log.Printf("Error marshalling JSON: %s", err)
			w.WriteHeader(500)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(200)
		w.Write(dat)
	}

}

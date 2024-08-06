package main

import (
	"encoding/json"
	"net/http"
	"strings"

	database "github.com/C1sc0Ram0s/Chirpy/internal/database"
)

type Chirp struct {
	Id   int    `json:"id"`
	Body string `json:"body"`
}

func (id *chirpIdConfig) handlerChirpsCreate(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Body string `json:"body"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters")
		return
	}

	const maxChirpLength = 140
	if len(params.Body) > maxChirpLength {
		respondWithError(w, http.StatusBadRequest, "Chirp is too long")
		return
	}

	badWords := map[string]struct{}{
		"kerfuffle": {},
		"sharbert":  {},
		"fornax":    {},
	}
	cleaned := getCleanedBody(params.Body, badWords)

	id.id++
	chirps := Chirp{
		Id:   id.id,
		Body: cleaned,
	}

	// Creates new database connection and writes a chirp to disk
	dbConnection, err := database.NewDB("database.json")
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Database connection failed")
	}
	dbConnection.CreateChirp(chirps.Body)
	respondWithJSON(w, 201, chirps)
}

func getCleanedBody(body string, badWords map[string]struct{}) string {
	words := strings.Split(body, " ")
	for i, word := range words {
		lowerWord := strings.ToLower(word)
		if _, exists := badWords[lowerWord]; exists {
			words[i] = "****"
		}
	}
	cleaned := strings.Join(words, " ")
	return cleaned
}

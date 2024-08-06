package main

import (
	"encoding/json"
	"net/http"

	"github.com/C1sc0Ram0s/Chirpy/internal/database"
)

type User struct {
	ID    int    `json:"id"`
	Email string `json:"email"`
}

func (userIdCfg *userIdConfig) handlerCreateUsers(w http.ResponseWriter, r *http.Request) {
	type parameter struct {
		Body string `json:"email"`
	}
	params := parameter{}

	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "error decoding parameter")
		return
	}

	userIdCfg.id++
	user := User{
		ID:    userIdCfg.id,
		Email: params.Body,
	}

	// Create a new DB connection and write user to disk
	dbConnection, err := database.NewDB("database.json")
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Database connection failed")
	}
	dbConnection.CreateUser(user.Email)
	respondWithJSON(w, 201, user)
}

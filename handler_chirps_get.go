package main

import (
	"fmt"
	"net/http"
	"strconv"

	database "github.com/C1sc0Ram0s/Chirpy/internal/database"
)

func handlerGetChirpID(w http.ResponseWriter, r *http.Request) {
	type returnVals struct {
		Id   int    `json:"id"`
		Body string `json:"body"`
	}
	dbConnection, err := database.NewDB("database.json")
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, fmt.Sprintf("%v", err))
		return
	}

	chirpId, err := strconv.Atoi(r.PathValue("chirpID"))
	if err != nil {
		respondWithError(w, http.StatusBadRequest, fmt.Sprintf("%v", err))
		return
	}

	chirp, err := dbConnection.GetChirps(chirpId)
	if err != nil {
		respondWithError(w, 404, fmt.Sprintf("%v", err))
		return
	}

	ret := returnVals{
		Id:   chirp[0].Id,
		Body: chirp[0].Body,
	}
	respondWithJSON(w, http.StatusOK, ret)

}

func handlerGetChirps(w http.ResponseWriter, r *http.Request) {
	dbConnection, err := database.NewDB("database.json")
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, fmt.Sprintf("%v", err))
	}

	chirps, err := dbConnection.GetChirps()
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, fmt.Sprintf("%v", err))
	}

	respondWithJSON(w, http.StatusOK, chirps)
}

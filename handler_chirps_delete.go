package main

import (
	"net/http"
	"strconv"

	"github.com/C1sc0Ram0s/Chirpy/internal/auth"
)

func (cfg *apiConfig) handlerChirpsDelete(w http.ResponseWriter, r *http.Request) {
	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "invalid token")
		return
	}

	userIDString, err := auth.ValidateJWT(token, cfg.Jwt)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "error validating token")
		return
	}
	userID, err := strconv.Atoi(userIDString)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "error converting userIDString to userID")
		return
	}

	chirpIDString := r.PathValue("chirpID")
	chirpID, err := strconv.Atoi(chirpIDString)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "invalid chirpID")
		return
	}

	err = cfg.DB.DeleteChirp(chirpID, userID)
	if err != nil {
		respondWithError(w, 403, "unauthorized chirp deletion")
		return
	}
	w.WriteHeader(204)
}

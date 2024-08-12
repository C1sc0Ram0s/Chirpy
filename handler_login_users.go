package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/C1sc0Ram0s/Chirpy/internal/auth"
)

func (cfg *apiConfig) handlerLogin(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Password string `json:"password"`
		Email    string `json:"email"`
	}
	type response struct {
		User
		Token        string `json:"token"`
		RefreshToken string `json:"refresh_token"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters")
		return
	}

	user, err := cfg.DB.GetUserByEmail(params.Email)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't get user")
		return
	}

	err = auth.CheckPasswordHash(params.Password, user.HashedPassword)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Invalid password")
		return
	}

	// JWT token creation
	token, err := auth.MakeJWT(user.ID, cfg.Jwt, time.Hour)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error creating JWT")
		return
	}

	// RefreshToken creation
	refreshToken, err := auth.MakeRefreshToken()
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error creating refresh token")
		return
	}

	err = cfg.DB.StoreRefreshToken(user.ID, refreshToken)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error storing refresh token")
		return
	}

	respondWithJSON(w, http.StatusOK, response{
		User: User{
			ID:          user.ID,
			Email:       user.Email,
			IsChirpyRed: user.IsChirpyRed,
		},
		Token:        token,
		RefreshToken: refreshToken,
	})
}

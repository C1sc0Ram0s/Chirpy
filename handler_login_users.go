package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/C1sc0Ram0s/Chirpy/internal/auth"
	"github.com/golang-jwt/jwt/v5"
)

func (cfg *apiConfig) handlerLogin(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Password   string `json:"password"`
		Email      string `json:"email"`
		Expiration int    `json:"expires_in_seconds"`
	}
	type response struct {
		User
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

	if params.Expiration == 0 {
		params.Expiration = 86400
	}

	// JWT token creation
	expiration_time := time.Duration(time.Duration(params.Expiration) * time.Second)
	claims := jwt.RegisteredClaims{
		Issuer:    "chirpy",
		IssuedAt:  jwt.NewNumericDate(time.Now()),
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(expiration_time)),
		Subject:   fmt.Sprintf("%d", user.ID),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	ss, err := token.SignedString([]byte(cfg.Jwt))
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't generate token")
		return
	}

	err = auth.CheckPasswordHash(params.Password, user.HashedPassword)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Invalid password")
		return
	}

	respondWithJSON(w, http.StatusOK, response{
		User: User{
			ID:       user.ID,
			Email:    user.Email,
			JwtToken: ss,
		},
	})
}

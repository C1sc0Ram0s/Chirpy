package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/C1sc0Ram0s/Chirpy/internal/auth"
	"github.com/golang-jwt/jwt/v5"
)

func (cfg *apiConfig) handlerUsersUpdate(w http.ResponseWriter, r *http.Request) {
	type Claim struct {
		jwt.RegisteredClaims
	}
	type parameters struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	type returnVals struct {
		ID    int    `json:"id"`
		Email string `json:"email"`
	}
	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "invalid parameters")
	}

	tokenString := r.Header.Get("Authorization")
	tokenString = strings.ReplaceAll(tokenString, "Bearer ", "")

	token, err := jwt.ParseWithClaims(tokenString, &Claim{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(cfg.Jwt), nil
	})
	if err != nil {
		respondWithError(w, 401, "Unauthorized token")
		return
	}
	if claims, ok := token.Claims.(*Claim); ok {
		userID, err := strconv.Atoi(claims.Subject)
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, fmt.Sprintf("%v", err))
			return
		}

		hashedPassword, err := auth.HashPassword(params.Password)
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, "error hashing password")
			return
		}

		user, err := cfg.DB.UpdateUser(userID, params.Email, hashedPassword)
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, "error updating user")
			return
		}

		ret := returnVals{
			ID:    userID,
			Email: user.Email,
		}
		respondWithJSON(w, http.StatusOK, ret)
	} else {
		respondWithError(w, http.StatusInternalServerError, "token.Claims not ok")
		return
	}

}

package main

import (
	"net/http"
	"time"

	"github.com/C1sc0Ram0s/Chirpy/internal/auth"
)

func (cfg *apiConfig) handlerRefresh(w http.ResponseWriter, r *http.Request) {
	type response struct {
		Token string `json:"token"`
	}

	refreshToken, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "error retreiving bearer token")
		return
	}

	user, err := cfg.DB.GetRefreshTokenUser(refreshToken)
	if err != nil {
		respondWithError(w, 401, "token is expired or doesn't exist")
		return
	}

	accessToken, err := auth.MakeJWT(
		user.ID,
		cfg.Jwt,
		time.Hour,
	)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "error validating token")
		return
	}

	respondWithJSON(w, http.StatusOK, response{
		Token: accessToken,
	})

}

func (cfg *apiConfig) handlerRevoke(w http.ResponseWriter, r *http.Request) {
	refreshToken, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "error finding token")
		return
	}

	err = cfg.DB.RevokeToken(refreshToken)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "error revoking session")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

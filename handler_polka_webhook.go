package main

import (
	"encoding/json"
	"net/http"

	"github.com/C1sc0Ram0s/Chirpy/internal/auth"
)

func (cfg apiConfig) handlerPolkaWebhook(w http.ResponseWriter, r *http.Request) {
	type parameter struct {
		Event string `json:"event"`
		Data  struct {
			UserID int `json:"user_id"`
		}
	}

	apiKey, err := auth.GetApiKey(r.Header)
	if err != nil || apiKey != cfg.PolkaKey {
		respondWithError(w, 401, "malformed header")
		return
	}

	params := parameter{}
	decoder := json.NewDecoder(r.Body)
	err = decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters")
		return
	}

	if params.Event != "user.upgraded" {
		w.WriteHeader(204)
		return
	}

	err = cfg.DB.GrantChirpyRed(params.Data.UserID)
	if err != nil {
		respondWithError(w, 404, "user can't be found")
		return
	}

	w.WriteHeader(204)
}

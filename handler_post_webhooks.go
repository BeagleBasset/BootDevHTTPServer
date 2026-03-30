package main

import (
	"net/http"
	"database/sql"
	"errors"

	"github.com/google/uuid"
	"github.com/BeagleBasset/BootDevHTTPServer/internal/auth"
)

func (cfg *apiConfig) handlerWebhooks(w http.ResponseWriter, r *http.Request) {
	type data struct {
		UserId uuid.UUID `json:"user_id"`
	}
	type body struct {
		Event string `json:"event"`
		Data  data   `json:"data"`
	}	

    apiKey, err := auth.GetAPIKey(r.Header)
    if err != nil || apiKey != cfg.polkaKey {
        w.WriteHeader(http.StatusUnauthorized)
        return
    }

	params, err := decodeJSON[body](r)
	if err != nil {
		respondWithError(w, 500, "Error to get body:", err)
		return
	}

	if params.Event != "user.upgraded" {
		w.WriteHeader(204)
		return
	}

	err = cfg.dbQueries.UpgradeUser(r.Context(), params.Data.UserId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			respondWithError(w, 404, "User not found", err)
			return
		}
		respondWithError(w, 500, "Couldn't update user", err)
		return
	}
	w.WriteHeader(204)

}

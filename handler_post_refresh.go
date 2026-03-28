package main

import (
	"net/http"
	"time"

	"github.com/BeagleBasset/BootDevHTTPServer/internal/auth"
)

func (cfg *apiConfig) handlerRefresh(w http.ResponseWriter, r *http.Request) {
	type response struct{
		Token string `json:"token"`
	}

	token, err := auth.GetBearerToken(r.Header)	
	if err != nil {
		respondWithError(w, 401, "Error to get token:", err)
		return
	}

	user, err := cfg.dbQueries.GetUserFromRefreshToken(r.Context(), token)
	if err != nil {
		respondWithError(w, 401, "Error to user:", err)
		return
	}

	newToken, err := auth.MakeJWT(user.ID, cfg.jwtSecret, time.Hour)
	if err != nil {
		respondWithError(w, 401, "Error to make new Token:", err)
		return
	}

	respondWithJSON(w, 200, response{
		Token: newToken,
	})
}

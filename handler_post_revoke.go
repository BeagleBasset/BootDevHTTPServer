package main

import (
	"net/http"

	"github.com/BeagleBasset/BootDevHTTPServer/internal/auth"
)

func (cfg *apiConfig) handlerRevoke(w http.ResponseWriter, r *http.Request) {
	token, err := auth.GetBearerToken(r.Header)	
	if err != nil {
		respondWithError(w, 401, "Error to get token:", err)
		return
	}

	err = cfg.dbQueries.RevokeRefreshToken(r.Context(), token)
	if err != nil {
		respondWithError(w, 500, "Error revoking token:", err)
		return
	}

	w.WriteHeader(204)
}

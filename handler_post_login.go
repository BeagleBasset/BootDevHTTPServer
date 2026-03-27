package main

import (
	"net/http"

	"github.com/BeagleBasset/BootDevHTTPServer/internal/auth"
)

func (cfg *apiConfig) handlerLogin(w http.ResponseWriter, r *http.Request) {
	params, err := decodeJSON[NewUser](r)
	if err != nil {
		respondWithError(w, 500, "Error decode request body:", err)
		return
	}

	user, err := cfg.dbQueries.GetUser(r.Context(), params.Email)
	if err != nil {
		respondWithError(w, 401, "Error in query:", err)
		return
	}

	isPasswordCorrect, err := auth.CheckPasswordHash(params.Password, user.HashedPassword)
	if err != nil {
		respondWithError(w, 401, "Error in hash check:", err)
		return
	}
	
	if isPasswordCorrect {
		respondWithJSON(w, 200, User{
			ID:        user.ID,
			CreatedAt: user.CreatedAt,
			UpdatedAt: user.UpdatedAt,
			Email:     user.Email,
		})
	} else {
	    respondWithError(w, 401, "Incorrect email or password", nil)
		return	
	}
}


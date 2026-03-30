package main

import (
	"net/http"
	
	"github.com/BeagleBasset/BootDevHTTPServer/internal/auth"
	"github.com/BeagleBasset/BootDevHTTPServer/internal/database"
)

func (cfg *apiConfig) handlerUpdateUser(w http.ResponseWriter, r *http.Request) {
	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, 401, "Baj van a radarral:", err)
		return
	}

	userId, err := auth.ValidateJWT(token, cfg.jwtSecret)
	if err != nil {
		respondWithError(w, 401, "Error to get userid:", err)
		return
	}

	params, err := decodeJSON[NewUser](r)
	if err != nil {
		respondWithError(w, 500, "Error decode request body:", err)
		return
	}
	
	hashedPwd, err := auth.HashPassword(params.Password)
	if err != nil {
		respondWithError(w, 500, "Error with hash:", err)
		return
	}

	user, err := cfg.dbQueries.UpdateUser(r.Context(), database.UpdateUserParams{
		Email: params.Email,
		HashedPassword: hashedPwd,
		ID: userId,
	})
	if err != nil {
		respondWithError(w, 500, "Error create user:", err)
		return
	}

	respondWithJSON(w, 200, User{
		ID:        user.ID,
		Email:     user.Email,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	})
}

package main

import (
	"net/http"
	
	"github.com/BeagleBasset/BootDevHTTPServer/internal/auth"
	"github.com/BeagleBasset/BootDevHTTPServer/internal/database"
)

func (cfg *apiConfig) handlerUsers(w http.ResponseWriter, r *http.Request) {
	params, err := decodeJSON[NewUser](r)
	if err != nil {
		respondWithError(w, 500, "Error decode request body:", err)
		return
	}

	hashedPassword, err := auth.HashPassword(params.Password)
	if err != nil {
		respondWithError(w, 500, "Error hashing password:", err)
		return
	}

	user, err := cfg.dbQueries.CreateUser(r.Context(), database.CreateUserParams{
		Email: params.Email,
		HashedPassword: hashedPassword,
	})
	if err != nil {
		respondWithError(w, 500, "Error create user:", err)
		return
	}

	resp := User{
		ID:			user.ID,
		CreatedAt: 	user.CreatedAt,
		UpdatedAt: 	user.UpdatedAt,
		Email: 		user.Email,
	}

	respondWithJSON(w, 201, resp)
}

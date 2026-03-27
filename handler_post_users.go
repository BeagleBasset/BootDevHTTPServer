package main

import (
	"net/http"
)

func (cfg *apiConfig) handlerUsers(w http.ResponseWriter, r *http.Request) {
	params, err := decodeJSON[Email](r)
	if err != nil {
		respondWithError(w, 500, "Error decode request body:", err)
		return
	}

	user, err := cfg.dbQueries.CreateUser(r.Context(), params.Email)
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

	respondWithJSON(w, 200, resp)
}

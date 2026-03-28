package main

import (
	"net/http"
	"time"

	"github.com/BeagleBasset/BootDevHTTPServer/internal/auth"
	"github.com/BeagleBasset/BootDevHTTPServer/internal/database"
)

func (cfg *apiConfig) handlerLogin(w http.ResponseWriter, r *http.Request) {
	type response struct {
        User
        Token string `json:"token"`
		RefreshToken string `json:"refresh_token"`
    }
	params, err := decodeJSON[NewUser](r)
	if err != nil {
		respondWithError(w, 500, "Error decode request body:", err)
		return
	}

	expirationTime := time.Duration(3600) * time.Second

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
	
	if !isPasswordCorrect {
	    respondWithError(w, 401, "Incorrect email or password", nil)
		return	
	}

	jwtToken, err := auth.MakeJWT(user.ID, cfg.jwtSecret, expirationTime)
	if err != nil {
		respondWithError(w, 401, "Error in make JWT:", err)
		return
	}

	refreshedToken := auth.MakeRefreshToken()
	_, err = cfg.dbQueries.CreateRefreshToken(r.Context(), database.CreateRefreshTokenParams{
		Token: refreshedToken,
		UserID: user.ID,
		ExpiresAt: time.Now().Add(time.Hour * 24 * 60),
	})
	if err != nil {
		respondWithError(w, 401, "Error create refresh token:", err)
		return
	}

	respondWithJSON(w, 200, response{ 
		User : User{
				ID:        user.ID,
				CreatedAt: user.CreatedAt,
				UpdatedAt: user.UpdatedAt,
				Email:     user.Email,
				},
		Token: jwtToken,
		RefreshToken: refreshedToken,
	})
}


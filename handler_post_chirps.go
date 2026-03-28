package main

import (
	"net/http"
	"strings"
	"errors"

    "github.com/BeagleBasset/BootDevHTTPServer/internal/database"
    "github.com/BeagleBasset/BootDevHTTPServer/internal/auth"
)

func cleanResponse(message string) string {
	badWords := map[string]struct{}{
		"kerfuffle": {},
		"sharbert":  {},
		"fornax":    {},
	}
	words := strings.Split(message, " ")

	for i, word := range words {
		word = strings.ToLower(word)
		_, ok := badWords[word]
		if ok {
			words[i] = "****"
		}
	}
	
	message = strings.Join(words, " ")
	return message
}


func validateChirp(body string) (string, error) {
    if len(body) > 140 {
        return "", errors.New("Chirp is too long")
    }
    cleaned := cleanResponse(body)
    return cleaned, nil
}

func (cfg *apiConfig) handlerChirps(w http.ResponseWriter, r *http.Request) {
	params, err := decodeJSON[NewChirp](r)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Couldn't decode parameters", err)
		return
	}

	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Error with token:", err)
		return
	}
	userID, err := auth.ValidateJWT(token, cfg.jwtSecret)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Error with userId:", err)
		return
	}

	chirp, err := validateChirp(params.Body)
	if err != nil {
		respondWithError(w, 500, "Error in chirp validation:", err)
		return
	}
	
	dbChirp, err := cfg.dbQueries.CreateChirp(r.Context(), database.CreateChirpParams{
		Body:   chirp,
		UserID: userID,
	})
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Error in query:", err)
		return
	}

	resp := Chirp{
		ID:			dbChirp.ID,
		CreatedAt: 	dbChirp.CreatedAt,
		UpdatedAt: 	dbChirp.UpdatedAt,
		Body: 		dbChirp.Body,
		UserID:		dbChirp.UserID,
	}

	respondWithJSON(w, 201, resp)
}

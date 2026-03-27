package main

import (
	"net/http"
)

func (cfg *apiConfig) handlerGetAllChirps(w http.ResponseWriter, r *http.Request) {
	dbChirps, err := cfg.dbQueries.GetAllChirps(r.Context())
	if err != nil {
		respondWithError(w, 500, "Hiba a chirp-ek lekérdezése közben", err)	
		return
	}

	resps := []Chirp{}
	for _, dbChirp := range dbChirps {
		resp := Chirp{
			ID:			dbChirp.ID,
			CreatedAt: 	dbChirp.CreatedAt,
			UpdatedAt: 	dbChirp.UpdatedAt,
			Body: 		dbChirp.Body,
			UserID:		dbChirp.UserID,
		}
		resps = append(resps, resp)
	}

	respondWithJSON(w, 200, resps)
}

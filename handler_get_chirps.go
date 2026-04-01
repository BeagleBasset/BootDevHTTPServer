package main

import (
	"net/http"
	"sort"

	"github.com/google/uuid"
)

func (cfg *apiConfig) handlerGetAllChirps(w http.ResponseWriter, r *http.Request) {
		sortParam := r.URL.Query().Get("sort")
	if sortParam == "" {
		sortParam = "asc"
	}

	authorIDString := r.URL.Query().Get("author_id")
	if authorIDString != "" {
		authorID, err := uuid.Parse(authorIDString)
		if err != nil {
			respondWithError(w, 400, "Error in parsing:", err)
			return
		}

		dbChirps, err := cfg.dbQueries.GetChirpFromUser(r.Context(), authorID)
		if err != nil {
			respondWithError(w, 500, "Error in query:", err)
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

	} else {
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

		sort.Slice(resps, func(i, j int) bool {
			if sortParam == "desc" {
				return resps[i].CreatedAt.After(resps[j].CreatedAt)
			}
			// default asc
			return resps[i].CreatedAt.Before(resps[j].CreatedAt)
		})
		respondWithJSON(w, 200, resps)
	}
}

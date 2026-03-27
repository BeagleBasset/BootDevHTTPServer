package main

import (
	"net/http"
	"log"
)

func (cfg *apiConfig) handlerReset(w http.ResponseWriter, r *http.Request) {
	if cfg.platform != "dev" {
		w.Header().Add("Content-Type", "text/plain; charset=utf-8")
		w.WriteHeader(http.StatusForbidden)
		w.Write([]byte("OK"))
	} else {
		cfg.fileserverHits.Store(0)
		err := cfg.dbQueries.Reset(r.Context())
		if err != nil {
			log.Printf("Error in Reset: %s", err)
			w.WriteHeader(500)
			return
		}

		w.Header().Add("Content-Type", "text/plain; charset=utf-8")
		w.WriteHeader(200)
		w.Write([]byte("OK"))
	}
}

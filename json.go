package main

import (
    "encoding/json"
    "log"
    "net/http"
)

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
    dat, err := json.Marshal(payload)
    if err != nil {
        log.Printf("Error marshalling JSON: %s", err)
        w.WriteHeader(500)
        return
    }
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(code)
    w.Write(dat)
}

func respondWithError(w http.ResponseWriter, code int, msg string, err error) {
    log.Printf("Error: %s - %v", msg, err)
    respondWithJSON(w, code, map[string]string{"error": msg})
}

func decodeJSON[T any](r *http.Request) (T, error) {
    var v T
    decoder := json.NewDecoder(r.Body)
    err := decoder.Decode(&v)
    return v, err
}

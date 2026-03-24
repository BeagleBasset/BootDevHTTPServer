package main

import (
	"fmt"
	"net/http"
	"sync/atomic"
)

type apiConfig struct {
	fileserverHits atomic.Int32
}

func main() {
	// Create a new ServeMux
	mux := http.NewServeMux()

	// Create a file server handler
	fileServer := http.FileServer(http.Dir("."))

	handler := http.StripPrefix("/app", fileServer)
	apiCfg := apiConfig{}

	// Register handler for root path
	mux.Handle("/app/", apiCfg.middlewareMetricsInc(handler))
	
	mux.HandleFunc("GET 	/healthz", 	handlerReadiness)
	mux.HandleFunc("GET  	/metrics", 	apiCfg.handlerMetrics)
	mux.HandleFunc("POST 	/reset", 	apiCfg.handlerReset)

	// Create a new Server
	server := http.Server{
		Addr:    ":8080",
		Handler: mux,
	}


	// Start the server
	server.ListenAndServe()
}

func handlerReadiness(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(200)
	w.Write([]byte("OK"))
}

func (cfg *apiConfig) handlerMetrics(w http.ResponseWriter, r *http.Request) {
	resp := fmt.Sprintf("Hits: %d", cfg.fileserverHits.Load())

	w.Header().Add("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(200)
	w.Write([]byte(resp))
}

func (cfg *apiConfig) handlerReset(w http.ResponseWriter, r *http.Request) {
	cfg.fileserverHits.Store(0)

	w.Header().Add("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(200)
	w.Write([]byte("OK"))
}

func (cfg *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cfg.fileserverHits.Add(1)
        next.ServeHTTP(w, r)
    })
}

package main

import (
	"net/http"
)

func main() {
	// Create a new ServeMux
	mux := http.NewServeMux()

	// Create a file server handler
	fileServer := http.FileServer(http.Dir("."))

	// Register handler for root path
	mux.Handle("/app/", http.StripPrefix("/app", fileServer))
	
	mux.HandleFunc("/healthz", handlerReadiness)

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

package main

import (
	"net/http"
)

func main() {
	// Create a new ServeMux
	serveMux := http.NewServeMux()

	// Create a file server handler
	fileServer := http.FileServer(http.Dir("."))

	// Register handler for root path
	serveMux.Handle("/", fileServer)

	// Create a new Server
	server := http.Server{
		Addr:    ":8080",
		Handler: serveMux,
	}

	// Start the server
	server.ListenAndServe()
}

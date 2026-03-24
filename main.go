package main

import "net/http"

func main() {
	// Create a new ServeMux
	serveMux := http.NewServeMux()

	// Create a new Server
	server := http.Server{
		Addr:    ":8080",
		Handler: serveMux,
	}

	// Start the server
	server.ListenAndServe()
}



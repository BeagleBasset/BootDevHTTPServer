package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"

	_ "github.com/lib/pq"
	"github.com/joho/godotenv"
    "github.com/BeagleBasset/BootDevHTTPServer/internal/database"
)


func main() {
	godotenv.Load()
	dbURL := os.Getenv("DB_URL")

	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Printf("Error open SQL database: %s", err)
		return
	}
	dbQueries := database.New(db)
	// Create a new ServeMux
	mux := http.NewServeMux()

	// Create a file server handler
	fileServer := http.FileServer(http.Dir("."))

	handler := http.StripPrefix("/app", fileServer)
	apiCfg := apiConfig{
		dbQueries: 	dbQueries,
		platform:	os.Getenv("PLATFORM"),
	}

	// Register handler for root path
	mux.Handle("/app/", apiCfg.middlewareMetricsInc(handler))
	
	mux.HandleFunc("GET 	/api/healthz",				handlerReadiness)
	mux.HandleFunc("GET  	/admin/metrics", 			apiCfg.handlerMetrics)
	mux.HandleFunc("GET  	/api/chirps", 				apiCfg.handlerGetAllChirps)
	mux.HandleFunc("GET  	/api/chirps/{chirpID}",		apiCfg.handlerGetChirp)
	mux.HandleFunc("POST 	/admin/reset", 				apiCfg.handlerReset)
	mux.HandleFunc("POST 	/api/users",		 		apiCfg.handlerUsers)
	mux.HandleFunc("POST 	/api/chirps",		 		apiCfg.handlerChirps)

	// Create a new Server
	server := http.Server{
		Addr:    ":8080",
		Handler: mux,
	}


	// Start the server
	server.ListenAndServe()
}


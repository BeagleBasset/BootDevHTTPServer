package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"sync/atomic"
	"time"

	_ "github.com/lib/pq"
	"github.com/joho/godotenv"
	"github.com/google/uuid"
    "github.com/BeagleBasset/BootDevHTTPServer/internal/database"
)

type apiConfig struct {
	fileserverHits 	atomic.Int32
	dbQueries      	*database.Queries
	platform		string
}

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
	mux.HandleFunc("POST 	/admin/reset", 				apiCfg.handlerReset)
	mux.HandleFunc("POST 	/api/validate_chirp", 		handlerValidateChirp)
	mux.HandleFunc("POST 	/api/users",		 		apiCfg.handlerUsers)

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

func handlerValidateChirp(w http.ResponseWriter, r *http.Request) {
	type validResponse struct {
		Valid bool `json:"valid"`
	}

	type errorResponse struct {
		Error string `json:"error"`
	}

	type parameters struct {
		Body string `json:"body"`
	}

	type cleanedResponse struct {
		CleanedBody string `json:"cleaned_body"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		log.Printf("Error in decoding: %s", err)
		w.WriteHeader(500)
		return
	}

	if len(params.Body) <= 140 {
		respText := cleanResponse(params.Body)
		respBody := cleanedResponse{CleanedBody: respText}
		dat, err := json.Marshal(respBody)
		if err != nil {
			log.Printf("Error marshalling JSON: %s", err)
			w.WriteHeader(500)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(200)
		w.Write(dat)
	} else {
		respBody := errorResponse{Error: "Chirp is too long"}
		dat, err := json.Marshal(respBody)
		if err != nil {
			log.Printf("Error marshalling JSON: %s", err)
			w.WriteHeader(500)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(400)
		w.Write(dat)
	}
}

func (cfg *apiConfig) handlerUsers(w http.ResponseWriter, r *http.Request) {
	type validResponse struct {
		ID        uuid.UUID `json:"id"`
		CreatedAt time.Time `json:"created_at"`
		UpdatedAt time.Time `json:"updated_at"`
		Email     string    `json:"email"`
	}

	type errorResponse struct {
		Error string `json:"error"`
	}

	type parameters struct {
		Email string `json:"email"`
	}

	params 	:= parameters{}
	decoder := json.NewDecoder(r.Body)
	err 	:= decoder.Decode(&params)
	if err != nil {
		log.Printf("Error decode request body: %s", err)
		w.WriteHeader(500)
		return
	}

	user, err := cfg.dbQueries.CreateUser(r.Context(), params.Email)
	if err != nil {
		log.Printf("Error create user: %s", err)
		w.WriteHeader(500)
		return
	}

	resp := validResponse{
		ID:			user.ID,
		CreatedAt: 	user.CreatedAt,
		UpdatedAt: 	user.UpdatedAt,
		Email: 		user.Email,
	}

	dat, err := json.Marshal(resp)
	if err != nil {
		log.Printf("Error marshalling JSON: %s", err)
		w.WriteHeader(500)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	w.Write(dat)
}
func (cfg *apiConfig) handlerMetrics(w http.ResponseWriter, r *http.Request) {
	resp := fmt.Sprintf(`
	<html>
  		<body>
    		<h1>Welcome, Chirpy Admin</h1>
    		<p>Chirpy has been visited %d times!</p>
  		</body>
	</html>`, cfg.fileserverHits.Load())

	w.Header().Add("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(200)
	w.Write([]byte(resp))
}

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

func (cfg *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cfg.fileserverHits.Add(1)
        next.ServeHTTP(w, r)
    })
}

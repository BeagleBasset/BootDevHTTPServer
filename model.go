package main

import (
	"time"
	"sync/atomic"

	"github.com/google/uuid"
    "github.com/BeagleBasset/BootDevHTTPServer/internal/database"
)

// main api struct
type apiConfig struct {
	fileserverHits 	atomic.Int32
	dbQueries      	*database.Queries
	platform		string
}

// Chirp struct, for chirp JSON respons
type Chirp struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Body      string    `json:"body"`
	UserID    uuid.UUID `json:"user_id"`
}

type NewChirp struct {
	Body string `json:"body"`
	UserID uuid.UUID `json:"user_id"`
}

type User struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Email     string    `json:"email"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}

type Email struct {
	Email string `json:"email"`
}

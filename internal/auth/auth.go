package auth

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/alexedwards/argon2id"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

func HashPassword(password string) (string, error) {
	hash, err := argon2id.CreateHash(password, argon2id.DefaultParams)
	if err != nil {
		return "", err
	}
	return hash, nil
}

func CheckPasswordHash(password, hash string) (bool, error) {
	match, err := argon2id.ComparePasswordAndHash(password, hash)
	if err != nil {
		return false, err
	}
	return match, nil
}

func MakeJWT(userID uuid.UUID, tokenSecret string, expiresIn time.Duration) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
		Issuer: "chirpy-access",
		IssuedAt: jwt.NewNumericDate(time.Now().UTC()),
		ExpiresAt: jwt.NewNumericDate(time.Now().UTC().Add(expiresIn)),
		Subject: userID.String(),})

	tokenString, err := token.SignedString([]byte(tokenSecret))

	if err != nil {
		return "", err
	} else {
		return tokenString, nil
	}
}

func ValidateJWT(tokenString, tokenSecret string) (uuid.UUID, error) {
	token, err := jwt.ParseWithClaims(tokenString, &jwt.RegisteredClaims{}, func(token *jwt.Token) (interface{}, error) {
			return []byte(tokenSecret), nil
		})
	if err != nil {
		return uuid.UUID{}, err
	}

	claims, ok := token.Claims.(*jwt.RegisteredClaims)
	if ok {
		userID, err := uuid.Parse(claims.Subject)
		if err != nil {
			return uuid.UUID{}, err 
		}
		return userID, nil 
	}

	return uuid.UUID{}, errors.New("Failed to validate")
}

func GetBearerToken(headers http.Header) (string, error){
	authorization := headers.Get("Authorization")
	authWords := strings.Split(authorization, " ")
	var token string
	var ErrNoAuthHeaderIncluded = errors.New("no auth header included in request")
	if len(authWords) != 2 {
		return "", ErrNoAuthHeaderIncluded
	}
	if authWords[0] != "Bearer" {
		return "", ErrNoAuthHeaderIncluded
	}

	token = authWords[1]
	if token != "" {
		return token, nil
	}	

	return "", ErrNoAuthHeaderIncluded
}

func MakeRefreshToken() string {
	var bits [32]byte
	_, err := rand.Read(bits[:])
	if err != nil {
		return ""
	}
	return hex.EncodeToString(bits[:])
}

package auth

import (
	"testing"
	"time"

	"github.com/google/uuid"
)


func TestMakeAndValidateJWT(t *testing.T) {
	userID := uuid.New()
	token, err := MakeJWT(userID, "secret", time.Hour)
	if err != nil {
		t.Fatalf("Failed to make JWT: %v", err)
	}
	validatedUserID, err := ValidateJWT(token, "secret")
	if err != nil {
		t.Fatalf("Failed to validate: %v", err)
	}
	if validatedUserID != userID {
		t.Errorf("expected %v got %v", userID, validatedUserID)
	}
}

func TestExpiredJWT(t *testing.T) {
	userID := uuid.New()
	token, err := MakeJWT(userID, "secret", -time.Hour)
	if err != nil {
		t.Fatalf("Failed to make JWT: %v", err)
	} 

	_, err = ValidateJWT(token, "secret")
	if err == nil {
		t.Errorf("expected an error for expired token, got nil")
	} 
}

func TestWrongSecret(t *testing.T) {
	userID := uuid.New()
	token, err := MakeJWT(userID, "secret1", time.Hour)
	if err != nil {
		t.Fatalf("Failed to make JWT: %v", err)
	} 

	_, err = ValidateJWT(token, "secret2")
	if err == nil {
		t.Errorf("expected an error for wrong secret, got nil")
	} 
}

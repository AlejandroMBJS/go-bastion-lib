// tests/jwt_test.go
package tests

import (
	"testing"
	"time"

	"github.com/alejandrombjs/go-bastion-lib/pkg/security"
)

func TestJWTGeneration(t *testing.T) {
	secret := "test-secret"
	sub := "user123"

	token, err := security.GenerateAccessToken(sub, time.Hour, secret, nil)
	if err != nil {
		t.Fatalf("Failed to generate token: %v", err)
	}

	if token == "" {
		t.Fatal("Token should not be empty")
	}
}

func TestJWTValidation(t *testing.T) {
	secret := "test-secret"
	sub := "user123"

	token, err := security.GenerateAccessToken(sub, time.Hour, secret, nil)
	if err != nil {
		t.Fatalf("Failed to generate token: %v", err)
	}

	claims, err := security.ParseAndValidateToken(token, secret)
	if err != nil {
		t.Fatalf("Failed to validate token: %v", err)
	}

	if claims.Subject != sub {
		t.Errorf("Expected subject %s, got %s", sub, claims.Subject)
	}
}

func TestJWTExpired(t *testing.T) {
	secret := "test-secret"
	sub := "user123"

	// Generate token that expired 1 hour ago
	token, err := security.GenerateAccessToken(sub, -time.Hour, secret, nil)
	if err != nil {
		t.Fatalf("Failed to generate token: %v", err)
	}

	_, err = security.ParseAndValidateToken(token, secret)
	if err != security.ErrExpiredToken {
		t.Errorf("Expected expired token error, got: %v", err)
	}
}

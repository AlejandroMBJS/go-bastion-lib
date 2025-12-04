// Package security provides security-related utilities.
package security

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var (
	ErrInvalidToken = errors.New("invalid token")
	ErrExpiredToken = errors.New("expired token")
)

// Claims represents JWT claims.
type Claims struct {
	Subject   string                 `json:"sub"`
	ExpiresAt time.Time              `json:"exp"`
	IssuedAt  time.Time              `json:"iat"`
	Scopes    []string               `json:"scopes,omitempty"`
	Extra     map[string]interface{} `json:"extra,omitempty"`
}

// GenerateAccessToken generates a new JWT access token.
func GenerateAccessToken(sub string, ttl time.Duration, secret string, extraClaims map[string]any) (string, error) {
	now := time.Now()

	claims := jwt.MapClaims{
		"sub": sub,
		"exp": now.Add(ttl).Unix(),
		"iat": now.Unix(),
	}

	// Add extra claims
	for k, v := range extraClaims {
		claims[k] = v
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}

// ParseAndValidateToken parses and validates a JWT token.
func ParseAndValidateToken(tokenString, secret string) (*Claims, error) {
	// Parse the token
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Validate signing method
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, ErrInvalidToken
		}
		return []byte(secret), nil
	})

	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			return nil, ErrExpiredToken
		}
		return nil, ErrInvalidToken
	}

	// Validate token
	if !token.Valid {
		return nil, ErrInvalidToken
	}

	// Extract claims
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, ErrInvalidToken
	}

	// Build Claims struct
	result := &Claims{}

	// Extract subject
	if sub, ok := claims["sub"].(string); ok {
		result.Subject = sub
	} else {
		return nil, ErrInvalidToken
	}

	// Extract expiration time
	if exp, ok := claims["exp"].(float64); ok {
		result.ExpiresAt = time.Unix(int64(exp), 0)
	}

	// Extract issued at time
	if iat, ok := claims["iat"].(float64); ok {
		result.IssuedAt = time.Unix(int64(iat), 0)
	}

	// Extract scopes
	if scopes, ok := claims["scopes"].([]interface{}); ok {
		result.Scopes = make([]string, len(scopes))
		for i, s := range scopes {
			if str, ok := s.(string); ok {
				result.Scopes[i] = str
			}
		}
	}

	// Extract extra claims
	result.Extra = make(map[string]interface{})
	for k, v := range claims {
		switch k {
		case "sub", "exp", "iat", "scopes":
			// Skip standard claims
		default:
			result.Extra[k] = v
		}
	}

	return result, nil
}

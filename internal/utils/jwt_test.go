package utils

import (
	"testing"

	"github.com/golang-jwt/jwt/v5"
)

func TestGenerateToken(t *testing.T) {
	tokenStr, err := GenerateToken("user-123", "moderator")
	if err != nil {
		t.Fatalf("failed to generate token: %v", err)
	}

	token, err := jwt.Parse(tokenStr, func(t *jwt.Token) (interface{}, error) {
		return []byte("secret-key"), nil
	})
	if err != nil || !token.Valid {
		t.Fatalf("token is invalid: %v", err)
	}

	claims := token.Claims.(jwt.MapClaims)
	if claims["user_id"] != "user-123" {
		t.Errorf("expected user_id to be 'user-123', got %v", claims["user_id"])
	}
	if claims["role"] != "moderator" {
		t.Errorf("expected role to be 'moderator', got %v", claims["role"])
	}
}

package auth

import (
	"testing"
)

func TestNewID(t *testing.T) {
	id := NewID()
	if id == "" {
		t.Error("NewID returned empty string")
	}
}

func TestGenerateAndValidateToken(t *testing.T) {
	j := NewJWT("test-secret-key")

	token, err := j.GenerateToken("user-1", "test@test.com", "venue")
	if err != nil {
		t.Fatalf("GenerateToken failed: %v", err)
	}

	if token == "" {
		t.Fatal("token is empty")
	}

	claims, err := j.ValidateToken(token)
	if err != nil {
		t.Fatalf("ValidateToken failed: %v", err)
	}

	if claims.UserID != "user-1" {
		t.Errorf("expected user-1, got %s", claims.UserID)
	}
	if claims.Email != "test@test.com" {
		t.Errorf("expected test@test.com, got %s", claims.Email)
	}
	if claims.Role != "venue" {
		t.Errorf("expected venue, got %s", claims.Role)
	}
}

func TestValidateInvalidToken(t *testing.T) {
	j := NewJWT("test-secret")

	_, err := j.ValidateToken("invalid-token")
	if err == nil {
		t.Error("expected error for invalid token")
	}
}

func TestValidateTokenWrongSecret(t *testing.T) {
	j1 := NewJWT("secret-1")
	j2 := NewJWT("secret-2")

	token, _ := j1.GenerateToken("user-1", "test@test.com", "artist")

	_, err := j2.ValidateToken(token)
	if err == nil {
		t.Error("expected error for wrong secret")
	}
}

package utils

import (
	"testing"
	"time"
)

func TestGenerateAndParseToken(t *testing.T) {
	secret := "test-secret-key"
	token, err := GenerateToken(42, "demo@sasivision.com", "admin", secret, time.Hour)
	if err != nil {
		t.Fatalf("GenerateToken: %v", err)
	}

	claims, err := ParseToken(token, secret)
	if err != nil {
		t.Fatalf("ParseToken: %v", err)
	}
	if claims.UserID != 42 {
		t.Fatalf("user_id = %d, want 42", claims.UserID)
	}
	if claims.Email != "demo@sasivision.com" {
		t.Fatalf("email = %q", claims.Email)
	}
	if claims.Role != "admin" {
		t.Fatalf("role = %q", claims.Role)
	}
}

func TestParseTokenInvalidSecret(t *testing.T) {
	token, err := GenerateToken(1, "a@b.com", "user", "secret-a", time.Hour)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := ParseToken(token, "secret-b"); err == nil {
		t.Fatal("expected error for invalid secret")
	}
}

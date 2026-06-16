package utils

import "testing"

func TestHashAndCheckPassword(t *testing.T) {
	hash, err := HashPassword("Sasivision123")
	if err != nil {
		t.Fatalf("HashPassword: %v", err)
	}
	if !CheckPassword("Sasivision123", hash) {
		t.Fatal("expected password to match hash")
	}
	if CheckPassword("wrong-password", hash) {
		t.Fatal("expected wrong password to fail")
	}
}

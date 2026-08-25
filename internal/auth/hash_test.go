package auth 

import "testing"

func TestHashPassword(t *testing.T) {
	password := "1234"

	// 1. Test hashing the password
	hash, err := HashPassword(password)
	if err != nil {
		t.Fatalf("expected no error hashing password, got %v", err)
	}
	if hash == "" {
		t.Fatal("expected hash string to not be empty")
	}

	// 2. Test checking with the correct password
	match, err := CheckPasswordHash(password, hash)
	if err != nil {
		t.Fatalf("expected no error checking correct password, got %v", err)
	}
	if !match {
		t.Errorf("expected password to match hash, but got false")
	}

	// 3. Test checking with an incorrect password
	match, err = CheckPasswordHash("wrongpassword", hash)
	if err != nil {
		t.Fatalf("expected no error checking wrong password, got %v", err)
	}
	if match {
		t.Errorf("expected incorrect password to return false, but got true")
	}
}
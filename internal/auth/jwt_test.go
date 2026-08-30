package auth

import (
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestJWT(t *testing.T) {
	tests := []struct {
		name string
		tokenSecret string
		validateSecret string
		userID uuid.UUID
		expiresIn time.Duration
		waitDuration time.Duration
		wantErr bool
	}{
		{name: "happy path", tokenSecret: "Trying_with_an_easy_token", validateSecret: "Trying_with_an_easy_token", userID: uuid.New(), expiresIn: 60 * time.Second, waitDuration: 0, wantErr: false},
		{name: "expired path", tokenSecret: "Trying_with_an_easy_token", validateSecret: "Trying_with_an_easy_token", userID: uuid.New(), expiresIn: 1 * time.Second, waitDuration: 2 * time.Second, wantErr: true},
		{name: "wrong validate token path", tokenSecret: "Trying_with_an_easy_token", validateSecret: "Trying_with_a_difficult_token", userID: uuid.New(), expiresIn: 60 * time.Second, waitDuration: 0, wantErr: true},
	}
	for _, tt := range tests {
		// 1. Test making JWT
		t.Run(tt.name, func(t *testing.T) {
			tokenString, err := MakeJWT(tt.userID, tt.tokenSecret, tt.expiresIn)
			if err != nil {
				t.Fatalf("expected no error making JWT, got %v", err)
			}
			
			time.Sleep(tt.waitDuration)

			// 2. Test validating JWT
			returnID, err := ValidateJWT(tokenString, tt.validateSecret)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateJWT() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr {
				return
			}
			if returnID != tt.userID {
				t.Fatalf("expected no difference between user and returned ID")
			}
		})
	}

}
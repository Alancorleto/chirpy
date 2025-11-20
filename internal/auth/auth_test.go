package auth

import (
	"net/http"
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestValidateJWT(t *testing.T) {
	userID := uuid.New()
	tokenSecret := "secret"
	tokenExpiresIn := 2 * time.Second

	tokenString, _ := MakeJWT(userID, tokenSecret, tokenExpiresIn)

	tests := []struct {
		name        string
		tokenString string
		tokenSecret string
		wantUserID  uuid.UUID
		wantErr     bool
	}{
		{
			name:        "Valid token",
			tokenString: tokenString,
			tokenSecret: "secret",
			wantUserID:  userID,
			wantErr:     false,
		},
		{
			name:        "Invalid token",
			tokenString: "invalid.token.string",
			tokenSecret: "secret",
			wantUserID:  uuid.Nil,
			wantErr:     true,
		},
		{
			name:        "Wrong secret",
			tokenString: tokenString,
			tokenSecret: "wrong_secret",
			wantUserID:  uuid.Nil,
			wantErr:     true,
		},
	}

	for _, testCase := range tests {
		gotUserID, err := ValidateJWT(testCase.tokenString, testCase.tokenSecret)
		if (err != nil) != testCase.wantErr {
			t.Errorf("Different error value from expected.\nexpected error: %v\nactual error: %v\n", testCase.wantErr, err)
			return
		}
		if gotUserID != testCase.wantUserID {
			t.Errorf("Different user ID from expected.\nexpected user ID: %v\nactual user ID: %v\n", testCase.wantUserID, gotUserID)
		}
	}
}

func TestGetBearerToken(t *testing.T) {
	tests := []struct {
		key       string
		value     string
		wantToken string
		wantErr   bool
	}{
		{
			key:       "Authorization",
			value:     "Bearer 123093450934",
			wantToken: "123093450934",
			wantErr:   false,
		},
		{
			key:       "Bad key",
			value:     "Bearer 123093450934",
			wantToken: "",
			wantErr:   true,
		},
		{
			key:       "Authorization",
			value:     "Bearer123093450934",
			wantToken: "",
			wantErr:   true,
		},
	}
	for _, testCase := range tests {
		header := http.Header{}
		header.Add(testCase.key, testCase.value)
		gotToken, err := GetBearerToken(header)
		if (err != nil) != testCase.wantErr {
			t.Errorf("Different error value from expected.\nexpected error: %v\nactual error: %v\n", testCase.wantErr, err)
			return
		}
		if gotToken != testCase.wantToken {
			t.Errorf("Different token string from expected.\nexpected: %v\nactual: %v\n", testCase.wantToken, gotToken)
		}
	}
}

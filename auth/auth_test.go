package auth

import (
	"os"
	"testing"
	"time"

	"golang.org/x/oauth2"
)

func TestSaveAndTokenFromFile(t *testing.T) {
	tmpFile := "test_token.json"
	defer os.Remove(tmpFile)

	token := &oauth2.Token{
		AccessToken:  "test-access-token",
		TokenType:    "Bearer",
		RefreshToken: "test-refresh-token",
		Expiry:       time.Now().Add(time.Hour),
	}

	// Test saveToken
	if err := saveToken(tmpFile, token); err != nil {
		t.Fatalf("saveToken failed: %v", err)
	}

	// Test tokenFromFile
	readToken, err := tokenFromFile(tmpFile)
	if err != nil {
		t.Fatalf("tokenFromFile failed: %v", err)
	}

	if readToken.AccessToken != token.AccessToken {
		t.Errorf("AccessToken mismatch: got %v, want %v", readToken.AccessToken, token.AccessToken)
	}
	if readToken.TokenType != token.TokenType {
		t.Errorf("TokenType mismatch: got %v, want %v", readToken.TokenType, token.TokenType)
	}
	if readToken.RefreshToken != token.RefreshToken {
		t.Errorf("RefreshToken mismatch: got %v, want %v", readToken.RefreshToken, token.RefreshToken)
	}
	if !readToken.Expiry.Equal(token.Expiry) {
		t.Errorf("Expiry mismatch: got %v, want %v", readToken.Expiry, token.Expiry)
	}
}

func TestTokenFromFileNotFound(t *testing.T) {
	_, err := tokenFromFile("nonexistent_token.json")
	if err == nil {
		t.Error("expected error for nonexistent file, got nil")
	}
}

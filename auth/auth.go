// Package auth handles OAuth2 authentication and token management for Google APIs.
package auth

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

// tokenFromFile reads an OAuth2 token from a file.
func tokenFromFile(file string) (*oauth2.Token, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer func() {
		cerr := f.Close()
		if cerr != nil {
			log.Printf("error closing token file: %v", cerr)
		}
	}()
	tok := &oauth2.Token{}
	if err := json.NewDecoder(f).Decode(tok); err != nil {
		return nil, err
	}
	return tok, nil
}

// saveToken saves an OAuth2 token to a file.
func saveToken(path string, token *oauth2.Token) error {
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer func() {
		cerr := f.Close()
		if cerr != nil {
			log.Printf("error closing token file: %v", cerr)
		}
	}()
	return json.NewEncoder(f).Encode(token)
}

// GetClient returns an authenticated HTTP client using credentials and token files.
func GetClient(ctx context.Context, credsPath, tokenPath, scope string) (*http.Client, error) {
	credBytes, err := os.ReadFile(credsPath)
	if err != nil {
		return nil, err
	}
	config, err := google.ConfigFromJSON(credBytes, scope)
	if err != nil {
		return nil, err
	}
	tok, err := tokenFromFile(tokenPath)
	if err != nil {
		// Token not found, get one from web
		url := config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
		log.Printf("Go to the following link in your browser then type the authorization code: \n%v\n", url)
		var code string
		if _, err := fmt.Scan(&code); err != nil {
			return nil, err
		}
		tok, err = config.Exchange(ctx, code)
		if err != nil {
			return nil, err
		}
		if err := saveToken(tokenPath, tok); err != nil {
			return nil, err
		}
	}
	return config.Client(ctx, tok), nil
}

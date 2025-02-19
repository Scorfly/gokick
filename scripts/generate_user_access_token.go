package main

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
)

// Put your scopes here.
var kickScopes = []string{
	"user:read",
}

var ENDPOINT = struct {
	AuthURL  string
	TokenURL string
}{
	AuthURL:  "https://id.kick.com/oauth/authorize",
	TokenURL: "https://id.kick.com/oauth/token",
}

// Put your redirect URL here.
const redirectURL = "http://localhost:3000/oauth/kick/callback"

// Client ID and Secret from environment variables.
var (
	kickClientID     = os.Getenv("KICK_CLIENT_ID")
	kickClientSecret = os.Getenv("KICK_CLIENT_SECRET")
)

// PKCE Helper Functions.
// Generate a random code verifier.
func generateCodeVerifier() (string, error) {
	buffer := make([]byte, 32)
	_, err := rand.Read(buffer)
	if err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(buffer), nil
}

// Generate a code challenge (SHA-256 hash of the verifier).
func generateCodeChallenge(verifier string) string {
	hash := sha256.Sum256([]byte(verifier))
	return base64.URLEncoding.EncodeToString(hash[:])
}

// Step 1: Generate a verifier, challenge, and send the user to the auth page.
func oauthKickHandler(w http.ResponseWriter, r *http.Request) {
	codeVerifier, err := generateCodeVerifier()
	if err != nil {
		http.Error(w, "Failed to generate code verifier", http.StatusInternalServerError)
		return
	}

	codeChallenge := generateCodeChallenge(codeVerifier)

	// Store the verifier in the state (not recommended for production).
	state := base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf(`{"codeVerifier":"%s"}`, codeVerifier)))

	// Build the authorization URL.
	authParams := url.Values{
		"client_id":             {kickClientID},
		"redirect_uri":          {redirectURL},
		"response_type":         {"code"},
		"scope":                 {strings.Join(kickScopes, " ")},
		"state":                 {state},
		"code_challenge":        {codeChallenge},
		"code_challenge_method": {"S256"},
	}

	authURL := fmt.Sprintf("%s?%s", ENDPOINT.AuthURL, authParams.Encode())
	http.Redirect(w, r, authURL, http.StatusFound)
}

// Step 2: Handle the redirect from the auth page.
func oauthKickCallbackHandler(w http.ResponseWriter, r *http.Request) {
	code := r.URL.Query().Get("code")
	state := r.URL.Query().Get("state")

	if code == "" {
		http.Error(w, "Missing authorization code", http.StatusBadRequest)
		return
	}

	// Decode the state to get the code verifier.
	stateBytes, err := base64.StdEncoding.DecodeString(state)
	if err != nil {
		http.Error(w, "Invalid state", http.StatusBadRequest)
		return
	}

	var stateData struct {
		CodeVerifier string `json:"codeVerifier"`
	}
	if err := json.Unmarshal(stateBytes, &stateData); err != nil {
		http.Error(w, "Invalid state data", http.StatusBadRequest)
		return
	}

	// Prepare the token request.
	tokenParams := url.Values{
		"grant_type":    {"authorization_code"},
		"client_id":     {kickClientID},
		"client_secret": {kickClientSecret},
		"code":          {code},
		"redirect_uri":  {redirectURL},
		"code_verifier": {stateData.CodeVerifier},
	}

	// Send the token request.
	resp, err := http.PostForm(ENDPOINT.TokenURL, tokenParams)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to request token: %s", err.Error()), http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	// Read and return the token response.
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		http.Error(w, "Failed to read token response", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_, err = w.Write(body)
	if err != nil {
		http.Error(w, "Failed to write token response", http.StatusInternalServerError)
	}
}

func main() {
	http.HandleFunc("/oauth/kick/", oauthKickHandler)
	http.HandleFunc("/oauth/kick/callback", oauthKickCallbackHandler)

	fmt.Println("Server running on http://localhost:3000")
	err := http.ListenAndServe(":3000", nil)
	if err != nil {
		fmt.Printf("Server failed to start: %s\n", err.Error())
		return
	}
}

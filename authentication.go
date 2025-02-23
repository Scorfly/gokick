package gokick

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"strings"
)

type TokenResponse struct {
	AccessToken  string `json:"access_token"`
	TokenType    string `json:"token_type"`
	ExpiresIn    int    `json:"expires_in"`
	Scope        string `json:"scope"`
	RefreshToken string `json:"refresh_token"`
}

func (c *Client) GetAuthorizeEndpoint(redirectURI, state, codeChallenge string, scope []Scope) (string, error) {
	scopes := make([]string, len(scope))
	for i, s := range scope {
		scopes[i] = url.QueryEscape(s.String())
	}

	return fmt.Sprintf(
		"%s/oauth/authorize?client_id=%s&response_type=code&redirect_uri=%s&state=%s&scope=%s&code_challenge=%s&code_challenge_method=S256",
		c.options.AuthBaseURL,
		c.options.ClientID,
		url.QueryEscape(redirectURI),
		state,
		strings.Join(scopes, "+"),
		codeChallenge,
	), nil
}

func (c *Client) GetToken(ctx context.Context, redirectURI, code, codeVerifier string) (TokenResponse, error) {
	if c.options.ClientID == "" {
		return TokenResponse{}, fmt.Errorf("client ID must be set on Client to refresh token")
	}

	if c.options.ClientSecret == "" {
		return TokenResponse{}, fmt.Errorf("client secret must be set on Client to refresh token")
	}

	formData := url.Values{}
	formData.Set("grant_type", "authorization_code")
	formData.Set("client_id", c.options.ClientID)
	formData.Set("client_secret", c.options.ClientSecret)
	formData.Set("code", code)
	formData.Set("redirect_uri", redirectURI)
	formData.Set("code_verifier", codeVerifier)

	response, err := makeAuthRequest[TokenResponse](
		ctx,
		c,
		http.MethodPost,
		"/oauth/token",
		http.StatusOK,
		strings.NewReader(formData.Encode()),
	)
	if err != nil {
		return TokenResponse{}, err
	}

	return response, nil
}

func (c *Client) RefreshToken(ctx context.Context, refreshToken string) (TokenResponse, error) {
	if c.options.ClientID == "" {
		return TokenResponse{}, fmt.Errorf("client ID must be set on Client to refresh token")
	}

	if c.options.ClientSecret == "" {
		return TokenResponse{}, fmt.Errorf("client secret must be set on Client to refresh token")
	}

	if refreshToken == "" {
		return TokenResponse{}, fmt.Errorf("refresh token must be defined")
	}

	formData := url.Values{}
	formData.Set("grant_type", "refresh_token")
	formData.Set("client_id", c.options.ClientID)
	formData.Set("client_secret", c.options.ClientSecret)
	formData.Set("refresh_token", refreshToken)

	response, err := makeAuthRequest[TokenResponse](
		ctx,
		c,
		http.MethodPost,
		"/oauth/token",
		http.StatusOK,
		strings.NewReader(formData.Encode()),
	)
	if err != nil {
		return TokenResponse{}, err
	}

	return response, nil
}

// TODO: Currently having issue : https://github.com/KickEngineering/KickDevDocs/issues/58
func (c *Client) RevokeToken(ctx context.Context, tokenType TokenType, token string) error {
	formData := url.Values{}
	formData.Set("token_hint_type", tokenType.String())
	formData.Set("token", token)

	_, err := makeAuthRequest[TokenResponse](
		ctx,
		c,
		http.MethodPost,
		"/oauth/revoke",
		http.StatusOK,
		strings.NewReader(formData.Encode()),
	)
	if err != nil {
		return err
	}

	return nil
}

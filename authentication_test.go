package gokick_test

import (
	"context"
	"fmt"
	"net/http"
	"testing"

	"github.com/scorfly/gokick"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetAuthorizeEndpointSuccess(t *testing.T) {
	kickClient, err := gokick.NewClient(&gokick.ClientOptions{
		ClientID: "client-id",
	})
	require.NoError(t, err)

	response, err := kickClient.GetAuthorizeEndpoint(
		"http://localhost:3000/oauth/kick/callback",
		"custom-state",
		"custom-code-challenge",
		[]gokick.Scope{gokick.ScopeUserRead, gokick.ScopeChannelRead},
	)
	require.NoError(t, err)
	assert.Equal(
		t,
		`https://id.kick.com/oauth/authorize?client_id=client-id&response_type=code&redirect_uri=`+
			`http%3A%2F%2Flocalhost%3A3000%2Foauth%2Fkick%2Fcallback&state=custom-state&scope=user%3Aread+channel%3Aread`+
			`&code_challenge=custom-code-challenge&code_challenge_method=S256`,
		response,
	)
}

func TestRefreshTokenError(t *testing.T) {
	t.Run("client ID missing", func(t *testing.T) {
		kickClient, err := gokick.NewClient(&gokick.ClientOptions{
			ClientSecret: "client-secret",
		})
		require.NoError(t, err)

		var ctx context.Context
		_, err = kickClient.RefreshToken(ctx, "refresh-token")
		require.EqualError(t, err, "client ID must be set on Client to refresh token")
	})

	t.Run("client secret missing", func(t *testing.T) {
		kickClient, err := gokick.NewClient(&gokick.ClientOptions{
			ClientID: "client-id",
		})
		require.NoError(t, err)

		var ctx context.Context
		_, err = kickClient.RefreshToken(ctx, "refresh-token")
		require.EqualError(t, err, "client secret must be set on Client to refresh token")
	})

	t.Run("refresh token empty", func(t *testing.T) {
		kickClient, err := gokick.NewClient(&gokick.ClientOptions{
			ClientID:     "client-id",
			ClientSecret: "client-secret",
		})
		require.NoError(t, err)

		var ctx context.Context
		_, err = kickClient.RefreshToken(ctx, "")
		require.EqualError(t, err, "refresh token must be defined")
	})

	t.Run("on new request", func(t *testing.T) {
		kickClient, err := gokick.NewClient(&gokick.ClientOptions{
			ClientID:     "client-id",
			ClientSecret: "client-secret",
		})
		require.NoError(t, err)

		var ctx context.Context
		_, err = kickClient.RefreshToken(ctx, "refresh-token")
		require.EqualError(t, err, "failed to create request: net/http: nil Context")
	})

	t.Run("timeout", func(t *testing.T) {
		kickClient := setupTimeoutMockAuthClient(t)

		_, err := kickClient.RefreshToken(context.Background(), "refresh-token")
		require.EqualError(t, err, `failed to make request: Post "https://id.kick.com/oauth/token": context deadline exceeded `+
			`(Client.Timeout exceeded while awaiting headers)`)
	})

	t.Run("unmarshal error response", func(t *testing.T) {
		kickClient := setupMockAuthClient(t, func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprint(w, `117`)
		})

		_, err := kickClient.RefreshToken(context.Background(), "refresh-token")

		assert.EqualError(t, err, `failed to unmarshal error response (KICK status code: 500 and body "117"): json: cannot unmarshal `+
			`number into Go value of type gokick.errorResponse`)
	})

	t.Run("unmarshal token response", func(t *testing.T) {
		kickClient := setupMockAuthClient(t, func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			fmt.Fprint(w, "117")
		})

		_, err := kickClient.RefreshToken(context.Background(), "refresh-token")

		assert.EqualError(t, err, `failed to unmarshal response body (KICK status code 200 and body "117"): json: cannot unmarshal `+
			`number into Go value of type gokick.TokenResponse`)
	})

	t.Run("reader failure", func(t *testing.T) {
		kickClient := setupMockAuthClient(t, func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Length", "10")
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprint(w, "")
		})

		_, err := kickClient.RefreshToken(context.Background(), "refresh-token")

		assert.EqualError(t, err, `failed to read response body (KICK status code 500): unexpected EOF`)
	})

	t.Run("with internal server error", func(t *testing.T) {
		kickClient := setupMockAuthClient(t, func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprint(w, `{"message":"internal server error", "data":null}`)
		})

		_, err := kickClient.RefreshToken(context.Background(), "refresh-token")

		var kickError gokick.Error
		require.ErrorAs(t, err, &kickError)
		assert.Equal(t, http.StatusInternalServerError, kickError.Code())
		assert.Equal(t, "internal server error", kickError.Message())
	})
}

func TestRefreshTokenSuccess(t *testing.T) {
	kickClient := setupMockAuthClient(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, `{
			"access_token":"new-access-token",
			"expires_in":7200,
			"refresh_token":"new-refresh-token",
			"scope":"user:read chat:write channel:read channel:write streamkey:read events:subscribe",
			"token_type":"Bearer"
		}`)
	})

	response, err := kickClient.RefreshToken(context.Background(), "refresh-token")
	require.NoError(t, err)
	assert.Equal(t, "new-access-token", response.AccessToken)
	assert.Equal(t, 7200, response.ExpiresIn)
	assert.Equal(t, "new-refresh-token", response.RefreshToken)
	assert.Equal(t, "user:read chat:write channel:read channel:write streamkey:read events:subscribe", response.Scope)
	assert.Equal(t, "Bearer", response.TokenType)
}

func TestRevokeTokenError(t *testing.T) {
	t.Run("on new request", func(t *testing.T) {
		kickClient, err := gokick.NewClient(&gokick.ClientOptions{
			ClientID:     "client-id",
			ClientSecret: "client-secret",
		})
		require.NoError(t, err)

		var ctx context.Context
		err = kickClient.RevokeToken(ctx, gokick.TokenTypeAccess, "access-token")
		require.EqualError(t, err, "failed to create request: net/http: nil Context")
	})

	t.Run("timeout", func(t *testing.T) {
		kickClient := setupTimeoutMockAuthClient(t)

		err := kickClient.RevokeToken(context.Background(), gokick.TokenTypeAccess, "access-token")
		require.EqualError(t, err, `failed to make request: Post "https://id.kick.com/oauth/revoke": context deadline exceeded `+
			`(Client.Timeout exceeded while awaiting headers)`)
	})

	t.Run("unmarshal error response", func(t *testing.T) {
		kickClient := setupMockAuthClient(t, func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprint(w, `117`)
		})

		err := kickClient.RevokeToken(context.Background(), gokick.TokenTypeAccess, "access-token")
		assert.EqualError(t, err, `failed to unmarshal error response (KICK status code: 500 and body "117"): json: cannot unmarshal `+
			`number into Go value of type gokick.errorResponse`)
	})

	t.Run("unmarshal token response", func(t *testing.T) {
		kickClient := setupMockAuthClient(t, func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			fmt.Fprint(w, "117")
		})

		err := kickClient.RevokeToken(context.Background(), gokick.TokenTypeAccess, "access-token")
		assert.EqualError(t, err, `failed to unmarshal response body (KICK status code 200 and body "117"): json: cannot unmarshal `+
			`number into Go value of type gokick.TokenResponse`)
	})

	t.Run("reader failure", func(t *testing.T) {
		kickClient := setupMockAuthClient(t, func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Length", "10")
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprint(w, "")
		})

		err := kickClient.RevokeToken(context.Background(), gokick.TokenTypeAccess, "access-token")
		assert.EqualError(t, err, `failed to read response body (KICK status code 500): unexpected EOF`)
	})

	t.Run("with internal server error", func(t *testing.T) {
		kickClient := setupMockAuthClient(t, func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprint(w, `{"message":"internal server error", "data":null}`)
		})

		err := kickClient.RevokeToken(context.Background(), gokick.TokenTypeAccess, "access-token")

		var kickError gokick.Error
		require.ErrorAs(t, err, &kickError)
		assert.Equal(t, http.StatusInternalServerError, kickError.Code())
		assert.Equal(t, "internal server error", kickError.Message())
	})
}

func TestRevokeTokenSuccess(t *testing.T) {
	testCases := map[string]struct {
		tokenType     gokick.TokenType
		token         string
		expectedPath  string
		expectedToken string
		tokenHintType string
	}{
		"access token": {
			tokenType:     gokick.TokenTypeAccess,
			token:         "access-token",
			expectedPath:  "/oauth/revoke",
			expectedToken: "access-token",
			tokenHintType: "access_token",
		},
		"refresh token": {
			tokenType:     gokick.TokenTypeRefresh,
			token:         "refresh-token",
			expectedPath:  "/oauth/revoke",
			expectedToken: "refresh-token",
			tokenHintType: "refresh_token",
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			kickClient := setupMockAuthClient(t, func(w http.ResponseWriter, r *http.Request) {
				assert.Equal(t, tc.expectedPath, r.URL.Path)
				assert.Equal(t, tc.expectedToken, r.FormValue("token"))
				assert.Equal(t, tc.tokenHintType, r.FormValue("token_hint_type"))

				w.WriteHeader(http.StatusOK)
				fmt.Fprint(w, `{}`)
			})

			err := kickClient.RevokeToken(context.Background(), tc.tokenType, tc.token)
			require.NoError(t, err)
		})
	}
}

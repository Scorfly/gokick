package gokick_test

import (
	"bytes"
	"context"
	"io"
	"net/http"
	"testing"

	"github.com/scorfly/gokick"
	"github.com/stretchr/testify/require"
)

func TestNewClientSuccess(t *testing.T) {
	client, err := gokick.NewClient(&gokick.ClientOptions{UserAccessToken: "access-token"})
	require.IsType(t, &gokick.Client{}, client)
	require.NoError(t, err)

	t.Run("SetAppAccessToken", func(t *testing.T) {
		client.SetAppAccessToken("app-token")
	})

	t.Run("SetUserAccessToken", func(t *testing.T) {
		client.SetUserAccessToken("user-access-token")
	})

	t.Run("SetUserRefreshToken", func(t *testing.T) {
		client.SetUserRefreshToken("user-refresh-token")
	})

	t.Run("OnUserAccessTokenRefreshed", func(t *testing.T) {
		client.OnUserAccessTokenRefreshed(func(a string, b string) {})
	})
}

type mockRoundTripper struct {
	code int
}

func (c *mockRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	return &http.Response{
		StatusCode: c.code,
		Body:       io.NopCloser(bytes.NewBufferString("OK")),
	}, nil
}

func TestClientRefreshTokenError(t *testing.T) {
	t.Run("on new request", func(t *testing.T) {
		client := &http.Client{
			Transport: &mockRoundTripper{
				code: http.StatusUnauthorized,
			},
		}

		kickClient, err := gokick.NewClient(&gokick.ClientOptions{
			UserAccessToken:  "access-token",
			ClientID:         "client-id",
			ClientSecret:     "client-secret",
			UserRefreshToken: "user-refresh-token",
			HTTPClient:       client,
		})
		require.NoError(t, err)

		_, err = kickClient.GetCategory(context.Background(), 117)
		require.EqualError(t, err, "failed to make request: failed to refresh token: failed to unmarshal error response "+
			"(KICK status code: 401 and body \"OK\"): invalid character 'O' looking for beginning of value")
	})
}

type mockRoundTripperRefreshTokenOK struct {
	code int
}

func (c *mockRoundTripperRefreshTokenOK) RoundTrip(req *http.Request) (*http.Response, error) {
	code := c.code
	c.code = http.StatusOK

	return &http.Response{
		StatusCode: code,
		Body: io.NopCloser(bytes.NewBufferString(`{
			"access_token":"access-token",
			"expires_in":7200,
			"token_type":"Bearer"
		}`)),
	}, nil
}

func TestClientRefreshTokenSuccess(t *testing.T) {
	t.Run("on new request", func(t *testing.T) {
		client := &http.Client{
			Transport: &mockRoundTripperRefreshTokenOK{
				code: http.StatusUnauthorized,
			},
		}

		kickClient, err := gokick.NewClient(&gokick.ClientOptions{
			UserAccessToken:  "access-token",
			ClientID:         "client-id",
			ClientSecret:     "client-secret",
			UserRefreshToken: "user-refresh-token",
			AppAccessToken:   "app-access-token",
			HTTPClient:       client,
		})
		require.NoError(t, err)

		kickClient.OnUserAccessTokenRefreshed(func(a string, b string) {})

		_, err = kickClient.GetCategory(context.Background(), 117)
		require.NoError(t, err)
	})
}

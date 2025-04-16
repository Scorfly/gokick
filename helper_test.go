package gokick_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/scorfly/gokick"
	"github.com/stretchr/testify/require"
)

func setupMockClient(t *testing.T, mockHandler http.HandlerFunc) *gokick.Client {
	t.Helper()

	server := httptest.NewServer(mockHandler)
	kickClient, err := gokick.NewClient(&gokick.ClientOptions{
		UserAccessToken: "access-token",
		APIBaseURL:      fmt.Sprintf("http://%s", server.Listener.Addr()),
	})
	require.NoError(t, err)

	t.Cleanup(func() { server.Close() })

	return kickClient
}

func setupMockAuthClient(t *testing.T, mockHandler http.HandlerFunc) *gokick.Client {
	t.Helper()

	server := httptest.NewServer(mockHandler)
	kickClient, err := gokick.NewClient(&gokick.ClientOptions{
		ClientID:     "client-id",
		ClientSecret: "client-secret",
		AuthBaseURL:  fmt.Sprintf("http://%s", server.Listener.Addr()),
	})
	require.NoError(t, err)

	t.Cleanup(func() { server.Close() })

	return kickClient
}

func setupTimeoutMockClient(t *testing.T) *gokick.Client {
	t.Helper()

	kickClient, err := gokick.NewClient(&gokick.ClientOptions{
		UserAccessToken: "access-token",
		HTTPClient:      &http.Client{Timeout: 1 * time.Nanosecond},
	})
	require.NoError(t, err)

	return kickClient
}

func setupTimeoutMockAuthClient(t *testing.T) *gokick.Client {
	t.Helper()

	kickClient, err := gokick.NewClient(&gokick.ClientOptions{
		ClientID:     "client-id",
		ClientSecret: "client-secret",
		HTTPClient:   &http.Client{Timeout: 1 * time.Nanosecond},
	})
	require.NoError(t, err)

	return kickClient
}

func stringPtr(s string) *string {
	return &s
}

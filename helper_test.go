package gokick_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/scorfly/gokick"
	"github.com/stretchr/testify/require"
)

func setupMockClient(t *testing.T, mockHandler http.HandlerFunc) *gokick.Client {
	t.Helper()

	server := httptest.NewServer(mockHandler)
	kickClient, err := gokick.NewClient(&http.Client{}, fmt.Sprintf("http://%s", server.Listener.Addr()), "access-token")
	require.NoError(t, err)

	t.Cleanup(func() { server.Close() })

	return kickClient
}

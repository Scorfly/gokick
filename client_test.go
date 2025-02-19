package gokick_test

import (
	"net/http"
	"testing"

	"github.com/scorfly/gokick"
	"github.com/stretchr/testify/require"
)

func TestNewClientSuccess(t *testing.T) {
	client, err := gokick.NewClient(&http.Client{}, "api-url")
	require.IsType(t, &gokick.Client{}, client)
	require.NoError(t, err)
}

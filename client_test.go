package gokick_test

import (
	"testing"

	"github.com/scorfly/gokick"
	"github.com/stretchr/testify/require"
)

func TestNewClientSuccess(t *testing.T) {
	client, err := gokick.NewClient(&gokick.ClientOptions{UserAccessToken: "access-token"})
	require.IsType(t, &gokick.Client{}, client)
	require.NoError(t, err)
}

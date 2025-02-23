package gokick_test

import (
	"fmt"
	"testing"

	"github.com/scorfly/gokick"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewTokenTypeError(t *testing.T) {
	testCases := map[string]string{
		"empty":         "",
		"not supported": "not supported",
	}

	for name, value := range testCases {
		t.Run(name, func(t *testing.T) {
			_, err := gokick.NewTokenType(value)
			assert.EqualError(t, err, fmt.Sprintf("unknown token type: %s", value))
		})
	}
}

func TestNewTokenTypeSuccess(t *testing.T) {
	testCases := map[string]gokick.TokenType{
		"access_token":  gokick.TokenTypeAccess,
		"refresh_token": gokick.TokenTypeRefresh,
	}

	for name, value := range testCases {
		t.Run(name, func(t *testing.T) {
			TokenType, err := gokick.NewTokenType(value.String())
			require.NoError(t, err)
			assert.Equal(t, TokenType, value)
		})
	}
}

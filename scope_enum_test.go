package gokick_test

import (
	"fmt"
	"testing"

	"github.com/scorfly/gokick"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewScopeError(t *testing.T) {
	testCases := map[string]string{
		"empty":         "",
		"not supported": "not supported",
	}

	for name, value := range testCases {
		t.Run(name, func(t *testing.T) {
			_, err := gokick.NewScope(value)
			assert.EqualError(t, err, fmt.Sprintf("unknown scope: %s", value))
		})
	}
}

func TestNewScopeSuccess(t *testing.T) {
	testCases := map[string]gokick.Scope{
		"user:read":        gokick.ScopeUserRead,
		"channel:read":     gokick.ScopeChannelRead,
		"channel:write":    gokick.ScopeChannelWrite,
		"chat:write":       gokick.ScopeChatWrite,
		"streamkey:read":   gokick.ScopeStremkeyRead,
		"events:subscribe": gokick.ScopeEventSubscribe,
		"moderation:ban":   gokick.ScopeModerationBan,
		"kicks:read":       gokick.ScopeKicksRead,
	}

	for name, value := range testCases {
		t.Run(name, func(t *testing.T) {
			Scope, err := gokick.NewScope(value.String())
			require.NoError(t, err)
			assert.Equal(t, Scope, value)
		})
	}
}

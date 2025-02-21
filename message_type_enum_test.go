package gokick_test

import (
	"fmt"
	"testing"

	"github.com/scorfly/gokick"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewMessageTypeError(t *testing.T) {
	testCases := map[string]string{
		"empty":         "",
		"not supported": "not supported",
	}

	for name, value := range testCases {
		t.Run(name, func(t *testing.T) {
			_, err := gokick.NewMessageType(value)
			assert.EqualError(t, err, fmt.Sprintf("unknown message type: %s", value))
		})
	}
}

func TestNewMessageTypeSuccess(t *testing.T) {
	testCases := map[string]gokick.MessageType{
		"bot":  gokick.MessageTypeBot,
		"user": gokick.MessageTypeUser,
	}

	for name, value := range testCases {
		t.Run(name, func(t *testing.T) {
			messageType, err := gokick.NewMessageType(value.String())
			require.NoError(t, err)
			assert.Equal(t, messageType, value)
		})
	}
}

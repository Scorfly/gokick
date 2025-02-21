package gokick_test

import (
	"fmt"
	"testing"

	"github.com/scorfly/gokick"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewSubscriptionMethodError(t *testing.T) {
	testCases := map[string]string{
		"empty":         "",
		"not supported": "not supported",
	}

	for name, value := range testCases {
		t.Run(name, func(t *testing.T) {
			_, err := gokick.NewSubscriptionMethod(value)
			assert.EqualError(t, err, fmt.Sprintf("unknown method: %s", value))
		})
	}
}

func TestNewSubscriptionMethodSuccess(t *testing.T) {
	testCases := map[string]gokick.SubscriptionMethod{
		"webhook": gokick.SubscriptionMethodWebhook,
	}

	for name, value := range testCases {
		t.Run(name, func(t *testing.T) {
			SubscriptionMethod, err := gokick.NewSubscriptionMethod(value.String())
			require.NoError(t, err)
			assert.Equal(t, SubscriptionMethod, value)
		})
	}
}

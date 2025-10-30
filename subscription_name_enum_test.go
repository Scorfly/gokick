package gokick_test

import (
	"fmt"
	"testing"

	"github.com/scorfly/gokick"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewSubscriptionNameError(t *testing.T) {
	testCases := map[string]string{
		"empty":         "",
		"not supported": "not supported",
	}

	for name, value := range testCases {
		t.Run(name, func(t *testing.T) {
			_, err := gokick.NewSubscriptionName(value)
			assert.EqualError(t, err, fmt.Sprintf("unknown name: %s", value))
		})
	}
}

func TestNewSubscriptionNameSuccess(t *testing.T) {
	testCases := map[string]gokick.SubscriptionName{
		"chat.message.sent":            gokick.SubscriptionNameChatMessage,
		"channel.followed":             gokick.SubscriptionNameChannelFollow,
		"channel.subscription.renewal": gokick.SubscriptionNameChannelSubscriptionRenewal,
		"channel.subscription.gifts":   gokick.SubscriptionNameChannelSubscriptionGifts,
		"channel.subscription.new":     gokick.SubscriptionNameChannelSubscriptionCreated,
		"livestream.status.updated":    gokick.SubscriptionNameLivestreamStatusUpdated,
		"livestream.metadata.updated":  gokick.SubscriptionNameLivestreamMetadataUpdated,
		"moderation.banned":            gokick.SubscriptionNameModerationBanned,
		"kicks.gifted":                 gokick.SubscriptionNameKicksGifted,
	}

	for name, value := range testCases {
		t.Run(name, func(t *testing.T) {
			subscriptionName, err := gokick.NewSubscriptionName(value.String())
			require.NoError(t, err)
			assert.Equal(t, subscriptionName, value)
		})
	}
}

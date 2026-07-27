package gokick_test

import (
	"fmt"
	"testing"

	"github.com/scorfly/gokick"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewChannelRewardRedemptionStatusError(t *testing.T) {
	testCases := map[string]string{
		"empty":         "",
		"not supported": "not supported",
	}

	for name, value := range testCases {
		t.Run(name, func(t *testing.T) {
			_, err := gokick.NewChannelRewardRedemptionStatus(value)
			assert.EqualError(t, err, fmt.Sprintf("unknown channel reward redemption status: %s", value))
		})
	}
}

func TestNewChannelRewardRedemptionStatusSuccess(t *testing.T) {
	testCases := map[string]gokick.ChannelRewardRedemptionStatus{
		"pending":  gokick.ChannelRewardRedemptionStatusPending,
		"accepted": gokick.ChannelRewardRedemptionStatusAccepted,
		"rejected": gokick.ChannelRewardRedemptionStatusRejected,
	}

	for name, value := range testCases {
		t.Run(name, func(t *testing.T) {
			status, err := gokick.NewChannelRewardRedemptionStatus(value.String())
			require.NoError(t, err)
			assert.Equal(t, value, status)
		})
	}
}

func TestNewFailedRedemptionReasonError(t *testing.T) {
	testCases := map[string]string{
		"empty":         "",
		"not supported": "not supported",
	}

	for name, value := range testCases {
		t.Run(name, func(t *testing.T) {
			_, err := gokick.NewFailedRedemptionReason(value)
			assert.EqualError(t, err, fmt.Sprintf("unknown failed redemption reason: %s", value))
		})
	}
}

func TestNewFailedRedemptionReasonSuccess(t *testing.T) {
	testCases := map[string]gokick.FailedRedemptionReason{
		"UNKNOWN":     gokick.FailedRedemptionReasonUnknown,
		"NOT_PENDING": gokick.FailedRedemptionReasonNotPending,
		"NOT_FOUND":   gokick.FailedRedemptionReasonNotFound,
		"NOT_OWNED":   gokick.FailedRedemptionReasonNotOwned,
	}

	for name, value := range testCases {
		t.Run(name, func(t *testing.T) {
			reason, err := gokick.NewFailedRedemptionReason(value.String())
			require.NoError(t, err)
			assert.Equal(t, value, reason)
		})
	}
}

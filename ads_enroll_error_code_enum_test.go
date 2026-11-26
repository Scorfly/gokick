package gokick_test

import (
	"fmt"
	"testing"

	"github.com/scorfly/gokick"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewAdsEnrollErrorCodeError(t *testing.T) {
	testCases := map[string]string{
		"empty":         "",
		"not supported": "not supported",
	}

	for name, value := range testCases {
		t.Run(name, func(t *testing.T) {
			_, err := gokick.NewAdsEnrollErrorCode(value)
			assert.EqualError(t, err, fmt.Sprintf("unknown ads enroll error code: %s", value))
		})
	}
}

func TestNewAdsEnrollErrorCodeSuccess(t *testing.T) {
	testCases := map[string]gokick.AdsEnrollErrorCode{
		"CHANNEL_NOT_FOUND":      gokick.AdsEnrollErrorCodeChannelNotFound,
		"CHANNEL_IS_LIVE":        gokick.AdsEnrollErrorCodeChannelIsLive,
		"CHANNEL_NEVER_STREAMED": gokick.AdsEnrollErrorCodeChannelNeverStreamed,
	}

	for name, value := range testCases {
		t.Run(name, func(t *testing.T) {
			code, err := gokick.NewAdsEnrollErrorCode(value.String())
			require.NoError(t, err)
			assert.Equal(t, value, code)
		})
	}
}

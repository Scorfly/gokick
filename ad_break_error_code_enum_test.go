package gokick_test

import (
	"fmt"
	"testing"

	"github.com/scorfly/gokick"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewAdBreakErrorCodeError(t *testing.T) {
	testCases := map[string]string{
		"empty":         "",
		"not supported": "not supported",
	}

	for name, value := range testCases {
		t.Run(name, func(t *testing.T) {
			_, err := gokick.NewAdBreakErrorCode(value)
			assert.EqualError(t, err, fmt.Sprintf("unknown ad break error code: %s", value))
		})
	}
}

func TestNewAdBreakErrorCodeSuccess(t *testing.T) {
	testCases := map[string]gokick.AdBreakErrorCode{
		"CHANNEL_IS_NOT_LIVE":   gokick.AdBreakErrorCodeChannelIsNotLive,
		"CHANNEL_NOT_OPTED_IN":  gokick.AdBreakErrorCodeChannelNotOptedIn,
		"ADS_BLOCKED":           gokick.AdBreakErrorCodeAdsBlocked,
		"ADS_RESTRICTED":        gokick.AdBreakErrorCodeAdsRestricted,
		"AD_BREAK_ACTIVE":       gokick.AdBreakErrorCodeAdBreakActive,
		"MAX_AD_BREAKS_REACHED": gokick.AdBreakErrorCodeMaxAdBreaksReached,
		"AUTO_ADS_ACTIVE":       gokick.AdBreakErrorCodeAutoAdsActive,
	}

	for name, value := range testCases {
		t.Run(name, func(t *testing.T) {
			code, err := gokick.NewAdBreakErrorCode(value.String())
			require.NoError(t, err)
			assert.Equal(t, value, code)
		})
	}
}

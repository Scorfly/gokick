package gokick

import "fmt"

// AdBreakErrorCode is the machine-readable code Kick returns in data when creating an ad break fails.
type AdBreakErrorCode int

const (
	AdBreakErrorCodeChannelIsNotLive   AdBreakErrorCode = iota // CHANNEL_IS_NOT_LIVE
	AdBreakErrorCodeChannelNotOptedIn                          // CHANNEL_NOT_OPTED_IN
	AdBreakErrorCodeAdsBlocked                                 // ADS_BLOCKED
	AdBreakErrorCodeAdsRestricted                              // ADS_RESTRICTED
	AdBreakErrorCodeAdBreakActive                              // AD_BREAK_ACTIVE
	AdBreakErrorCodeMaxAdBreaksReached                         // MAX_AD_BREAKS_REACHED
	AdBreakErrorCodeAutoAdsActive                              // AUTO_ADS_ACTIVE
)

// NewAdBreakErrorCode parses a Kick ad-break error code.
func NewAdBreakErrorCode(code string) (AdBreakErrorCode, error) {
	switch code {
	case "CHANNEL_IS_NOT_LIVE":
		return AdBreakErrorCodeChannelIsNotLive, nil
	case "CHANNEL_NOT_OPTED_IN":
		return AdBreakErrorCodeChannelNotOptedIn, nil
	case "ADS_BLOCKED":
		return AdBreakErrorCodeAdsBlocked, nil
	case "ADS_RESTRICTED":
		return AdBreakErrorCodeAdsRestricted, nil
	case "AD_BREAK_ACTIVE":
		return AdBreakErrorCodeAdBreakActive, nil
	case "MAX_AD_BREAKS_REACHED":
		return AdBreakErrorCodeMaxAdBreaksReached, nil
	case "AUTO_ADS_ACTIVE":
		return AdBreakErrorCodeAutoAdsActive, nil
	default:
		return 0, fmt.Errorf("unknown ad break error code: %s", code)
	}
}

func (c AdBreakErrorCode) String() string {
	switch c {
	case AdBreakErrorCodeChannelIsNotLive:
		return "CHANNEL_IS_NOT_LIVE"
	case AdBreakErrorCodeChannelNotOptedIn:
		return "CHANNEL_NOT_OPTED_IN"
	case AdBreakErrorCodeAdsBlocked:
		return "ADS_BLOCKED"
	case AdBreakErrorCodeAdsRestricted:
		return "ADS_RESTRICTED"
	case AdBreakErrorCodeAdBreakActive:
		return "AD_BREAK_ACTIVE"
	case AdBreakErrorCodeMaxAdBreaksReached:
		return "MAX_AD_BREAKS_REACHED"
	case AdBreakErrorCodeAutoAdsActive:
		return "AUTO_ADS_ACTIVE"
	default:
		return "unknown"
	}
}

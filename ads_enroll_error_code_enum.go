package gokick

import "fmt"

// AdsEnrollErrorCode is the machine-readable code Kick returns in data when enrolling in ads fails.
type AdsEnrollErrorCode int

const (
	AdsEnrollErrorCodeChannelNotFound      AdsEnrollErrorCode = iota // CHANNEL_NOT_FOUND
	AdsEnrollErrorCodeChannelIsLive                                  // CHANNEL_IS_LIVE
	AdsEnrollErrorCodeChannelNeverStreamed                           // CHANNEL_NEVER_STREAMED
)

// NewAdsEnrollErrorCode parses a Kick ads-enrollment error code.
func NewAdsEnrollErrorCode(code string) (AdsEnrollErrorCode, error) {
	switch code {
	case "CHANNEL_NOT_FOUND":
		return AdsEnrollErrorCodeChannelNotFound, nil
	case "CHANNEL_IS_LIVE":
		return AdsEnrollErrorCodeChannelIsLive, nil
	case "CHANNEL_NEVER_STREAMED":
		return AdsEnrollErrorCodeChannelNeverStreamed, nil
	default:
		return 0, fmt.Errorf("unknown ads enroll error code: %s", code)
	}
}

func (c AdsEnrollErrorCode) String() string {
	switch c {
	case AdsEnrollErrorCodeChannelNotFound:
		return "CHANNEL_NOT_FOUND"
	case AdsEnrollErrorCodeChannelIsLive:
		return "CHANNEL_IS_LIVE"
	case AdsEnrollErrorCodeChannelNeverStreamed:
		return "CHANNEL_NEVER_STREAMED"
	default:
		return "unknown"
	}
}

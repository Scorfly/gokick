package gokick

import "fmt"

type ChannelRewardRedemptionStatus int

const (
	ChannelRewardRedemptionStatusPending  ChannelRewardRedemptionStatus = iota // pending
	ChannelRewardRedemptionStatusAccepted                                      // accepted
	ChannelRewardRedemptionStatusRejected                                      // rejected
)

func NewChannelRewardRedemptionStatus(status string) (ChannelRewardRedemptionStatus, error) {
	switch status {
	case "pending":
		return ChannelRewardRedemptionStatusPending, nil
	case "accepted":
		return ChannelRewardRedemptionStatusAccepted, nil
	case "rejected":
		return ChannelRewardRedemptionStatusRejected, nil
	default:
		return 0, fmt.Errorf("unknown channel reward redemption status: %s", status)
	}
}

func (s ChannelRewardRedemptionStatus) String() string {
	switch s {
	case ChannelRewardRedemptionStatusPending:
		return "pending"
	case ChannelRewardRedemptionStatusAccepted:
		return "accepted"
	case ChannelRewardRedemptionStatusRejected:
		return "rejected"
	default:
		return "unknown"
	}
}

type FailedRedemptionReason int

const (
	FailedRedemptionReasonUnknown    FailedRedemptionReason = iota // UNKNOWN
	FailedRedemptionReasonNotPending                               // NOT_PENDING
	FailedRedemptionReasonNotFound                                 // NOT_FOUND
	FailedRedemptionReasonNotOwned                                 // NOT_OWNED
)

func NewFailedRedemptionReason(reason string) (FailedRedemptionReason, error) {
	switch reason {
	case "UNKNOWN":
		return FailedRedemptionReasonUnknown, nil
	case "NOT_PENDING":
		return FailedRedemptionReasonNotPending, nil
	case "NOT_FOUND":
		return FailedRedemptionReasonNotFound, nil
	case "NOT_OWNED":
		return FailedRedemptionReasonNotOwned, nil
	default:
		return 0, fmt.Errorf("unknown failed redemption reason: %s", reason)
	}
}

func (r FailedRedemptionReason) String() string {
	switch r {
	case FailedRedemptionReasonUnknown:
		return "UNKNOWN"
	case FailedRedemptionReasonNotPending:
		return "NOT_PENDING"
	case FailedRedemptionReasonNotFound:
		return "NOT_FOUND"
	case FailedRedemptionReasonNotOwned:
		return "NOT_OWNED"
	default:
		return "unknown"
	}
}

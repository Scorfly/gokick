package gokick

import "fmt"

type SubscriptionName int

const (
	SubscriptionNameChatMessage                SubscriptionName = iota // chat.message.sent
	SubscriptionNameChannelFollow                                      // channel.followed
	SubscriptionNameChannelSubscriptionRenewal                         // channel.subscription.renewal
	SubscriptionNameChannelSubscriptionGifts                           // channel.subscription.gifts
	SubscriptionNameChannelSubscriptionCreated                         // channel.subscription.new
	SubscriptionNameLivestreamStatusUpdated                            // livestream.status.updated
	SubscriptionNameLivestreamMetadataUpdated                          // livestream.metadata.updated
	SubscriptionNameModerationBanned                                   // moderation.banned
	SubscriptionNameKicksGifted                                        // kicks.gifted
)

func NewSubscriptionName(name string) (SubscriptionName, error) {
	switch name {
	case "chat.message.sent":
		return SubscriptionNameChatMessage, nil
	case "channel.followed":
		return SubscriptionNameChannelFollow, nil
	case "channel.subscription.renewal":
		return SubscriptionNameChannelSubscriptionRenewal, nil
	case "channel.subscription.gifts":
		return SubscriptionNameChannelSubscriptionGifts, nil
	case "channel.subscription.new":
		return SubscriptionNameChannelSubscriptionCreated, nil
	case "livestream.status.updated":
		return SubscriptionNameLivestreamStatusUpdated, nil
	case "livestream.metadata.updated":
		return SubscriptionNameLivestreamMetadataUpdated, nil
	case "moderation.banned":
		return SubscriptionNameModerationBanned, nil
	case "kicks.gifted":
		return SubscriptionNameKicksGifted, nil
	default:
		return 0, fmt.Errorf("unknown name: %s", name)
	}
}

func (s SubscriptionName) String() string {
	switch s {
	case SubscriptionNameChatMessage:
		return "chat.message.sent"
	case SubscriptionNameChannelFollow:
		return "channel.followed"
	case SubscriptionNameChannelSubscriptionRenewal:
		return "channel.subscription.renewal"
	case SubscriptionNameChannelSubscriptionGifts:
		return "channel.subscription.gifts"
	case SubscriptionNameChannelSubscriptionCreated:
		return "channel.subscription.new"
	case SubscriptionNameLivestreamStatusUpdated:
		return "livestream.status.updated"
	case SubscriptionNameLivestreamMetadataUpdated:
		return "livestream.metadata.updated"
	case SubscriptionNameModerationBanned:
		return "moderation.banned"
	case SubscriptionNameKicksGifted:
		return "kicks.gifted"
	default:
		return "unknown"
	}
}

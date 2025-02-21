package gokick

import "fmt"

type SubscriptionMethod int

const (
	SubscriptionMethodWebhook SubscriptionMethod = iota // webhook
)

func NewSubscriptionMethod(method string) (SubscriptionMethod, error) {
	switch method {
	case "webhook":
		return SubscriptionMethodWebhook, nil
	default:
		return 0, fmt.Errorf("unknown method: %s", method)
	}
}

func (s SubscriptionMethod) String() string {
	switch s {
	case SubscriptionMethodWebhook:
		return "webhook"
	default:
		return "unknown"
	}
}

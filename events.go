package gokick

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

type (
	EventsResponseWrapper              Response[[]EventResponse]
	EventResponseWrapper               Response[EventResponse]
	CreateSubscriptionsResponseWrapper Response[[]CreateSubscriptionResponse]
	CreateSubscriptionResponseWrapper  Response[CreateSubscriptionResponse]
)

type EventResponse struct {
	AppID             string `json:"app_id"`
	BroadcasterUserID int    `json:"broadcaster_user_id"`
	CreatedAt         string `json:"created_at"`
	Event             string `json:"event"`
	ID                string `json:"id"`
	Method            string `json:"method"`
	UpdatedAt         string `json:"updated_at"`
	Version           int    `json:"version"`
}

type CreateSubscriptionResponse struct {
	Error          string `json:"error"`
	Name           string `json:"name"`
	SubscriptionID string `json:"subscription_id"`
	Version        int    `json:"version"`
}

func (c *Client) GetSubscriptions(ctx context.Context) (EventsResponseWrapper, error) {
	response, err := makeRequest[[]EventResponse](
		ctx,
		c,
		http.MethodGet,
		"/public/v1/events/subscriptions",
		http.StatusOK,
		http.NoBody,
	)
	if err != nil {
		return EventsResponseWrapper{}, err
	}

	return EventsResponseWrapper(response), nil
}

type (
	SubscriptionMethod int
)

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

type (
	SubscriptionName int
)

const (
	SubscriptionNameChatMessage                 SubscriptionName = iota // chat.message.sent
	SubscriptionNameChannelFollow                                       // channel.followed
	SubscriptionNameChannelSubscriptionRenewal                          // channel.subscription.renewal
	SubscriptionNameChannelSubscriptionGifts                            // channel.subscription.gifts
	SubscriptionNameChannelSubscriptionfCreated                         // channel.subscription.new
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
		return SubscriptionNameChannelSubscriptionfCreated, nil
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
	case SubscriptionNameChannelSubscriptionfCreated:
		return "channel.subscription.new"
	default:
		return "unknown"
	}
}

type SubscriptionRequest struct {
	Name    SubscriptionName `json:"name"`
	Version int              `json:"version"`
}

func (c *Client) CreateSubscriptions(
	ctx context.Context,
	method SubscriptionMethod,
	subscriptions []SubscriptionRequest,
) (CreateSubscriptionsResponseWrapper, error) {
	type postBodyRequestSubscription struct {
		Name    string `json:"name"`
		Version int    `json:"version"`
	}

	type postBodyRequest struct {
		Method string                        `json:"method"`
		Events []postBodyRequestSubscription `json:"events"`
	}

	events := make([]postBodyRequestSubscription, len(subscriptions))
	for i := range subscriptions {
		events[i] = postBodyRequestSubscription{
			Name:    subscriptions[i].Name.String(),
			Version: subscriptions[i].Version,
		}
	}

	r := postBodyRequest{
		Method: method.String(),
		Events: events,
	}

	body, err := json.Marshal(r)
	if err != nil {
		return CreateSubscriptionsResponseWrapper{}, fmt.Errorf("failed to marshal body: %v", err)
	}

	response, err := makeRequest[[]CreateSubscriptionResponse](
		ctx,
		c,
		http.MethodPost,
		"/public/v1/events/subscriptions",
		http.StatusCreated,
		bytes.NewReader(body),
	)
	if err != nil {
		return CreateSubscriptionsResponseWrapper{}, err
	}

	return CreateSubscriptionsResponseWrapper(response), nil
}

package gokick

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
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

type SubscriptionRequest struct {
	Name    SubscriptionName `json:"name"`
	Version int              `json:"version"`
}

func (c *Client) CreateSubscriptions(
	ctx context.Context,
	method SubscriptionMethod,
	subscriptions []SubscriptionRequest,
	broadcasterUserID *int,
) (CreateSubscriptionsResponseWrapper, error) {
	type postBodyRequestSubscription struct {
		Name    string `json:"name"`
		Version int    `json:"version"`
	}

	type postBodyRequest struct {
		Method            string                        `json:"method"`
		Events            []postBodyRequestSubscription `json:"events"`
		BroadcasterUserID int                           `json:"broadcaster_user_id,omitempty"`
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

	if broadcasterUserID != nil {
		r.BroadcasterUserID = *broadcasterUserID
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
		http.StatusOK,
		bytes.NewReader(body),
	)
	if err != nil {
		return CreateSubscriptionsResponseWrapper{}, err
	}

	return CreateSubscriptionsResponseWrapper(response), nil
}

type SubscriptionToDeleteFilter struct {
	queryParams url.Values
}

func NewSubscriptionToDeleteFilter() SubscriptionToDeleteFilter {
	return SubscriptionToDeleteFilter{queryParams: make(url.Values)}
}

func (f SubscriptionToDeleteFilter) SetIDs(ids []string) SubscriptionToDeleteFilter {
	for i := range ids {
		f.queryParams.Add("id", ids[i])
	}

	return f
}

func (f SubscriptionToDeleteFilter) ToQueryString() string {
	if len(f.queryParams) == 0 {
		return ""
	}

	return "?" + f.queryParams.Encode()
}

func (c *Client) DeleteSubscriptions(ctx context.Context, filter SubscriptionToDeleteFilter) (EmptyResponse, error) {
	_, err := makeRequest[EmptyResponse](
		ctx,
		c,
		http.MethodDelete,
		fmt.Sprintf("/public/v1/events/subscriptions%s", filter.ToQueryString()),
		http.StatusNoContent,
		http.NoBody,
	)
	if err != nil {
		return EmptyResponse{}, err
	}

	return EmptyResponse{}, nil
}

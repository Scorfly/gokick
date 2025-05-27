package gokick

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

type BanUserResponseWrapper Response[BanUserResponse]

type BanUserResponse struct{}

func (c *Client) BanUser(
	ctx context.Context,
	broadcasterUserID int,
	userID int,
	duration *int,
	reason *string,
) (BanUserResponseWrapper, error) {
	type postBodyRequest struct {
		BroadcasterUserID int    `json:"broadcaster_user_id"`
		Duration          int    `json:"duration,omitempty"`
		Reason            string `json:"reason,omitempty"`
		UserID            int    `json:"user_id"`
	}

	r := postBodyRequest{
		BroadcasterUserID: broadcasterUserID,
		UserID:            userID,
	}

	if duration != nil {
		r.Duration = *duration
	}

	if reason != nil {
		r.Reason = *reason
	}

	body, err := json.Marshal(r)
	if err != nil {
		return BanUserResponseWrapper{}, fmt.Errorf("failed to marshal body: %v", err)
	}

	response, err := makeRequest[BanUserResponse](
		ctx,
		c,
		http.MethodPost,
		"/public/v1/moderation/bans",
		http.StatusOK,
		bytes.NewReader(body),
	)
	if err != nil {
		return BanUserResponseWrapper{}, err
	}

	return BanUserResponseWrapper(response), nil
}

func (c *Client) UnbanUser(
	ctx context.Context,
	broadcasterUserID int,
	userID int,
) (BanUserResponseWrapper, error) {
	type postBodyRequest struct {
		BroadcasterUserID int `json:"broadcaster_user_id"`
		UserID            int `json:"user_id"`
	}

	r := postBodyRequest{
		BroadcasterUserID: broadcasterUserID,
		UserID:            userID,
	}

	body, err := json.Marshal(r)
	if err != nil {
		return BanUserResponseWrapper{}, fmt.Errorf("failed to marshal body: %v", err)
	}

	response, err := makeRequest[BanUserResponse](
		ctx,
		c,
		http.MethodDelete,
		"/public/v1/moderation/bans",
		http.StatusOK,
		bytes.NewReader(body),
	)
	if err != nil {
		return BanUserResponseWrapper{}, err
	}

	return BanUserResponseWrapper(response), nil
}

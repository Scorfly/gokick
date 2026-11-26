package gokick

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

// AdBreakResponse is the payload of a successfully accepted ad break.
type AdBreakResponse struct {
	AdBreakID         string `json:"ad_break_id"`
	ID                string `json:"id"`
	RemainingAdBreaks int    `json:"remaining_ad_breaks"`
}

// CreateAdBreakRequest starts an ad break on the authenticated broadcaster's channel.
// ID must be a client-generated UUID. BreakDurationSeconds must be between 7 and 90.
type CreateAdBreakRequest struct {
	BreakDurationSeconds int    `json:"break_duration_seconds"`
	ID                   string `json:"id"`
}

// CreateAdBreak starts an ad break on the authenticated broadcaster's channel.
// Kick responds with 202 Accepted. Requires the ads:write scope.
// When Kick rejects the break, Error.Description holds the machine-readable code.
// Parse it with NewAdBreakErrorCode.
func (c *Client) CreateAdBreak(ctx context.Context, req CreateAdBreakRequest) (Response[AdBreakResponse], error) {
	body, err := json.Marshal(req)
	if err != nil {
		return Response[AdBreakResponse]{}, fmt.Errorf("failed to marshal body: %v", err)
	}

	response, err := makeRequest[AdBreakResponse](
		ctx,
		c,
		http.MethodPost,
		"/public/v1/ads/ad-break",
		http.StatusAccepted,
		bytes.NewReader(body),
	)
	if err != nil {
		return Response[AdBreakResponse]{}, err
	}

	return response, nil
}

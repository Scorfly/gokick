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

// AdBreakLimits is how many ad breaks a channel may run in a period.
type AdBreakLimits struct {
	Max           int `json:"max"`
	PeriodSeconds int `json:"period_seconds"`
}

// AdBreakStatusItem is one ad break returned by GetAdBreakStatus.
// AdBreakID is nil until Kick accepts the break.
type AdBreakStatusItem struct {
	AdBreakID            *string `json:"ad_break_id,omitempty"`
	BreakDurationSeconds int     `json:"break_duration_seconds"`
	CreatedAt            string  `json:"created_at"`
	ID                   string  `json:"id"`
	Source               string  `json:"source"`
	Status               string  `json:"status"`
}

// AdBreakStatusResponse is the payload of GET /public/v1/ads/ad-break-status.
type AdBreakStatusResponse struct {
	AdBreaks          []AdBreakStatusItem `json:"ad_breaks"`
	AdsBlocked        bool                `json:"ads_blocked"`
	AutoAdsActive     bool                `json:"auto_ads_active"`
	Limits            AdBreakLimits       `json:"limits"`
	OptedIn           bool                `json:"opted_in"`
	RemainingAdBreaks int                 `json:"remaining_ad_breaks"`
}

// GetAdBreakStatus returns ad break status for the authenticated broadcaster's channel.
// Requires the ads:read or ads:write scope.
func (c *Client) GetAdBreakStatus(ctx context.Context) (Response[AdBreakStatusResponse], error) {
	response, err := makeRequest[AdBreakStatusResponse](
		ctx,
		c,
		http.MethodGet,
		"/public/v1/ads/ad-break-status",
		http.StatusOK,
		http.NoBody,
	)
	if err != nil {
		return Response[AdBreakStatusResponse]{}, err
	}

	return response, nil
}

// EnrollInAds enrolls the authenticated broadcaster's channel in ads.
// Kick responds with 200 OK and an empty payload. Requires the ads:write scope.
// When Kick rejects enrollment, Error.Description holds the machine-readable code.
// Parse it with NewAdsEnrollErrorCode.
func (c *Client) EnrollInAds(ctx context.Context) (EmptyResponse, error) {
	_, err := makeRequest[EmptyResponse](
		ctx,
		c,
		http.MethodPost,
		"/public/v1/ads/enroll",
		http.StatusOK,
		http.NoBody,
	)
	if err != nil {
		return EmptyResponse{}, err
	}

	return EmptyResponse{}, nil
}

package gokick

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
)

type (
	DropsClaimsResponseWrapper Response[DropsClaimsData]
)

// DropsClaimResponse is a single Drops reward claim.
type DropsClaimResponse struct {
	CampaignID     string `json:"campaign_id"`
	ClaimID        string `json:"claim_id"`
	CreatedAt      string `json:"created_at"`
	ExternalID     string `json:"external_id"`
	ExternalStatus string `json:"external_status"`
	RewardID       string `json:"reward_id"`
	UpdatedAt      string `json:"updated_at"`
	UserID         int    `json:"user_id"`
}

// DropsClaimsData is the payload of GET /public/v1/drops/claims
// (claims list plus a cursor for the next page).
type DropsClaimsData struct {
	Claims []DropsClaimResponse `json:"claims"`
	Cursor string               `json:"cursor"`
}

// DropsClaimUpdate is one claim status update for PATCH /public/v1/drops/claims.
type DropsClaimUpdate struct {
	ClaimID        string `json:"claim_id"`
	ExternalStatus string `json:"external_status"`
}

// DropsClaimsFilter builds query parameters for GET /public/v1/drops/claims.
// Only OAuth apps associated with a Kick organization can call this endpoint.
type DropsClaimsFilter struct {
	queryParams url.Values
}

func NewDropsClaimsFilter() DropsClaimsFilter {
	return DropsClaimsFilter{queryParams: make(url.Values)}
}

func (f DropsClaimsFilter) SetCampaignID(campaignID string) DropsClaimsFilter {
	f.queryParams.Set("campaign_id", campaignID)
	return f
}

func (f DropsClaimsFilter) SetLimit(limit int) DropsClaimsFilter {
	f.queryParams.Set("limit", strconv.Itoa(limit))
	return f
}

func (f DropsClaimsFilter) SetCursor(cursor string) DropsClaimsFilter {
	f.queryParams.Set("cursor", cursor)
	return f
}

func (f DropsClaimsFilter) SetUserID(userID int) DropsClaimsFilter {
	f.queryParams.Set("user_id", strconv.Itoa(userID))
	return f
}

func (f DropsClaimsFilter) SetClaimID(claimID string) DropsClaimsFilter {
	f.queryParams.Set("claim_id", claimID)
	return f
}

func (f DropsClaimsFilter) SetExternalStatus(status string) DropsClaimsFilter {
	f.queryParams.Set("external_status", status)
	return f
}

func (f DropsClaimsFilter) ToQueryString() string {
	if len(f.queryParams) == 0 {
		return ""
	}
	return "?" + f.queryParams.Encode()
}

// GetDropsClaims retrieves Drops reward claims for campaigns owned by the organization
// linked to the app. Requires an app access token from an organization-associated OAuth app.
func (c *Client) GetDropsClaims(ctx context.Context, filter DropsClaimsFilter) (DropsClaimsResponseWrapper, error) {
	response, err := makeRequest[DropsClaimsData](
		ctx,
		c,
		http.MethodGet,
		fmt.Sprintf("/public/v1/drops/claims%s", filter.ToQueryString()),
		http.StatusOK,
		http.NoBody,
	)
	if err != nil {
		return DropsClaimsResponseWrapper{}, err
	}

	return DropsClaimsResponseWrapper(response), nil
}

// UpdateDropsClaims updates external_status for up to 100 claims.
// Requires an app access token from an organization-associated OAuth app.
func (c *Client) UpdateDropsClaims(ctx context.Context, claims []DropsClaimUpdate) (EmptyResponse, error) {
	body, err := json.Marshal(struct {
		Claims []DropsClaimUpdate `json:"claims"`
	}{Claims: claims})
	if err != nil {
		return EmptyResponse{}, fmt.Errorf("failed to marshal body: %v", err)
	}

	_, err = makeRequest[EmptyResponse](
		ctx,
		c,
		http.MethodPatch,
		"/public/v1/drops/claims",
		http.StatusNoContent,
		bytes.NewReader(body),
	)
	if err != nil {
		return EmptyResponse{}, err
	}

	return EmptyResponse{}, nil
}

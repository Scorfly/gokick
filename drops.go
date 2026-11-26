package gokick

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
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

// DropsClaimEvent is the body Kick POSTs when a viewer claims a Drops reward.
// ExternalID is empty when the reward has no external reference.
type DropsClaimEvent struct {
	ClaimID    string `json:"claim_id"`
	UserID     int    `json:"user_id"`
	CampaignID string `json:"campaign_id"`
	RewardID   string `json:"reward_id"`
	ExternalID string `json:"external_id"`
}

// dropsClaimRequestMeta reads the headers Kick sends on a Drops claim webhook.
// Those names differ from event-subscription webhooks: Kick-Event-Message-ID and Kick-Event-Timestamp.
// Go canonicalizes the message ID header to Kick-Event-Message-Id.
func dropsClaimRequestMeta(h http.Header) (signature, messageID, timestamp string) {
	signature = h.Get("Kick-Event-Signature")
	messageID = h.Get("Kick-Event-Message-Id")
	timestamp = h.Get("Kick-Event-Timestamp")
	return signature, messageID, timestamp
}

// GetDropsClaimFromRequest verifies a Drops claim webhook and parses its body.
// The signature is message ID, timestamp, and the raw body, joined by dots, same as other Kick webhooks.
// The public key is fetched from GET /public/v1/public-key unless DefaultEventPublicKey is set.
// There is no Kick-Event-Type header on this webhook.
func GetDropsClaimFromRequest(request *http.Request) (DropsClaimEvent, error) {
	if request == nil {
		return DropsClaimEvent{}, errors.New("request cannot be nil")
	}

	signature, messageID, timestamp := dropsClaimRequestMeta(request.Header)

	body, err := io.ReadAll(request.Body)
	if err != nil {
		return DropsClaimEvent{}, fmt.Errorf("failed to read body: %v", err)
	}

	if !SkipSignatureValidation {
		signed := []byte(fmt.Sprintf("%s.%s.%s", messageID, timestamp, body))
		err = checkEventSignature(request.Context(), signed, []byte(signature))
		if err != nil {
			return DropsClaimEvent{}, wrapEventSignatureError(err, "failed to verify drops claim")
		}
	}

	var claim DropsClaimEvent
	err = json.Unmarshal(body, &claim)
	if err != nil {
		return DropsClaimEvent{}, fmt.Errorf("failed to unmarshal drops claim: %v", err)
	}

	return claim, nil
}

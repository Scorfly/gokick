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
	ChannelsResponseWrapper                       Response[[]ChannelResponse]
	ChannelResponseWrapper                        Response[ChannelResponse]
	ChannelRewardsResponseWrapper                 Response[[]ChannelRewardResponse]
	ChannelRewardRedemptionsResponseWrapper       PaginatedResponse[[]RedemptionsByReward]
	FailedChannelRewardRedemptionsResponseWrapper Response[[]FailedRedemption]
)

type StreamResponse struct {
	CustomTags  []string `json:"custom_tags"`
	Key         string   `json:"key"`
	URL         string   `json:"url"`
	IsLive      bool     `json:"is_live"`
	IsMature    bool     `json:"is_mature"`
	Language    string   `json:"language"`
	StartTime   string   `json:"start_time"`
	Thumbnail   string   `json:"thumbnail"`
	ViewerCount int      `json:"viewer_count"`
}

type ChannelResponse struct {
	ActiveSubscribersCount   int              `json:"active_subscribers_count"`
	BannerPicture            string           `json:"banner_picture"`
	BroadcasterUserID        int              `json:"broadcaster_user_id"`
	CanceledSubscribersCount int              `json:"canceled_subscribers_count"`
	Category                 CategoryResponse `json:"category"`
	ChannelDescription       string           `json:"channel_description"`
	Slug                     string           `json:"slug"`
	Stream                   StreamResponse   `json:"stream"`
	StreamTitle              string           `json:"stream_title"`
}

type ChannelRewardResponse struct {
	BackgroundColor                   *string `json:"background_color,omitempty"`
	Cost                              int     `json:"cost"`
	Description                       *string `json:"description,omitempty"`
	ID                                string  `json:"id"`
	IsEnabled                         *bool   `json:"is_enabled,omitempty"`
	IsPaused                          *bool   `json:"is_paused,omitempty"`
	IsUserInputRequired               *bool   `json:"is_user_input_required,omitempty"`
	ShouldRedemptionsSkipRequestQueue *bool   `json:"should_redemptions_skip_request_queue,omitempty"`
	Title                             string  `json:"title"`
}

type MinimalChannelReward struct {
	CanManage   *bool   `json:"can_manage,omitempty"`
	Cost        *int    `json:"cost,omitempty"`
	Description *string `json:"description,omitempty"`
	ID          string  `json:"id"`
	IsDeleted   *bool   `json:"is_deleted,omitempty"`
	Title       string  `json:"title"`
}

type RedemptionUserInfo struct {
	UserID int `json:"user_id"`
}

type ChannelRewardRedemption struct {
	ID         string             `json:"id"`
	RedeemedAt string             `json:"redeemed_at"`
	Redeemer   RedemptionUserInfo `json:"redeemer"`
	Status     string             `json:"status"`
	UserInput  string             `json:"user_input"`
}

type RedemptionsByReward struct {
	Redemptions []ChannelRewardRedemption `json:"redemptions"`
	Reward      MinimalChannelReward      `json:"reward"`
}

type FailedRedemption struct {
	ID     string `json:"id"`
	Reason string `json:"reason"`
}

type ChannelListFilter struct {
	queryParams url.Values
}

func NewChannelListFilter() ChannelListFilter {
	return ChannelListFilter{queryParams: make(url.Values)}
}

func (f ChannelListFilter) SetBroadcasterUserIDs(ids []int) ChannelListFilter {
	for i := range ids {
		f.queryParams.Add("broadcaster_user_id", fmt.Sprintf("%d", ids[i]))
	}

	return f
}

func (f ChannelListFilter) SetSlug(slugs []string) ChannelListFilter {
	for i := range slugs {
		f.queryParams.Add("slug", slugs[i])
	}

	return f
}

func (f ChannelListFilter) ToQueryString() string {
	if len(f.queryParams) == 0 {
		return ""
	}

	return "?" + f.queryParams.Encode()
}

func (c *Client) GetChannels(ctx context.Context, filter ChannelListFilter) (ChannelsResponseWrapper, error) {
	response, err := makeRequest[[]ChannelResponse](
		ctx,
		c,
		http.MethodGet,
		fmt.Sprintf("/public/v1/channels%s", filter.ToQueryString()),
		http.StatusOK,
		http.NoBody,
	)
	if err != nil {
		return ChannelsResponseWrapper{}, err
	}

	return ChannelsResponseWrapper(response), nil
}

func (c *Client) UpdateStreamTitle(ctx context.Context, title string) (EmptyResponse, error) {
	type patchBodyRequest struct {
		StreamTitle string `json:"stream_title"`
	}

	body, err := json.Marshal(patchBodyRequest{StreamTitle: title})
	if err != nil {
		return EmptyResponse{}, fmt.Errorf("failed to marshal body: %v", err)
	}

	_, err = makeRequest[EmptyResponse](
		ctx,
		c,
		http.MethodPatch,
		"/public/v1/channels",
		http.StatusNoContent,
		bytes.NewReader(body),
	)
	if err != nil {
		return EmptyResponse{}, err
	}

	return EmptyResponse{}, nil
}

func (c *Client) UpdateStreamCategory(ctx context.Context, categoryID int) (EmptyResponse, error) {
	type patchBodyRequest struct {
		CategoryID int `json:"category_id"`
	}

	body, err := json.Marshal(patchBodyRequest{CategoryID: categoryID})
	if err != nil {
		return EmptyResponse{}, fmt.Errorf("failed to marshal body: %v", err)
	}

	_, err = makeRequest[EmptyResponse](
		ctx,
		c,
		http.MethodPatch,
		"/public/v1/channels",
		http.StatusNoContent,
		bytes.NewReader(body),
	)
	if err != nil {
		return EmptyResponse{}, err
	}

	return EmptyResponse{}, nil
}

func (c *Client) UpdateStreamTags(ctx context.Context, tags []string) (EmptyResponse, error) {
	type patchBodyRequest struct {
		Tags []string `json:"custom_tags"`
	}

	body, err := json.Marshal(patchBodyRequest{Tags: tags})
	if err != nil {
		return EmptyResponse{}, fmt.Errorf("failed to marshal body: %v", err)
	}

	_, err = makeRequest[EmptyResponse](
		ctx,
		c,
		http.MethodPatch,
		"/public/v1/channels",
		http.StatusNoContent,
		bytes.NewReader(body),
	)
	if err != nil {
		return EmptyResponse{}, err
	}

	return EmptyResponse{}, nil
}

func (c *Client) GetChannelRewards(ctx context.Context) (ChannelRewardsResponseWrapper, error) {
	response, err := makeRequest[[]ChannelRewardResponse](
		ctx,
		c,
		http.MethodGet,
		"/public/v1/channels/rewards",
		http.StatusOK,
		http.NoBody,
	)
	if err != nil {
		return ChannelRewardsResponseWrapper{}, err
	}

	return ChannelRewardsResponseWrapper(response), nil
}

type CreateChannelRewardRequest struct {
	BackgroundColor                   *string `json:"background_color,omitempty"`
	Cost                              int     `json:"cost"`
	Description                       *string `json:"description,omitempty"`
	IsEnabled                         *bool   `json:"is_enabled,omitempty"`
	IsUserInputRequired               *bool   `json:"is_user_input_required,omitempty"`
	ShouldRedemptionsSkipRequestQueue *bool   `json:"should_redemptions_skip_request_queue,omitempty"`
	Title                             string  `json:"title"`
}

func (c *Client) CreateChannelReward(ctx context.Context, req CreateChannelRewardRequest) (Response[ChannelRewardResponse], error) {
	body, err := json.Marshal(req)
	if err != nil {
		return Response[ChannelRewardResponse]{}, fmt.Errorf("failed to marshal body: %v", err)
	}

	response, err := makeRequest[ChannelRewardResponse](
		ctx,
		c,
		http.MethodPost,
		"/public/v1/channels/rewards",
		http.StatusOK,
		bytes.NewReader(body),
	)
	if err != nil {
		return Response[ChannelRewardResponse]{}, err
	}

	return response, nil
}

type UpdateChannelRewardRequest struct {
	BackgroundColor                   *string `json:"background_color,omitempty"`
	Cost                              *int    `json:"cost,omitempty"`
	Description                       *string `json:"description,omitempty"`
	IsEnabled                         *bool   `json:"is_enabled,omitempty"`
	IsPaused                          *bool   `json:"is_paused,omitempty"`
	IsUserInputRequired               *bool   `json:"is_user_input_required,omitempty"`
	ShouldRedemptionsSkipRequestQueue *bool   `json:"should_redemptions_skip_request_queue,omitempty"`
	Title                             *string `json:"title,omitempty"`
}

func (c *Client) UpdateChannelReward(
	ctx context.Context,
	id string,
	req UpdateChannelRewardRequest,
) (Response[ChannelRewardResponse], error) {
	body, err := json.Marshal(req)
	if err != nil {
		return Response[ChannelRewardResponse]{}, fmt.Errorf("failed to marshal body: %v", err)
	}

	response, err := makeRequest[ChannelRewardResponse](
		ctx,
		c,
		http.MethodPatch,
		fmt.Sprintf("/public/v1/channels/rewards/%s", id),
		http.StatusOK,
		bytes.NewReader(body),
	)
	if err != nil {
		return Response[ChannelRewardResponse]{}, err
	}

	return response, nil
}

func (c *Client) DeleteChannelReward(ctx context.Context, id string) (EmptyResponse, error) {
	_, err := makeRequest[EmptyResponse](
		ctx,
		c,
		http.MethodDelete,
		fmt.Sprintf("/public/v1/channels/rewards/%s", id),
		http.StatusNoContent,
		http.NoBody,
	)
	if err != nil {
		return EmptyResponse{}, err
	}

	return EmptyResponse{}, nil
}

type ChannelRewardRedemptionListFilter struct {
	cursor   string
	rewardID string
	status   *ChannelRewardRedemptionStatus
	ids      []string
}

func NewChannelRewardRedemptionListFilter() ChannelRewardRedemptionListFilter {
	return ChannelRewardRedemptionListFilter{}
}

// SetRewardID filters redemptions for a specific reward.
func (f ChannelRewardRedemptionListFilter) SetRewardID(rewardID string) ChannelRewardRedemptionListFilter {
	f.rewardID = rewardID
	return f
}

// SetStatus filters by redemption status (defaults to pending on the API when omitted).
func (f ChannelRewardRedemptionListFilter) SetStatus(status ChannelRewardRedemptionStatus) ChannelRewardRedemptionListFilter {
	f.status = &status
	return f
}

// AddID adds a redemption ID filter. When IDs are set, other filters must not be used.
func (f ChannelRewardRedemptionListFilter) AddID(id string) ChannelRewardRedemptionListFilter {
	f.ids = append(f.ids, id)
	return f
}

// SetCursor sets the pagination cursor.
func (f ChannelRewardRedemptionListFilter) SetCursor(cursor string) ChannelRewardRedemptionListFilter {
	f.cursor = cursor
	return f
}

func (f ChannelRewardRedemptionListFilter) ToQueryString() string {
	v := url.Values{}
	if f.cursor != "" {
		v.Set("cursor", f.cursor)
	}
	if f.rewardID != "" {
		v.Set("reward_id", f.rewardID)
	}
	if f.status != nil {
		v.Set("status", f.status.String())
	}
	for _, id := range f.ids {
		v.Add("id", id)
	}
	if len(v) == 0 {
		return ""
	}
	return "?" + v.Encode()
}

func (c *Client) GetChannelRewardRedemptions(
	ctx context.Context,
	filter ChannelRewardRedemptionListFilter,
) (ChannelRewardRedemptionsResponseWrapper, error) {
	response, err := makePaginatedRequest[[]RedemptionsByReward](
		ctx,
		c,
		http.MethodGet,
		fmt.Sprintf("/public/v1/channels/rewards/redemptions%s", filter.ToQueryString()),
		http.StatusOK,
		http.NoBody,
	)
	if err != nil {
		return ChannelRewardRedemptionsResponseWrapper{}, err
	}

	return ChannelRewardRedemptionsResponseWrapper(response), nil
}

func (c *Client) AcceptChannelRewardRedemptions(
	ctx context.Context,
	ids []string,
) (FailedChannelRewardRedemptionsResponseWrapper, error) {
	body, err := json.Marshal(struct {
		IDs []string `json:"ids"`
	}{IDs: ids})
	if err != nil {
		return FailedChannelRewardRedemptionsResponseWrapper{}, fmt.Errorf("failed to marshal body: %v", err)
	}

	response, err := makeRequest[[]FailedRedemption](
		ctx,
		c,
		http.MethodPost,
		"/public/v1/channels/rewards/redemptions/accept",
		http.StatusOK,
		bytes.NewReader(body),
	)
	if err != nil {
		return FailedChannelRewardRedemptionsResponseWrapper{}, err
	}

	return FailedChannelRewardRedemptionsResponseWrapper(response), nil
}

func (c *Client) RejectChannelRewardRedemptions(
	ctx context.Context,
	ids []string,
) (FailedChannelRewardRedemptionsResponseWrapper, error) {
	body, err := json.Marshal(struct {
		IDs []string `json:"ids"`
	}{IDs: ids})
	if err != nil {
		return FailedChannelRewardRedemptionsResponseWrapper{}, fmt.Errorf("failed to marshal body: %v", err)
	}

	response, err := makeRequest[[]FailedRedemption](
		ctx,
		c,
		http.MethodPost,
		"/public/v1/channels/rewards/redemptions/reject",
		http.StatusOK,
		bytes.NewReader(body),
	)
	if err != nil {
		return FailedChannelRewardRedemptionsResponseWrapper{}, err
	}

	return FailedChannelRewardRedemptionsResponseWrapper(response), nil
}

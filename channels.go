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
	ChannelsResponseWrapper       Response[[]ChannelResponse]
	ChannelResponseWrapper        Response[ChannelResponse]
	ChannelRewardsResponseWrapper Response[[]ChannelRewardResponse]
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
	BannerPicture      string           `json:"banner_picture"`
	BroadcasterUserID  int              `json:"broadcaster_user_id"`
	Category           CategoryResponse `json:"category"`
	ChannelDescription string           `json:"channel_description"`
	Slug               string           `json:"slug"`
	Stream             StreamResponse   `json:"stream"`
	StreamTitle        string           `json:"stream_title"`
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

func (c *Client) UpdateCustomTags(ctx context.Context, tags []string) (EmptyResponse, error) {
	type patchBodyRequest struct {
		CustomTags []string `json:"custom_tags"`
	}

	body, err := json.Marshal(patchBodyRequest{CustomTags: tags})
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

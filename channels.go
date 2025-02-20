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
	ChannelsResponseWrapper Response[[]ChannelResponse]
	ChannelResponseWrapper  Response[ChannelResponse]
)

type ChannelResponse struct {
	BannerPicture      string           `json:"banner_picture"`
	BroadcasterUserID  int              `json:"broadcaster_user_id"`
	Category           CategoryResponse `json:"category"`
	ChannelDescription string           `json:"channel_description"`
	Slug               string           `json:"slug"`
	Stream             struct {
		Key string `json:"key"`
		URL string `json:"url"`
	} `json:"stream"`
	StreamTitle string `json:"stream_title"`
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
		http.StatusOK,
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
		http.StatusOK,
		bytes.NewReader(body),
	)
	if err != nil {
		return EmptyResponse{}, err
	}

	return EmptyResponse{}, nil
}

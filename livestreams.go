package gokick

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
)

type (
	LivestreamsResponseWrapper Response[[]LivestreamResponse]
	LivestreamResponseWrapper  Response[LivestreamResponse]
)

type LivestreamResponse struct {
	BroadcasterUserID int                `json:"broadcaster_user_id"`
	Category          CategoryResponseV2 `json:"category"`
	ChannelID         int                `json:"channel_id"`
	HasMatureContent  bool               `json:"has_mature_content"`
	Language          string             `json:"language"`
	Slug              string             `json:"slug"`
	StartedAt         string             `json:"started_at"`
	StreamTitle       string             `json:"stream_title"`
	ThumbnailURL      string             `json:"thumbnail_url"`
	ViewerCount       int                `json:"viewer_count"`
}

type LivestreamListFilter struct {
	queryParams url.Values
}

func NewLivestreamListFilter() LivestreamListFilter {
	return LivestreamListFilter{queryParams: make(url.Values)}
}

func (f LivestreamListFilter) SetBroadcasterUserIDs(id int) LivestreamListFilter {
	f.queryParams.Add("broadcaster_user_id", fmt.Sprintf("%d", id))

	return f
}

func (f LivestreamListFilter) SetCategoryID(id int) LivestreamListFilter {
	f.queryParams.Add("category_id", fmt.Sprintf("%d", id))

	return f
}

func (f LivestreamListFilter) SetLanguage(lang string) LivestreamListFilter {
	f.queryParams.Add("language", lang)

	return f
}

func (f LivestreamListFilter) SetLimit(limit int) LivestreamListFilter {
	f.queryParams.Add("limit", fmt.Sprintf("%d", limit))

	return f
}

func (f LivestreamListFilter) SetSort(sort LivestreamSort) LivestreamListFilter {
	f.queryParams.Add("sort", sort.String())

	return f
}

func (f LivestreamListFilter) ToQueryString() string {
	if len(f.queryParams) == 0 {
		return ""
	}

	return "?" + f.queryParams.Encode()
}

func (c *Client) GetLivestreams(ctx context.Context, filter LivestreamListFilter) (LivestreamsResponseWrapper, error) {
	response, err := makeRequest[[]LivestreamResponse](
		ctx,
		c,
		http.MethodGet,
		fmt.Sprintf("/public/v1/livestreams%s", filter.ToQueryString()),
		http.StatusOK,
		http.NoBody,
	)
	if err != nil {
		return LivestreamsResponseWrapper{}, err
	}

	return LivestreamsResponseWrapper(response), nil
}

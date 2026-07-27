package gokick

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
)

type (
	LivestreamsResponseWrapper      Response[[]LivestreamResponse]
	LivestreamResponseWrapper       Response[LivestreamResponse]
	LivestreamsV2ResponseWrapper    PaginatedResponse[[]LivestreamV2Response]
	UsersLivestreamsResponseWrapper Response[[]LivestreamV2Response]
	LivestreamStatsResponseWrapper  Response[LivestreamStatsResponse]
)

type LivestreamV2Response struct {
	ID               string                    `json:"id"`
	BroadcasterUser  LivestreamBroadcasterUser `json:"broadcaster_user"`
	Category         LivestreamCategory        `json:"category"`
	Channel          LivestreamChannel         `json:"channel"`
	HasMatureContent bool                      `json:"has_mature_content"`
	LanguageCode     string                    `json:"language_code"`
	StartedAt        string                    `json:"started_at"`
	Tags             []string                  `json:"tags"`
	Thumbnail        string                    `json:"thumbnail"`
	Title            string                    `json:"title"`
	ViewerCount      int                       `json:"viewer_count"`
}

type LivestreamBroadcasterUser struct {
	ID             int    `json:"id"`
	ProfilePicture string `json:"profile_picture"`
	Username       string `json:"username"`
}

type LivestreamCategory struct {
	ID        int    `json:"id"`
	Name      string `json:"name"`
	Thumbnail string `json:"thumbnail"`
}

type LivestreamChannel struct {
	Slug string `json:"slug"`
}

type LivestreamResponse struct {
	BroadcasterUserID int              `json:"broadcaster_user_id"`
	Category          CategoryResponse `json:"category"`
	ChannelID         int              `json:"channel_id"`
	CustomTags        []string         `json:"custom_tags"`
	HasMatureContent  bool             `json:"has_mature_content"`
	Language          string           `json:"language"`
	ProfilePicture    string           `json:"profile_picture"`
	Slug              string           `json:"slug"`
	StartedAt         string           `json:"started_at"`
	StreamTitle       string           `json:"stream_title"`
	Thumbnail         string           `json:"thumbnail"`
	ViewerCount       int              `json:"viewer_count"`
}

type LivestreamStatsResponse struct {
	TotalCount int `json:"total_count"`
}

// LivestreamListFilter builds query parameters for the deprecated GET /public/v1/livestreams endpoint.
//
// Deprecated: use LivestreamV2ListFilter with GetLivestreamsV2, or UsersLivestreamsFilter with GetUsersLivestreams.
type LivestreamListFilter struct {
	queryParams url.Values
}

// NewLivestreamListFilter returns a filter for the deprecated v1 livestreams list.
//
// Deprecated: use NewLivestreamV2ListFilter or NewUsersLivestreamsFilter instead.
func NewLivestreamListFilter() LivestreamListFilter {
	return LivestreamListFilter{queryParams: make(url.Values)}
}

// SetBroadcasterUserIDs filters by broadcaster user IDs (v1 only, up to 50).
//
// Deprecated: use GetUsersLivestreams with UsersLivestreamsFilter for up to 100 user IDs.
func (f LivestreamListFilter) SetBroadcasterUserIDs(ids []int) LivestreamListFilter {
	for i := range ids {
		f.queryParams.Add("broadcaster_user_id", fmt.Sprintf("%d", ids[i]))
	}

	return f
}

func (f LivestreamListFilter) SetCategoryID(id int) LivestreamListFilter {
	f.queryParams.Add("category_id", fmt.Sprintf("%d", id))

	return f
}

// SetLanguage filters by BCP 47 language tag (v1 query param "language").
//
// Deprecated: use LivestreamV2ListFilter.AddLanguageCode with GetLivestreamsV2.
func (f LivestreamListFilter) SetLanguage(lang string) LivestreamListFilter {
	f.queryParams.Add("language", lang)

	return f
}

func (f LivestreamListFilter) SetLimit(limit int) LivestreamListFilter {
	f.queryParams.Add("limit", fmt.Sprintf("%d", limit))

	return f
}

// SetSort sets sort order (v1 only: viewer_count or started_at).
//
// Deprecated: v2 livestreams are sorted oldest-to-newest by the API; this parameter is not supported on v2.
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

type LivestreamV2ListFilter struct {
	cursor        string
	limit         *int
	categoryIDs   []int
	languageCodes []string
}

func NewLivestreamV2ListFilter() LivestreamV2ListFilter {
	return LivestreamV2ListFilter{}
}

// SetCursor sets the pagination cursor (v2).
func (f LivestreamV2ListFilter) SetCursor(cursor string) LivestreamV2ListFilter {
	f.cursor = cursor
	return f
}

// SetLimit sets the page size (1–1000, API default 100).
func (f LivestreamV2ListFilter) SetLimit(limit int) LivestreamV2ListFilter {
	f.limit = &limit
	return f
}

// AddCategoryID adds a category filter (up to 25 category IDs).
func (f LivestreamV2ListFilter) AddCategoryID(id int) LivestreamV2ListFilter {
	f.categoryIDs = append(f.categoryIDs, id)
	return f
}

// AddLanguageCode adds a BCP 47 language filter (up to 25 language codes).
func (f LivestreamV2ListFilter) AddLanguageCode(code string) LivestreamV2ListFilter {
	f.languageCodes = append(f.languageCodes, code)
	return f
}

func (f LivestreamV2ListFilter) ToQueryString() string {
	v := url.Values{}
	if f.cursor != "" {
		v.Set("cursor", f.cursor)
	}
	if f.limit != nil {
		v.Set("limit", strconv.Itoa(*f.limit))
	}
	for _, id := range f.categoryIDs {
		v.Add("category_id", strconv.Itoa(id))
	}
	for _, code := range f.languageCodes {
		v.Add("language_code", code)
	}
	if len(v) == 0 {
		return ""
	}
	return "?" + v.Encode()
}

type UsersLivestreamsFilter struct {
	userIDs []int
}

func NewUsersLivestreamsFilter() UsersLivestreamsFilter {
	return UsersLivestreamsFilter{}
}

// AddUserID adds a broadcaster user ID (required; up to 100 user IDs).
func (f UsersLivestreamsFilter) AddUserID(id int) UsersLivestreamsFilter {
	f.userIDs = append(f.userIDs, id)
	return f
}

func (f UsersLivestreamsFilter) ToQueryString() string {
	v := url.Values{}
	for _, id := range f.userIDs {
		v.Add("user_id", strconv.Itoa(id))
	}
	if len(v) == 0 {
		return ""
	}
	return "?" + v.Encode()
}

// GetLivestreams returns active livestreams from the deprecated v1 endpoint.
//
// Deprecated: use GetLivestreamsV2 or GetUsersLivestreams instead.
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

func (c *Client) GetLivestreamsV2(ctx context.Context, filter LivestreamV2ListFilter) (LivestreamsV2ResponseWrapper, error) {
	response, err := makePaginatedRequest[[]LivestreamV2Response](
		ctx,
		c,
		http.MethodGet,
		fmt.Sprintf("/public/v2/livestreams%s", filter.ToQueryString()),
		http.StatusOK,
		http.NoBody,
	)
	if err != nil {
		return LivestreamsV2ResponseWrapper{}, err
	}

	return LivestreamsV2ResponseWrapper(response), nil
}

func (c *Client) GetUsersLivestreams(ctx context.Context, filter UsersLivestreamsFilter) (UsersLivestreamsResponseWrapper, error) {
	response, err := makeRequest[[]LivestreamV2Response](
		ctx,
		c,
		http.MethodGet,
		fmt.Sprintf("/public/v1/users/livestreams%s", filter.ToQueryString()),
		http.StatusOK,
		http.NoBody,
	)
	if err != nil {
		return UsersLivestreamsResponseWrapper{}, err
	}

	return UsersLivestreamsResponseWrapper(response), nil
}

func (c *Client) GetLivestreamsStats(ctx context.Context) (LivestreamStatsResponseWrapper, error) {
	response, err := makeRequest[LivestreamStatsResponse](
		ctx,
		c,
		http.MethodGet,
		"/public/v1/livestreams/stats",
		http.StatusOK,
		http.NoBody,
	)
	if err != nil {
		return LivestreamStatsResponseWrapper{}, err
	}

	return LivestreamStatsResponseWrapper(response), nil
}

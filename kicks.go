package gokick

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
)

type (
	KicksLeaderboardResponseWrapper Response[KicksLeaderboardResponse]
)

type KicksLeaderboardResponse struct {
	Lifetime []KicksLeaderboardEntry `json:"lifetime"`
	Month    []KicksLeaderboardEntry `json:"month"`
	Week     []KicksLeaderboardEntry `json:"week"`
}

type KicksLeaderboardEntry struct {
	GiftedAmount int    `json:"gifted_amount"`
	Rank         int    `json:"rank"`
	UserID       int    `json:"user_id"`
	Username     string `json:"username"`
}

type KicksLeaderboardFilter struct {
	queryParams url.Values
}

func NewKicksLeaderboardFilter() KicksLeaderboardFilter {
	return KicksLeaderboardFilter{queryParams: make(url.Values)}
}

func (f KicksLeaderboardFilter) SetTop(top int) KicksLeaderboardFilter {
	f.queryParams.Set("top", fmt.Sprintf("%d", top))
	return f
}

func (f KicksLeaderboardFilter) ToQueryString() string {
	if len(f.queryParams) == 0 {
		return ""
	}

	return "?" + f.queryParams.Encode()
}

func (c *Client) GetKicksLeaderboard(ctx context.Context, filter KicksLeaderboardFilter) (KicksLeaderboardResponseWrapper, error) {
	response, err := makeRequest[KicksLeaderboardResponse](
		ctx,
		c,
		http.MethodGet,
		fmt.Sprintf("/public/v1/kicks/leaderboard%s", filter.ToQueryString()),
		http.StatusOK,
		http.NoBody,
	)
	if err != nil {
		return KicksLeaderboardResponseWrapper{}, err
	}

	return KicksLeaderboardResponseWrapper(response), nil
}

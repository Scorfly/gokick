package gokick_test

import (
	"context"
	"fmt"
	"net/http"
	"testing"

	"github.com/scorfly/gokick"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewKicksLeaderboardFilterSuccess(t *testing.T) {
	testCases := map[string]struct {
		filter              gokick.KicksLeaderboardFilter
		expectedQueryString string
	}{
		"default": {
			filter:              gokick.NewKicksLeaderboardFilter(),
			expectedQueryString: "",
		},
		"with top parameter": {
			filter:              gokick.NewKicksLeaderboardFilter().SetTop(10),
			expectedQueryString: "?top=10",
		},
		"with maximum top parameter": {
			filter:              gokick.NewKicksLeaderboardFilter().SetTop(100),
			expectedQueryString: "?top=100",
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			assert.Equal(t, tc.expectedQueryString, tc.filter.ToQueryString())
		})
	}
}

func TestGetKicksLeaderboardError(t *testing.T) {
	t.Run("on new request", func(t *testing.T) {
		kickClient, err := gokick.NewClient(&gokick.ClientOptions{AppAccessToken: "access-token"})
		require.NoError(t, err)

		var ctx context.Context
		_, err = kickClient.GetKicksLeaderboard(ctx, gokick.NewKicksLeaderboardFilter())
		require.EqualError(t, err, "failed to create request: net/http: nil Context")
	})

	t.Run("timeout", func(t *testing.T) {
		kickClient := setupTimeoutMockClient(t)

		ctx := context.Background()
		_, err := kickClient.GetKicksLeaderboard(ctx, gokick.NewKicksLeaderboardFilter())
		require.Error(t, err)
		require.Contains(t, err.Error(), "failed to make request:")
	})
}

func TestGetKicksLeaderboardSuccess(t *testing.T) {
	mockResponse := `{
		"data": {
			"lifetime": [
				{
					"gifted_amount": 500,
					"rank": 1,
					"user_id": 123,
					"username": "top-gifter"
				},
				{
					"gifted_amount": 250,
					"rank": 2,
					"user_id": 456,
					"username": "second-place"
				}
			],
			"month": [
				{
					"gifted_amount": 200,
					"rank": 1,
					"user_id": 789,
					"username": "monthly-leader"
				}
			],
			"week": [
				{
					"gifted_amount": 100,
					"rank": 1,
					"user_id": 321,
					"username": "weekly-champion"
				}
			]
		},
		"message": "Success"
	}`

	testCases := map[string]struct {
		filter      gokick.KicksLeaderboardFilter
		expectedURL string
	}{
		"no filter": {
			filter:      gokick.NewKicksLeaderboardFilter(),
			expectedURL: "/public/v1/kicks/leaderboard",
		},
		"with top filter": {
			filter:      gokick.NewKicksLeaderboardFilter().SetTop(10),
			expectedURL: "/public/v1/kicks/leaderboard?top=10",
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			kickClient := setupMockClient(t, func(w http.ResponseWriter, r *http.Request) {
				assert.Equal(t, http.MethodGet, r.Method)

				actualURL := r.URL.Path
				if r.URL.RawQuery != "" {
					actualURL += "?" + r.URL.RawQuery
				}
				assert.Equal(t, tc.expectedURL, actualURL)
				assert.Contains(t, r.Header.Get("Authorization"), "Bearer")

				w.WriteHeader(http.StatusOK)
				fmt.Fprint(w, mockResponse)
			})

			ctx := context.Background()
			response, err := kickClient.GetKicksLeaderboard(ctx, tc.filter)

			require.NoError(t, err)
			assert.Len(t, response.Result.Lifetime, 2)
			assert.Equal(t, 123, response.Result.Lifetime[0].UserID)
			assert.Equal(t, "top-gifter", response.Result.Lifetime[0].Username)
			assert.Equal(t, 500, response.Result.Lifetime[0].GiftedAmount)
			assert.Equal(t, 1, response.Result.Lifetime[0].Rank)

			assert.Len(t, response.Result.Month, 1)
			assert.Equal(t, 789, response.Result.Month[0].UserID)
			assert.Equal(t, "monthly-leader", response.Result.Month[0].Username)

			assert.Len(t, response.Result.Week, 1)
			assert.Equal(t, 321, response.Result.Week[0].UserID)
			assert.Equal(t, "weekly-champion", response.Result.Week[0].Username)
		})
	}
}

func TestGetKicksLeaderboardAPIError(t *testing.T) {
	kickClient := setupMockClient(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, `{
			"data": null,
			"message": "Invalid top parameter. Must be between 1 and 100."
		}`)
	})

	ctx := context.Background()
	_, err := kickClient.GetKicksLeaderboard(ctx, gokick.NewKicksLeaderboardFilter().SetTop(101))

	require.Error(t, err)
	kickError, ok := err.(gokick.Error)
	require.True(t, ok)
	assert.Equal(t, 400, kickError.Code())
	assert.Equal(t, "Invalid top parameter. Must be between 1 and 100.", kickError.Message())
}

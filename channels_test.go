package gokick_test

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	"github.com/scorfly/gokick"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewChannelListFilterSuccess(t *testing.T) {
	testCases := map[string]struct {
		filter              gokick.ChannelListFilter
		expectedQueryString string
	}{
		"default": {
			filter:              gokick.NewChannelListFilter(),
			expectedQueryString: "",
		},
		"with broadcaster_user_id query": {
			filter:              gokick.NewChannelListFilter().SetBroadcasterUserIDs([]int{118, 218}),
			expectedQueryString: "?broadcaster_user_id=118&broadcaster_user_id=218",
		},
		"with slug query": {
			filter:              gokick.NewChannelListFilter().SetSlug([]string{"slug_1", "slug_2"}),
			expectedQueryString: "?slug=slug_1&slug=slug_2",
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			assert.Equal(t, tc.expectedQueryString, tc.filter.ToQueryString())
		})
	}
}

func TestGetChannelsError(t *testing.T) {
	t.Run("on new request", func(t *testing.T) {
		kickClient, err := gokick.NewClient(&gokick.ClientOptions{UserAccessToken: "access-token"})
		require.NoError(t, err)

		var ctx context.Context
		_, err = kickClient.GetChannels(ctx, gokick.NewChannelListFilter())
		require.EqualError(t, err, "failed to create request: net/http: nil Context")
	})

	t.Run("timeout", func(t *testing.T) {
		kickClient := setupTimeoutMockClient(t)

		_, err := kickClient.GetChannels(context.Background(), gokick.NewChannelListFilter())
		require.EqualError(t, err, `failed to make request: Get "https://api.kick.com/public/v1/channels": context deadline exceeded `+
			`(Client.Timeout exceeded while awaiting headers)`)
	})

	t.Run("unmarshal error response", func(t *testing.T) {
		kickClient := setupMockClient(t, func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprint(w, `117`)
		})

		_, err := kickClient.GetChannels(context.Background(), gokick.NewChannelListFilter())

		assert.EqualError(t, err, `failed to unmarshal error response (KICK status code: 500 and body "117"): json: cannot unmarshal `+
			`number into Go value of type gokick.errorResponse`)
	})

	t.Run("unmarshal channels response", func(t *testing.T) {
		kickClient := setupMockClient(t, func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			fmt.Fprint(w, "117")
		})

		_, err := kickClient.GetChannels(context.Background(), gokick.NewChannelListFilter())

		assert.EqualError(t, err, `failed to unmarshal response body (KICK status code 200 and body "117"): json: cannot unmarshal `+
			`number into Go value of type gokick.successResponse[[]github.com/scorfly/gokick.ChannelResponse]`)
	})

	t.Run("reader failure", func(t *testing.T) {
		kickClient := setupMockClient(t, func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Length", "10")
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprint(w, "")
		})

		_, err := kickClient.GetChannels(context.Background(), gokick.NewChannelListFilter())

		assert.EqualError(t, err, `failed to read response body (KICK status code 500): unexpected EOF`)
	})

	t.Run("with internal server error", func(t *testing.T) {
		kickClient := setupMockClient(t, func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprint(w, `{"message":"internal server error", "data":null}`)
		})

		_, err := kickClient.GetChannels(context.Background(), gokick.NewChannelListFilter())

		var kickError gokick.Error
		require.ErrorAs(t, err, &kickError)
		assert.Equal(t, http.StatusInternalServerError, kickError.Code())
		assert.Equal(t, "internal server error", kickError.Message())
	})
}

func TestGetChannelsSuccess(t *testing.T) {
	t.Run("without result", func(t *testing.T) {
		kickClient := setupMockClient(t, func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			fmt.Fprint(w, `{"message":"success", "data":[]}`)
		})

		channelsResponse, err := kickClient.GetChannels(context.Background(), gokick.NewChannelListFilter())
		require.NoError(t, err)
		assert.Empty(t, channelsResponse.Result)
	})

	t.Run("with result", func(t *testing.T) {
		kickClient := setupMockClient(t, func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			fmt.Fprint(w, `{"message":"success", "data":[{
				"active_subscribers_count": 42,
				"banner_picture": "banner picture",
				"broadcaster_user_id": 117,
				"canceled_subscribers_count": 3,
				"category": {
					"id": 1,
					"name": "category name",
					"thumbnail": "category thumbnail"
				},
				"channel_description": "channel description",
				"slug": "slug",
				"stream": {
					"key": "stream key",
					"url": "stream URL",
					"custom_tags": ["tag1", "tag2"]
				},
				"stream_title": "stream title"
			}]}`)
		})

		channelsResponse, err := kickClient.GetChannels(context.Background(), gokick.NewChannelListFilter())
		require.NoError(t, err)
		require.Len(t, channelsResponse.Result, 1)
		assert.Equal(t, 42, channelsResponse.Result[0].ActiveSubscribersCount)
		assert.Equal(t, "banner picture", channelsResponse.Result[0].BannerPicture)
		assert.Equal(t, 117, channelsResponse.Result[0].BroadcasterUserID)
		assert.Equal(t, 3, channelsResponse.Result[0].CanceledSubscribersCount)
		assert.Equal(t, 1, channelsResponse.Result[0].Category.ID)
		assert.Equal(t, "category name", channelsResponse.Result[0].Category.Name)
		assert.Equal(t, "category thumbnail", channelsResponse.Result[0].Category.Thumbnail)
		assert.Equal(t, "channel description", channelsResponse.Result[0].ChannelDescription)
		assert.Equal(t, "slug", channelsResponse.Result[0].Slug)
		assert.Equal(t, "stream key", channelsResponse.Result[0].Stream.Key)
		assert.Equal(t, "stream URL", channelsResponse.Result[0].Stream.URL)
		assert.Equal(t, []string{"tag1", "tag2"}, channelsResponse.Result[0].Stream.CustomTags)
		assert.Equal(t, "stream title", channelsResponse.Result[0].StreamTitle)
	})
}

func TestUpdateStreamTitleError(t *testing.T) {
	t.Run("on new request", func(t *testing.T) {
		kickClient, err := gokick.NewClient(&gokick.ClientOptions{UserAccessToken: "access-token"})
		require.NoError(t, err)

		var ctx context.Context
		_, err = kickClient.UpdateStreamTitle(ctx, "new stream title")
		require.EqualError(t, err, "failed to create request: net/http: nil Context")
	})

	t.Run("timeout", func(t *testing.T) {
		kickClient := setupTimeoutMockClient(t)

		_, err := kickClient.UpdateStreamTitle(context.Background(), "new stream title")
		require.EqualError(t, err, `failed to make request: Patch "https://api.kick.com/public/v1/channels": context deadline exceeded `+
			`(Client.Timeout exceeded while awaiting headers)`)
	})

	t.Run("unmarshal error response", func(t *testing.T) {
		kickClient := setupMockClient(t, func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprint(w, `117`)
		})

		_, err := kickClient.UpdateStreamTitle(context.Background(), "new stream title")

		assert.EqualError(t, err, `failed to unmarshal error response (KICK status code: 500 and body "117"): json: cannot unmarshal `+
			`number into Go value of type gokick.errorResponse`)
	})

	t.Run("unmarshal token response", func(t *testing.T) {
		kickClient := setupMockClient(t, func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			fmt.Fprint(w, "117")
		})

		_, err := kickClient.UpdateStreamTitle(context.Background(), "new stream title")

		assert.EqualError(t, err, `failed to unmarshal error response (KICK status code: 200 and body "117"): json: cannot unmarshal`+
			` number into Go value of type gokick.errorResponse`)
	})

	t.Run("reader failure", func(t *testing.T) {
		kickClient := setupMockClient(t, func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Length", "10")
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprint(w, "")
		})

		_, err := kickClient.UpdateStreamTitle(context.Background(), "new stream title")

		assert.EqualError(t, err, `failed to read response body (KICK status code 500): unexpected EOF`)
	})

	t.Run("with internal server error", func(t *testing.T) {
		kickClient := setupMockClient(t, func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprint(w, `{"message":"internal server error", "data":null}`)
		})

		_, err := kickClient.UpdateStreamTitle(context.Background(), "new stream title")

		var kickError gokick.Error
		require.ErrorAs(t, err, &kickError)
		assert.Equal(t, http.StatusInternalServerError, kickError.Code())
		assert.Equal(t, "internal server error", kickError.Message())
	})
}

func TestUpdateStreamTitleSuccess(t *testing.T) {
	kickClient := setupMockClient(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNoContent)
	})

	_, err := kickClient.UpdateStreamTitle(context.Background(), "new stream title")
	require.NoError(t, err)
}

func TestUpdateStreamCategoryError(t *testing.T) {
	t.Run("on new request", func(t *testing.T) {
		kickClient, err := gokick.NewClient(&gokick.ClientOptions{UserAccessToken: "access-token"})
		require.NoError(t, err)

		var ctx context.Context
		_, err = kickClient.UpdateStreamCategory(ctx, 117)
		require.EqualError(t, err, "failed to create request: net/http: nil Context")
	})

	t.Run("timeout", func(t *testing.T) {
		kickClient := setupTimeoutMockClient(t)

		_, err := kickClient.UpdateStreamCategory(context.Background(), 117)
		require.EqualError(t, err, `failed to make request: Patch "https://api.kick.com/public/v1/channels": context deadline exceeded `+
			`(Client.Timeout exceeded while awaiting headers)`)
	})

	t.Run("unmarshal error response", func(t *testing.T) {
		kickClient := setupMockClient(t, func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprint(w, `117`)
		})

		_, err := kickClient.UpdateStreamCategory(context.Background(), 117)

		assert.EqualError(t, err, `failed to unmarshal error response (KICK status code: 500 and body "117"): json: cannot unmarshal `+
			`number into Go value of type gokick.errorResponse`)
	})

	t.Run("unmarshal token response", func(t *testing.T) {
		kickClient := setupMockClient(t, func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			fmt.Fprint(w, "117")
		})

		_, err := kickClient.UpdateStreamCategory(context.Background(), 117)

		assert.EqualError(t, err, `failed to unmarshal error response (KICK status code: 200 and body "117"): json: cannot unmarshal `+
			`number into Go value of type gokick.errorResponse`)
	})

	t.Run("reader failure", func(t *testing.T) {
		kickClient := setupMockClient(t, func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Length", "10")
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprint(w, "")
		})

		_, err := kickClient.UpdateStreamCategory(context.Background(), 117)

		assert.EqualError(t, err, `failed to read response body (KICK status code 500): unexpected EOF`)
	})

	t.Run("with internal server error", func(t *testing.T) {
		kickClient := setupMockClient(t, func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprint(w, `{"message":"internal server error", "data":null}`)
		})

		_, err := kickClient.UpdateStreamCategory(context.Background(), 117)

		var kickError gokick.Error
		require.ErrorAs(t, err, &kickError)
		assert.Equal(t, http.StatusInternalServerError, kickError.Code())
		assert.Equal(t, "internal server error", kickError.Message())
	})
}

func TestUpdateStreamCategorySuccess(t *testing.T) {
	kickClient := setupMockClient(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNoContent)
	})

	_, err := kickClient.UpdateStreamCategory(context.Background(), 117)
	require.NoError(t, err)
}

func TestUpdateStreamTagsError(t *testing.T) {
	t.Run("on new request", func(t *testing.T) {
		kickClient, err := gokick.NewClient(&gokick.ClientOptions{UserAccessToken: "access-token"})
		require.NoError(t, err)

		var ctx context.Context
		_, err = kickClient.UpdateStreamTags(ctx, []string{"tag1", "tag2"})
		require.EqualError(t, err, "failed to create request: net/http: nil Context")
	})

	t.Run("timeout", func(t *testing.T) {
		kickClient := setupTimeoutMockClient(t)

		_, err := kickClient.UpdateStreamTags(context.Background(), []string{"tag1", "tag2"})
		require.EqualError(t, err, `failed to make request: Patch "https://api.kick.com/public/v1/channels": context deadline exceeded `+
			`(Client.Timeout exceeded while awaiting headers)`)
	})

	t.Run("unmarshal error response", func(t *testing.T) {
		kickClient := setupMockClient(t, func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprint(w, `117`)
		})

		_, err := kickClient.UpdateStreamTags(context.Background(), []string{"tag1", "tag2"})

		assert.EqualError(t, err, `failed to unmarshal error response (KICK status code: 500 and body "117"): json: cannot unmarshal `+
			`number into Go value of type gokick.errorResponse`)
	})

	t.Run("unmarshal token response", func(t *testing.T) {
		kickClient := setupMockClient(t, func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			fmt.Fprint(w, "117")
		})

		_, err := kickClient.UpdateStreamTags(context.Background(), []string{"tag1", "tag2"})

		assert.EqualError(t, err, `failed to unmarshal error response (KICK status code: 200 and body "117"): json: cannot unmarshal `+
			`number into Go value of type gokick.errorResponse`)
	})

	t.Run("reader failure", func(t *testing.T) {
		kickClient := setupMockClient(t, func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Length", "10")
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprint(w, "")
		})

		_, err := kickClient.UpdateStreamTags(context.Background(), []string{"tag1", "tag2"})

		assert.EqualError(t, err, `failed to read response body (KICK status code 500): unexpected EOF`)
	})

	t.Run("with internal server error", func(t *testing.T) {
		kickClient := setupMockClient(t, func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprint(w, `{"message":"internal server error", "data":null}`)
		})

		_, err := kickClient.UpdateStreamTags(context.Background(), []string{"tag1", "tag2"})

		var kickError gokick.Error
		require.ErrorAs(t, err, &kickError)
		assert.Equal(t, http.StatusInternalServerError, kickError.Code())
		assert.Equal(t, "internal server error", kickError.Message())
	})
}

func TestUpdateStreamTagsSuccess(t *testing.T) {
	kickClient := setupMockClient(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNoContent)
	})

	_, err := kickClient.UpdateStreamTags(context.Background(), []string{"tag1", "tag2"})
	require.NoError(t, err)
}

func TestGetChannelRewardsError(t *testing.T) {
	t.Run("on new request", func(t *testing.T) {
		kickClient, err := gokick.NewClient(&gokick.ClientOptions{UserAccessToken: "access-token"})
		require.NoError(t, err)

		var ctx context.Context
		_, err = kickClient.GetChannelRewards(ctx)
		require.EqualError(t, err, "failed to create request: net/http: nil Context")
	})

	t.Run("timeout", func(t *testing.T) {
		kickClient := setupTimeoutMockClient(t)

		_, err := kickClient.GetChannelRewards(context.Background())
		require.EqualError(t, err, `failed to make request: Get "https://api.kick.com/public/v1/channels/rewards": context deadline exceeded `+
			`(Client.Timeout exceeded while awaiting headers)`)
	})

	t.Run("unmarshal error response", func(t *testing.T) {
		kickClient := setupMockClient(t, func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprint(w, `117`)
		})

		_, err := kickClient.GetChannelRewards(context.Background())

		assert.EqualError(t, err, `failed to unmarshal error response (KICK status code: 500 and body "117"): json: cannot unmarshal `+
			`number into Go value of type gokick.errorResponse`)
	})

	t.Run("unmarshal rewards response", func(t *testing.T) {
		kickClient := setupMockClient(t, func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			fmt.Fprint(w, "117")
		})

		_, err := kickClient.GetChannelRewards(context.Background())

		assert.EqualError(t, err, `failed to unmarshal response body (KICK status code 200 and body "117"): json: cannot unmarshal `+
			`number into Go value of type gokick.successResponse[[]github.com/scorfly/gokick.ChannelRewardResponse]`)
	})

	t.Run("reader failure", func(t *testing.T) {
		kickClient := setupMockClient(t, func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Length", "10")
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprint(w, "")
		})

		_, err := kickClient.GetChannelRewards(context.Background())

		assert.EqualError(t, err, `failed to read response body (KICK status code 500): unexpected EOF`)
	})

	t.Run("with internal server error", func(t *testing.T) {
		kickClient := setupMockClient(t, func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprint(w, `{"message":"internal server error", "data":null}`)
		})

		_, err := kickClient.GetChannelRewards(context.Background())

		var kickError gokick.Error
		require.ErrorAs(t, err, &kickError)
		assert.Equal(t, http.StatusInternalServerError, kickError.Code())
		assert.Equal(t, "internal server error", kickError.Message())
	})
}

func TestGetChannelRewardsSuccess(t *testing.T) {
	t.Run("without result", func(t *testing.T) {
		kickClient := setupMockClient(t, func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			fmt.Fprint(w, `{"message":"text", "data":[]}`)
		})

		rewardsResponse, err := kickClient.GetChannelRewards(context.Background())
		require.NoError(t, err)
		assert.Empty(t, rewardsResponse.Result)
	})

	t.Run("with result", func(t *testing.T) {
		kickClient := setupMockClient(t, func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			fmt.Fprint(w, `{"message":"text", "data":[{
				"background_color": "#00e701",
				"cost": 100,
				"description": "Request a song by providing a URL",
				"id": "01HZ8X9K2M4N6P8Q0R2S4T6V8W0Y2Z4",
				"is_enabled": true,
				"is_paused": false,
				"is_user_input_required": false,
				"should_redemptions_skip_request_queue": false,
				"title": "Song Request"
			}]}`)
		})

		rewardsResponse, err := kickClient.GetChannelRewards(context.Background())
		require.NoError(t, err)
		require.Len(t, rewardsResponse.Result, 1)
		require.NotNil(t, rewardsResponse.Result[0].BackgroundColor)
		assert.Equal(t, "#00e701", *rewardsResponse.Result[0].BackgroundColor)
		assert.Equal(t, 100, rewardsResponse.Result[0].Cost)
		require.NotNil(t, rewardsResponse.Result[0].Description)
		assert.Equal(t, "Request a song by providing a URL", *rewardsResponse.Result[0].Description)
		assert.Equal(t, "01HZ8X9K2M4N6P8Q0R2S4T6V8W0Y2Z4", rewardsResponse.Result[0].ID)
		require.NotNil(t, rewardsResponse.Result[0].IsEnabled)
		assert.True(t, *rewardsResponse.Result[0].IsEnabled)
		require.NotNil(t, rewardsResponse.Result[0].IsPaused)
		assert.False(t, *rewardsResponse.Result[0].IsPaused)
		require.NotNil(t, rewardsResponse.Result[0].IsUserInputRequired)
		assert.False(t, *rewardsResponse.Result[0].IsUserInputRequired)
		require.NotNil(t, rewardsResponse.Result[0].ShouldRedemptionsSkipRequestQueue)
		assert.False(t, *rewardsResponse.Result[0].ShouldRedemptionsSkipRequestQueue)
		assert.Equal(t, "Song Request", rewardsResponse.Result[0].Title)
	})

	t.Run("with only required fields", func(t *testing.T) {
		kickClient := setupMockClient(t, func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			fmt.Fprint(w, `{"message":"text", "data":[{
				"cost": 50,
				"id": "01HZ8X9K2M4N6P8Q0R2S4T6V8W0Y2Z4",
				"title": "Basic Reward"
			}]}`)
		})

		rewardsResponse, err := kickClient.GetChannelRewards(context.Background())
		require.NoError(t, err)
		require.Len(t, rewardsResponse.Result, 1)
		assert.Nil(t, rewardsResponse.Result[0].BackgroundColor)
		assert.Equal(t, 50, rewardsResponse.Result[0].Cost)
		assert.Nil(t, rewardsResponse.Result[0].Description)
		assert.Equal(t, "01HZ8X9K2M4N6P8Q0R2S4T6V8W0Y2Z4", rewardsResponse.Result[0].ID)
		assert.Nil(t, rewardsResponse.Result[0].IsEnabled)
		assert.Nil(t, rewardsResponse.Result[0].IsPaused)
		assert.Nil(t, rewardsResponse.Result[0].IsUserInputRequired)
		assert.Nil(t, rewardsResponse.Result[0].ShouldRedemptionsSkipRequestQueue)
		assert.Equal(t, "Basic Reward", rewardsResponse.Result[0].Title)
	})
}

func TestCreateChannelRewardError(t *testing.T) {
	t.Run("on new request", func(t *testing.T) {
		kickClient, err := gokick.NewClient(&gokick.ClientOptions{UserAccessToken: "access-token"})
		require.NoError(t, err)

		var ctx context.Context
		_, err = kickClient.CreateChannelReward(ctx, gokick.CreateChannelRewardRequest{
			Cost:  100,
			Title: "Song Request",
		})
		require.EqualError(t, err, "failed to create request: net/http: nil Context")
	})

	t.Run("timeout", func(t *testing.T) {
		kickClient := setupTimeoutMockClient(t)

		_, err := kickClient.CreateChannelReward(context.Background(), gokick.CreateChannelRewardRequest{
			Cost:  100,
			Title: "Song Request",
		})
		require.EqualError(t, err, `failed to make request: Post "https://api.kick.com/public/v1/channels/rewards": context deadline exceeded `+
			`(Client.Timeout exceeded while awaiting headers)`)
	})

	t.Run("unmarshal error response", func(t *testing.T) {
		kickClient := setupMockClient(t, func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprint(w, `117`)
		})

		_, err := kickClient.CreateChannelReward(context.Background(), gokick.CreateChannelRewardRequest{
			Cost:  100,
			Title: "Song Request",
		})

		assert.EqualError(t, err, `failed to unmarshal error response (KICK status code: 500 and body "117"): json: cannot unmarshal `+
			`number into Go value of type gokick.errorResponse`)
	})

	t.Run("unmarshal reward response", func(t *testing.T) {
		kickClient := setupMockClient(t, func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			fmt.Fprint(w, "117")
		})

		_, err := kickClient.CreateChannelReward(context.Background(), gokick.CreateChannelRewardRequest{
			Cost:  100,
			Title: "Song Request",
		})

		assert.EqualError(t, err, `failed to unmarshal response body (KICK status code 200 and body "117"): json: cannot unmarshal `+
			`number into Go value of type gokick.successResponse[github.com/scorfly/gokick.ChannelRewardResponse]`)
	})

	t.Run("reader failure", func(t *testing.T) {
		kickClient := setupMockClient(t, func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Length", "10")
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprint(w, "")
		})

		_, err := kickClient.CreateChannelReward(context.Background(), gokick.CreateChannelRewardRequest{
			Cost:  100,
			Title: "Song Request",
		})

		assert.EqualError(t, err, `failed to read response body (KICK status code 500): unexpected EOF`)
	})

	t.Run("with internal server error", func(t *testing.T) {
		kickClient := setupMockClient(t, func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprint(w, `{"message":"internal server error", "data":null}`)
		})

		_, err := kickClient.CreateChannelReward(context.Background(), gokick.CreateChannelRewardRequest{
			Cost:  100,
			Title: "Song Request",
		})

		var kickError gokick.Error
		require.ErrorAs(t, err, &kickError)
		assert.Equal(t, http.StatusInternalServerError, kickError.Code())
		assert.Equal(t, "internal server error", kickError.Message())
	})
}

func TestCreateChannelRewardSuccess(t *testing.T) {
	t.Run("with all fields", func(t *testing.T) {
		backgroundColor := "#fff111"
		description := "Request a song by providing a URL"
		isEnabled := true
		isUserInputRequired := false
		shouldRedemptionsSkipRequestQueue := false

		kickClient := setupMockClient(t, func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, http.MethodPost, r.Method)
			assert.Equal(t, "/public/v1/channels/rewards", r.URL.Path)
			assert.Equal(t, "application/json", r.Header.Get("Content-Type"))
			assert.Equal(t, "Bearer access-token", r.Header.Get("Authorization"))

			var body map[string]interface{}
			err := json.NewDecoder(r.Body).Decode(&body)
			assert.NoError(t, err)

			assert.Equal(t, "#fff111", body["background_color"])
			assert.InEpsilon(t, float64(100), body["cost"], 0)
			assert.Equal(t, "Request a song by providing a URL", body["description"])
			assert.Equal(t, true, body["is_enabled"])
			assert.Equal(t, false, body["is_user_input_required"])
			assert.Equal(t, false, body["should_redemptions_skip_request_queue"])
			assert.Equal(t, "Song Request", body["title"])

			w.WriteHeader(http.StatusOK)
			fmt.Fprint(w, `{
				"data": {
					"background_color": "#00e701",
					"cost": 100,
					"description": "Request a song by providing a URL",
					"id": "01HZ8X9K2M4N6P8Q0R2S4T6V8W0Y2Z4",
					"is_enabled": true,
					"is_paused": false,
					"is_user_input_required": false,
					"should_redemptions_skip_request_queue": false,
					"title": "Song Request"
				},
				"message": "text"
			}`)
		})

		rewardResponse, err := kickClient.CreateChannelReward(context.Background(), gokick.CreateChannelRewardRequest{
			BackgroundColor:                   &backgroundColor,
			Cost:                              100,
			Description:                       &description,
			IsEnabled:                         &isEnabled,
			IsUserInputRequired:               &isUserInputRequired,
			ShouldRedemptionsSkipRequestQueue: &shouldRedemptionsSkipRequestQueue,
			Title:                             "Song Request",
		})
		require.NoError(t, err)
		require.NotNil(t, rewardResponse.Result.BackgroundColor)
		assert.Equal(t, "#00e701", *rewardResponse.Result.BackgroundColor)
		assert.Equal(t, 100, rewardResponse.Result.Cost)
		require.NotNil(t, rewardResponse.Result.Description)
		assert.Equal(t, "Request a song by providing a URL", *rewardResponse.Result.Description)
		assert.Equal(t, "01HZ8X9K2M4N6P8Q0R2S4T6V8W0Y2Z4", rewardResponse.Result.ID)
		require.NotNil(t, rewardResponse.Result.IsEnabled)
		assert.True(t, *rewardResponse.Result.IsEnabled)
		require.NotNil(t, rewardResponse.Result.IsPaused)
		assert.False(t, *rewardResponse.Result.IsPaused)
		require.NotNil(t, rewardResponse.Result.IsUserInputRequired)
		assert.False(t, *rewardResponse.Result.IsUserInputRequired)
		require.NotNil(t, rewardResponse.Result.ShouldRedemptionsSkipRequestQueue)
		assert.False(t, *rewardResponse.Result.ShouldRedemptionsSkipRequestQueue)
		assert.Equal(t, "Song Request", rewardResponse.Result.Title)
	})

	t.Run("with only required fields", func(t *testing.T) {
		kickClient := setupMockClient(t, func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, http.MethodPost, r.Method)
			assert.Equal(t, "/public/v1/channels/rewards", r.URL.Path)
			assert.Equal(t, "application/json", r.Header.Get("Content-Type"))
			assert.Equal(t, "Bearer access-token", r.Header.Get("Authorization"))

			var body map[string]interface{}
			err := json.NewDecoder(r.Body).Decode(&body)
			assert.NoError(t, err)

			assert.InEpsilon(t, float64(100), body["cost"], 0)
			assert.Equal(t, "Song Request", body["title"])
			_, hasBackgroundColor := body["background_color"]
			assert.False(t, hasBackgroundColor, "background_color should not be present")
			_, hasDescription := body["description"]
			assert.False(t, hasDescription, "description should not be present")
			_, hasIsEnabled := body["is_enabled"]
			assert.False(t, hasIsEnabled, "is_enabled should not be present")
			_, hasIsUserInputRequired := body["is_user_input_required"]
			assert.False(t, hasIsUserInputRequired, "is_user_input_required should not be present")
			_, hasShouldRedemptionsSkipRequestQueue := body["should_redemptions_skip_request_queue"]
			assert.False(t, hasShouldRedemptionsSkipRequestQueue, "should_redemptions_skip_request_queue should not be present")

			w.WriteHeader(http.StatusOK)
			fmt.Fprint(w, `{
				"data": {
					"background_color": "#00e701",
					"cost": 100,
					"description": "Request a song by providing a URL",
					"id": "01HZ8X9K2M4N6P8Q0R2S4T6V8W0Y2Z4",
					"is_enabled": true,
					"is_paused": false,
					"is_user_input_required": false,
					"should_redemptions_skip_request_queue": false,
					"title": "Song Request"
				},
				"message": "text"
			}`)
		})

		rewardResponse, err := kickClient.CreateChannelReward(context.Background(), gokick.CreateChannelRewardRequest{
			Cost:  100,
			Title: "Song Request",
		})
		require.NoError(t, err)
		require.NotNil(t, rewardResponse.Result.BackgroundColor)
		assert.Equal(t, "#00e701", *rewardResponse.Result.BackgroundColor)
		assert.Equal(t, 100, rewardResponse.Result.Cost)
		require.NotNil(t, rewardResponse.Result.Description)
		assert.Equal(t, "Request a song by providing a URL", *rewardResponse.Result.Description)
		assert.Equal(t, "01HZ8X9K2M4N6P8Q0R2S4T6V8W0Y2Z4", rewardResponse.Result.ID)
		require.NotNil(t, rewardResponse.Result.IsEnabled)
		assert.True(t, *rewardResponse.Result.IsEnabled)
		require.NotNil(t, rewardResponse.Result.IsPaused)
		assert.False(t, *rewardResponse.Result.IsPaused)
		require.NotNil(t, rewardResponse.Result.IsUserInputRequired)
		assert.False(t, *rewardResponse.Result.IsUserInputRequired)
		require.NotNil(t, rewardResponse.Result.ShouldRedemptionsSkipRequestQueue)
		assert.False(t, *rewardResponse.Result.ShouldRedemptionsSkipRequestQueue)
		assert.Equal(t, "Song Request", rewardResponse.Result.Title)
	})
}

func TestUpdateChannelRewardError(t *testing.T) {
	t.Run("on new request", func(t *testing.T) {
		kickClient, err := gokick.NewClient(&gokick.ClientOptions{UserAccessToken: "access-token"})
		require.NoError(t, err)

		var ctx context.Context
		_, err = kickClient.UpdateChannelReward(ctx, "reward-id", gokick.UpdateChannelRewardRequest{})
		require.EqualError(t, err, "failed to create request: net/http: nil Context")
	})

	t.Run("timeout", func(t *testing.T) {
		kickClient := setupTimeoutMockClient(t)

		_, err := kickClient.UpdateChannelReward(context.Background(), "reward-id", gokick.UpdateChannelRewardRequest{})
		require.EqualError(t, err, `failed to make request: Patch "https://api.kick.com/public/v1/channels/rewards/reward-id": `+
			`context deadline exceeded (Client.Timeout exceeded while awaiting headers)`)
	})

	t.Run("unmarshal error response", func(t *testing.T) {
		kickClient := setupMockClient(t, func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprint(w, `117`)
		})

		_, err := kickClient.UpdateChannelReward(context.Background(), "reward-id", gokick.UpdateChannelRewardRequest{})

		assert.EqualError(t, err, `failed to unmarshal error response (KICK status code: 500 and body "117"): json: cannot unmarshal `+
			`number into Go value of type gokick.errorResponse`)
	})

	t.Run("unmarshal token response", func(t *testing.T) {
		kickClient := setupMockClient(t, func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			fmt.Fprint(w, "117")
		})

		_, err := kickClient.UpdateChannelReward(context.Background(), "reward-id", gokick.UpdateChannelRewardRequest{})

		assert.EqualError(t, err, `failed to unmarshal response body (KICK status code 200 and body "117"): json: cannot unmarshal `+
			`number into Go value of type gokick.successResponse[github.com/scorfly/gokick.ChannelRewardResponse]`)
	})

	t.Run("reader failure", func(t *testing.T) {
		kickClient := setupMockClient(t, func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Length", "10")
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprint(w, "")
		})

		_, err := kickClient.UpdateChannelReward(context.Background(), "reward-id", gokick.UpdateChannelRewardRequest{})

		assert.EqualError(t, err, `failed to read response body (KICK status code 500): unexpected EOF`)
	})

	t.Run("with internal server error", func(t *testing.T) {
		kickClient := setupMockClient(t, func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprint(w, `{"message":"internal server error", "data":null}`)
		})

		_, err := kickClient.UpdateChannelReward(context.Background(), "reward-id", gokick.UpdateChannelRewardRequest{})

		var kickError gokick.Error
		require.ErrorAs(t, err, &kickError)
		assert.Equal(t, http.StatusInternalServerError, kickError.Code())
		assert.Equal(t, "internal server error", kickError.Message())
	})
}

func TestUpdateChannelRewardSuccess(t *testing.T) {
	t.Run("with all fields", func(t *testing.T) {
		backgroundColor := "#00e701"
		cost := 100
		description := "Updated description"
		isEnabled := true
		isPaused := false
		isUserInputRequired := false
		shouldRedemptionsSkipRequestQueue := false
		title := "Updated Title"

		kickClient := setupMockClient(t, func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, http.MethodPatch, r.Method)
			assert.Equal(t, "/public/v1/channels/rewards/reward-id", r.URL.Path)
			assert.Equal(t, "application/json", r.Header.Get("Content-Type"))
			assert.Equal(t, "Bearer access-token", r.Header.Get("Authorization"))

			var body map[string]interface{}
			err := json.NewDecoder(r.Body).Decode(&body)
			assert.NoError(t, err)

			assert.Equal(t, "#00e701", body["background_color"])
			assert.InEpsilon(t, float64(100), body["cost"], 0)
			assert.Equal(t, "Updated description", body["description"])
			assert.Equal(t, true, body["is_enabled"])
			assert.Equal(t, false, body["is_paused"])
			assert.Equal(t, false, body["is_user_input_required"])
			assert.Equal(t, false, body["should_redemptions_skip_request_queue"])
			assert.Equal(t, "Updated Title", body["title"])

			w.WriteHeader(http.StatusOK)
			fmt.Fprint(w, `{
				"data": {
					"background_color": "#00e701",
					"cost": 100,
					"description": "Request a song by providing a URL",
					"id": "01HZ8X9K2M4N6P8Q0R2S4T6V8W0Y2Z4",
					"is_enabled": true,
					"is_paused": false,
					"is_user_input_required": false,
					"should_redemptions_skip_request_queue": false,
					"title": "Song Request"
				},
				"message": "text"
			}`)
		})

		rewardResponse, err := kickClient.UpdateChannelReward(context.Background(), "reward-id", gokick.UpdateChannelRewardRequest{
			BackgroundColor:                   &backgroundColor,
			Cost:                              &cost,
			Description:                       &description,
			IsEnabled:                         &isEnabled,
			IsPaused:                          &isPaused,
			IsUserInputRequired:               &isUserInputRequired,
			ShouldRedemptionsSkipRequestQueue: &shouldRedemptionsSkipRequestQueue,
			Title:                             &title,
		})
		require.NoError(t, err)
		require.NotNil(t, rewardResponse.Result.BackgroundColor)
		assert.Equal(t, "#00e701", *rewardResponse.Result.BackgroundColor)
		assert.Equal(t, 100, rewardResponse.Result.Cost)
		require.NotNil(t, rewardResponse.Result.Description)
		assert.Equal(t, "Request a song by providing a URL", *rewardResponse.Result.Description)
		assert.Equal(t, "01HZ8X9K2M4N6P8Q0R2S4T6V8W0Y2Z4", rewardResponse.Result.ID)
		require.NotNil(t, rewardResponse.Result.IsEnabled)
		assert.True(t, *rewardResponse.Result.IsEnabled)
		require.NotNil(t, rewardResponse.Result.IsPaused)
		assert.False(t, *rewardResponse.Result.IsPaused)
		require.NotNil(t, rewardResponse.Result.IsUserInputRequired)
		assert.False(t, *rewardResponse.Result.IsUserInputRequired)
		require.NotNil(t, rewardResponse.Result.ShouldRedemptionsSkipRequestQueue)
		assert.False(t, *rewardResponse.Result.ShouldRedemptionsSkipRequestQueue)
		assert.Equal(t, "Song Request", rewardResponse.Result.Title)
	})

	t.Run("with partial fields", func(t *testing.T) {
		cost := 200
		title := "New Title"

		kickClient := setupMockClient(t, func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, http.MethodPatch, r.Method)
			assert.Equal(t, "/public/v1/channels/rewards/reward-id", r.URL.Path)
			assert.Equal(t, "application/json", r.Header.Get("Content-Type"))
			assert.Equal(t, "Bearer access-token", r.Header.Get("Authorization"))

			var body map[string]interface{}
			err := json.NewDecoder(r.Body).Decode(&body)
			assert.NoError(t, err)

			assert.InEpsilon(t, float64(200), body["cost"], 0)
			assert.Equal(t, "New Title", body["title"])
			_, hasBackgroundColor := body["background_color"]
			assert.False(t, hasBackgroundColor, "background_color should not be present")
			_, hasDescription := body["description"]
			assert.False(t, hasDescription, "description should not be present")
			_, hasIsEnabled := body["is_enabled"]
			assert.False(t, hasIsEnabled, "is_enabled should not be present")
			_, hasIsPaused := body["is_paused"]
			assert.False(t, hasIsPaused, "is_paused should not be present")
			_, hasIsUserInputRequired := body["is_user_input_required"]
			assert.False(t, hasIsUserInputRequired, "is_user_input_required should not be present")
			_, hasShouldRedemptionsSkipRequestQueue := body["should_redemptions_skip_request_queue"]
			assert.False(t, hasShouldRedemptionsSkipRequestQueue, "should_redemptions_skip_request_queue should not be present")

			w.WriteHeader(http.StatusOK)
			fmt.Fprint(w, `{
				"data": {
					"background_color": "#00e701",
					"cost": 200,
					"description": "Request a song by providing a URL",
					"id": "01HZ8X9K2M4N6P8Q0R2S4T6V8W0Y2Z4",
					"is_enabled": true,
					"is_paused": false,
					"is_user_input_required": false,
					"should_redemptions_skip_request_queue": false,
					"title": "New Title"
				},
				"message": "text"
			}`)
		})

		rewardResponse, err := kickClient.UpdateChannelReward(context.Background(), "reward-id", gokick.UpdateChannelRewardRequest{
			Cost:  &cost,
			Title: &title,
		})
		require.NoError(t, err)
		assert.Equal(t, 200, rewardResponse.Result.Cost)
		assert.Equal(t, "New Title", rewardResponse.Result.Title)
		assert.Equal(t, "01HZ8X9K2M4N6P8Q0R2S4T6V8W0Y2Z4", rewardResponse.Result.ID)
	})

	t.Run("with empty request", func(t *testing.T) {
		kickClient := setupMockClient(t, func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, http.MethodPatch, r.Method)
			assert.Equal(t, "/public/v1/channels/rewards/reward-id", r.URL.Path)
			assert.Equal(t, "application/json", r.Header.Get("Content-Type"))
			assert.Equal(t, "Bearer access-token", r.Header.Get("Authorization"))

			var body map[string]interface{}
			err := json.NewDecoder(r.Body).Decode(&body)
			assert.NoError(t, err)

			assert.Empty(t, body)

			w.WriteHeader(http.StatusOK)
			fmt.Fprint(w, `{
				"data": {
					"background_color": "#00e701",
					"cost": 100,
					"description": "Request a song by providing a URL",
					"id": "01HZ8X9K2M4N6P8Q0R2S4T6V8W0Y2Z4",
					"is_enabled": true,
					"is_paused": false,
					"is_user_input_required": false,
					"should_redemptions_skip_request_queue": false,
					"title": "Song Request"
				},
				"message": "text"
			}`)
		})

		rewardResponse, err := kickClient.UpdateChannelReward(context.Background(), "reward-id", gokick.UpdateChannelRewardRequest{})
		require.NoError(t, err)
		assert.Equal(t, "01HZ8X9K2M4N6P8Q0R2S4T6V8W0Y2Z4", rewardResponse.Result.ID)
	})
}

func TestDeleteChannelRewardError(t *testing.T) {
	t.Run("on new request", func(t *testing.T) {
		kickClient, err := gokick.NewClient(&gokick.ClientOptions{UserAccessToken: "access-token"})
		require.NoError(t, err)

		var ctx context.Context
		_, err = kickClient.DeleteChannelReward(ctx, "reward-id")
		require.EqualError(t, err, "failed to create request: net/http: nil Context")
	})

	t.Run("timeout", func(t *testing.T) {
		kickClient := setupTimeoutMockClient(t)

		_, err := kickClient.DeleteChannelReward(context.Background(), "reward-id")
		require.EqualError(t, err, `failed to make request: Delete "https://api.kick.com/public/v1/channels/rewards/reward-id": `+
			`context deadline exceeded (Client.Timeout exceeded while awaiting headers)`)
	})

	t.Run("unmarshal error response", func(t *testing.T) {
		kickClient := setupMockClient(t, func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprint(w, `117`)
		})

		_, err := kickClient.DeleteChannelReward(context.Background(), "reward-id")

		assert.EqualError(t, err, `failed to unmarshal error response (KICK status code: 500 and body "117"): json: cannot unmarshal `+
			`number into Go value of type gokick.errorResponse`)
	})

	t.Run("unmarshal token response", func(t *testing.T) {
		kickClient := setupMockClient(t, func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			fmt.Fprint(w, "117")
		})

		_, err := kickClient.DeleteChannelReward(context.Background(), "reward-id")

		assert.EqualError(t, err, `failed to unmarshal error response (KICK status code: 200 and body "117"): json: cannot unmarshal `+
			`number into Go value of type gokick.errorResponse`)
	})

	t.Run("reader failure", func(t *testing.T) {
		kickClient := setupMockClient(t, func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Length", "10")
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprint(w, "")
		})

		_, err := kickClient.DeleteChannelReward(context.Background(), "reward-id")

		assert.EqualError(t, err, `failed to read response body (KICK status code 500): unexpected EOF`)
	})

	t.Run("with internal server error", func(t *testing.T) {
		kickClient := setupMockClient(t, func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprint(w, `{"message":"internal server error", "data":null}`)
		})

		_, err := kickClient.DeleteChannelReward(context.Background(), "reward-id")

		var kickError gokick.Error
		require.ErrorAs(t, err, &kickError)
		assert.Equal(t, http.StatusInternalServerError, kickError.Code())
		assert.Equal(t, "internal server error", kickError.Message())
	})
}

func TestDeleteChannelRewardSuccess(t *testing.T) {
	t.Run("successful deletion", func(t *testing.T) {
		kickClient := setupMockClient(t, func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, http.MethodDelete, r.Method)
			assert.Equal(t, "/public/v1/channels/rewards/reward-id", r.URL.Path)
			assert.Equal(t, "application/json", r.Header.Get("Content-Type"))
			assert.Equal(t, "Bearer access-token", r.Header.Get("Authorization"))

			w.WriteHeader(http.StatusNoContent)
		})

		_, err := kickClient.DeleteChannelReward(context.Background(), "reward-id")
		require.NoError(t, err)
	})
}

func TestNewChannelRewardRedemptionListFilterSuccess(t *testing.T) {
	testCases := map[string]struct {
		filter              gokick.ChannelRewardRedemptionListFilter
		expectedQueryString string
	}{
		"default": {
			filter:              gokick.NewChannelRewardRedemptionListFilter(),
			expectedQueryString: "",
		},
		"with reward_id": {
			filter:              gokick.NewChannelRewardRedemptionListFilter().SetRewardID("01HZREWARD"),
			expectedQueryString: "?reward_id=01HZREWARD",
		},
		"with status": {
			filter: gokick.NewChannelRewardRedemptionListFilter().
				SetStatus(gokick.ChannelRewardRedemptionStatusAccepted),
			expectedQueryString: "?status=accepted",
		},
		"with cursor": {
			filter:              gokick.NewChannelRewardRedemptionListFilter().SetCursor("abc123"),
			expectedQueryString: "?cursor=abc123",
		},
		"with single id": {
			filter:              gokick.NewChannelRewardRedemptionListFilter().AddID("01HZRED1"),
			expectedQueryString: "?id=01HZRED1",
		},
		"with multiple ids": {
			filter: gokick.NewChannelRewardRedemptionListFilter().
				AddID("01HZRED1").
				AddID("01HZRED2"),
			expectedQueryString: "?id=01HZRED1&id=01HZRED2",
		},
		"with reward_id status and cursor": {
			filter: gokick.NewChannelRewardRedemptionListFilter().
				SetRewardID("01HZREWARD").
				SetStatus(gokick.ChannelRewardRedemptionStatusPending).
				SetCursor("next"),
			expectedQueryString: "?cursor=next&reward_id=01HZREWARD&status=pending",
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			assert.Equal(t, tc.expectedQueryString, tc.filter.ToQueryString())
		})
	}
}

func TestGetChannelRewardRedemptionsError(t *testing.T) {
	t.Run("on new request", func(t *testing.T) {
		kickClient, err := gokick.NewClient(&gokick.ClientOptions{UserAccessToken: "access-token"})
		require.NoError(t, err)

		var ctx context.Context
		_, err = kickClient.GetChannelRewardRedemptions(ctx, gokick.NewChannelRewardRedemptionListFilter())
		require.EqualError(t, err, "failed to create request: net/http: nil Context")
	})

	t.Run("timeout", func(t *testing.T) {
		kickClient := setupTimeoutMockClient(t)

		_, err := kickClient.GetChannelRewardRedemptions(context.Background(), gokick.NewChannelRewardRedemptionListFilter())
		require.EqualError(t, err, `failed to make request: Get "https://api.kick.com/public/v1/channels/rewards/redemptions": `+
			`context deadline exceeded (Client.Timeout exceeded while awaiting headers)`)
	})

	t.Run("unmarshal error response", func(t *testing.T) {
		kickClient := setupMockClient(t, func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprint(w, `117`)
		})

		_, err := kickClient.GetChannelRewardRedemptions(context.Background(), gokick.NewChannelRewardRedemptionListFilter())

		assert.EqualError(t, err, `failed to unmarshal error response (KICK status code: 500 and body "117"): json: cannot unmarshal `+
			`number into Go value of type gokick.errorResponse`)
	})

	t.Run("unmarshal redemptions response", func(t *testing.T) {
		kickClient := setupMockClient(t, func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			fmt.Fprint(w, "117")
		})

		_, err := kickClient.GetChannelRewardRedemptions(context.Background(), gokick.NewChannelRewardRedemptionListFilter())

		assert.EqualError(t, err, `failed to unmarshal response body (KICK status code 200 and body "117"): json: cannot unmarshal `+
			`number into Go value of type gokick.PaginatedResponse[[]github.com/scorfly/gokick.RedemptionsByReward]`)
	})

	t.Run("reader failure", func(t *testing.T) {
		kickClient := setupMockClient(t, func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Length", "10")
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprint(w, "")
		})

		_, err := kickClient.GetChannelRewardRedemptions(context.Background(), gokick.NewChannelRewardRedemptionListFilter())

		assert.EqualError(t, err, `failed to read response body (KICK status code 500): unexpected EOF`)
	})

	t.Run("with internal server error", func(t *testing.T) {
		kickClient := setupMockClient(t, func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprint(w, `{"message":"internal server error", "data":null}`)
		})

		_, err := kickClient.GetChannelRewardRedemptions(context.Background(), gokick.NewChannelRewardRedemptionListFilter())

		var kickError gokick.Error
		require.ErrorAs(t, err, &kickError)
		assert.Equal(t, http.StatusInternalServerError, kickError.Code())
		assert.Equal(t, "internal server error", kickError.Message())
	})
}

func TestGetChannelRewardRedemptionsSuccess(t *testing.T) {
	t.Run("when result is empty", func(t *testing.T) {
		kickClient := setupMockClient(t, func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			fmt.Fprint(w, `{"data":[], "message":"text", "pagination":{"next_cursor":""}}`)
		})

		response, err := kickClient.GetChannelRewardRedemptions(context.Background(), gokick.NewChannelRewardRedemptionListFilter())
		require.NoError(t, err)
		assert.Empty(t, response.Result)
		assert.Empty(t, response.Pagination.NextCursor)
	})

	t.Run("when result is filled", func(t *testing.T) {
		kickClient := setupMockClient(t, func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, http.MethodGet, r.Method)
			assert.Equal(t, "/public/v1/channels/rewards/redemptions", r.URL.Path)
			assert.Equal(t, "01HZREWARD", r.URL.Query().Get("reward_id"))
			assert.Equal(t, "pending", r.URL.Query().Get("status"))
			assert.Equal(t, "next-page", r.URL.Query().Get("cursor"))
			assert.Equal(t, "Bearer access-token", r.Header.Get("Authorization"))

			w.WriteHeader(http.StatusOK)
			fmt.Fprint(w, `{
				"data": [{
					"reward": {
						"id": "01HZREWARD",
						"title": "Song Request",
						"cost": 100,
						"description": "Request a song",
						"can_manage": true,
						"is_deleted": false
					},
					"redemptions": [{
						"id": "01HZRED1",
						"redeemed_at": "2025-02-21T23:23:36Z",
						"redeemer": {"user_id": 117},
						"status": "pending",
						"user_input": "https://example.com/song"
					}]
				}],
				"message": "text",
				"pagination": {"next_cursor": "cursor-2"}
			}`)
		})

		response, err := kickClient.GetChannelRewardRedemptions(
			context.Background(),
			gokick.NewChannelRewardRedemptionListFilter().
				SetRewardID("01HZREWARD").
				SetStatus(gokick.ChannelRewardRedemptionStatusPending).
				SetCursor("next-page"),
		)
		require.NoError(t, err)
		require.Len(t, response.Result, 1)
		assert.Equal(t, "cursor-2", response.Pagination.NextCursor)

		reward := response.Result[0].Reward
		assert.Equal(t, "01HZREWARD", reward.ID)
		assert.Equal(t, "Song Request", reward.Title)
		require.NotNil(t, reward.Cost)
		assert.Equal(t, 100, *reward.Cost)
		require.NotNil(t, reward.Description)
		assert.Equal(t, "Request a song", *reward.Description)
		require.NotNil(t, reward.CanManage)
		assert.True(t, *reward.CanManage)
		require.NotNil(t, reward.IsDeleted)
		assert.False(t, *reward.IsDeleted)

		require.Len(t, response.Result[0].Redemptions, 1)
		redemption := response.Result[0].Redemptions[0]
		assert.Equal(t, "01HZRED1", redemption.ID)
		assert.Equal(t, "2025-02-21T23:23:36Z", redemption.RedeemedAt)
		assert.Equal(t, 117, redemption.Redeemer.UserID)
		assert.Equal(t, "pending", redemption.Status)
		assert.Equal(t, "https://example.com/song", redemption.UserInput)
	})

	t.Run("when filtering by redemption ids", func(t *testing.T) {
		kickClient := setupMockClient(t, func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, []string{"01HZRED1", "01HZRED2"}, r.URL.Query()["id"])

			w.WriteHeader(http.StatusOK)
			fmt.Fprint(w, `{
				"data": [{
					"reward": {"id": "01HZREWARD", "title": "Song Request"},
					"redemptions": []
				}],
				"message": "text",
				"pagination": {"next_cursor": ""}
			}`)
		})

		response, err := kickClient.GetChannelRewardRedemptions(
			context.Background(),
			gokick.NewChannelRewardRedemptionListFilter().AddID("01HZRED1").AddID("01HZRED2"),
		)
		require.NoError(t, err)
		require.Len(t, response.Result, 1)
		assert.Equal(t, "01HZREWARD", response.Result[0].Reward.ID)
		assert.Nil(t, response.Result[0].Reward.CanManage)
		assert.Nil(t, response.Result[0].Reward.Cost)
		assert.Nil(t, response.Result[0].Reward.Description)
		assert.Nil(t, response.Result[0].Reward.IsDeleted)
		assert.Empty(t, response.Result[0].Redemptions)
	})
}

func TestAcceptChannelRewardRedemptionsError(t *testing.T) {
	t.Run("on new request", func(t *testing.T) {
		kickClient, err := gokick.NewClient(&gokick.ClientOptions{UserAccessToken: "access-token"})
		require.NoError(t, err)

		var ctx context.Context
		_, err = kickClient.AcceptChannelRewardRedemptions(ctx, []string{"01HZRED1"})
		require.EqualError(t, err, "failed to create request: net/http: nil Context")
	})

	t.Run("timeout", func(t *testing.T) {
		kickClient := setupTimeoutMockClient(t)

		_, err := kickClient.AcceptChannelRewardRedemptions(context.Background(), []string{"01HZRED1"})
		require.EqualError(t, err, `failed to make request: Post "https://api.kick.com/public/v1/channels/rewards/redemptions/accept": `+
			`context deadline exceeded (Client.Timeout exceeded while awaiting headers)`)
	})

	t.Run("unmarshal error response", func(t *testing.T) {
		kickClient := setupMockClient(t, func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprint(w, `117`)
		})

		_, err := kickClient.AcceptChannelRewardRedemptions(context.Background(), []string{"01HZRED1"})

		assert.EqualError(t, err, `failed to unmarshal error response (KICK status code: 500 and body "117"): json: cannot unmarshal `+
			`number into Go value of type gokick.errorResponse`)
	})

	t.Run("unmarshal response", func(t *testing.T) {
		kickClient := setupMockClient(t, func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			fmt.Fprint(w, "117")
		})

		_, err := kickClient.AcceptChannelRewardRedemptions(context.Background(), []string{"01HZRED1"})

		assert.EqualError(t, err, `failed to unmarshal response body (KICK status code 200 and body "117"): json: cannot unmarshal `+
			`number into Go value of type gokick.successResponse[[]github.com/scorfly/gokick.FailedRedemption]`)
	})

	t.Run("reader failure", func(t *testing.T) {
		kickClient := setupMockClient(t, func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Length", "10")
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprint(w, "")
		})

		_, err := kickClient.AcceptChannelRewardRedemptions(context.Background(), []string{"01HZRED1"})

		assert.EqualError(t, err, `failed to read response body (KICK status code 500): unexpected EOF`)
	})

	t.Run("with internal server error", func(t *testing.T) {
		kickClient := setupMockClient(t, func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprint(w, `{"message":"internal server error", "data":null}`)
		})

		_, err := kickClient.AcceptChannelRewardRedemptions(context.Background(), []string{"01HZRED1"})

		var kickError gokick.Error
		require.ErrorAs(t, err, &kickError)
		assert.Equal(t, http.StatusInternalServerError, kickError.Code())
		assert.Equal(t, "internal server error", kickError.Message())
	})
}

func TestAcceptChannelRewardRedemptionsSuccess(t *testing.T) {
	t.Run("when all succeed", func(t *testing.T) {
		kickClient := setupMockClient(t, func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, http.MethodPost, r.Method)
			assert.Equal(t, "/public/v1/channels/rewards/redemptions/accept", r.URL.Path)
			assert.Equal(t, "application/json", r.Header.Get("Content-Type"))
			assert.Equal(t, "Bearer access-token", r.Header.Get("Authorization"))

			var body map[string]interface{}
			err := json.NewDecoder(r.Body).Decode(&body)
			assert.NoError(t, err)
			assert.Equal(t, []interface{}{"01HZRED1", "01HZRED2"}, body["ids"])

			w.WriteHeader(http.StatusOK)
			fmt.Fprint(w, `{"data":[], "message":"text"}`)
		})

		response, err := kickClient.AcceptChannelRewardRedemptions(context.Background(), []string{"01HZRED1", "01HZRED2"})
		require.NoError(t, err)
		assert.Empty(t, response.Result)
	})

	t.Run("when some fail", func(t *testing.T) {
		kickClient := setupMockClient(t, func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			fmt.Fprint(w, `{
				"data": [
					{"id": "01HZRED1", "reason": "NOT_PENDING"},
					{"id": "01HZRED2", "reason": "NOT_FOUND"}
				],
				"message": "text"
			}`)
		})

		response, err := kickClient.AcceptChannelRewardRedemptions(context.Background(), []string{"01HZRED1", "01HZRED2"})
		require.NoError(t, err)
		require.Len(t, response.Result, 2)
		assert.Equal(t, "01HZRED1", response.Result[0].ID)
		assert.Equal(t, "NOT_PENDING", response.Result[0].Reason)
		assert.Equal(t, "01HZRED2", response.Result[1].ID)
		assert.Equal(t, "NOT_FOUND", response.Result[1].Reason)
	})
}

func TestRejectChannelRewardRedemptionsError(t *testing.T) {
	t.Run("on new request", func(t *testing.T) {
		kickClient, err := gokick.NewClient(&gokick.ClientOptions{UserAccessToken: "access-token"})
		require.NoError(t, err)

		var ctx context.Context
		_, err = kickClient.RejectChannelRewardRedemptions(ctx, []string{"01HZRED1"})
		require.EqualError(t, err, "failed to create request: net/http: nil Context")
	})

	t.Run("timeout", func(t *testing.T) {
		kickClient := setupTimeoutMockClient(t)

		_, err := kickClient.RejectChannelRewardRedemptions(context.Background(), []string{"01HZRED1"})
		require.EqualError(t, err, `failed to make request: Post "https://api.kick.com/public/v1/channels/rewards/redemptions/reject": `+
			`context deadline exceeded (Client.Timeout exceeded while awaiting headers)`)
	})

	t.Run("unmarshal error response", func(t *testing.T) {
		kickClient := setupMockClient(t, func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprint(w, `117`)
		})

		_, err := kickClient.RejectChannelRewardRedemptions(context.Background(), []string{"01HZRED1"})

		assert.EqualError(t, err, `failed to unmarshal error response (KICK status code: 500 and body "117"): json: cannot unmarshal `+
			`number into Go value of type gokick.errorResponse`)
	})

	t.Run("unmarshal response", func(t *testing.T) {
		kickClient := setupMockClient(t, func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			fmt.Fprint(w, "117")
		})

		_, err := kickClient.RejectChannelRewardRedemptions(context.Background(), []string{"01HZRED1"})

		assert.EqualError(t, err, `failed to unmarshal response body (KICK status code 200 and body "117"): json: cannot unmarshal `+
			`number into Go value of type gokick.successResponse[[]github.com/scorfly/gokick.FailedRedemption]`)
	})

	t.Run("reader failure", func(t *testing.T) {
		kickClient := setupMockClient(t, func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Length", "10")
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprint(w, "")
		})

		_, err := kickClient.RejectChannelRewardRedemptions(context.Background(), []string{"01HZRED1"})

		assert.EqualError(t, err, `failed to read response body (KICK status code 500): unexpected EOF`)
	})

	t.Run("with internal server error", func(t *testing.T) {
		kickClient := setupMockClient(t, func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprint(w, `{"message":"internal server error", "data":null}`)
		})

		_, err := kickClient.RejectChannelRewardRedemptions(context.Background(), []string{"01HZRED1"})

		var kickError gokick.Error
		require.ErrorAs(t, err, &kickError)
		assert.Equal(t, http.StatusInternalServerError, kickError.Code())
		assert.Equal(t, "internal server error", kickError.Message())
	})
}

func TestRejectChannelRewardRedemptionsSuccess(t *testing.T) {
	t.Run("when all succeed", func(t *testing.T) {
		kickClient := setupMockClient(t, func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, http.MethodPost, r.Method)
			assert.Equal(t, "/public/v1/channels/rewards/redemptions/reject", r.URL.Path)
			assert.Equal(t, "application/json", r.Header.Get("Content-Type"))
			assert.Equal(t, "Bearer access-token", r.Header.Get("Authorization"))

			var body map[string]interface{}
			err := json.NewDecoder(r.Body).Decode(&body)
			assert.NoError(t, err)
			assert.Equal(t, []interface{}{"01HZRED1"}, body["ids"])

			w.WriteHeader(http.StatusOK)
			fmt.Fprint(w, `{"data":[], "message":"text"}`)
		})

		response, err := kickClient.RejectChannelRewardRedemptions(context.Background(), []string{"01HZRED1"})
		require.NoError(t, err)
		assert.Empty(t, response.Result)
	})

	t.Run("when some fail", func(t *testing.T) {
		kickClient := setupMockClient(t, func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			fmt.Fprint(w, `{
				"data": [
					{"id": "01HZRED1", "reason": "NOT_OWNED"},
					{"id": "01HZRED2", "reason": "UNKNOWN"}
				],
				"message": "text"
			}`)
		})

		response, err := kickClient.RejectChannelRewardRedemptions(context.Background(), []string{"01HZRED1", "01HZRED2"})
		require.NoError(t, err)
		require.Len(t, response.Result, 2)
		assert.Equal(t, "01HZRED1", response.Result[0].ID)
		assert.Equal(t, "NOT_OWNED", response.Result[0].Reason)
		assert.Equal(t, "01HZRED2", response.Result[1].ID)
		assert.Equal(t, "UNKNOWN", response.Result[1].Reason)
	})
}

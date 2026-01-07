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
				"banner_picture": "banner picture",
				"broadcaster_user_id": 117,
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
		assert.Equal(t, "banner picture", channelsResponse.Result[0].BannerPicture)
		assert.Equal(t, 117, channelsResponse.Result[0].BroadcasterUserID)
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

func TestUpdateCustomTagsError(t *testing.T) {
	t.Run("on new request", func(t *testing.T) {
		kickClient, err := gokick.NewClient(&gokick.ClientOptions{UserAccessToken: "access-token"})
		require.NoError(t, err)

		var ctx context.Context
		_, err = kickClient.UpdateCustomTags(ctx, []string{"pog", "chad"})
		require.EqualError(t, err, "failed to create request: net/http: nil Context")
	})

	t.Run("timeout", func(t *testing.T) {
		kickClient := setupTimeoutMockClient(t)

		_, err := kickClient.UpdateCustomTags(context.Background(), []string{"pog", "chad"})
		require.EqualError(t, err, `failed to make request: Patch "https://api.kick.com/public/v1/channels": context deadline exceeded `+
			`(Client.Timeout exceeded while awaiting headers)`)
	})

	t.Run("unmarshal error response", func(t *testing.T) {
		kickClient := setupMockClient(t, func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprint(w, `117`)
		})

		_, err := kickClient.UpdateCustomTags(context.Background(), []string{"pog", "chad"})

		assert.EqualError(t, err, `failed to unmarshal error response (KICK status code: 500 and body "117"): json: cannot unmarshal `+
			`number into Go value of type gokick.errorResponse`)
	})

	t.Run("unmarshal token response", func(t *testing.T) {
		kickClient := setupMockClient(t, func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			fmt.Fprint(w, "117")
		})

		_, err := kickClient.UpdateCustomTags(context.Background(), []string{"pog", "chad"})

		assert.EqualError(t, err, `failed to unmarshal error response (KICK status code: 200 and body "117"): json: cannot unmarshal `+
			`number into Go value of type gokick.errorResponse`)
	})

	t.Run("reader failure", func(t *testing.T) {
		kickClient := setupMockClient(t, func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Length", "10")
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprint(w, "")
		})

		_, err := kickClient.UpdateCustomTags(context.Background(), []string{"pog", "chad"})

		assert.EqualError(t, err, `failed to read response body (KICK status code 500): unexpected EOF`)
	})

	t.Run("with internal server error", func(t *testing.T) {
		kickClient := setupMockClient(t, func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprint(w, `{"message":"internal server error", "data":null}`)
		})

		_, err := kickClient.UpdateCustomTags(context.Background(), []string{"pog", "chad"})

		var kickError gokick.Error
		require.ErrorAs(t, err, &kickError)
		assert.Equal(t, http.StatusInternalServerError, kickError.Code())
		assert.Equal(t, "internal server error", kickError.Message())
	})
}

func TestUpdateCustomTagsSuccess(t *testing.T) {
	t.Run("with tags", func(t *testing.T) {
		kickClient := setupMockClient(t, func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, http.MethodPatch, r.Method)
			assert.Equal(t, "/public/v1/channels", r.URL.Path)
			assert.Equal(t, "application/json", r.Header.Get("Content-Type"))
			assert.Equal(t, "Bearer access-token", r.Header.Get("Authorization"))

			var body map[string]interface{}
			err := json.NewDecoder(r.Body).Decode(&body)
			assert.NoError(t, err)

			customTags, ok := body["custom_tags"].([]interface{})
			assert.True(t, ok, "custom_tags should be present in request body")
			assert.Equal(t, []interface{}{"pog", "chad"}, customTags)

			w.WriteHeader(http.StatusNoContent)
		})

		_, err := kickClient.UpdateCustomTags(context.Background(), []string{"pog", "chad"})
		require.NoError(t, err)
	})

	t.Run("with empty tags", func(t *testing.T) {
		kickClient := setupMockClient(t, func(w http.ResponseWriter, r *http.Request) {
			var body map[string]interface{}
			err := json.NewDecoder(r.Body).Decode(&body)
			assert.NoError(t, err)

			customTags, ok := body["custom_tags"].([]interface{})
			assert.True(t, ok, "custom_tags should be present in request body")
			assert.Empty(t, customTags)

			w.WriteHeader(http.StatusNoContent)
		})

		_, err := kickClient.UpdateCustomTags(context.Background(), []string{})
		require.NoError(t, err)
	})
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

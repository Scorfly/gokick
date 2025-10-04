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

func TestNewLivestreamListFilterSuccess(t *testing.T) {
	testCases := map[string]struct {
		filter              gokick.LivestreamListFilter
		expectedQueryString string
	}{
		"default": {
			filter:              gokick.NewLivestreamListFilter(),
			expectedQueryString: "",
		},
		"with broadcaster user ID": {
			filter:              gokick.NewLivestreamListFilter().SetBroadcasterUserIDs(118),
			expectedQueryString: "?broadcaster_user_id=118",
		},
		"with category ID": {
			filter:              gokick.NewLivestreamListFilter().SetCategoryID(218),
			expectedQueryString: "?category_id=218",
		},
		"with language": {
			filter:              gokick.NewLivestreamListFilter().SetLanguage("fr"),
			expectedQueryString: "?language=fr",
		},
		"with limit": {
			filter:              gokick.NewLivestreamListFilter().SetLimit(117),
			expectedQueryString: "?limit=117",
		},
		"with sort by viewer count": {
			filter:              gokick.NewLivestreamListFilter().SetSort(gokick.LivestreamSortViewerCount),
			expectedQueryString: "?sort=viewer_count",
		},
		"with sort by started at": {
			filter:              gokick.NewLivestreamListFilter().SetSort(gokick.LivestreamSortStartedAt),
			expectedQueryString: "?sort=started_at",
		},
		"with all params": {
			filter: gokick.NewLivestreamListFilter().
				SetBroadcasterUserIDs(118).
				SetCategoryID(218).
				SetLanguage("fr").
				SetLimit(117).
				SetSort(gokick.LivestreamSortStartedAt),
			expectedQueryString: "?broadcaster_user_id=118&category_id=218&language=fr&limit=117&sort=started_at",
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			assert.Equal(t, tc.expectedQueryString, tc.filter.ToQueryString())
		})
	}
}

func TestGetLivestreamsError(t *testing.T) {
	t.Run("on new request", func(t *testing.T) {
		kickClient, err := gokick.NewClient(&gokick.ClientOptions{UserAccessToken: "access-token"})
		require.NoError(t, err)

		var ctx context.Context
		_, err = kickClient.GetLivestreams(ctx, gokick.NewLivestreamListFilter())
		require.EqualError(t, err, "failed to create request: net/http: nil Context")
	})

	t.Run("timeout", func(t *testing.T) {
		kickClient := setupTimeoutMockClient(t)

		_, err := kickClient.GetLivestreams(context.Background(), gokick.NewLivestreamListFilter())
		require.EqualError(t, err, `failed to make request: Get "https://api.kick.com/public/v1/livestreams": context deadline exceeded `+
			`(Client.Timeout exceeded while awaiting headers)`)
	})

	t.Run("unmarshal error response", func(t *testing.T) {
		kickClient := setupMockClient(t, func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprint(w, `117`)
		})

		_, err := kickClient.GetLivestreams(context.Background(), gokick.NewLivestreamListFilter())

		assert.EqualError(t, err, `failed to unmarshal error response (KICK status code: 500 and body "117"): json: cannot unmarshal `+
			`number into Go value of type gokick.errorResponse`)
	})

	t.Run("unmarshal Livestreams response", func(t *testing.T) {
		kickClient := setupMockClient(t, func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			fmt.Fprint(w, "117")
		})

		_, err := kickClient.GetLivestreams(context.Background(), gokick.NewLivestreamListFilter())

		assert.EqualError(t, err, `failed to unmarshal response body (KICK status code 200 and body "117"): json: cannot unmarshal `+
			`number into Go value of type gokick.successResponse[[]github.com/scorfly/gokick.LivestreamResponse]`)
	})

	t.Run("reader failure", func(t *testing.T) {
		kickClient := setupMockClient(t, func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Length", "10")
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprint(w, "")
		})

		_, err := kickClient.GetLivestreams(context.Background(), gokick.NewLivestreamListFilter())

		assert.EqualError(t, err, `failed to read response body (KICK status code 500): unexpected EOF`)
	})

	t.Run("with internal server error", func(t *testing.T) {
		kickClient := setupMockClient(t, func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprint(w, `{"message":"internal server error", "data":null}`)
		})

		_, err := kickClient.GetLivestreams(context.Background(), gokick.NewLivestreamListFilter())

		var kickError gokick.Error
		require.ErrorAs(t, err, &kickError)
		assert.Equal(t, http.StatusInternalServerError, kickError.Code())
		assert.Equal(t, "internal server error", kickError.Message())
	})
}

func TestGetLivestreamsSuccess(t *testing.T) {
	t.Run("without result", func(t *testing.T) {
		kickClient := setupMockClient(t, func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			fmt.Fprint(w, `{"message":"success", "data":[]}`)
		})

		LivestreamsResponse, err := kickClient.GetLivestreams(context.Background(), gokick.NewLivestreamListFilter())
		require.NoError(t, err)
		assert.Empty(t, LivestreamsResponse.Result)
	})

	t.Run("with result", func(t *testing.T) {
		kickClient := setupMockClient(t, func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			fmt.Fprint(w, `{"message":"success", "data":[{
				"broadcaster_user_id": 219,
				"category": {
					"id": 123,
					"thumbnail": "category image url",
					"name": "category name"
				},
				"channel_id": 198,
				"has_mature_content": true,
				"language": "fr",
				"slug": "slug",
				"started_at": "started_at",
				"stream_title": "stream_title",
				"thumbnail": "thumbnail_url",
				"viewer_count": 167
				}
  			]}`)
		})

		LivestreamsResponse, err := kickClient.GetLivestreams(context.Background(), gokick.NewLivestreamListFilter())
		require.NoError(t, err)
		require.Len(t, LivestreamsResponse.Result, 1)
		assert.Equal(t, 219, LivestreamsResponse.Result[0].BroadcasterUserID)
		assert.Equal(t, 123, LivestreamsResponse.Result[0].Category.ID)
		assert.Equal(t, "category name", LivestreamsResponse.Result[0].Category.Name)
		assert.Equal(t, "category image url", LivestreamsResponse.Result[0].Category.Thumbnail)
		assert.Equal(t, 198, LivestreamsResponse.Result[0].ChannelID)
		assert.True(t, LivestreamsResponse.Result[0].HasMatureContent)
		assert.Equal(t, "fr", LivestreamsResponse.Result[0].Language)
		assert.Equal(t, "slug", LivestreamsResponse.Result[0].Slug)
		assert.Equal(t, "started_at", LivestreamsResponse.Result[0].StartedAt)
		assert.Equal(t, "stream_title", LivestreamsResponse.Result[0].StreamTitle)
		assert.Equal(t, "thumbnail_url", LivestreamsResponse.Result[0].Thumbnail)
		assert.Equal(t, 167, LivestreamsResponse.Result[0].ViewerCount)
	})
}

func TestGetLivestreamsStatsError(t *testing.T) {
	t.Run("on new request", func(t *testing.T) {
		kickClient, err := gokick.NewClient(&gokick.ClientOptions{UserAccessToken: "access-token"})
		require.NoError(t, err)

		var ctx context.Context
		_, err = kickClient.GetLivestreamsStats(ctx)
		require.EqualError(t, err, "failed to create request: net/http: nil Context")
	})

	t.Run("timeout", func(t *testing.T) {
		kickClient := setupTimeoutMockClient(t)

		_, err := kickClient.GetLivestreamsStats(context.Background())
		require.EqualError(t, err, `failed to make request: Get "https://api.kick.com/public/v1/livestreams/stats": context deadline exceeded `+
			`(Client.Timeout exceeded while awaiting headers)`)
	})

	t.Run("unmarshal error response", func(t *testing.T) {
		kickClient := setupMockClient(t, func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprint(w, `117`)
		})

		_, err := kickClient.GetLivestreamsStats(context.Background())

		assert.EqualError(t, err, `failed to unmarshal error response (KICK status code: 500 and body "117"): json: cannot unmarshal `+
			`number into Go value of type gokick.errorResponse`)
	})

	t.Run("unmarshal Livestreams response", func(t *testing.T) {
		kickClient := setupMockClient(t, func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			fmt.Fprint(w, "117")
		})

		_, err := kickClient.GetLivestreamsStats(context.Background())

		assert.EqualError(t, err, `failed to unmarshal response body (KICK status code 200 and body "117"): json: cannot unmarshal `+
			`number into Go value of type gokick.successResponse[github.com/scorfly/gokick.LivestreamStatsResponse]`)
	})

	t.Run("reader failure", func(t *testing.T) {
		kickClient := setupMockClient(t, func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Length", "10")
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprint(w, "")
		})

		_, err := kickClient.GetLivestreamsStats(context.Background())

		assert.EqualError(t, err, `failed to read response body (KICK status code 500): unexpected EOF`)
	})

	t.Run("with internal server error", func(t *testing.T) {
		kickClient := setupMockClient(t, func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprint(w, `{"message":"internal server error", "data":null}`)
		})

		_, err := kickClient.GetLivestreamsStats(context.Background())

		var kickError gokick.Error
		require.ErrorAs(t, err, &kickError)
		assert.Equal(t, http.StatusInternalServerError, kickError.Code())
		assert.Equal(t, "internal server error", kickError.Message())
	})
}

func TestGetLivestreamsStatsSuccess(t *testing.T) {
	kickClient := setupMockClient(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, `{"message":"success", "data":{
				"broadcaster_user_id": 219,
				"category": {
					"id": 123,
					"thumbnail": "category image url",
					"name": "category name"
				},
				"channel_id": 198,
				"has_mature_content": true,
				"language": "fr",
				"slug": "slug",
				"started_at": "started_at",
				"stream_title": "stream_title",
				"thumbnail": "thumbnail_url",
				"viewer_count": 167
				}
  			}`)
	})

	LivestreamsResponse, err := kickClient.GetLivestreamsStats(context.Background())
	require.NoError(t, err)
	assert.Equal(t, 219, LivestreamsResponse.Result.BroadcasterUserID)
	assert.Equal(t, 123, LivestreamsResponse.Result.Category.ID)
	assert.Equal(t, "category name", LivestreamsResponse.Result.Category.Name)
	assert.Equal(t, "category image url", LivestreamsResponse.Result.Category.Thumbnail)
	assert.Equal(t, 198, LivestreamsResponse.Result.ChannelID)
	assert.True(t, LivestreamsResponse.Result.HasMatureContent)
	assert.Equal(t, "fr", LivestreamsResponse.Result.Language)
	assert.Equal(t, "slug", LivestreamsResponse.Result.Slug)
	assert.Equal(t, "started_at", LivestreamsResponse.Result.StartedAt)
	assert.Equal(t, "stream_title", LivestreamsResponse.Result.StreamTitle)
	assert.Equal(t, "thumbnail_url", LivestreamsResponse.Result.Thumbnail)
	assert.Equal(t, 167, LivestreamsResponse.Result.ViewerCount)
}

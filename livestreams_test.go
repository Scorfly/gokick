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
			filter:              gokick.NewLivestreamListFilter().SetBroadcasterUserIDs([]int{118}),
			expectedQueryString: "?broadcaster_user_id=118",
		},
		"with multiple broadcaster user IDs": {
			filter:              gokick.NewLivestreamListFilter().SetBroadcasterUserIDs([]int{118, 218}),
			expectedQueryString: "?broadcaster_user_id=118&broadcaster_user_id=218",
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
				SetBroadcasterUserIDs([]int{118}).
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
				"custom_tags": ["tag1", "tag2"],
				"has_mature_content": true,
				"language": "fr",
				"profile_picture": "profile_picture_url",
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
		assert.Equal(t, []string{"tag1", "tag2"}, LivestreamsResponse.Result[0].CustomTags)
		assert.True(t, LivestreamsResponse.Result[0].HasMatureContent)
		assert.Equal(t, "fr", LivestreamsResponse.Result[0].Language)
		assert.Equal(t, "profile_picture_url", LivestreamsResponse.Result[0].ProfilePicture)
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
		fmt.Fprint(w, `{"message":"success", "data":{"total_count": 100}}`)
	})

	LivestreamsResponse, err := kickClient.GetLivestreamsStats(context.Background())
	require.NoError(t, err)
	assert.Equal(t, 100, LivestreamsResponse.Result.TotalCount)
}

func TestNewLivestreamV2ListFilterSuccess(t *testing.T) {
	testCases := map[string]struct {
		filter              gokick.LivestreamV2ListFilter
		expectedQueryString string
	}{
		"default": {
			filter:              gokick.NewLivestreamV2ListFilter(),
			expectedQueryString: "",
		},
		"with cursor": {
			filter:              gokick.NewLivestreamV2ListFilter().SetCursor("abc123"),
			expectedQueryString: "?cursor=abc123",
		},
		"with limit": {
			filter:              gokick.NewLivestreamV2ListFilter().SetLimit(50),
			expectedQueryString: "?limit=50",
		},
		"with single category ID": {
			filter:              gokick.NewLivestreamV2ListFilter().AddCategoryID(15),
			expectedQueryString: "?category_id=15",
		},
		"with multiple category IDs": {
			filter:              gokick.NewLivestreamV2ListFilter().AddCategoryID(15).AddCategoryID(20),
			expectedQueryString: "?category_id=15&category_id=20",
		},
		"with language code": {
			filter:              gokick.NewLivestreamV2ListFilter().AddLanguageCode("fr"),
			expectedQueryString: "?language_code=fr",
		},
		"with multiple language codes": {
			filter:              gokick.NewLivestreamV2ListFilter().AddLanguageCode("en").AddLanguageCode("fr"),
			expectedQueryString: "?language_code=en&language_code=fr",
		},
		"with all params": {
			filter: gokick.NewLivestreamV2ListFilter().
				SetCursor("next").
				SetLimit(100).
				AddCategoryID(15).
				AddLanguageCode("en"),
			expectedQueryString: "?category_id=15&cursor=next&language_code=en&limit=100",
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			assert.Equal(t, tc.expectedQueryString, tc.filter.ToQueryString())
		})
	}
}

func TestNewUsersLivestreamsFilterSuccess(t *testing.T) {
	testCases := map[string]struct {
		filter              gokick.UsersLivestreamsFilter
		expectedQueryString string
	}{
		"default": {
			filter:              gokick.NewUsersLivestreamsFilter(),
			expectedQueryString: "",
		},
		"with single user ID": {
			filter:              gokick.NewUsersLivestreamsFilter().AddUserID(118),
			expectedQueryString: "?user_id=118",
		},
		"with multiple user IDs": {
			filter:              gokick.NewUsersLivestreamsFilter().AddUserID(118).AddUserID(219),
			expectedQueryString: "?user_id=118&user_id=219",
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			assert.Equal(t, tc.expectedQueryString, tc.filter.ToQueryString())
		})
	}
}

func TestGetLivestreamsV2Error(t *testing.T) {
	t.Run("on new request", func(t *testing.T) {
		kickClient, err := gokick.NewClient(&gokick.ClientOptions{UserAccessToken: "access-token"})
		require.NoError(t, err)

		var ctx context.Context
		_, err = kickClient.GetLivestreamsV2(ctx, gokick.NewLivestreamV2ListFilter())
		require.EqualError(t, err, "failed to create request: net/http: nil Context")
	})

	t.Run("timeout", func(t *testing.T) {
		kickClient := setupTimeoutMockClient(t)

		_, err := kickClient.GetLivestreamsV2(context.Background(), gokick.NewLivestreamV2ListFilter())
		require.EqualError(t, err, `failed to make request: Get "https://api.kick.com/public/v2/livestreams": context deadline exceeded `+
			`(Client.Timeout exceeded while awaiting headers)`)
	})

	t.Run("unmarshal error response", func(t *testing.T) {
		kickClient := setupMockClient(t, func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprint(w, `117`)
		})

		_, err := kickClient.GetLivestreamsV2(context.Background(), gokick.NewLivestreamV2ListFilter())

		assert.EqualError(t, err, `failed to unmarshal error response (KICK status code: 500 and body "117"): json: cannot unmarshal `+
			`number into Go value of type gokick.errorResponse`)
	})

	t.Run("unmarshal livestreams response", func(t *testing.T) {
		kickClient := setupMockClient(t, func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			fmt.Fprint(w, "117")
		})

		_, err := kickClient.GetLivestreamsV2(context.Background(), gokick.NewLivestreamV2ListFilter())

		assert.Contains(t, err.Error(), `failed to unmarshal response body (KICK status code 200 and body "117")`)
	})

	t.Run("reader failure", func(t *testing.T) {
		kickClient := setupMockClient(t, func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Length", "10")
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprint(w, "")
		})

		_, err := kickClient.GetLivestreamsV2(context.Background(), gokick.NewLivestreamV2ListFilter())

		assert.EqualError(t, err, `failed to read response body (KICK status code 500): unexpected EOF`)
	})

	t.Run("with internal server error", func(t *testing.T) {
		kickClient := setupMockClient(t, func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprint(w, `{"message":"internal server error", "data":null}`)
		})

		_, err := kickClient.GetLivestreamsV2(context.Background(), gokick.NewLivestreamV2ListFilter())

		var kickError gokick.Error
		require.ErrorAs(t, err, &kickError)
		assert.Equal(t, http.StatusInternalServerError, kickError.Code())
		assert.Equal(t, "internal server error", kickError.Message())
	})
}

func TestGetLivestreamsV2Success(t *testing.T) {
	t.Run("without result", func(t *testing.T) {
		kickClient := setupMockClient(t, func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			fmt.Fprint(w, `{"message":"success", "data":[],"pagination":{"next_cursor":""}}`)
		})

		response, err := kickClient.GetLivestreamsV2(context.Background(), gokick.NewLivestreamV2ListFilter())
		require.NoError(t, err)
		assert.Empty(t, response.Result)
		assert.Empty(t, response.Pagination.NextCursor)
	})

	t.Run("with result", func(t *testing.T) {
		kickClient := setupMockClient(t, func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			fmt.Fprint(w, `{"message":"success", "data":[{
				"id": "550e8400-e29b-41d4-a716-446655440000",
				"broadcaster_user": {
					"id": 219,
					"profile_picture": "profile_picture_url",
					"username": "streamer"
				},
				"category": {
					"id": 123,
					"name": "category name",
					"thumbnail": "category image url"
				},
				"channel": {
					"slug": "slug"
				},
				"has_mature_content": true,
				"language_code": "fr",
				"started_at": "started_at",
				"tags": ["tag1", "tag2"],
				"thumbnail": "thumbnail_url",
				"title": "stream title",
				"viewer_count": 167
			}], "pagination":{"next_cursor":"next-page"}}`)
		})

		response, err := kickClient.GetLivestreamsV2(context.Background(), gokick.NewLivestreamV2ListFilter())
		require.NoError(t, err)
		require.Len(t, response.Result, 1)
		assert.Equal(t, "550e8400-e29b-41d4-a716-446655440000", response.Result[0].ID)
		assert.Equal(t, 219, response.Result[0].BroadcasterUser.ID)
		assert.Equal(t, "profile_picture_url", response.Result[0].BroadcasterUser.ProfilePicture)
		assert.Equal(t, "streamer", response.Result[0].BroadcasterUser.Username)
		assert.Equal(t, 123, response.Result[0].Category.ID)
		assert.Equal(t, "category name", response.Result[0].Category.Name)
		assert.Equal(t, "category image url", response.Result[0].Category.Thumbnail)
		assert.Equal(t, "slug", response.Result[0].Channel.Slug)
		assert.True(t, response.Result[0].HasMatureContent)
		assert.Equal(t, "fr", response.Result[0].LanguageCode)
		assert.Equal(t, "started_at", response.Result[0].StartedAt)
		assert.Equal(t, []string{"tag1", "tag2"}, response.Result[0].Tags)
		assert.Equal(t, "thumbnail_url", response.Result[0].Thumbnail)
		assert.Equal(t, "stream title", response.Result[0].Title)
		assert.Equal(t, 167, response.Result[0].ViewerCount)
		assert.Equal(t, "next-page", response.Pagination.NextCursor)
	})
}

func TestGetUsersLivestreamsError(t *testing.T) {
	t.Run("on new request", func(t *testing.T) {
		kickClient, err := gokick.NewClient(&gokick.ClientOptions{UserAccessToken: "access-token"})
		require.NoError(t, err)

		var ctx context.Context
		_, err = kickClient.GetUsersLivestreams(ctx, gokick.NewUsersLivestreamsFilter().AddUserID(118))
		require.EqualError(t, err, "failed to create request: net/http: nil Context")
	})

	t.Run("timeout", func(t *testing.T) {
		kickClient := setupTimeoutMockClient(t)

		_, err := kickClient.GetUsersLivestreams(context.Background(), gokick.NewUsersLivestreamsFilter().AddUserID(118))
		require.EqualError(t, err, `failed to make request: Get "https://api.kick.com/public/v1/users/livestreams?user_id=118": `+
			`context deadline exceeded (Client.Timeout exceeded while awaiting headers)`)
	})

	t.Run("unmarshal error response", func(t *testing.T) {
		kickClient := setupMockClient(t, func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprint(w, `117`)
		})

		_, err := kickClient.GetUsersLivestreams(context.Background(), gokick.NewUsersLivestreamsFilter().AddUserID(118))

		assert.EqualError(t, err, `failed to unmarshal error response (KICK status code: 500 and body "117"): json: cannot unmarshal `+
			`number into Go value of type gokick.errorResponse`)
	})

	t.Run("unmarshal livestreams response", func(t *testing.T) {
		kickClient := setupMockClient(t, func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			fmt.Fprint(w, "117")
		})

		_, err := kickClient.GetUsersLivestreams(context.Background(), gokick.NewUsersLivestreamsFilter().AddUserID(118))

		assert.EqualError(t, err, `failed to unmarshal response body (KICK status code 200 and body "117"): json: cannot unmarshal `+
			`number into Go value of type gokick.successResponse[[]github.com/scorfly/gokick.LivestreamV2Response]`)
	})

	t.Run("reader failure", func(t *testing.T) {
		kickClient := setupMockClient(t, func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Length", "10")
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprint(w, "")
		})

		_, err := kickClient.GetUsersLivestreams(context.Background(), gokick.NewUsersLivestreamsFilter().AddUserID(118))

		assert.EqualError(t, err, `failed to read response body (KICK status code 500): unexpected EOF`)
	})

	t.Run("with internal server error", func(t *testing.T) {
		kickClient := setupMockClient(t, func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprint(w, `{"message":"internal server error", "data":null}`)
		})

		_, err := kickClient.GetUsersLivestreams(context.Background(), gokick.NewUsersLivestreamsFilter().AddUserID(118))

		var kickError gokick.Error
		require.ErrorAs(t, err, &kickError)
		assert.Equal(t, http.StatusInternalServerError, kickError.Code())
		assert.Equal(t, "internal server error", kickError.Message())
	})
}

func TestGetUsersLivestreamsSuccess(t *testing.T) {
	t.Run("without result", func(t *testing.T) {
		kickClient := setupMockClient(t, func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			fmt.Fprint(w, `{"message":"success", "data":[]}`)
		})

		response, err := kickClient.GetUsersLivestreams(context.Background(), gokick.NewUsersLivestreamsFilter().AddUserID(118))
		require.NoError(t, err)
		assert.Empty(t, response.Result)
	})

	t.Run("with result", func(t *testing.T) {
		kickClient := setupMockClient(t, func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			fmt.Fprint(w, `{"message":"success", "data":[{
				"id": "550e8400-e29b-41d4-a716-446655440000",
				"broadcaster_user": {
					"id": 219,
					"profile_picture": "profile_picture_url",
					"username": "streamer"
				},
				"category": {
					"id": 123,
					"name": "category name",
					"thumbnail": "category image url"
				},
				"channel": {
					"slug": "slug"
				},
				"has_mature_content": false,
				"language_code": "en",
				"started_at": "2025-04-01T14:38:29Z",
				"tags": [],
				"thumbnail": "thumbnail_url",
				"title": "stream title",
				"viewer_count": 42
			}]}`)
		})

		response, err := kickClient.GetUsersLivestreams(context.Background(), gokick.NewUsersLivestreamsFilter().AddUserID(219))
		require.NoError(t, err)
		require.Len(t, response.Result, 1)
		assert.Equal(t, "550e8400-e29b-41d4-a716-446655440000", response.Result[0].ID)
		assert.Equal(t, 219, response.Result[0].BroadcasterUser.ID)
		assert.Equal(t, "stream title", response.Result[0].Title)
		assert.Equal(t, 42, response.Result[0].ViewerCount)
	})
}

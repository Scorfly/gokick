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

func TestNewCategoryListFilterSuccess(t *testing.T) {
	testCases := map[string]struct {
		filter              gokick.CategoryListFilter
		expectedQueryString string
	}{
		"default": {
			filter:              gokick.NewCategoryListFilter(),
			expectedQueryString: "",
		},
		"with cursor": {
			filter:              gokick.NewCategoryListFilter().SetCursor("abc123"),
			expectedQueryString: "?cursor=abc123",
		},
		"with limit": {
			filter:              gokick.NewCategoryListFilter().SetLimit(50),
			expectedQueryString: "?limit=50",
		},
		"with single name": {
			filter:              gokick.NewCategoryListFilter().AddName("gaming"),
			expectedQueryString: "?name=gaming",
		},
		"with two names": {
			filter:              gokick.NewCategoryListFilter().AddName("a").AddName("b"),
			expectedQueryString: "?name=a%2Cb",
		},
		"with tag": {
			filter:              gokick.NewCategoryListFilter().AddTag("fps"),
			expectedQueryString: "?tag=fps",
		},
		"with id": {
			filter:              gokick.NewCategoryListFilter().AddID(117),
			expectedQueryString: "?id=117",
		},
		"cursor limit and id": {
			filter: gokick.NewCategoryListFilter().
				SetCursor("next").
				SetLimit(10).
				AddID(5),
			expectedQueryString: "?cursor=next&id=5&limit=10",
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			assert.Equal(t, tc.expectedQueryString, tc.filter.ToQueryString())
		})
	}
}

func TestGetCategoriesError(t *testing.T) {
	t.Run("on new request", func(t *testing.T) {
		kickClient, err := gokick.NewClient(&gokick.ClientOptions{UserAccessToken: "access-token"})
		require.NoError(t, err)

		var ctx context.Context
		_, err = kickClient.GetCategories(ctx, gokick.NewCategoryListFilter())
		require.EqualError(t, err, "failed to create request: net/http: nil Context")
	})

	t.Run("timeout", func(t *testing.T) {
		kickClient := setupTimeoutMockClient(t)

		_, err := kickClient.GetCategories(context.Background(), gokick.NewCategoryListFilter())
		require.EqualError(t, err, `failed to make request: Get "https://api.kick.com/public/v2/categories": context deadline exceeded `+
			`(Client.Timeout exceeded while awaiting headers)`)
	})

	t.Run("unmarshal error response", func(t *testing.T) {
		kickClient := setupMockClient(t, func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprint(w, `117`)
		})

		_, err := kickClient.GetCategories(context.Background(), gokick.NewCategoryListFilter())

		assert.EqualError(t, err, `failed to unmarshal error response (KICK status code: 500 and body "117"): json: cannot unmarshal `+
			`number into Go value of type gokick.errorResponse`)
	})

	t.Run("unmarshal categories response", func(t *testing.T) {
		kickClient := setupMockClient(t, func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			fmt.Fprint(w, "117")
		})

		_, err := kickClient.GetCategories(context.Background(), gokick.NewCategoryListFilter())

		assert.Contains(t, err.Error(), `failed to unmarshal response body (KICK status code 200 and body "117")`)
	})

	t.Run("reader failure", func(t *testing.T) {
		kickClient := setupMockClient(t, func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Length", "10")
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprint(w, "")
		})

		_, err := kickClient.GetCategories(context.Background(), gokick.NewCategoryListFilter())

		assert.EqualError(t, err, `failed to read response body (KICK status code 500): unexpected EOF`)
	})

	t.Run("with internal server error", func(t *testing.T) {
		kickClient := setupMockClient(t, func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprint(w, `{"message":"internal server error", "data":null}`)
		})

		_, err := kickClient.GetCategories(context.Background(), gokick.NewCategoryListFilter())

		var kickError gokick.Error
		require.ErrorAs(t, err, &kickError)
		assert.Equal(t, http.StatusInternalServerError, kickError.Code())
		assert.Equal(t, "internal server error", kickError.Message())
	})
}

func TestGetCategoriesSuccess(t *testing.T) {
	t.Run("without result", func(t *testing.T) {
		kickClient := setupMockClient(t, func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			fmt.Fprint(w, `{"message":"success", "data":[],"pagination":{"next_cursor":""}}`)
		})

		categoriesResponse, err := kickClient.GetCategories(context.Background(), gokick.NewCategoryListFilter())
		require.NoError(t, err)
		assert.Empty(t, categoriesResponse.Result)
		assert.Empty(t, categoriesResponse.Pagination.NextCursor)
	})

	t.Run("with result", func(t *testing.T) {
		kickClient := setupMockClient(t, func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			fmt.Fprint(w, `{"message":"success", "data":[
				{"id":1, "name":"Music", "tags":["t1"], "thumbnail":"music"},
				{"id":2, "name":"Comedy", "thumbnail":"comedy"}
			], "pagination":{"next_cursor":"next-page"}}`)
		})

		categoriesResponse, err := kickClient.GetCategories(context.Background(), gokick.NewCategoryListFilter())
		require.NoError(t, err)
		require.Len(t, categoriesResponse.Result, 2)
		assert.Equal(t, 1, categoriesResponse.Result[0].ID)
		assert.Equal(t, "Music", categoriesResponse.Result[0].Name)
		assert.Equal(t, []string{"t1"}, categoriesResponse.Result[0].Tags)
		assert.Equal(t, "music", categoriesResponse.Result[0].Thumbnail)
		assert.Equal(t, 2, categoriesResponse.Result[1].ID)
		assert.Equal(t, "Comedy", categoriesResponse.Result[1].Name)
		assert.Equal(t, "comedy", categoriesResponse.Result[1].Thumbnail)
		assert.Equal(t, "next-page", categoriesResponse.Pagination.NextCursor)
	})
}

func TestGetCategoryError(t *testing.T) {
	t.Run("on new request", func(t *testing.T) {
		kickClient, err := gokick.NewClient(&gokick.ClientOptions{UserAccessToken: "access-token"})
		require.NoError(t, err)

		var ctx context.Context
		_, err = kickClient.GetCategory(ctx, 117)
		require.EqualError(t, err, "failed to create request: net/http: nil Context")
	})

	t.Run("timeout", func(t *testing.T) {
		kickClient := setupTimeoutMockClient(t)

		_, err := kickClient.GetCategory(context.Background(), 117)
		wantErr := `failed to make request: Get "https://api.kick.com/public/v2/categories?id=117&limit=1": ` +
			`context deadline exceeded (Client.Timeout exceeded while awaiting headers)`
		require.EqualError(t, err, wantErr)
	})

	t.Run("unmarshal error response", func(t *testing.T) {
		kickClient := setupMockClient(t, func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprint(w, `117`)
		})

		_, err := kickClient.GetCategory(context.Background(), 117)

		assert.EqualError(t, err, `failed to unmarshal error response (KICK status code: 500 and body "117"): json: cannot unmarshal `+
			`number into Go value of type gokick.errorResponse`)
	})

	t.Run("unmarshal categories response", func(t *testing.T) {
		kickClient := setupMockClient(t, func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			fmt.Fprint(w, "117")
		})

		_, err := kickClient.GetCategory(context.Background(), 117)

		assert.Contains(t, err.Error(), `failed to unmarshal response body (KICK status code 200 and body "117")`)
	})

	t.Run("reader failure", func(t *testing.T) {
		kickClient := setupMockClient(t, func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Length", "10")
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprint(w, "")
		})

		_, err := kickClient.GetCategory(context.Background(), 117)

		assert.EqualError(t, err, `failed to read response body (KICK status code 500): unexpected EOF`)
	})

	t.Run("with internal server error", func(t *testing.T) {
		kickClient := setupMockClient(t, func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprint(w, `{"message":"internal server error", "data":null}`)
		})

		_, err := kickClient.GetCategory(context.Background(), 117)

		var kickError gokick.Error
		require.ErrorAs(t, err, &kickError)
		assert.Equal(t, http.StatusInternalServerError, kickError.Code())
		assert.Equal(t, "internal server error", kickError.Message())
	})

	t.Run("empty result", func(t *testing.T) {
		kickClient := setupMockClient(t, func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			fmt.Fprint(w, `{"message":"success", "data":[],"pagination":{"next_cursor":""}}`)
		})

		_, err := kickClient.GetCategory(context.Background(), 999)
		require.EqualError(t, err, "category id 999: empty result")
	})
}

func TestGetCategorySuccess(t *testing.T) {
	kickClient := setupMockClient(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, `{
			"message":"success",
			"data":[{"id":117, "name":"Hubert", "tags":["x"], "thumbnail":"Bonisseur de La Bath"}],
			"pagination":{"next_cursor":""}
		}`)
	})

	categoryResponse, err := kickClient.GetCategory(context.Background(), 117)
	require.NoError(t, err)
	assert.Equal(t, 117, categoryResponse.Result.ID)
	assert.Equal(t, "Hubert", categoryResponse.Result.Name)
	assert.Equal(t, []string{"x"}, categoryResponse.Result.Tags)
	assert.Equal(t, "Bonisseur de La Bath", categoryResponse.Result.Thumbnail)
}

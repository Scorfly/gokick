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

func TestBanUserError(t *testing.T) {
	t.Run("on new request", func(t *testing.T) {
		kickClient, err := gokick.NewClient(&gokick.ClientOptions{UserAccessToken: "access-token"})
		require.NoError(t, err)

		var ctx context.Context
		_, err = kickClient.BanUser(ctx, 1234, 345, nil, nil)
		require.EqualError(t, err, "failed to create request: net/http: nil Context")
	})

	t.Run("timeout", func(t *testing.T) {
		kickClient := setupTimeoutMockClient(t)

		_, err := kickClient.BanUser(context.Background(), 1234, 345, nil, nil)
		require.EqualError(t, err, `failed to make request: Post "https://api.kick.com/public/v1/moderation/bans": context deadline exceeded `+
			`(Client.Timeout exceeded while awaiting headers)`)
	})

	t.Run("unmarshal error response", func(t *testing.T) {
		kickClient := setupMockClient(t, func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprint(w, `117`)
		})

		_, err := kickClient.BanUser(context.Background(), 1234, 345, nil, nil)
		require.EqualError(t, err, `failed to unmarshal error response (KICK status code: 500 and body "117"): json: cannot unmarshal `+
			`number into Go value of type gokick.errorResponse`)
	})

	t.Run("unmarshal token response", func(t *testing.T) {
		kickClient := setupMockClient(t, func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			fmt.Fprint(w, "117")
		})

		_, err := kickClient.BanUser(context.Background(), 1234, 345, nil, nil)

		assert.EqualError(t, err, `failed to unmarshal response body (KICK status code 200 and body "117"): json: cannot unmarshal `+
			`number into Go value of type gokick.successResponse[github.com/scorfly/gokick.BanUserResponse]`)
	})

	t.Run("reader failure", func(t *testing.T) {
		kickClient := setupMockClient(t, func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Length", "10")
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprint(w, "")
		})

		_, err := kickClient.BanUser(context.Background(), 1234, 345, nil, nil)

		assert.EqualError(t, err, `failed to read response body (KICK status code 500): unexpected EOF`)
	})

	t.Run("with internal server error", func(t *testing.T) {
		kickClient := setupMockClient(t, func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprint(w, `{"message":"internal server error", "data":null}`)
		})

		_, err := kickClient.BanUser(context.Background(), 1234, 345, nil, nil)

		var kickError gokick.Error
		require.ErrorAs(t, err, &kickError)
		assert.Equal(t, http.StatusInternalServerError, kickError.Code())
		assert.Equal(t, "internal server error", kickError.Message())
	})
}

func TestBanUserSuccess(t *testing.T) {
	testCases := map[string]struct {
		reason   *string
		duration *int
	}{
		"default": {},
		"with reason": {
			reason: stringPtr("ban reason"),
		},
		"with duration": {
			duration: intPtr(123),
		},
		"with reason and duration": {
			reason:   stringPtr("ban reason"),
			duration: intPtr(123),
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			kickClient := setupMockClient(t, func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
				fmt.Fprint(w, `{
					"message":"success",
					"data":{}
				}`)
			})

			_, err := kickClient.BanUser(context.Background(), 1234, 345, tc.duration, tc.reason)
			require.NoError(t, err)
		})
	}
}

func TestUnbanUserError(t *testing.T) {
	t.Run("on new request", func(t *testing.T) {
		kickClient, err := gokick.NewClient(&gokick.ClientOptions{UserAccessToken: "access-token"})
		require.NoError(t, err)

		var ctx context.Context
		_, err = kickClient.UnbanUser(ctx, 1234, 345)
		require.EqualError(t, err, "failed to create request: net/http: nil Context")
	})

	t.Run("timeout", func(t *testing.T) {
		kickClient := setupTimeoutMockClient(t)

		_, err := kickClient.UnbanUser(context.Background(), 1234, 345)
		require.EqualError(t, err, `failed to make request: Delete "https://api.kick.com/public/v1/moderation/bans": context deadline exceeded `+
			`(Client.Timeout exceeded while awaiting headers)`)
	})

	t.Run("unmarshal error response", func(t *testing.T) {
		kickClient := setupMockClient(t, func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprint(w, `117`)
		})

		_, err := kickClient.UnbanUser(context.Background(), 1234, 345)
		require.EqualError(t, err, `failed to unmarshal error response (KICK status code: 500 and body "117"): json: cannot unmarshal `+
			`number into Go value of type gokick.errorResponse`)
	})

	t.Run("unmarshal token response", func(t *testing.T) {
		kickClient := setupMockClient(t, func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			fmt.Fprint(w, "117")
		})

		_, err := kickClient.UnbanUser(context.Background(), 1234, 345)

		assert.EqualError(t, err, `failed to unmarshal response body (KICK status code 200 and body "117"): json: cannot unmarshal `+
			`number into Go value of type gokick.successResponse[github.com/scorfly/gokick.BanUserResponse]`)
	})

	t.Run("reader failure", func(t *testing.T) {
		kickClient := setupMockClient(t, func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Length", "10")
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprint(w, "")
		})

		_, err := kickClient.UnbanUser(context.Background(), 1234, 345)

		assert.EqualError(t, err, `failed to read response body (KICK status code 500): unexpected EOF`)
	})

	t.Run("with internal server error", func(t *testing.T) {
		kickClient := setupMockClient(t, func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprint(w, `{"message":"internal server error", "data":null}`)
		})

		_, err := kickClient.UnbanUser(context.Background(), 1234, 345)

		var kickError gokick.Error
		require.ErrorAs(t, err, &kickError)
		assert.Equal(t, http.StatusInternalServerError, kickError.Code())
		assert.Equal(t, "internal server error", kickError.Message())
	})
}

func TestUnbanUserSuccess(t *testing.T) {
	kickClient := setupMockClient(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, `{
			"message":"success",
			"data":{}
		}`)
	})

	_, err := kickClient.UnbanUser(context.Background(), 1234, 345)
	require.NoError(t, err)
}

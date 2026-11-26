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

func TestCreateAdBreakError(t *testing.T) {
	req := gokick.CreateAdBreakRequest{
		BreakDurationSeconds: 30,
		ID:                   "123e4567-e89b-12d3-a456-426614174000",
	}

	t.Run("on new request", func(t *testing.T) {
		kickClient, err := gokick.NewClient(&gokick.ClientOptions{UserAccessToken: "access-token"})
		require.NoError(t, err)

		var ctx context.Context
		_, err = kickClient.CreateAdBreak(ctx, req)
		require.EqualError(t, err, "failed to create request: net/http: nil Context")
	})

	t.Run("timeout", func(t *testing.T) {
		kickClient := setupTimeoutMockClient(t)

		_, err := kickClient.CreateAdBreak(context.Background(), req)
		require.EqualError(t, err, `failed to make request: Post "https://api.kick.com/public/v1/ads/ad-break": context deadline exceeded `+
			`(Client.Timeout exceeded while awaiting headers)`)
	})

	t.Run("unmarshal error response", func(t *testing.T) {
		kickClient := setupMockClient(t, func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprint(w, `117`)
		})

		_, err := kickClient.CreateAdBreak(context.Background(), req)
		assert.EqualError(t, err, `failed to unmarshal error response (KICK status code: 500 and body "117"): json: cannot unmarshal `+
			`number into Go value of type gokick.errorResponse`)
	})

	t.Run("unmarshal ad break response", func(t *testing.T) {
		kickClient := setupMockClient(t, func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusAccepted)
			fmt.Fprint(w, "117")
		})

		_, err := kickClient.CreateAdBreak(context.Background(), req)
		assert.EqualError(t, err, `failed to unmarshal response body (KICK status code 202 and body "117"): json: cannot unmarshal `+
			`number into Go value of type gokick.successResponse[github.com/scorfly/gokick.AdBreakResponse]`)
	})

	t.Run("reader failure", func(t *testing.T) {
		kickClient := setupMockClient(t, func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Length", "10")
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprint(w, "")
		})

		_, err := kickClient.CreateAdBreak(context.Background(), req)
		assert.EqualError(t, err, `failed to read response body (KICK status code 500): unexpected EOF`)
	})

	t.Run("with internal server error", func(t *testing.T) {
		kickClient := setupMockClient(t, func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprint(w, `{"message":"internal server error", "data":null}`)
		})

		_, err := kickClient.CreateAdBreak(context.Background(), req)

		var kickError gokick.Error
		require.ErrorAs(t, err, &kickError)
		assert.Equal(t, http.StatusInternalServerError, kickError.Code())
		assert.Equal(t, "internal server error", kickError.Message())
		assert.Empty(t, kickError.Description())
	})

	errorCases := map[string]struct {
		status int
		code   string
	}{
		"channel is not live": {
			status: http.StatusBadRequest,
			code:   "CHANNEL_IS_NOT_LIVE",
		},
		"channel not opted in": {
			status: http.StatusBadRequest,
			code:   "CHANNEL_NOT_OPTED_IN",
		},
		"ads blocked": {
			status: http.StatusBadRequest,
			code:   "ADS_BLOCKED",
		},
		"ads restricted": {
			status: http.StatusBadRequest,
			code:   "ADS_RESTRICTED",
		},
		"ad break active": {
			status: http.StatusConflict,
			code:   "AD_BREAK_ACTIVE",
		},
		"max ad breaks reached": {
			status: http.StatusTooManyRequests,
			code:   "MAX_AD_BREAKS_REACHED",
		},
		"auto ads active": {
			status: http.StatusConflict,
			code:   "AUTO_ADS_ACTIVE",
		},
	}

	for name, testCase := range errorCases {
		t.Run(name, func(t *testing.T) {
			kickClient := setupMockClient(t, func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(testCase.status)
				fmt.Fprintf(w, `{"message":"Bad request","data":%q}`, testCase.code)
			})

			_, err := kickClient.CreateAdBreak(context.Background(), req)

			var kickError gokick.Error
			require.ErrorAs(t, err, &kickError)
			assert.Equal(t, testCase.status, kickError.Code())
			assert.Equal(t, "Bad request", kickError.Message())
			assert.Equal(t, testCase.code, kickError.Description())

			code, parseErr := gokick.NewAdBreakErrorCode(kickError.Description())
			require.NoError(t, parseErr)
			assert.Equal(t, testCase.code, code.String())
		})
	}
}

func TestCreateAdBreakSuccess(t *testing.T) {
	kickClient := setupMockClient(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPost, r.Method)
		assert.Equal(t, "/public/v1/ads/ad-break", r.URL.Path)
		assert.Equal(t, "application/json", r.Header.Get("Content-Type"))
		assert.Equal(t, "Bearer access-token", r.Header.Get("Authorization"))

		var body map[string]interface{}
		err := json.NewDecoder(r.Body).Decode(&body)
		assert.NoError(t, err)
		assert.InEpsilon(t, float64(30), body["break_duration_seconds"], 0)
		assert.Equal(t, "123e4567-e89b-12d3-a456-426614174000", body["id"])

		w.WriteHeader(http.StatusAccepted)
		fmt.Fprint(w, `{
			"data": {
				"ad_break_id": "Xq3LmP9aB2cD",
				"id": "123e4567-e89b-12d3-a456-426614174000",
				"remaining_ad_breaks": 2
			},
			"message": "Accepted"
		}`)
	})

	response, err := kickClient.CreateAdBreak(context.Background(), gokick.CreateAdBreakRequest{
		BreakDurationSeconds: 30,
		ID:                   "123e4567-e89b-12d3-a456-426614174000",
	})
	require.NoError(t, err)
	assert.Equal(t, "Xq3LmP9aB2cD", response.Result.AdBreakID)
	assert.Equal(t, "123e4567-e89b-12d3-a456-426614174000", response.Result.ID)
	assert.Equal(t, 2, response.Result.RemainingAdBreaks)
}

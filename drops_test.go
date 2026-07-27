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

func TestNewDropsClaimsFilterSuccess(t *testing.T) {
	testCases := map[string]struct {
		filter              gokick.DropsClaimsFilter
		expectedQueryString string
	}{
		"default": {
			filter:              gokick.NewDropsClaimsFilter(),
			expectedQueryString: "",
		},
		"with campaign_id": {
			filter:              gokick.NewDropsClaimsFilter().SetCampaignID("01JAXK8N4QWRTY5PM7ZEBVJDS8S"),
			expectedQueryString: "?campaign_id=01JAXK8N4QWRTY5PM7ZEBVJDS8S",
		},
		"with limit": {
			filter:              gokick.NewDropsClaimsFilter().SetLimit(10),
			expectedQueryString: "?limit=10",
		},
		"with cursor": {
			filter:              gokick.NewDropsClaimsFilter().SetCursor("01K0TDDR08Q5ZNWXK92VRXD7GW"),
			expectedQueryString: "?cursor=01K0TDDR08Q5ZNWXK92VRXD7GW",
		},
		"with user_id": {
			filter:              gokick.NewDropsClaimsFilter().SetUserID(1234),
			expectedQueryString: "?user_id=1234",
		},
		"with claim_id": {
			filter:              gokick.NewDropsClaimsFilter().SetClaimID("01JAXK8N4QWRTY5PM7ZEBVGH2S"),
			expectedQueryString: "?claim_id=01JAXK8N4QWRTY5PM7ZEBVGH2S",
		},
		"with external_status": {
			filter:              gokick.NewDropsClaimsFilter().SetExternalStatus("initiated"),
			expectedQueryString: "?external_status=initiated",
		},
		"with all filters": {
			filter: gokick.NewDropsClaimsFilter().
				SetCampaignID("01JAXK8N4QWRTY5PM7ZEBVJDS8S").
				SetLimit(10).
				SetCursor("01K0TDDR08Q5ZNWXK92VRXD7GW").
				SetUserID(1234).
				SetExternalStatus("initiated"),
			expectedQueryString: `?campaign_id=01JAXK8N4QWRTY5PM7ZEBVJDS8S&cursor=01K0TDDR08Q5ZNWXK92VRXD7GW&` +
				`external_status=initiated&limit=10&user_id=1234`,
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			assert.Equal(t, tc.expectedQueryString, tc.filter.ToQueryString())
		})
	}
}

func TestGetDropsClaimsError(t *testing.T) {
	t.Run("on new request", func(t *testing.T) {
		kickClient, err := gokick.NewClient(&gokick.ClientOptions{AppAccessToken: "access-token"})
		require.NoError(t, err)

		var ctx context.Context
		_, err = kickClient.GetDropsClaims(ctx, gokick.NewDropsClaimsFilter())
		require.EqualError(t, err, "failed to create request: net/http: nil Context")
	})

	t.Run("timeout", func(t *testing.T) {
		kickClient := setupTimeoutMockClient(t)

		_, err := kickClient.GetDropsClaims(context.Background(), gokick.NewDropsClaimsFilter())
		require.EqualError(t, err, `failed to make request: Get "https://api.kick.com/public/v1/drops/claims": `+
			`context deadline exceeded (Client.Timeout exceeded while awaiting headers)`)
	})

	t.Run("unmarshal error response", func(t *testing.T) {
		kickClient := setupMockClient(t, func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprint(w, `117`)
		})

		_, err := kickClient.GetDropsClaims(context.Background(), gokick.NewDropsClaimsFilter())

		assert.EqualError(t, err, `failed to unmarshal error response (KICK status code: 500 and body "117"): json: cannot unmarshal `+
			`number into Go value of type gokick.errorResponse`)
	})

	t.Run("unmarshal claims response", func(t *testing.T) {
		kickClient := setupMockClient(t, func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			fmt.Fprint(w, "117")
		})

		_, err := kickClient.GetDropsClaims(context.Background(), gokick.NewDropsClaimsFilter())

		assert.EqualError(t, err, `failed to unmarshal response body (KICK status code 200 and body "117"): json: cannot unmarshal `+
			`number into Go value of type gokick.successResponse[github.com/scorfly/gokick.DropsClaimsData]`)
	})

	t.Run("unauthorized response", func(t *testing.T) {
		kickClient := setupMockClient(t, func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusUnauthorized)
			fmt.Fprint(w, `{"message":"Unauthorized", "data":null}`)
		})

		_, err := kickClient.GetDropsClaims(context.Background(), gokick.NewDropsClaimsFilter())

		var kickError gokick.Error
		require.ErrorAs(t, err, &kickError)
		assert.Equal(t, http.StatusUnauthorized, kickError.Code())
		assert.Equal(t, "Unauthorized", kickError.Message())
	})
}

func TestGetDropsClaimsSuccess(t *testing.T) {
	t.Run("when result is empty", func(t *testing.T) {
		kickClient := setupMockClient(t, func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			fmt.Fprint(w, `{"data":{"claims":[],"cursor":""}, "message":"OK"}`)
		})

		response, err := kickClient.GetDropsClaims(context.Background(), gokick.NewDropsClaimsFilter())
		require.NoError(t, err)
		assert.Empty(t, response.Result.Claims)
		assert.Empty(t, response.Result.Cursor)
	})

	t.Run("when result is filled", func(t *testing.T) {
		kickClient := setupMockClient(t, func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, http.MethodGet, r.Method)
			assert.Equal(t, "/public/v1/drops/claims", r.URL.Path)
			assert.Equal(t, "01JAXK8N4QWRTY5PM7ZEBVJDS8S", r.URL.Query().Get("campaign_id"))
			assert.Equal(t, "10", r.URL.Query().Get("limit"))
			assert.Equal(t, "01K0TDDR08Q5ZNWXK92VRXD7GW", r.URL.Query().Get("cursor"))
			assert.Equal(t, "1234", r.URL.Query().Get("user_id"))
			assert.Equal(t, "initiated", r.URL.Query().Get("external_status"))
			assert.Equal(t, "Bearer access-token", r.Header.Get("Authorization"))

			w.WriteHeader(http.StatusOK)
			fmt.Fprint(w, `{
				"data": {
					"claims": [{
						"claim_id": "01JAXK8N4QWRTY5PM7ZEBVGH2S",
						"user_id": 1234,
						"campaign_id": "01JAXK8N4QWRTY5PM7ZEBVJDS8S",
						"reward_id": "01JAXK8N4QWRTY5PM7ZEBVJDS9T",
						"external_id": "ext-123",
						"external_status": "initiated",
						"created_at": "2026-01-15T12:00:00Z",
						"updated_at": "2026-01-15T12:00:00Z"
					}],
					"cursor": "01K0TDDR08Q5ZNWXK92SDH7SDH"
				},
				"message": "OK"
			}`)
		})

		response, err := kickClient.GetDropsClaims(
			context.Background(),
			gokick.NewDropsClaimsFilter().
				SetCampaignID("01JAXK8N4QWRTY5PM7ZEBVJDS8S").
				SetLimit(10).
				SetCursor("01K0TDDR08Q5ZNWXK92VRXD7GW").
				SetUserID(1234).
				SetExternalStatus("initiated"),
		)
		require.NoError(t, err)
		require.Len(t, response.Result.Claims, 1)
		assert.Equal(t, "01JAXK8N4QWRTY5PM7ZEBVGH2S", response.Result.Claims[0].ClaimID)
		assert.Equal(t, 1234, response.Result.Claims[0].UserID)
		assert.Equal(t, "01JAXK8N4QWRTY5PM7ZEBVJDS8S", response.Result.Claims[0].CampaignID)
		assert.Equal(t, "01JAXK8N4QWRTY5PM7ZEBVJDS9T", response.Result.Claims[0].RewardID)
		assert.Equal(t, "ext-123", response.Result.Claims[0].ExternalID)
		assert.Equal(t, "initiated", response.Result.Claims[0].ExternalStatus)
		assert.Equal(t, "2026-01-15T12:00:00Z", response.Result.Claims[0].CreatedAt)
		assert.Equal(t, "2026-01-15T12:00:00Z", response.Result.Claims[0].UpdatedAt)
		assert.Equal(t, "01K0TDDR08Q5ZNWXK92SDH7SDH", response.Result.Cursor)
	})

	t.Run("when filtering by claim_id", func(t *testing.T) {
		kickClient := setupMockClient(t, func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "01JAXK8N4QWRTY5PM7ZEBVGH2S", r.URL.Query().Get("claim_id"))

			w.WriteHeader(http.StatusOK)
			fmt.Fprint(w, `{"data":{"claims":[],"cursor":""}, "message":"OK"}`)
		})

		_, err := kickClient.GetDropsClaims(
			context.Background(),
			gokick.NewDropsClaimsFilter().SetClaimID("01JAXK8N4QWRTY5PM7ZEBVGH2S"),
		)
		require.NoError(t, err)
	})
}

func TestUpdateDropsClaimsError(t *testing.T) {
	claims := []gokick.DropsClaimUpdate{
		{ClaimID: "01KAAFHJ2PNXS48NG8XXPGWKCZ", ExternalStatus: "processed"},
	}

	t.Run("on new request", func(t *testing.T) {
		kickClient, err := gokick.NewClient(&gokick.ClientOptions{AppAccessToken: "access-token"})
		require.NoError(t, err)

		var ctx context.Context
		_, err = kickClient.UpdateDropsClaims(ctx, claims)
		require.EqualError(t, err, "failed to create request: net/http: nil Context")
	})

	t.Run("timeout", func(t *testing.T) {
		kickClient := setupTimeoutMockClient(t)

		_, err := kickClient.UpdateDropsClaims(context.Background(), claims)
		require.EqualError(t, err, `failed to make request: Patch "https://api.kick.com/public/v1/drops/claims": `+
			`context deadline exceeded (Client.Timeout exceeded while awaiting headers)`)
	})

	t.Run("unmarshal error response", func(t *testing.T) {
		kickClient := setupMockClient(t, func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprint(w, `117`)
		})

		_, err := kickClient.UpdateDropsClaims(context.Background(), claims)

		assert.EqualError(t, err, `failed to unmarshal error response (KICK status code: 500 and body "117"): json: cannot unmarshal `+
			`number into Go value of type gokick.errorResponse`)
	})

	t.Run("unauthorized response", func(t *testing.T) {
		kickClient := setupMockClient(t, func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusUnauthorized)
			fmt.Fprint(w, `{"message":"Unauthorized", "data":null}`)
		})

		_, err := kickClient.UpdateDropsClaims(context.Background(), claims)

		var kickError gokick.Error
		require.ErrorAs(t, err, &kickError)
		assert.Equal(t, http.StatusUnauthorized, kickError.Code())
		assert.Equal(t, "Unauthorized", kickError.Message())
	})
}

func TestUpdateDropsClaimsSuccess(t *testing.T) {
	kickClient := setupMockClient(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPatch, r.Method)
		assert.Equal(t, "/public/v1/drops/claims", r.URL.Path)
		assert.Equal(t, "application/json", r.Header.Get("Content-Type"))
		assert.Equal(t, "Bearer access-token", r.Header.Get("Authorization"))

		var body map[string]interface{}
		err := json.NewDecoder(r.Body).Decode(&body)
		assert.NoError(t, err)
		assert.Equal(t, []interface{}{
			map[string]interface{}{
				"claim_id":        "01KAAFHJ2PNXS48NG8XXPGWKCZ",
				"external_status": "processed",
			},
			map[string]interface{}{
				"claim_id":        "01KAAFHJ2PNXS48NG8XXPGWKCA",
				"external_status": "initiated",
			},
		}, body["claims"])

		w.WriteHeader(http.StatusNoContent)
	})

	response, err := kickClient.UpdateDropsClaims(context.Background(), []gokick.DropsClaimUpdate{
		{ClaimID: "01KAAFHJ2PNXS48NG8XXPGWKCZ", ExternalStatus: "processed"},
		{ClaimID: "01KAAFHJ2PNXS48NG8XXPGWKCA", ExternalStatus: "initiated"},
	})
	require.NoError(t, err)
	assert.Equal(t, gokick.EmptyResponse{}, response)
}

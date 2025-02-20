package gokick_test

import (
	"context"
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/scorfly/gokick"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetSubscriptionsError(t *testing.T) {
	t.Run("on new request", func(t *testing.T) {
		kickClient, err := gokick.NewClient(&http.Client{}, "", "access-token")
		require.NoError(t, err)

		var ctx context.Context
		_, err = kickClient.GetSubscriptions(ctx)
		require.EqualError(t, err, "failed to create request: net/http: nil Context")
	})

	t.Run("timeout", func(t *testing.T) {
		kickClient, err := gokick.NewClient(&http.Client{Timeout: 1 * time.Nanosecond}, "", "access-token")
		require.NoError(t, err)

		_, err = kickClient.GetSubscriptions(context.Background())
		require.EqualError(t, err, `failed to make request: Get "https://api.kick.com/public/v1/events/subscriptions": context deadline `+
			`exceeded (Client.Timeout exceeded while awaiting headers)`)
	})

	t.Run("unmarshal error response", func(t *testing.T) {
		kickClient := setupMockClient(t, func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprint(w, `117`)
		})

		_, err := kickClient.GetSubscriptions(context.Background())

		assert.EqualError(t, err, `failed to unmarshal error response (KICK status code: 500 and body "117"): json: cannot unmarshal `+
			`number into Go value of type gokick.errorResponse`)
	})

	t.Run("unmarshal users response", func(t *testing.T) {
		kickClient := setupMockClient(t, func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			fmt.Fprint(w, "117")
		})

		_, err := kickClient.GetSubscriptions(context.Background())

		assert.EqualError(t, err, `failed to unmarshal response body (KICK status code 200 and body "117"): json: cannot unmarshal `+
			`number into Go value of type gokick.successResponse[[]github.com/scorfly/gokick.EventResponse]`)
	})

	t.Run("reader failure", func(t *testing.T) {
		kickClient := setupMockClient(t, func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Length", "10")
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprint(w, "")
		})

		_, err := kickClient.GetSubscriptions(context.Background())

		assert.EqualError(t, err, `failed to read response body (KICK status code 500): unexpected EOF`)
	})

	t.Run("with internal server error", func(t *testing.T) {
		kickClient := setupMockClient(t, func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprint(w, `{"message":"internal server error", "data":null}`)
		})

		_, err := kickClient.GetSubscriptions(context.Background())

		var kickError gokick.Error
		require.ErrorAs(t, err, &kickError)
		assert.Equal(t, http.StatusInternalServerError, kickError.Code())
		assert.Equal(t, "internal server error", kickError.Message())
	})
}

func TestGetSubscriptionsSuccess(t *testing.T) {
	t.Run("without result", func(t *testing.T) {
		kickClient := setupMockClient(t, func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			fmt.Fprint(w, `{"message":"success", "data":[]}`)
		})

		usersResponse, err := kickClient.GetSubscriptions(context.Background())
		require.NoError(t, err)
		assert.Empty(t, usersResponse.Result)
	})

	t.Run("with result", func(t *testing.T) {
		kickClient := setupMockClient(t, func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			fmt.Fprint(w, `{"message":"success", "data":[{
				"app_id": "app id",
				"broadcaster_user_id": 111,
				"created_at": "created at",
				"event": "event",
				"id": "id",
				"method": "method",
				"updated_at": "updated at",
				"version": 222
			}]}`)
		})

		response, err := kickClient.GetSubscriptions(context.Background())
		require.NoError(t, err)
		require.Len(t, response.Result, 1)
		assert.Equal(t, "app id", response.Result[0].AppID)
		assert.Equal(t, 111, response.Result[0].BroadcasterUserID)
		assert.Equal(t, "created at", response.Result[0].CreatedAt)
		assert.Equal(t, "event", response.Result[0].Event)
		assert.Equal(t, "id", response.Result[0].ID)
		assert.Equal(t, "method", response.Result[0].Method)
		assert.Equal(t, "updated at", response.Result[0].UpdatedAt)
		assert.Equal(t, 222, response.Result[0].Version)
	})
}

func TestNewSubscriptionMethodError(t *testing.T) {
	testCases := map[string]string{
		"empty":         "",
		"not supported": "not supported",
	}

	for name, value := range testCases {
		t.Run(name, func(t *testing.T) {
			_, err := gokick.NewSubscriptionMethod(value)
			assert.EqualError(t, err, fmt.Sprintf("unknown method: %s", value))
		})
	}
}

func TestNewSubscriptionMethodSuccess(t *testing.T) {
	testCases := map[string]gokick.SubscriptionMethod{
		"webhook": gokick.SubscriptionMethodWebhook,
	}

	for name, value := range testCases {
		t.Run(name, func(t *testing.T) {
			SubscriptionMethod, err := gokick.NewSubscriptionMethod(value.String())
			require.NoError(t, err)
			assert.Equal(t, SubscriptionMethod, value)
		})
	}
}

func TestNewSubscriptionNameError(t *testing.T) {
	testCases := map[string]string{
		"empty":         "",
		"not supported": "not supported",
	}

	for name, value := range testCases {
		t.Run(name, func(t *testing.T) {
			_, err := gokick.NewSubscriptionName(value)
			assert.EqualError(t, err, fmt.Sprintf("unknown name: %s", value))
		})
	}
}

func TestNewSubscriptionNameSuccess(t *testing.T) {
	testCases := map[string]gokick.SubscriptionName{
		"chat.message.sent":            gokick.SubscriptionNameChatMessage,
		"channel.followed":             gokick.SubscriptionNameChannelFollow,
		"channel.subscription.renewal": gokick.SubscriptionNameChannelSubscriptionRenewal,
		"channel.subscription.gifts":   gokick.SubscriptionNameChannelSubscriptionGifts,
		"channel.subscription.new":     gokick.SubscriptionNameChannelSubscriptionfCreated,
	}

	for name, value := range testCases {
		t.Run(name, func(t *testing.T) {
			SubscriptionName, err := gokick.NewSubscriptionName(value.String())
			require.NoError(t, err)
			assert.Equal(t, SubscriptionName, value)
		})
	}
}

func TestCreateSubscriptionsError(t *testing.T) {
	t.Run("on new request", func(t *testing.T) {
		kickClient, err := gokick.NewClient(&http.Client{}, "", "access-token")
		require.NoError(t, err)

		var ctx context.Context
		_, err = kickClient.CreateSubscriptions(ctx, gokick.SubscriptionMethodWebhook, []gokick.SubscriptionRequest{})
		require.EqualError(t, err, "failed to create request: net/http: nil Context")
	})

	t.Run("timeout", func(t *testing.T) {
		kickClient, err := gokick.NewClient(&http.Client{Timeout: 1 * time.Nanosecond}, "", "access-token")
		require.NoError(t, err)

		_, err = kickClient.CreateSubscriptions(context.Background(), gokick.SubscriptionMethodWebhook, []gokick.SubscriptionRequest{})
		require.EqualError(t, err, `failed to make request: Post "https://api.kick.com/public/v1/events/subscriptions": context deadline `+
			`exceeded (Client.Timeout exceeded while awaiting headers)`)
	})

	t.Run("unmarshal error response", func(t *testing.T) {
		kickClient := setupMockClient(t, func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprint(w, `117`)
		})

		_, err := kickClient.CreateSubscriptions(context.Background(), gokick.SubscriptionMethodWebhook, []gokick.SubscriptionRequest{})

		assert.EqualError(t, err, `failed to unmarshal error response (KICK status code: 500 and body "117"): json: cannot unmarshal `+
			`number into Go value of type gokick.errorResponse`)
	})

	t.Run("unmarshal token response", func(t *testing.T) {
		kickClient := setupMockClient(t, func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			fmt.Fprint(w, "117")
		})

		_, err := kickClient.CreateSubscriptions(context.Background(), gokick.SubscriptionMethodWebhook, []gokick.SubscriptionRequest{})

		assert.EqualError(t, err, `failed to unmarshal error response (KICK status code: 200 and body "117"): json: cannot unmarshal `+
			`number into Go value of type gokick.errorResponse`)
	})

	t.Run("reader failure", func(t *testing.T) {
		kickClient := setupMockClient(t, func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Length", "10")
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprint(w, "")
		})

		_, err := kickClient.CreateSubscriptions(context.Background(), gokick.SubscriptionMethodWebhook, []gokick.SubscriptionRequest{})

		assert.EqualError(t, err, `failed to read response body (KICK status code 500): unexpected EOF`)
	})

	t.Run("with internal server error", func(t *testing.T) {
		kickClient := setupMockClient(t, func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprint(w, `{"message":"internal server error", "data":null}`)
		})

		_, err := kickClient.CreateSubscriptions(context.Background(), gokick.SubscriptionMethodWebhook, []gokick.SubscriptionRequest{})

		var kickError gokick.Error
		require.ErrorAs(t, err, &kickError)
		assert.Equal(t, http.StatusInternalServerError, kickError.Code())
		assert.Equal(t, "internal server error", kickError.Message())
	})
}

func TestCreateSubscriptionsSuccess(t *testing.T) {
	kickClient := setupMockClient(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusCreated)
		fmt.Fprint(w, `{
			"message":"success",
			"data":[{"error": "error","name": "name","subscription_id": "subscription id","version": 1}]
		}`)
	})

	response, err := kickClient.CreateSubscriptions(context.Background(), gokick.SubscriptionMethodWebhook, []gokick.SubscriptionRequest{
		{
			Name:    gokick.SubscriptionNameChatMessage,
			Version: 1,
		},
		{
			Name:    gokick.SubscriptionNameChannelFollow,
			Version: 1,
		},
		{
			Name:    gokick.SubscriptionNameChannelSubscriptionRenewal,
			Version: 1,
		},
		{
			Name:    gokick.SubscriptionNameChannelSubscriptionGifts,
			Version: 1,
		},
		{
			Name:    gokick.SubscriptionNameChannelSubscriptionfCreated,
			Version: 1,
		},
	})

	require.NoError(t, err)
	assert.Len(t, response.Result, 1)
	assert.Equal(t, "error", response.Result[0].Error)
	assert.Equal(t, "name", response.Result[0].Name)
	assert.Equal(t, "subscription id", response.Result[0].SubscriptionID)
	assert.Equal(t, 1, response.Result[0].Version)
}

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

func TestSendChatMessageError(t *testing.T) {
	t.Run("on new request", func(t *testing.T) {
		kickClient, err := gokick.NewClient(&gokick.ClientOptions{UserAccessToken: "access-token"})
		require.NoError(t, err)

		var ctx context.Context
		_, err = kickClient.SendChatMessage(ctx, nil, "message", nil, gokick.MessageTypeBot)
		require.EqualError(t, err, "failed to create request: net/http: nil Context")
	})

	t.Run("timeout", func(t *testing.T) {
		kickClient := setupTimeoutMockClient(t)

		_, err := kickClient.SendChatMessage(context.Background(), nil, "message", nil, gokick.MessageTypeBot)
		require.EqualError(t, err, `failed to make request: Post "https://api.kick.com/public/v1/chat": context deadline exceeded `+
			`(Client.Timeout exceeded while awaiting headers)`)
	})

	t.Run("unmarshal error response", func(t *testing.T) {
		kickClient := setupMockClient(t, func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprint(w, `117`)
		})

		_, err := kickClient.SendChatMessage(context.Background(), nil, "message", nil, gokick.MessageTypeBot)

		assert.EqualError(t, err, `failed to unmarshal error response (KICK status code: 500 and body "117"): json: cannot unmarshal `+
			`number into Go value of type gokick.errorResponse`)
	})

	t.Run("unmarshal token response", func(t *testing.T) {
		kickClient := setupMockClient(t, func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			fmt.Fprint(w, "117")
		})

		_, err := kickClient.SendChatMessage(context.Background(), nil, "message", nil, gokick.MessageTypeBot)

		assert.EqualError(t, err, `failed to unmarshal response body (KICK status code 200 and body "117"): json: cannot unmarshal `+
			`number into Go value of type gokick.successResponse[github.com/scorfly/gokick.ChatResponse]`)
	})

	t.Run("reader failure", func(t *testing.T) {
		kickClient := setupMockClient(t, func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Length", "10")
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprint(w, "")
		})

		_, err := kickClient.SendChatMessage(context.Background(), nil, "message", nil, gokick.MessageTypeBot)

		assert.EqualError(t, err, `failed to read response body (KICK status code 500): unexpected EOF`)
	})

	t.Run("with internal server error", func(t *testing.T) {
		kickClient := setupMockClient(t, func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprint(w, `{"message":"internal server error", "data":null}`)
		})

		_, err := kickClient.SendChatMessage(context.Background(), nil, "message", nil, gokick.MessageTypeBot)

		var kickError gokick.Error
		require.ErrorAs(t, err, &kickError)
		assert.Equal(t, http.StatusInternalServerError, kickError.Code())
		assert.Equal(t, "internal server error", kickError.Message())
	})
}

func TestSendChatMessageSuccess(t *testing.T) {
	testCases := map[string]struct {
		broadcasterUserID   *int
		replyToMessageID    *string
		expectedQueryString string
	}{
		"without optional params": {},
		"with reply to message ID": {
			replyToMessageID: stringPtr("message-id"),
		},
		"with broadcaster user ID": {
			broadcasterUserID: intPtr(12345),
		},
		"with all optional params": {
			replyToMessageID:  stringPtr("message-id"),
			broadcasterUserID: intPtr(12345),
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			kickClient := setupMockClient(t, func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
				fmt.Fprint(w, `{
					"message":"success",
					"data":{"is_sent":true, "message_id":"message id"}
				}`)
			})

			resonse, err := kickClient.SendChatMessage(
				context.Background(),
				tc.broadcasterUserID,
				"message",
				tc.replyToMessageID,
				gokick.MessageTypeBot,
			)
			require.NoError(t, err)
			assert.True(t, resonse.Result.IsSent)
			assert.Equal(t, "message id", resonse.Result.MessageID)
		})
	}
}

func TestDeleteChatMessageError(t *testing.T) {
	t.Run("on new request", func(t *testing.T) {
		kickClient, err := gokick.NewClient(&gokick.ClientOptions{UserAccessToken: "access-token"})
		require.NoError(t, err)

		var ctx context.Context
		_, err = kickClient.DeleteChatMessage(ctx, "message-id")
		require.EqualError(t, err, "failed to create request: net/http: nil Context")
	})

	t.Run("timeout", func(t *testing.T) {
		kickClient := setupTimeoutMockClient(t)

		_, err := kickClient.DeleteChatMessage(context.Background(), "message-id")
		require.EqualError(t, err, `failed to make request: Delete "https://api.kick.com/public/v1/chat/message-id": `+
			`context deadline exceeded (Client.Timeout exceeded while awaiting headers)`)
	})

	t.Run("unmarshal error response", func(t *testing.T) {
		kickClient := setupMockClient(t, func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprint(w, `117`)
		})

		_, err := kickClient.DeleteChatMessage(context.Background(), "message-id")

		assert.EqualError(t, err, `failed to unmarshal error response (KICK status code: 500 and body "117"): json: cannot unmarshal `+
			`number into Go value of type gokick.errorResponse`)
	})

	t.Run("unmarshal token response", func(t *testing.T) {
		kickClient := setupMockClient(t, func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			fmt.Fprint(w, "117")
		})

		_, err := kickClient.DeleteChatMessage(context.Background(), "message-id")

		assert.EqualError(t, err, `failed to unmarshal error response (KICK status code: 200 and body "117"): json: cannot unmarshal `+
			`number into Go value of type gokick.errorResponse`)
	})

	t.Run("reader failure", func(t *testing.T) {
		kickClient := setupMockClient(t, func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Length", "10")
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprint(w, "")
		})

		_, err := kickClient.DeleteChatMessage(context.Background(), "message-id")

		assert.EqualError(t, err, `failed to read response body (KICK status code 500): unexpected EOF`)
	})

	t.Run("with internal server error", func(t *testing.T) {
		kickClient := setupMockClient(t, func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprint(w, `{"message":"internal server error", "data":null}`)
		})

		_, err := kickClient.DeleteChatMessage(context.Background(), "message-id")

		var kickError gokick.Error
		require.ErrorAs(t, err, &kickError)
		assert.Equal(t, http.StatusInternalServerError, kickError.Code())
		assert.Equal(t, "internal server error", kickError.Message())
	})
}

func TestDeleteChatMessageSuccess(t *testing.T) {
	t.Run("successful deletion", func(t *testing.T) {
		kickClient := setupMockClient(t, func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, http.MethodDelete, r.Method)
			assert.Equal(t, "/public/v1/chat/message-id", r.URL.Path)
			assert.Equal(t, "application/json", r.Header.Get("Content-Type"))
			assert.Equal(t, "Bearer access-token", r.Header.Get("Authorization"))

			w.WriteHeader(http.StatusNoContent)
		})

		_, err := kickClient.DeleteChatMessage(context.Background(), "message-id")
		require.NoError(t, err)
	})
}

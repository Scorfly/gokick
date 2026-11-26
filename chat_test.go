package gokick_test

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/scorfly/gokick"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestValidateChatMessageContent(t *testing.T) {
	t.Run("within grapheme limit", func(t *testing.T) {
		require.NoError(t, gokick.ValidateChatMessageContent(""))
		require.NoError(t, gokick.ValidateChatMessageContent(strings.Repeat("a", gokick.ChatMessageContentMaxGraphemes)))
		require.NoError(t, gokick.ValidateChatMessageContent(strings.Repeat("é", gokick.ChatMessageContentMaxGraphemes)))
	})

	t.Run("combining mark counts with the base character", func(t *testing.T) {
		cluster := "e\u0301"
		require.NoError(t, gokick.ValidateChatMessageContent(strings.Repeat(cluster, gokick.ChatMessageContentMaxGraphemes)))

		err := gokick.ValidateChatMessageContent(strings.Repeat(cluster, gokick.ChatMessageContentMaxGraphemes+1))
		require.ErrorIs(t, err, gokick.ErrChatMessageContentTooLong)
	})

	t.Run("over grapheme limit", func(t *testing.T) {
		err := gokick.ValidateChatMessageContent(strings.Repeat("a", gokick.ChatMessageContentMaxGraphemes+1))
		require.ErrorIs(t, err, gokick.ErrChatMessageContentTooLong)
	})

	t.Run("within byte limit", func(t *testing.T) {
		flag := "\U0001F1EB\U0001F1F7"
		require.NoError(t, gokick.ValidateChatMessageContent(strings.Repeat(flag, gokick.ChatMessageContentMaxBytes/len(flag))))
	})

	t.Run("over byte limit", func(t *testing.T) {
		// One family emoji is a single grapheme of several UTF-8 bytes, so the byte limit is reached first.
		family := "\U0001F468\u200D\U0001F469\u200D\U0001F467\u200D\U0001F466"
		fit := gokick.ChatMessageContentMaxBytes / len(family)
		require.Less(t, fit, gokick.ChatMessageContentMaxGraphemes)

		require.NoError(t, gokick.ValidateChatMessageContent(strings.Repeat(family, fit)))
		err := gokick.ValidateChatMessageContent(strings.Repeat(family, fit+1))
		require.ErrorIs(t, err, gokick.ErrChatMessageContentTooLong)
	})
}

func TestSendChatMessageError(t *testing.T) {
	t.Run("content too long", func(t *testing.T) {
		var called bool
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			called = true
			w.WriteHeader(http.StatusOK)
		}))
		t.Cleanup(server.Close)

		client, err := gokick.NewClient(&gokick.ClientOptions{
			UserAccessToken: "access-token",
			APIBaseURL:      fmt.Sprintf("http://%s", server.Listener.Addr()),
		})
		require.NoError(t, err)

		_, err = client.SendChatMessage(
			context.Background(),
			nil,
			strings.Repeat("x", gokick.ChatMessageContentMaxGraphemes+1),
			nil,
			gokick.MessageTypeBot,
		)
		require.ErrorIs(t, err, gokick.ErrChatMessageContentTooLong)
		assert.False(t, called)
	})

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

	t.Run("at max content length", func(t *testing.T) {
		kickClient := setupMockClient(t, func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			fmt.Fprint(w, `{
				"message":"success",
				"data":{"is_sent":true, "message_id":"message id"}
			}`)
		})

		msg := strings.Repeat("x", gokick.ChatMessageContentMaxGraphemes)
		response, err := kickClient.SendChatMessage(
			context.Background(),
			nil,
			msg,
			nil,
			gokick.MessageTypeBot,
		)
		require.NoError(t, err)
		assert.True(t, response.Result.IsSent)
	})
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

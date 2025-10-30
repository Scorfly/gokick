package gokick_test

import (
	"errors"
	"net/http"
	"strings"
	"testing"

	"github.com/scorfly/gokick"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetEventFromRequestError(t *testing.T) {
	t.Run("request not set", func(t *testing.T) {
		_, err := gokick.GetEventFromRequest(nil)
		require.EqualError(t, err, "request cannot be nil")
	})

	t.Run("invalid subscription name", func(t *testing.T) {
		req, err := http.NewRequest("GET", "https://domain.tld", strings.NewReader(""))
		require.NoError(t, err)
		req.Header.Set("X-Event-Subscription", "invalid")

		_, err = gokick.GetEventFromRequest(req)
		require.EqualError(t, err, "failed to parse subscription name: unknown name: invalid")
	})

	t.Run("invalid body", func(t *testing.T) {
		req, err := http.NewRequest("GET", "https://domain.tld", faultyReader{})
		require.NoError(t, err)
		req.Header.Set("X-Event-Subscription", "chat.message.sent")

		_, err = gokick.GetEventFromRequest(req)
		require.EqualError(t, err, "failed to read body: read error")
	})

	t.Run("cannot parse event", func(t *testing.T) {
		req, err := http.NewRequest("GET", "https://domain.tld", strings.NewReader(""))
		require.NoError(t, err)
		req.Header.Set("X-Event-Subscription", "chat.message.sent")

		_, err = gokick.GetEventFromRequest(req)
		require.EqualError(t, err, "failed to verify event validity: failed to verify signature: crypto/rsa: verification error")
	})
}

func TestGetEventFromRequestSuccess(t *testing.T) {
	skipSignatureValidation(t)

	req, err := http.NewRequest("GET", "https://domain.tld", strings.NewReader("{}"))
	require.NoError(t, err)
	req.Header.Set("X-Event-Subscription", "chat.message.sent")
	req.Header.Set("X-Event-Version", "1")
	req.Header.Set("X-Event-Signature", "signature")
	req.Header.Set("X-Event-Message-Id", "message ID")
	req.Header.Set("X-Event-Timestamp", "2025-02-21T23:23:36Z")

	event, err := gokick.GetEventFromRequest(req)
	require.NoError(t, err)
	assert.IsType(t, &gokick.ChatMessageEvent{}, event)
}

func TestValidateEventError(t *testing.T) {
	previousKey := gokick.DefaultEventPublicKey
	t.Cleanup(func() { gokick.DefaultEventPublicKey = previousKey })

	gokick.DefaultEventPublicKey = "invalid key"

	headers := http.Header{}
	headers.Set("Kick-Event-Message-Id", "msg123")

	valid := gokick.ValidateEvent(headers, []byte("body"))
	assert.False(t, valid)
}

func TestValidateEventSuccess(t *testing.T) {
	skipSignatureValidation(t)

	headers := http.Header{}
	headers.Set("Kick-Event-Message-Id", "msg123")

	valid := gokick.ValidateEvent(
		headers,
		[]byte(`{"message_id":"bb9832e4-e865-48f4-a0c3-392f78bf3b1a","broadcaster":{"is_anonymous":false,"user_id":721956,`+
			`"username":"Scorfly","is_verified":false,"profile_picture":"https://files.kick.com/images/user/721956/profile_image/`+
			`conversion/44a9f1fb-0498-47b5-820e-ef9399fd23d4-fullsize.webp","channel_slug":"scorfly"},"sender":{"is_anonymous":false,`+
			`"user_id":721956,"username":"Scorfly","is_verified":false,"profile_picture":"https://files.kick.com/images/user/721956/`+
			`profile_image/conversion/44a9f1fb-0498-47b5-820e-ef9399fd23d4-fullsize.webp","channel_slug":"scorfly"},`+
			`"content":"coucou","emotes":null}`),
	)
	assert.True(t, valid)
}

func TestValidateAndParseEventError(t *testing.T) {
	t.Run("failed to decode public key", func(t *testing.T) {
		previousKey := gokick.DefaultEventPublicKey
		t.Cleanup(func() { gokick.DefaultEventPublicKey = previousKey })

		gokick.DefaultEventPublicKey = "invalid key"
		_, err := gokick.ValidateAndParseEvent(
			gokick.SubscriptionNameChatMessage,
			"1",
			"signature",
			"message ID",
			"timestamp",
			"body",
		)
		require.EqualError(t, err, "failed to parse public key: failed to decode public key")
	})

	t.Run("key is not public", func(t *testing.T) {
		previousKey := gokick.DefaultEventPublicKey
		t.Cleanup(func() { gokick.DefaultEventPublicKey = previousKey })

		gokick.DefaultEventPublicKey = `-----BEGIN RSA PRIVATE KEY-----
MIIBOgIBAAJBAKj34GkxFhD90vcNLYLInFEX6Ppy1tPf9Cnzj4p4WGeKLs1Pt8Qu
KUpRKfFLfRYC9AIKjbJTWit+CqvjWYzvQwECAwEAAQJAIJLixBy2qpFoS4DSmoEm
o3qGy0t6z09AIJtH+5OeRV1be+N4cDYJKffGzDa88vQENZiRm0GRq6a+HPGQMd2k
TQIhAKMSvzIBnni7ot/OSie2TmJLY4SwTQAevXysE2RbFDYdAiEBCUEaRQnMnbp7
9mxDXDf6AU0cN/RPBjb9qSHDcWZHGzUCIG2Es59z8ugGrDY+pxLQnwfotadxd+Uy
v/Ow5T0q5gIJAiEAyS4RaI9YG8EWx/2w0T67ZUVAw8eOMB6BIUg0Xcu+3okCIBOs
/5OiPgoTdSy7bcF9IGpSE8ZgGKzgYQVZeN97YE00
-----END RSA PRIVATE KEY-----`
		_, err := gokick.ValidateAndParseEvent(
			gokick.SubscriptionNameChatMessage,
			"1",
			"signature",
			"message ID",
			"timestamp",
			"body",
		)
		require.EqualError(t, err, "failed to parse public key: not public key")
	})

	t.Run("failed to parse public key", func(t *testing.T) {
		previousKey := gokick.DefaultEventPublicKey
		t.Cleanup(func() { gokick.DefaultEventPublicKey = previousKey })

		gokick.DefaultEventPublicKey = `-----BEGIN PUBLIC KEY-----
-----END PUBLIC KEY-----`
		_, err := gokick.ValidateAndParseEvent(
			gokick.SubscriptionNameChatMessage,
			"1",
			"signature",
			"message ID",
			"timestamp",
			"body",
		)
		require.EqualError(t, err, "failed to parse public key: failed to parse public key: asn1: syntax error: sequence truncated")
	})

	t.Run("failed to ensure validy of signature", func(t *testing.T) {
		_, err := gokick.ValidateAndParseEvent(
			gokick.SubscriptionNameChatMessage,
			"1",
			"signature",
			"message ID",
			"timestamp",
			"body",
		)
		require.EqualError(t, err, "failed to verify event validity: failed to decode signature: illegal base64 data at input byte 8")
	})

	t.Run("failed to ensure verify PKCS1v15", func(t *testing.T) {
		_, err := gokick.ValidateAndParseEvent(
			gokick.SubscriptionNameChatMessage,
			"1",
			"b3NzMTE3",
			"message ID",
			"timestamp",
			"body",
		)
		require.EqualError(t, err, "failed to verify event validity: failed to verify signature: crypto/rsa: verification error")
	})

	t.Run("with invalid body", func(t *testing.T) {
		skipSignatureValidation(t)

		_, err := gokick.ValidateAndParseEvent(
			gokick.SubscriptionNameChatMessage,
			"2",
			"signature",
			"message ID",
			"2025-02-21T23:23:36Z",
			`invalid JSON`,
		)
		require.EqualError(t, err, "failed to unmarshal event: invalid character 'i' looking for beginning of value")
	})
}

func TestValidateAndParseEventSuccess(t *testing.T) {
	t.Run("with new chat message event detailed", func(t *testing.T) {
		skipSignatureValidation(t)

		event, err := gokick.ValidateAndParseEvent(
			gokick.SubscriptionNameChatMessage,
			"1",
			"signature",
			"message ID",
			"2025-02-21T23:23:36Z",
			`{"message_id":"bb9832e4-e865-48f4-a0c3-392f78bf3b1a","broadcaster":{"is_anonymous":false,"user_id":721956,`+
				`"username":"Scorfly","is_verified":false,"profile_picture":"https://files.kick.com/images/user/721956/profile_image/`+
				`conversion/44a9f1fb-0498-47b5-820e-ef9399fd23d4-fullsize.webp","channel_slug":"scorfly"},"sender":{"is_anonymous":false,`+
				`"user_id":721956,"username":"Scorfly","is_verified":false,"profile_picture":"https://files.kick.com/images/user/721956/`+
				`profile_image/conversion/44a9f1fb-0498-47b5-820e-ef9399fd23d4-fullsize.webp","channel_slug":"scorfly"},`+
				`"content":"coucou","emotes":null}`,
		)
		require.NoError(t, err)
		assert.IsType(t, &gokick.ChatMessageEvent{}, event)
		assert.Equal(t, "bb9832e4-e865-48f4-a0c3-392f78bf3b1a", event.(*gokick.ChatMessageEvent).MessageID)
		assert.Equal(t, "Scorfly", event.(*gokick.ChatMessageEvent).Broadcaster.Username)
		assert.Equal(t, "Scorfly", event.(*gokick.ChatMessageEvent).Sender.Username)
		assert.Equal(t, "coucou", event.(*gokick.ChatMessageEvent).Content)
		assert.Nil(t, event.(*gokick.ChatMessageEvent).Emotes)
	})

	t.Run("with new kicks gifted event detailed", func(t *testing.T) {
		skipSignatureValidation(t)

		event, err := gokick.ValidateAndParseEvent(
			gokick.SubscriptionNameKicksGifted,
			"1",
			"signature",
			"message ID",
			"2025-10-20T04:00:08.634Z",
			`{
				"broadcaster": {
					"user_id": 123456789,
					"username": "broadcaster_name",
					"is_verified": true,
					"profile_picture": "https://example.com/broadcaster_avatar.jpg",
					"channel_slug": "broadcaster_channel"
				},
				"sender": {
					"user_id": 987654321,
					"username": "gift_sender",
					"is_verified": false,
					"profile_picture": "https://example.com/sender_avatar.jpg",
					"channel_slug": "gift_sender_channel"
				},
				"gift": {
					"amount": 100,
					"name": "Full Send",
					"type": "BASIC",
					"tier": "BASIC",
					"message": "w"
				},
				"created_at": "2025-10-20T04:00:08.634Z"
			}`,
		)
		require.NoError(t, err)
		assert.IsType(t, &gokick.KicksGiftedEvent{}, event)

		kicksEvent := event.(*gokick.KicksGiftedEvent)
		assert.Equal(t, 123456789, kicksEvent.Broadcaster.UserID)
		assert.Equal(t, "broadcaster_name", kicksEvent.Broadcaster.Username)
		assert.True(t, kicksEvent.Broadcaster.IsVerified)
		assert.Equal(t, "broadcaster_channel", kicksEvent.Broadcaster.ChannelSlug)

		assert.Equal(t, 987654321, kicksEvent.Sender.UserID)
		assert.Equal(t, "gift_sender", kicksEvent.Sender.Username)
		assert.False(t, kicksEvent.Sender.IsVerified)
		assert.Equal(t, "gift_sender_channel", kicksEvent.Sender.ChannelSlug)

		assert.Equal(t, 100, kicksEvent.Gift.Amount)
		assert.Equal(t, "Full Send", kicksEvent.Gift.Name)
		assert.Equal(t, "BASIC", kicksEvent.Gift.Type)
		assert.Equal(t, "BASIC", kicksEvent.Gift.Tier)
		assert.Equal(t, "w", kicksEvent.Gift.Message)
		assert.Equal(t, "2025-10-20T04:00:08.634Z", kicksEvent.CreatedAt)
	})

	t.Run("with new chat message event details with unexisting version", func(t *testing.T) {
		skipSignatureValidation(t)

		event, err := gokick.ValidateAndParseEvent(
			gokick.SubscriptionNameChatMessage,
			"-1",
			"signature",
			"message ID",
			"2025-02-21T23:23:36Z",
			`{"message_id":"bb9832e4-e865-48f4-a0c3-392f78bf3b1a","broadcaster":{"is_anonymous":false,"user_id":721956,`+
				`"username":"Scorfly","is_verified":false,"profile_picture":"https://files.kick.com/images/user/721956/profile_image/`+
				`conversion/44a9f1fb-0498-47b5-820e-ef9399fd23d4-fullsize.webp","channel_slug":"scorfly"},"sender":{"is_anonymous":false,`+
				`"user_id":721956,"username":"Scorfly","is_verified":false,"profile_picture":"https://files.kick.com/images/user/721956/`+
				`profile_image/conversion/44a9f1fb-0498-47b5-820e-ef9399fd23d4-fullsize.webp","channel_slug":"scorfly"},"content":"coucou",`+
				`"emotes":null}`,
		)
		require.NoError(t, err)
		assert.IsType(t, map[string]interface{}{}, event)
		assert.Equal(t, "bb9832e4-e865-48f4-a0c3-392f78bf3b1a", event.(map[string]interface{})["message_id"])
		assert.Equal(t, "coucou", event.(map[string]interface{})["content"])
	})

	t.Run("all events with version", func(t *testing.T) {
		testCases := map[string]struct {
			subscription gokick.SubscriptionName
			version      string
			expectedType interface{}
		}{
			"with new chat message version 1": {
				subscription: gokick.SubscriptionNameChatMessage,
				version:      "1",
				expectedType: &gokick.ChatMessageEvent{},
			},
			"with new follower version 1": {
				subscription: gokick.SubscriptionNameChannelFollow,
				version:      "1",
				expectedType: &gokick.ChannelFollowEvent{},
			},
			"with new subscription renewal version 1": {
				subscription: gokick.SubscriptionNameChannelSubscriptionRenewal,
				version:      "1",
				expectedType: &gokick.ChannelSubscriptionRenewalEvent{},
			},
			"with new subscription gifts version 1": {
				subscription: gokick.SubscriptionNameChannelSubscriptionGifts,
				version:      "1",
				expectedType: &gokick.ChannelSubscriptionGiftsEvent{},
			},
			"with new subscription created version 1": {
				subscription: gokick.SubscriptionNameChannelSubscriptionCreated,
				version:      "1",
				expectedType: &gokick.ChannelSubscriptionCreatedEvent{},
			},
			"with new livestream status updated version 1": {
				subscription: gokick.SubscriptionNameLivestreamStatusUpdated,
				version:      "1",
				expectedType: &gokick.LivestreamStatusUpdatedEvent{},
			},
			"with new livestream metadata updated version 1": {
				subscription: gokick.SubscriptionNameLivestreamMetadataUpdated,
				version:      "1",
				expectedType: &gokick.LivestreamMetadataUpdatedEvent{},
			},
			"with new banned user updated version 1": {
				subscription: gokick.SubscriptionNameModerationBanned,
				version:      "1",
				expectedType: &gokick.ModerationBannedEvent{},
			},
			"with new kicks gifted version 1": {
				subscription: gokick.SubscriptionNameKicksGifted,
				version:      "1",
				expectedType: &gokick.KicksGiftedEvent{},
			},
		}

		for name, testCase := range testCases {
			t.Run(name, func(t *testing.T) {
				skipSignatureValidation(t)

				event, err := gokick.ValidateAndParseEvent(
					testCase.subscription,
					testCase.version,
					"signature",
					"message ID",
					"2025-02-21T23:23:36Z",
					`{}`,
				)
				require.NoError(t, err)
				assert.IsType(t, testCase.expectedType, event)
			})
		}
	})
}

type faultyReader struct{}

func (r faultyReader) Read(p []byte) (n int, err error) {
	return 0, errors.New("read error")
}

func skipSignatureValidation(t *testing.T) {
	t.Helper()

	previousKey := gokick.SkipSignatureValidation
	t.Cleanup(func() { gokick.SkipSignatureValidation = previousKey })
	gokick.SkipSignatureValidation = true
}

package gokick

import (
	"bytes"
	"crypto"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"encoding/pem"
	"errors"
	"fmt"
	"io"
	"net/http"
)

type Badge struct {
	Text  string `json:"text"`
	Type  string `json:"type"`
	Count int    `json:"count"`
}

type IdentityEvent struct {
	UsernameColor string  `json:"username_color"`
	Badges        []Badge `json:"badges"`
}

type UserEvent struct {
	IsAnonymous    bool          `json:"is_anonymous"`
	UserID         int           `json:"user_id"`
	Username       string        `json:"username"`
	IsVerified     bool          `json:"is_verified"`
	ProfilePicture string        `json:"profile_picture"`
	ChannelSlug    string        `json:"channel_slug"`
	Identity       IdentityEvent `json:"identity"`
}

type ChatMessageEmotesEvent struct {
	EmoteID   int `json:"emote_id"`
	Positions []struct {
		Start int `json:"s"`
		End   int `json:"e"`
	} `json:"positions"`
}

type ChatMessageEvent struct {
	MessageID string `json:"message_id"`
	RepliesTo struct {
		MessageID string    `json:"message_id"`
		Sender    UserEvent `json:"sender"`
		Content   string    `json:"content"`
	} `json:"replies_to"`
	Broadcaster UserEvent                `json:"broadcaster"`
	Sender      UserEvent                `json:"sender"`
	Content     string                   `json:"content"`
	Emotes      []ChatMessageEmotesEvent `json:"emotes"`
	CreatedAt   string                   `json:"created_at"`
}

type ChannelFollowEvent struct {
	Broadcaster UserEvent `json:"broadcaster"`
	Follower    UserEvent `json:"follower"`
}

type ChannelSubscriptionRenewalEvent struct {
	Broadcaster UserEvent `json:"broadcaster"`
	Subscriber  UserEvent `json:"subscriber"`
	Duration    int       `json:"duration"`
	CreatedAt   string    `json:"created_at"`
	ExpiresAt   string    `json:"expires_at"`
}

type ChannelSubscriptionGiftsEvent struct {
	Broadcaster UserEvent   `json:"broadcaster"`
	Gifter      UserEvent   `json:"gifter"`
	Giftees     []UserEvent `json:"giftees"`
	CreatedAt   string      `json:"created_at"`
	ExpiresAt   string      `json:"expires_at"`
}

type ChannelSubscriptionCreatedEvent struct {
	Broadcaster UserEvent `json:"broadcaster"`
	Subscriber  UserEvent `json:"subscriber"`
	Duration    int       `json:"duration"`
	CreatedAt   string    `json:"created_at"`
	ExpiresAt   string    `json:"expires_at"`
}

type LivestreamStatusUpdatedEvent struct {
	Broadcaster UserEvent `json:"broadcaster"`
	IsLive      bool      `json:"is_live"`
	Title       string    `json:"title"`
	StartedAt   string    `json:"started_at"`
	EndedAt     string    `json:"ended_at"`
}

type LivestreamMetadataUpdatedEvent struct {
	Broadcaster UserEvent `json:"broadcaster"`
	Metadata    struct {
		Title            string `json:"title"`
		Language         string `json:"language"`
		HasMatureContent bool   `json:"has_mature_content"`
		Category         struct {
			ID        string `json:"id"`
			Name      string `json:"name"`
			Thumbnail string `json:"thumbnail"`
		} `json:"category"`
	} `json:"metadata"`
}

type ModerationBannedEvent struct {
	Broadcaster UserEvent `json:"broadcaster"`
	Moderator   UserEvent `json:"moderator"`
	BannedUser  UserEvent `json:"banned_user"`
	Metadata    struct {
		Reason    string `json:"reason"`
		CreatedAt string `json:"created_at"`
		ExpiresAt string `json:"expires_at"`
	} `json:"reason"`
}

type KicksGiftedEvent struct {
	Broadcaster UserEvent `json:"broadcaster"`
	Sender      UserEvent `json:"sender"`
	Gift        struct {
		Amount  int    `json:"amount"`
		Name    string `json:"name"`
		Type    string `json:"type"`
		Tier    string `json:"tier"`
		Message string `json:"message"`
	} `json:"gift"`
	CreatedAt string `json:"created_at"`
}

// I set it as public to be able to change it in tests.
// It's not a good practice to do so, but it's the only way to do it for now.
var DefaultEventPublicKey = `-----BEGIN PUBLIC KEY-----
MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEAq/+l1WnlRrGSolDMA+A8
6rAhMbQGmQ2SapVcGM3zq8ANXjnhDWocMqfWcTd95btDydITa10kDvHzw9WQOqp2
MZI7ZyrfzJuz5nhTPCiJwTwnEtWft7nV14BYRDHvlfqPUaZ+1KR4OCaO/wWIk/rQ
L/TjY0M70gse8rlBkbo2a8rKhu69RQTRsoaf4DVhDPEeSeI5jVrRDGAMGL3cGuyY
6CLKGdjVEM78g3JfYOvDU/RvfqD7L89TZ3iN94jrmWdGz34JNlEI5hqK8dd7C5EF
BEbZ5jgB8s8ReQV8H+MkuffjdAj3ajDDX3DOJMIut1lBrUVD1AaSrGCKHooWoL2e
twIDAQAB
-----END PUBLIC KEY-----`

// I set it to be able to tests without having to sign the events.
// It's not a good practice to do so, but it's the only way to do it for now.
// Do not override it in production !
var SkipSignatureValidation = false

func GetEventFromRequest(request *http.Request) (interface{}, error) {
	if request == nil {
		return nil, errors.New("request cannot be nil")
	}

	subscriptionName, err := NewSubscriptionName(request.Header.Get("X-Event-Subscription"))
	if err != nil {
		return nil, fmt.Errorf("failed to parse subscription name: %v", err)
	}

	version := request.Header.Get("X-Event-Version")
	eventSignature := request.Header.Get("X-Event-Signature")
	messageID := request.Header.Get("X-Event-Message-Id")
	timestamp := request.Header.Get("X-Event-Timestamp")

	body, err := io.ReadAll(request.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read body: %v", err)
	}

	return ValidateAndParseEvent(
		subscriptionName,
		version,
		eventSignature,
		messageID,
		timestamp,
		string(body),
	)
}

func ValidateEvent(
	header http.Header,
	body []byte,
) bool {
	if !SkipSignatureValidation {
		payload := struct {
			signature string
			messageID string
			timestamp string
		}{
			signature: header.Get("Kick-Event-Signature"),
			messageID: header.Get("Kick-Event-Message-Id"),
			timestamp: header.Get("Kick-Event-Message-Timestamp"),
		}

		signature := bytes.Join([][]byte{
			[]byte(payload.messageID),
			[]byte(payload.timestamp),
			body,
		}, []byte("."))

		publicKey, _ := parsePublicKey([]byte(DefaultEventPublicKey))

		err := verifyEventValidity(&publicKey, signature, []byte(payload.signature))
		if err != nil {
			return false
		}
	}

	return true
}

func ValidateAndParseEvent(
	subscriptionName SubscriptionName,
	version string,
	eventSignature string,
	messageID string,
	timestamp string,
	body string,
) (interface{}, error) {
	if !SkipSignatureValidation {
		signature := []byte(fmt.Sprintf("%s.%s.%s", messageID, timestamp, body))

		publicKey, err := parsePublicKey([]byte(DefaultEventPublicKey))
		if err != nil {
			return nil, fmt.Errorf("failed to parse public key: %v", err)
		}

		err = verifyEventValidity(&publicKey, signature, []byte(eventSignature))
		if err != nil {
			return nil, fmt.Errorf("failed to verify event validity: %v", err)
		}
	}

	var event interface{}
	if versionConstructor, ok := eventConstructors[subscriptionName]; ok {
		if constructor, ok := versionConstructor[version]; ok {
			event = constructor()
		}
	}

	err := json.Unmarshal([]byte(body), &event)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal event: %v", err)
	}

	return event, nil
}

func parsePublicKey(key []byte) (rsa.PublicKey, error) {
	block, _ := pem.Decode(key)
	if block == nil {
		return rsa.PublicKey{}, errors.New("failed to decode public key")
	}

	if block.Type != "PUBLIC KEY" {
		return rsa.PublicKey{}, errors.New("not public key")
	}

	parsed, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return rsa.PublicKey{}, fmt.Errorf("failed to parse public key: %v", err)
	}

	publicKey, ok := parsed.(*rsa.PublicKey)
	if !ok {
		return rsa.PublicKey{}, errors.New("not expected public key interface")
	}

	return *publicKey, nil
}

func verifyEventValidity(publicKey *rsa.PublicKey, body []byte, signature []byte) error {
	decoded := make([]byte, base64.StdEncoding.DecodedLen(len(signature)))

	n, err := base64.StdEncoding.Decode(decoded, signature)
	if err != nil {
		return fmt.Errorf("failed to decode signature: %v", err)
	}

	signature = decoded[:n]
	hashed := sha256.Sum256(body)

	err = rsa.VerifyPKCS1v15(publicKey, crypto.SHA256, hashed[:], signature)
	if err != nil {
		return fmt.Errorf("failed to verify signature: %v", err)
	}

	return nil
}

type eventConstructor func() interface{}

var eventConstructors = map[SubscriptionName]map[string]eventConstructor{
	SubscriptionNameChatMessage: {
		"1": func() interface{} { return new(ChatMessageEvent) },
	},
	SubscriptionNameChannelFollow: {
		"1": func() interface{} { return new(ChannelFollowEvent) },
	},
	SubscriptionNameChannelSubscriptionRenewal: {
		"1": func() interface{} { return new(ChannelSubscriptionRenewalEvent) },
	},
	SubscriptionNameChannelSubscriptionGifts: {
		"1": func() interface{} { return new(ChannelSubscriptionGiftsEvent) },
	},
	SubscriptionNameChannelSubscriptionCreated: {
		"1": func() interface{} { return new(ChannelSubscriptionCreatedEvent) },
	},
	SubscriptionNameLivestreamStatusUpdated: {
		"1": func() interface{} { return new(LivestreamStatusUpdatedEvent) },
	},
	SubscriptionNameLivestreamMetadataUpdated: {
		"1": func() interface{} { return new(LivestreamMetadataUpdatedEvent) },
	},
	SubscriptionNameModerationBanned: {
		"1": func() interface{} { return new(ModerationBannedEvent) },
	},
	SubscriptionNameKicksGifted: {
		"1": func() interface{} { return new(KicksGiftedEvent) },
	},
}

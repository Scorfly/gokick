package gokick

import (
	"crypto"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"encoding/pem"
	"errors"
	"fmt"
)

type BroadcasterEvent struct {
	IsAnonymous    bool   `json:"is_anonymous"`
	UserID         int    `json:"user_id"`
	Username       string `json:"username"`
	IsVerified     bool   `json:"is_verified"`
	ProfilePicture string `json:"profile_picture"`
	ChannelSlug    string `json:"channel_slug"`
}

type ChatMessageEmotesEvent struct {
	EmoteID   int `json:"emote_id"`
	Positions []struct {
		Start int `json:"s"`
		End   int `json:"e"`
	} `json:"positions"`
}

type ChatMessageEvent struct {
	MessageID   string                   `json:"message_id"`
	Broadcaster BroadcasterEvent         `json:"broadcaster"`
	Sender      BroadcasterEvent         `json:"sender"`
	Content     string                   `json:"content"`
	Emotes      []ChatMessageEmotesEvent `json:"emotes"`
}

type ChannelFollowEvent struct {
	Broadcaster BroadcasterEvent `json:"broadcaster"`
	Follower    BroadcasterEvent `json:"follower"`
}

type ChannelSubscriptionRenewalEvent struct {
	Broadcaster BroadcasterEvent `json:"broadcaster"`
	Subscriber  BroadcasterEvent `json:"subscriber"`
	Duration    int              `json:"duration"`
	CreatedAt   string           `json:"created_at"`
}

type ChannelSubscriptionGiftsEvent struct {
	Broadcaster BroadcasterEvent   `json:"broadcaster"`
	Gifter      BroadcasterEvent   `json:"gifter"`
	Giftees     []BroadcasterEvent `json:"giftees"`
	CreatedAt   string             `json:"created_at"`
}

type ChannelSubscriptionCreatedEvent struct {
	Broadcaster BroadcasterEvent `json:"broadcaster"`
	Subscriber  BroadcasterEvent `json:"subscriber"`
	Duration    int              `json:"duration"`
	CreatedAt   string           `json:"created_at"`
}

// I set it as public to be able to change it in tests.
// It’s not a good practice to do so, but it’s the only way to do it for now.
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
// It’s not a good practice to do so, but it’s the only way to do it for now.
// Do not override it in production !
var SkipSignatureValidation = false

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
}

package gokick

import (
	"context"
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
	"sync"
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
	IsAnonymous    bool           `json:"is_anonymous"`
	UserID         int            `json:"user_id"`
	Username       string         `json:"username"`
	IsVerified     bool           `json:"is_verified"`
	ProfilePicture string         `json:"profile_picture"`
	ChannelSlug    string         `json:"channel_slug"`
	Identity       *IdentityEvent `json:"identity"`
}

type ChatMessageEmotesEvent struct {
	EmoteID   json.Number `json:"emote_id"`
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
			ID        int    `json:"id"`
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
	} `json:"metadata"`
}

type KicksGiftedEvent struct {
	Broadcaster UserEvent `json:"broadcaster"`
	Sender      UserEvent `json:"sender"`
	Gift        struct {
		Amount            int    `json:"amount"`
		Name              string `json:"name"`
		Type              string `json:"type"`
		Tier              string `json:"tier"`
		Message           string `json:"message"`
		PinnedTimeSeconds int    `json:"pinned_time_seconds"`
	} `json:"gift"`
	CreatedAt string `json:"created_at"`
}

type ChannelRewardRedemptionUpdatedEvent struct {
	ID         string `json:"id"`
	UserInput  string `json:"user_input"`
	Status     string `json:"status"`
	RedeemedAt string `json:"redeemed_at"`
	Reward     struct {
		ID          string `json:"id"`
		Title       string `json:"title"`
		Cost        int    `json:"cost"`
		Description string `json:"description"`
	} `json:"reward"`
	Redeemer    UserEvent `json:"redeemer"`
	Broadcaster UserEvent `json:"broadcaster"`
}

// DefaultEventPublicKey overrides the webhook public key when it is not empty.
// Leave it empty so verification loads the current key from GET /public/v1/public-key.
// Kick can rotate that key; a hardcoded copy fails after a rotation.
// Tests set this to avoid calling the API.
var DefaultEventPublicKey string

// EventPublicKeyBaseURL overrides the API origin used to fetch the webhook public key.
// Empty uses https://api.kick.com.
var EventPublicKeyBaseURL string

var (
	fetchedEventPublicKey   string
	fetchedEventPublicKeyMu sync.Mutex
)

// ClearFetchedEventPublicKey drops the cached webhook public key.
// The next verification loads a new key from GET /public/v1/public-key.
func ClearFetchedEventPublicKey() {
	fetchedEventPublicKeyMu.Lock()
	fetchedEventPublicKey = ""
	fetchedEventPublicKeyMu.Unlock()
}

// I set it to be able to tests without having to sign the events.
// It's not a good practice to do so, but it's the only way to do it for now.
// Do not override it in production !
var SkipSignatureValidation = false

// webhookRequestMeta reads Kick webhook headers per https://docs.kick.com/events/event-types.md,
// with fallback to legacy X-Event-* names used by older examples.
func webhookRequestMeta(h http.Header) (eventName, version, signature, messageID, timestamp string) {
	eventName = h.Get("Kick-Event-Type")
	if eventName == "" {
		eventName = h.Get("X-Event-Subscription")
	}
	version = h.Get("Kick-Event-Version")
	if version == "" {
		version = h.Get("X-Event-Version")
	}
	signature = h.Get("Kick-Event-Signature")
	if signature == "" {
		signature = h.Get("X-Event-Signature")
	}
	messageID = h.Get("Kick-Event-Message-Id")
	if messageID == "" {
		messageID = h.Get("X-Event-Message-Id")
	}
	timestamp = h.Get("Kick-Event-Message-Timestamp")
	if timestamp == "" {
		timestamp = h.Get("X-Event-Timestamp")
	}
	return eventName, version, signature, messageID, timestamp
}

func GetEventFromRequest(request *http.Request) (interface{}, error) {
	if request == nil {
		return nil, errors.New("request cannot be nil")
	}

	eventName, version, eventSignature, messageID, timestamp := webhookRequestMeta(request.Header)

	subscriptionName, err := NewSubscriptionName(eventName)
	if err != nil {
		return nil, fmt.Errorf("failed to parse subscription name: %v", err)
	}

	body, err := io.ReadAll(request.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read body: %v", err)
	}

	return validateAndParseEventWithContext(
		request.Context(),
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
		_, _, sig, msgID, ts := webhookRequestMeta(header)
		signed := []byte(fmt.Sprintf("%s.%s.%s", msgID, ts, body))
		err := checkEventSignature(context.Background(), signed, []byte(sig))
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
	return validateAndParseEventWithContext(
		context.Background(),
		subscriptionName,
		version,
		eventSignature,
		messageID,
		timestamp,
		body,
	)
}

func validateAndParseEventWithContext(
	ctx context.Context,
	subscriptionName SubscriptionName,
	version string,
	eventSignature string,
	messageID string,
	timestamp string,
	body string,
) (interface{}, error) {
	if !SkipSignatureValidation {
		signed := []byte(fmt.Sprintf("%s.%s.%s", messageID, timestamp, body))
		err := checkEventSignature(ctx, signed, []byte(eventSignature))
		if err != nil {
			return nil, wrapEventSignatureError(err, "failed to verify event validity")
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
		return signatureMismatchError{err: err}
	}

	return nil
}

type signatureMismatchError struct {
	err error
}

func (e signatureMismatchError) Error() string {
	return fmt.Sprintf("failed to verify signature: %v", e.err)
}

func (e signatureMismatchError) Unwrap() error {
	return e.err
}

type eventKeyError struct {
	op  string
	err error
}

func (e *eventKeyError) Error() string {
	return e.err.Error()
}

func (e *eventKeyError) Unwrap() error {
	return e.err
}

func checkEventSignature(ctx context.Context, signed, signature []byte) error {
	pemKey, err := eventPublicKeyPEM(ctx, false)
	if err != nil {
		return err
	}

	verifyErr := verifyPEM(pemKey, signed, signature)
	if verifyErr == nil || DefaultEventPublicKey != "" || !isSignatureMismatch(verifyErr) {
		return verifyErr
	}

	refreshed, refreshErr := eventPublicKeyPEM(ctx, true)
	if refreshErr != nil || refreshed == pemKey {
		return verifyErr
	}

	return verifyPEM(refreshed, signed, signature)
}

func isSignatureMismatch(err error) bool {
	var mismatch signatureMismatchError
	return errors.As(err, &mismatch)
}

func verifyPEM(pemKey string, signed, signature []byte) error {
	publicKey, err := parsePublicKey([]byte(pemKey))
	if err != nil {
		return &eventKeyError{op: "parse", err: err}
	}

	return verifyEventValidity(&publicKey, signed, signature)
}

func eventPublicKeyPEM(ctx context.Context, refresh bool) (string, error) {
	if DefaultEventPublicKey != "" {
		return DefaultEventPublicKey, nil
	}

	fetchedEventPublicKeyMu.Lock()
	defer fetchedEventPublicKeyMu.Unlock()

	if !refresh && fetchedEventPublicKey != "" {
		return fetchedEventPublicKey, nil
	}

	pemKey, err := fetchEventPublicKey(ctx)
	if err != nil {
		return "", &eventKeyError{op: "fetch", err: err}
	}

	fetchedEventPublicKey = pemKey
	return pemKey, nil
}

func fetchEventPublicKey(ctx context.Context) (string, error) {
	client, err := NewClient(&ClientOptions{APIBaseURL: EventPublicKeyBaseURL})
	if err != nil {
		return "", err
	}

	response, err := client.GetPublicKey(ctx)
	if err != nil {
		return "", err
	}

	if response.Result.PublicKey == "" {
		return "", errors.New("public key is empty")
	}

	return response.Result.PublicKey, nil
}

func wrapEventSignatureError(err error, verifyPrefix string) error {
	var keyErr *eventKeyError
	if errors.As(err, &keyErr) {
		switch keyErr.op {
		case "parse":
			return fmt.Errorf("failed to parse public key: %v", keyErr.err)
		case "fetch":
			return fmt.Errorf("failed to fetch public key: %v", keyErr.err)
		}
	}

	return fmt.Errorf("%s: %v", verifyPrefix, err)
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
	SubscriptionNameChannelRewardRedemptionUpdated: {
		"1": func() interface{} { return new(ChannelRewardRedemptionUpdatedEvent) },
	},
}

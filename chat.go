package gokick

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/rivo/uniseg"
)

// ChatMessageContentMaxGraphemes is the maximum number of user-perceived characters
// Kick accepts in a chat message. Combining marks belong to the preceding character.
const ChatMessageContentMaxGraphemes = 500

// ChatMessageContentMaxBytes is the maximum UTF-8 size Kick accepts for chat message content.
const ChatMessageContentMaxBytes = 2048

// ChatMessageContentMaxRunes is the grapheme-cluster limit.
//
// Deprecated: Kick counts user-perceived characters, not Unicode code points.
// Use ChatMessageContentMaxGraphemes.
const ChatMessageContentMaxRunes = ChatMessageContentMaxGraphemes

// ErrChatMessageContentTooLong is returned when content exceeds ChatMessageContentMaxGraphemes
// grapheme clusters or ChatMessageContentMaxBytes UTF-8 bytes.
var ErrChatMessageContentTooLong = errors.New(
	"gokick: chat message content exceeds 500 grapheme clusters or 2048 UTF-8 bytes",
)

// ValidateChatMessageContent returns ErrChatMessageContentTooLong when content exceeds
// Kick's grapheme-cluster or UTF-8 byte limit.
func ValidateChatMessageContent(content string) error {
	if len(content) > ChatMessageContentMaxBytes {
		return ErrChatMessageContentTooLong
	}
	if uniseg.GraphemeClusterCount(content) > ChatMessageContentMaxGraphemes {
		return ErrChatMessageContentTooLong
	}
	return nil
}

type ChatResponseWrapper Response[ChatResponse]

type ChatResponse struct {
	IsSent    bool   `json:"is_sent"`
	MessageID string `json:"message_id"`
}

func (c *Client) SendChatMessage(
	ctx context.Context,
	broadcasterUserID *int,
	content string,
	replyToMessageID *string,
	messageType MessageType,
) (ChatResponseWrapper, error) {
	err := ValidateChatMessageContent(content)
	if err != nil {
		return ChatResponseWrapper{}, err
	}

	type postBodyRequest struct {
		BroadcasterUserID int    `json:"broadcaster_user_id,omitempty"`
		Content           string `json:"content"`
		ReplyToMessageID  string `json:"reply_to_message_id,omitempty"`
		Type              string `json:"type"`
	}

	r := postBodyRequest{
		Content: content,
		Type:    messageType.String(),
	}

	if replyToMessageID != nil {
		r.ReplyToMessageID = *replyToMessageID
	}

	if broadcasterUserID != nil {
		r.BroadcasterUserID = *broadcasterUserID
	}

	body, err := json.Marshal(r)
	if err != nil {
		return ChatResponseWrapper{}, fmt.Errorf("failed to marshal body: %v", err)
	}

	response, err := makeRequest[ChatResponse](
		ctx,
		c,
		http.MethodPost,
		"/public/v1/chat",
		http.StatusOK,
		bytes.NewReader(body),
	)
	if err != nil {
		return ChatResponseWrapper{}, err
	}

	return ChatResponseWrapper(response), nil
}

func (c *Client) DeleteChatMessage(ctx context.Context, messageID string) (EmptyResponse, error) {
	_, err := makeRequest[EmptyResponse](
		ctx,
		c,
		http.MethodDelete,
		fmt.Sprintf("/public/v1/chat/%s", messageID),
		http.StatusNoContent,
		http.NoBody,
	)
	if err != nil {
		return EmptyResponse{}, err
	}

	return EmptyResponse{}, nil
}

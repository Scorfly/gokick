package gokick

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

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

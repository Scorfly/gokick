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
	broadcasterUserID int,
	content string,
	messageType MessageType,
) (ChatResponseWrapper, error) {
	type postBodyRequest struct {
		BroadcasterUserID int    `json:"broadcaster_user_id"`
		Content           string `json:"content"`
		Type              string `json:"type"`
	}

	r := postBodyRequest{
		BroadcasterUserID: broadcasterUserID,
		Content:           content,
		Type:              messageType.String(),
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

package gokick

import (
	"fmt"
)

type MessageType int

const (
	MessageTypeUser MessageType = iota // user
	MessageTypeBot                     // bot
)

func NewMessageType(messageType string) (MessageType, error) {
	switch messageType {
	case "user":
		return MessageTypeUser, nil
	case "bot":
		return MessageTypeBot, nil
	default:
		return 0, fmt.Errorf("unknown message type: %s", messageType)
	}
}

func (m MessageType) String() string {
	switch m {
	case MessageTypeUser:
		return "user"
	case MessageTypeBot:
		return "bot"
	default:
		return "unknown"
	}
}

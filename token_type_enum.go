package gokick

import "fmt"

type TokenType int

const (
	TokenTypeAccess  TokenType = iota // access_token
	TokenTypeRefresh                  // refresh_token
)

func NewTokenType(tokenType string) (TokenType, error) {
	switch tokenType {
	case "access_token":
		return TokenTypeAccess, nil
	case "refresh_token":
		return TokenTypeRefresh, nil
	default:
		return 0, fmt.Errorf("unknown token type: %s", tokenType)
	}
}

func (s TokenType) String() string {
	switch s {
	case TokenTypeAccess:
		return "access_token"
	case TokenTypeRefresh:
		return "refresh_token"
	default:
		return "unknown"
	}
}

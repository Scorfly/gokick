package gokick

import "fmt"

type Scope int

const (
	ScopeUserRead       Scope = iota // user:read
	ScopeChannelRead                 // channel:read
	ScopeChannelWrite                // channel:write
	ScopeChatWrite                   // chat:write
	ScopeStremkeyRead                // streamkey:read
	ScopeEventSubscribe              // events:subscribe
	ScopeModerationBan               // moderation:ban
	ScopeKicksRead                   // kicks:read
)

func NewScope(scope string) (Scope, error) {
	switch scope {
	case "user:read":
		return ScopeUserRead, nil
	case "channel:read":
		return ScopeChannelRead, nil
	case "channel:write":
		return ScopeChannelWrite, nil
	case "chat:write":
		return ScopeChatWrite, nil
	case "streamkey:read":
		return ScopeStremkeyRead, nil
	case "events:subscribe":
		return ScopeEventSubscribe, nil
	case "moderation:ban":
		return ScopeModerationBan, nil
	case "kicks:read":
		return ScopeKicksRead, nil
	default:
		return 0, fmt.Errorf("unknown scope: %s", scope)
	}
}

func (s Scope) String() string {
	switch s {
	case ScopeUserRead:
		return "user:read"
	case ScopeChannelRead:
		return "channel:read"
	case ScopeChannelWrite:
		return "channel:write"
	case ScopeChatWrite:
		return "chat:write"
	case ScopeStremkeyRead:
		return "streamkey:read"
	case ScopeEventSubscribe:
		return "events:subscribe"
	case ScopeModerationBan:
		return "moderation:ban"
	case ScopeKicksRead:
		return "kicks:read"
	default:
		return "unknown"
	}
}

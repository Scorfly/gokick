package gokick

import "fmt"

type Scope int

const (
	ScopeUserRead                    Scope = iota // user:read
	ScopeChannelRead                              // channel:read
	ScopeChannelWrite                             // channel:write
	ScopeChannelRewardsRead                       // channel:rewards:read
	ScopeChannelRewardsWrite                      // channel:rewards:write
	ScopeChatWrite                                // chat:write
	ScopeStreamkeyRead                            // streamkey:read
	ScopeEventSubscribe                           // events:subscribe
	ScopeModerationBan                            // moderation:ban
	ScopeModerationChatMessageManage              // moderation:chat_message:manage
	ScopeKicksRead                                // kicks:read
)

func NewScope(scope string) (Scope, error) {
	switch scope {
	case "user:read":
		return ScopeUserRead, nil
	case "channel:read":
		return ScopeChannelRead, nil
	case "channel:write":
		return ScopeChannelWrite, nil
	case "channel:rewards:read":
		return ScopeChannelRewardsRead, nil
	case "channel:rewards:write":
		return ScopeChannelRewardsWrite, nil
	case "chat:write":
		return ScopeChatWrite, nil
	case "streamkey:read":
		return ScopeStreamkeyRead, nil
	case "events:subscribe":
		return ScopeEventSubscribe, nil
	case "moderation:ban":
		return ScopeModerationBan, nil
	case "moderation:chat_message:manage":
		return ScopeModerationChatMessageManage, nil
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
	case ScopeChannelRewardsRead:
		return "channel:rewards:read"
	case ScopeChannelRewardsWrite:
		return "channel:rewards:write"
	case ScopeChatWrite:
		return "chat:write"
	case ScopeStreamkeyRead:
		return "streamkey:read"
	case ScopeEventSubscribe:
		return "events:subscribe"
	case ScopeModerationBan:
		return "moderation:ban"
	case ScopeModerationChatMessageManage:
		return "moderation:chat_message:manage"
	case ScopeKicksRead:
		return "kicks:read"
	default:
		return "unknown"
	}
}

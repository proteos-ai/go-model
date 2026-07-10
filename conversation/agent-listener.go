package conversationmodel

import (
	"time"

	"go.proteos.ai/model/common"
)

// AgentListener is a rule binding inbound messages to an agent: "when a message
// matching <trigger> arrives on <connection|conversation>, drive agent
// <agent_key> and send its reply back". Exactly one of ConnectionId /
// ConversationId is set (DB CHECK): connection-bound listeners cover a whole
// integration (a Slack workspace), conversation-bound ones a single thread (the
// meeting-companion case). The dispatcher acts as ActingUser — that user needs
// agent-sessions:write + messages:write FGA grants (documented setup step).
type AgentListener struct {
	Id             string `json:"id"`
	OrgId          string `json:"org_id"`
	ConnectionId   string `json:"connection_id"`
	ConversationId string `json:"conversation_id"`
	Name           string `json:"name" sortable:""`
	// AgentKey references an agent-service agent by its immutable key.
	AgentKey    string                   `json:"agent_key" sortable:""`
	TriggerType AgentListenerTriggerType `json:"trigger_type"`
	// TriggerConfig is the typed, per-trigger configuration (a tagged union keyed
	// by TriggerType — see agent-listener-trigger.go): channel → ChannelConfig,
	// keyword → KeywordConfig, always/mention carry none. Serializes to the bare
	// variant ({external_channel_id} / {phrases} / {}); TriggerType discriminates.
	TriggerConfig AgentListenerTriggerConfig `json:"trigger_config,omitempty"`
	// ActingUser is the platform user the dispatcher acts as when driving the
	// agent — a common.UserRef ({type,id}) so a non-person actor (agent/api) can
	// own a listener later; it needs agent-sessions:write + messages:write grants.
	ActingUser common.UserRef `json:"acting_user"`
	IsEnabled  bool           `json:"is_enabled" sortable:""`
	// Priority breaks ties when several listeners match one message: highest wins,
	// exactly one listener dispatches.
	Priority  int            `json:"priority" sortable:""`
	CreatedAt time.Time      `json:"created_at" sortable:""`
	CreatedBy common.UserRef `json:"created_by"`
	UpdatedAt time.Time      `json:"updated_at" sortable:""`
	UpdatedBy common.UserRef `json:"updated_by"`
}

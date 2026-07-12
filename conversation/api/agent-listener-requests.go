package conversationapi

import (
	"go.proteos.ai/model/common"
	conversationmodel "go.proteos.ai/model/conversation"
)

// CreateAgentListenerRequest binds inbound messages to an agent. Exactly one of
// ConnectionId / ConversationId must be set (enforced by logic, mirrored by a DB
// CHECK). ActingUserId is a bare user id (the platform user the dispatcher acts
// as — it needs agent-sessions:write + messages:write grants); the service wraps
// it into a common.UserRef (type=person). TriggerConfig is the raw per-trigger
// parameters; the service validates + types it against TriggerType. IsEnabled is
// a tri-state pointer — omitted defaults to TRUE (an omitted create should NOT
// yield a silently-disabled listener).
type CreateAgentListenerRequest struct {
	ConnectionId   string                                     `json:"connection_id"`
	ConversationId string                                     `json:"conversation_id"`
	Name           string                                     `json:"name" validate:"required"`
	AgentKey       string                                     `json:"agent_key" validate:"required"`
	TriggerType    conversationmodel.AgentListenerTriggerType `json:"trigger_type" validate:"required"`
	TriggerConfig  map[string]any                             `json:"trigger_config"`
	ActingUserId   string                                     `json:"acting_user_id" validate:"required"`
	// WakePhrase gates the listener behind a wake word: when set, the listener
	// stays dormant (trigger suppressed) until a message containing the phrase
	// wakes it by starting a session. Empty ⇒ no gate.
	WakePhrase string `json:"wake_phrase"`
	IsEnabled  *bool  `json:"is_enabled,omitempty"`
	// IsAutoForwardAgentRepliesEnabled is a tri-state pointer — omitted defaults to
	// TRUE (an omitted create must keep the historic auto-forward behavior). TRUE ⇒
	// the platform posts the agent's reply automatically; FALSE ⇒ the agent replies
	// itself via the send_message/reply tool.
	IsAutoForwardAgentRepliesEnabled *bool `json:"is_auto_forward_agent_replies_enabled,omitempty"`
	Priority                         int   `json:"priority"`
}

type UpdateAgentListenerRequest struct {
	Name          *string                                     `json:"name,omitempty"`
	AgentKey      *string                                     `json:"agent_key,omitempty"`
	TriggerType   *conversationmodel.AgentListenerTriggerType `json:"trigger_type,omitempty"`
	TriggerConfig *map[string]any                             `json:"trigger_config,omitempty"`
	ActingUserId  *string                                     `json:"acting_user_id,omitempty"`
	// WakePhrase is a tri-state pointer: nil leaves the stored value untouched,
	// "" clears the gate, a value sets it.
	WakePhrase *string `json:"wake_phrase,omitempty"`
	IsEnabled  *bool   `json:"is_enabled,omitempty"`
	// IsAutoForwardAgentRepliesEnabled is a tri-state pointer: nil leaves the stored
	// value untouched, true/false sets it.
	IsAutoForwardAgentRepliesEnabled *bool `json:"is_auto_forward_agent_replies_enabled,omitempty"`
	Priority                         *int  `json:"priority,omitempty"`
}

type GetManyAgentListenersQuery struct {
	ConnectionId   *string `json:"connection_id" form:"connection_id" db:"connection_id"`
	ConversationId *string `json:"conversation_id" form:"conversation_id" db:"conversation_id"`
	AgentKey       *string `json:"agent_key" form:"agent_key" db:"agent_key"`
	IsEnabled      *bool   `json:"is_enabled" form:"is_enabled" db:"is_enabled"`
	common.Pagination
	common.Sorting
}

type GetManyAgentListenersResponse struct {
	Meta common.ResponseMeta               `json:"meta"`
	Data []conversationmodel.AgentListener `json:"data"`
}

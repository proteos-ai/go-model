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
// meeting-companion case). A connection-bound listener may be narrowed to one
// room via RoomId. When several listeners are eligible for a message, the most
// specific scope wins (conversation > room > connection); Priority orders
// within a tier. The dispatcher acts as the resolved acting user (see
// ActingUserMode) — that user needs agent-sessions:write + messages:write FGA
// grants (documented setup step).
type AgentListener struct {
	Id             string `json:"id"`
	OrgId          string `json:"org_id"`
	ConnectionId   string `json:"connection_id"`
	ConversationId string `json:"conversation_id"`
	// RoomId optionally NARROWS a connection-bound listener to one room — the
	// internal room-table row id, matched against Conversation.RoomId at
	// dispatch. Empty ⇒ the whole connection. Never set on conversation-bound
	// listeners (DB CHECK). No FK: a pruned-and-re-minted room gets a new id and
	// the listener then silently stops matching new threads (a conversation
	// stamped room-less at ingest is backfilled by the next inbound message).
	RoomId string `json:"room_id,omitempty"`
	Name   string `json:"name" sortable:""`
	// AgentKey references an agent-service agent by its immutable key.
	AgentKey    string                   `json:"agent_key" sortable:""`
	TriggerType AgentListenerTriggerType `json:"trigger_type"`
	// TriggerConfig is the typed, per-trigger configuration (a tagged union keyed
	// by TriggerType — see agent-listener-trigger.go): channel → ChannelConfig,
	// keyword → KeywordConfig, always/mention carry none. Serializes to the bare
	// variant ({external_channel_id} / {phrases} / {}); TriggerType discriminates.
	TriggerConfig AgentListenerTriggerConfig `json:"trigger_config,omitempty"`
	// WakePhrase, when non-empty, turns the listener into a dormant agent: its
	// configured trigger is suppressed until a session is connected to the
	// conversation, and the ONLY thing that starts that session is an inbound
	// message whose text contains the wake phrase (case-insensitive substring).
	// Once awake, the trigger drives the existing session normally. Empty ⇒ no
	// gate (the trigger fires as always). Orthogonal to TriggerType — it composes
	// on top of whichever trigger matched (a meeting "Hey Ava" companion is
	// trigger=always + wake_phrase="hey ava").
	WakePhrase string `json:"wake_phrase,omitempty"`
	// AcknowledgementType + AcknowledgementConfig configure the immediate
	// acknowledgement the DISPATCHER (never the agent) places on the triggering
	// message once the agent turn is queued — an emoji reaction or a short text
	// — so people see the agent is on it while it thinks. Tagged union like
	// TriggerType/TriggerConfig; "" ⇒ no acknowledgement, config nil.
	AcknowledgementType   AgentListenerAcknowledgementType   `json:"acknowledgement_type"`
	AcknowledgementConfig AgentListenerAcknowledgementConfig `json:"acknowledgement_config,omitempty"`
	// ActingUserMode says where the dispatcher takes its acting user from:
	// defined ⇒ ActingUser below (historic behavior, zero-value default so
	// pre-mode rows and stale cache entries keep working); inferred ⇒ the
	// triggering message sender's resolved platform user, with ActingUser as
	// OPTIONAL fallback — sender unresolved and no fallback ⇒ the dispatch is
	// skipped (no session, no acknowledgement, no fall-through to another
	// listener).
	ActingUserMode AgentListenerActingUserMode `json:"acting_user_mode"`
	// ActingUser is the platform user the dispatcher acts as when driving the
	// agent (mode=defined), or the optional fallback when the sender has no
	// platform user (mode=inferred; zero-value = no fallback). A common.UserRef
	// ({type,id}) so a non-person actor (agent/api) can own a listener later;
	// whoever the dispatcher acts as needs agent-sessions:write +
	// messages:write grants.
	ActingUser common.UserRef `json:"acting_user"`
	IsEnabled  bool           `json:"is_enabled" sortable:""`
	// IsAutoForwardAgentRepliesEnabled controls how the agent's reply reaches the
	// conversation. TRUE (the default) ⇒ the dispatcher folds the agent's
	// agent.message text and posts it back automatically. FALSE ⇒ the platform
	// forwards nothing; the agent is told (in the session preamble) to reply
	// itself via the conversations send_message/reply tool. Distinct from
	// IsEnabled, which turns the whole listener on/off.
	IsAutoForwardAgentRepliesEnabled bool `json:"is_auto_forward_agent_replies_enabled" sortable:""`
	// Priority breaks ties when several listeners match one message: highest wins,
	// exactly one listener dispatches.
	Priority  int            `json:"priority" sortable:""`
	CreatedAt time.Time      `json:"created_at" sortable:""`
	CreatedBy common.UserRef `json:"created_by"`
	UpdatedAt time.Time      `json:"updated_at" sortable:""`
	UpdatedBy common.UserRef `json:"updated_by"`
}

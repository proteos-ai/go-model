package conversationmodel

import "time"

// ConversationSession is the dispatcher's conversation↔agent-session map: the
// first dispatched message of a conversation creates an agent-service session
// and records it here; every later message in that conversation drives the same
// session, so the agent keeps its context. Keyed (org_id, conversation_id);
// concurrent dispatchers race-safely upsert ON CONFLICT DO NOTHING and re-read.
type ConversationSession struct {
	OrgId          string    `json:"org_id"`
	ConversationId string    `json:"conversation_id"`
	SessionId      string    `json:"session_id"`
	AgentKey       string    `json:"agent_key"`
	ListenerId     string    `json:"listener_id"`
	CreatedAt      time.Time `json:"created_at"`
}

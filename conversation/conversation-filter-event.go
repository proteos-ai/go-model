package conversationmodel

import "time"

// ConversationFilterEvent is one entry in the append-only audit trail of
// messages dropped at ingest by a conversation filter — the trace that makes
// silent-drop safe: it answers "why didn't X's email arrive" and feeds the
// per-rule RecentEventCount rollup. Content-free by design (routing/identity
// facts only, no body, no subject) so the audit stays GDPR-light; rows are
// purged after the retention window (see the filter service's purge runner).
// Append-only: never updated, so no updated_*/created_by audit — OccurredAt is
// the single timestamp.
type ConversationFilterEvent struct {
	Id                   string                 `json:"id"`
	OrgId                string                 `json:"org_id"`
	ConnectionId         string                 `json:"connection_id"`
	ConversationFilterId string                 `json:"conversation_filter_id"`
	FilterType           ConversationFilterType `json:"filter_type"`
	// Reason is the matched class + action, e.g. domain_block,
	// internal_conversations, automated_block — the drop explanation.
	Reason string `json:"reason"`
	// SenderKind/SenderValue identify who was dropped, in canonical address
	// terms (the sender's channel-primary ContactAddressKey).
	SenderKind        ContactAddressKind `json:"sender_kind,omitempty"`
	SenderValue       string             `json:"sender_value,omitempty"`
	ExternalMessageId string             `json:"external_message_id,omitempty"`
	OccurredAt        time.Time          `json:"occurred_at" sortable:""`
}

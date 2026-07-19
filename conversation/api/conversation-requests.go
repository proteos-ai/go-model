package conversationapi

import (
	"go.proteos.ai/model/common"
	conversationmodel "go.proteos.ai/model/conversation"
)

// UpdateConversationRequest covers the user-editable surface of a conversation;
// everything else (external key, channel, timestamps) is owned by ingest.
type UpdateConversationRequest struct {
	Subject *string `json:"subject,omitempty"`
	// Summary is the conversation's markdown summary (auto-generated for
	// meetings, freely editable).
	Summary  *string                               `json:"summary,omitempty"`
	Status   *conversationmodel.ConversationStatus `json:"status,omitempty"`
	Metadata *map[string]any                       `json:"metadata,omitempty"`
}

type GetManyConversationsQuery struct {
	Channel      *string `json:"channel" form:"channel" db:"channel"`
	Status       *string `json:"status" form:"status" db:"status"`
	ConnectionId *string `json:"connection_id" form:"connection_id" db:"connection_id"`
	RoomId       *string `json:"room_id" form:"room_id" db:"room_id"`
	// ParentConversationId filters a conversation's thread children (Slack
	// threads forked from one main conversation).
	ParentConversationId *string `json:"parent_conversation_id" form:"parent_conversation_id" db:"parent_conversation_id"`
	// ContactId filters threads whose roster contains the resolved person — the
	// contact stream (subsumes the retired participant_external_id filter:
	// the expansion below covers raw-identifier matching too). Post-contacts
	// snapshots match on the roster's contact_id; pre-contacts rows are covered
	// by the server-side expansion below. No db tag (jsonb containment,
	// repository-applied).
	ContactId *string `json:"contact_id" form:"contact_id"`
	// ContactAddressExternalIds is SERVER-FILLED (never bound from the wire —
	// no form tag): the contact's address identifiers (raw + canonical values),
	// OR-ed into the roster containment so pre-contacts snapshots (which carry
	// no contact_id) still match the person stream.
	ContactAddressExternalIds []string `json:"-" form:"-"`
	// Include opts into expensive read projections; the only value today is
	// "messages_summary" (root/latest message, totals, repliers). No db tag —
	// handled by the repository, not the generic filter mapper.
	Include *string `json:"include" form:"include"`
	common.Pagination
	common.Sorting
}

// IncludeMessagesSummary is the only recognized ?include= value: opt into the
// root/latest message projection on conversation lists (hub stream rows).
const IncludeMessagesSummary = "messages_summary"

type GetManyConversationsResponse struct {
	Meta common.ResponseMeta              `json:"meta"`
	Data []conversationmodel.Conversation `json:"data"`
}

// UnreadCounts aggregates the requesting user's unread conversations (unit:
// conversations, not messages) — the topbar envelope total and the per-channel
// switcher badges.
type UnreadCounts struct {
	Total    int                               `json:"total"`
	Channels map[conversationmodel.Channel]int `json:"channels"`
}

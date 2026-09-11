package conversationmodel

import (
	"time"

	"go.proteos.ai/model/common"
)

// ConversationBrief is guidance prepared AHEAD of a conversation — "talk to
// these people about this": the conversation sibling of a draft message. An
// agent (or a person) writes what the conversation should cover, a human
// opens the brief (a phone brief opens the softphone, which shows it before
// and during the call), and the brief binds to the conversation it was used
// on so the conversation can later be judged against it.
//
// Binding is always EXPLICIT — the human's button press carries the brief id
// into the call (the softphone hands it to the provider, the call-start
// callback attaches it), or an attach call names the conversation. Nothing
// matches briefs to conversations by contact or number, and any number of
// briefs may be open for the same people, channel and type at once: a brief
// is a prepared task, not a rule. Records that want a brief (an activity)
// store its id; every UI surface renders from the id.
//
// v1 is text only. Structured topics (and the adherence rating computed from
// the transcript) are the planned extension — they will land as additional
// fields, never as a replacement of Content.
type ConversationBrief struct {
	Id    string `json:"id"`
	OrgId string `json:"org_id"`
	// Subject is the brief's human label ("Renewal call with Dana") — the
	// headline wherever the brief is listed; empty allowed.
	Subject string `json:"subject" sortable:""`
	// ContactIds are the people the conversation is with (≥ 1, deduplicated).
	// Serializes as an array, never null.
	ContactIds []string `json:"contact_ids"`
	// Channel says how the conversation should happen (phone opens the
	// softphone); "" = unspecified — the brief is read-only guidance.
	Channel Channel `json:"channel,omitempty" sortable:""`
	// TypeKey is the ConversationType the conversation is meant to be; "" =
	// untyped. Stamped onto the conversation when the brief attaches.
	TypeKey string `json:"type_key,omitempty"`
	// ConversationId is the conversation the brief was used on; "" while
	// prepared or discarded.
	ConversationId string                  `json:"conversation_id,omitempty"`
	Status         ConversationBriefStatus `json:"status" sortable:""`
	// Content is the guidance text (markdown allowed).
	Content string `json:"content"`
	// UsedAt stamps the moment the brief attached to a conversation.
	UsedAt    *time.Time     `json:"used_at,omitempty"`
	CreatedAt time.Time      `json:"created_at" sortable:""`
	CreatedBy common.UserRef `json:"created_by"`
	UpdatedAt time.Time      `json:"updated_at" sortable:""`
	UpdatedBy common.UserRef `json:"updated_by"`
}

package conversationmodel

import (
	"time"

	"go.proteos.ai/model/common"
)

// Conversation is one thread of communication on any channel: a Slack thread,
// an email thread, a meeting, or an ad-hoc captured exchange. ConnectionId is
// nullable — adhoc/meeting conversations exist without an integration. The
// external key is unique per connection so webhook replays upsert instead of
// duplicating; adhoc conversations carry an empty external key and stay outside
// that constraint.
type Conversation struct {
	Id           string  `json:"id"`
	OrgId        string  `json:"org_id"`
	ConnectionId string  `json:"connection_id"`
	Channel      Channel `json:"channel" sortable:""`
	// ExternalConversationId is the integration-side thread identity (Slack
	// "channel:thread_ts", Gmail threadId). Empty for adhoc/meeting, and empty on
	// an outbound-originated conversation until the connector sends the first
	// message and the minted thread id is stamped on. The origination target
	// (channel id / To address) is a transient send input (SendOptions.Recipient),
	// NOT stored here — a conversation's durable identity is this thread id (where)
	// plus Participants (who) plus Subject.
	ExternalConversationId string             `json:"external_conversation_id"`
	Subject                string             `json:"subject" sortable:""`
	Status                 ConversationStatus `json:"status" sortable:""`
	// RoomId links a room-borne thread (a Slack channel conversation) to its
	// room directory row; empty for DMs, email, meeting, adhoc. Stamped at
	// ingest; no FK — a pruned room leaves the id dangling by design.
	RoomId string `json:"room_id,omitempty"`
	// ParentConversationId links a child conversation to the main conversation
	// it forked from (Slack: the channel's main conversation; provider-side the
	// child is a native Slack thread); empty for main/DM/email/meeting/adhoc.
	// No FK — empty-string convention.
	ParentConversationId string `json:"parent_conversation_id,omitempty"`
	// StartedByMessageId is OUR message id of the child conversation's root
	// message, which lives in the parent conversation (a message belongs to
	// exactly one conversation); empty when the root is unknown (the child's
	// first message was ingested before its root was ever seen).
	StartedByMessageId string `json:"started_by_message_id,omitempty"`
	// StartedAt/EndedAt are the real time bounds of time-bounded media (meeting,
	// adhoc); nil for open-ended chat threads.
	StartedAt *time.Time `json:"started_at,omitempty"`
	EndedAt   *time.Time `json:"ended_at,omitempty"`
	// Participants is the conversation's complete roster — everyone in the thread,
	// INCLUDING the connected account itself (each entry the account owns is
	// marked IsSelf). It is the canonical "who" of a channel thread, populated for
	// every channel (email, Slack DM, MPIM) as well as meeting/adhoc, and is the
	// queryable person dimension: a person filter is a containment match
	// (participants @> '[{"external_id":…}]'::jsonb, GIN-indexed). Deduped by
	// external_id; grows as new senders/recipients appear on later messages.
	Participants []ConversationParticipant `json:"participants,omitempty"`
	// Transcription artifacts (source/normalized audio, rendered transcript) do
	// NOT live here — they belong to the Transcription that materialized this
	// conversation, linked the other way via Transcription.ConversationId. A
	// conversation may have several transcriptions (segmented meeting, appended
	// audio); fetch them with GET /transcriptions?conversation_id=<id>.
	//
	// LastMessageAt drives the inbox ordering (last_message_at DESC).
	LastMessageAt time.Time      `json:"last_message_at" sortable:""`
	Metadata      map[string]any `json:"metadata"`
	// Messages summary — read-time projection over the message table, present
	// only when a list/get opted in via ?include=messages_summary. Lets a hub
	// stream render each conversation row (root on Slack, latest on email,
	// reply affordance) without a per-conversation message fetch.
	RootMessage   *MessagePreview `json:"root_message,omitempty"`
	LatestMessage *MessagePreview `json:"latest_message,omitempty"`
	MessagesTotal int             `json:"messages_total,omitempty"`
	// LastRepliers are up to 3 most recent distinct senders excluding the root
	// sender — the stream's stacked mini avatars.
	LastRepliers []ParticipantRef `json:"last_repliers,omitempty"`
	// Read state of the requesting user (never another user's) — projected onto
	// every authenticated read. LastReadAt nil = never opened. IsUnread is the
	// server-computed inbox predicate: activity after the marker AND the latest
	// message is inbound (own outbound sends don't re-flag the thread).
	LastReadAt *time.Time     `json:"last_read_at,omitempty"`
	IsUnread   bool           `json:"is_unread"`
	CreatedAt  time.Time      `json:"created_at" sortable:""`
	CreatedBy  common.UserRef `json:"created_by"`
	UpdatedAt  time.Time      `json:"updated_at" sortable:""`
	UpdatedBy  common.UserRef `json:"updated_by"`
}

// MessagePreview is the compact snapshot of one message riding on conversation
// reads (?include=messages_summary): enough to render a stream row — who, the
// first ~500 chars of text, direction/status for outbound identity — without
// the full content blocks.
type MessagePreview struct {
	Sender     ParticipantRef   `json:"sender"`
	Text       string           `json:"text"`
	Direction  MessageDirection `json:"direction"`
	Status     MessageStatus    `json:"status"`
	OccurredAt time.Time        `json:"occurred_at"`
}

// ParticipantRef is the inline snapshot of a person on a message/reaction —
// the integration-side identity plus, when resolved, the platform user. It is
// hydrated from the participant directory at ingest (a local read, never a
// provider call), so reads carry rich identity with no extra hop. Name
// is always populated when known so UIs and agents have something to show;
// Email rides along when the directory has it and backs the platform-user
// resolution (email match).
type ParticipantRef struct {
	ExternalId   string          `json:"external_id,omitempty"`
	Name         string          `json:"name"`
	Email        string          `json:"email,omitempty"`
	PlatformUser *common.UserRef `json:"platform_user,omitempty"`
}

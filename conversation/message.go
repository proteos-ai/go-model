package conversationmodel

import (
	"time"

	"go.proteos.ai/model/common"
)

// Message is what a participant said or sent — THE platform Message (the bus
// envelope is events.Event since the events rename). Inbound messages are
// deduped on ExternalMessageId per connection; outbound messages carry the
// pending→sent/failed lifecycle. OccurredAt is when it was SAID (real utterance
// or import time for meeting/adhoc/backfill; defaults to created_at otherwise)
// and is the list ordering — not created_at, which is merely when the row landed.
type Message struct {
	Id             string           `json:"id"`
	OrgId          string           `json:"org_id"`
	ConversationId string           `json:"conversation_id"`
	ConnectionId   string           `json:"connection_id"`
	Channel        Channel          `json:"channel"`
	Direction      MessageDirection `json:"direction"`
	// ExternalMessageId is the integration-side message identity (Slack event ts,
	// Gmail message id). UNIQUE per org+connection where non-empty AND inbound —
	// the webhook-replay dedupe. Empty for outbound until the connector returns one.
	ExternalMessageId string     `json:"external_message_id"`
	Sender            ContactRef `json:"sender"`
	// Recipients is who the message was addressed to, each tagged with its role
	// (to/cc/bcc). Email carries the full To/Cc list (and Bcc on messages WE
	// sent); other channels list the direct target(s) as `to`. Stored as one
	// jsonb column; empty for messages with no explicit addressees.
	Recipients []MessageRecipient `json:"recipients,omitempty"`
	Content    []ContentBlock     `json:"content"`
	Status     MessageStatus      `json:"status"`
	// ReplyToMessageId records the message this one was targeted at (the
	// thread-mode send anchor). Persisted so a draft remembers its target until
	// the human sends it (the provider-side anchor is re-resolved at send time);
	// stamped on directly-sent targeted replies too, as provenance. Empty for
	// everything else.
	ReplyToMessageId string    `json:"reply_to_message_id,omitempty"`
	OccurredAt       time.Time `json:"occurred_at" sortable:""`
	// Reactions is a READ-TIME projection aggregated from the message_reaction
	// edges (never a stored column): empty for a message nobody reacted to and
	// for channels without reactions — that IS the graceful degradation.
	Reactions []Reaction `json:"reactions,omitempty"`
	// Attachments is a READ-TIME projection of the message's attachment rows
	// (bytes live in storage-service); absent for messages without files.
	Attachments []Attachment `json:"attachments,omitempty"`
	// Child-conversation projection (READ-TIME, per-conversation lists only):
	// the conversation rooted at this message and its message count — the "N
	// replies" affordance. Populated via conversation.started_by_message_id;
	// zero-valued for messages no child conversation forked from.
	ChildConversationId string `json:"child_conversation_id,omitempty"`
	ChildMessagesTotal  int    `json:"child_messages_total,omitempty"`
	// Metadata carries channel-specific extras; for spoken messages (materialized
	// transcription turns): start_ms, end_ms, confidence, attribution_source.
	Metadata map[string]any `json:"metadata"`
	// Error holds the connector failure detail when status=failed.
	Error     string         `json:"error"`
	CreatedAt time.Time      `json:"created_at" sortable:""`
	CreatedBy common.UserRef `json:"created_by"`
	UpdatedAt time.Time      `json:"updated_at" sortable:""`
	UpdatedBy common.UserRef `json:"updated_by"`
}

// ContentBlock is one unit of message content — what was SAID: text, or html
// (the sanitized rich body of an email, stored ALONGSIDE the text block,
// never instead of it). Text stays the canonical content every consumer
// (agents, previews, search) reads; html is display-only fidelity for UIs
// that can render it in isolation. Files a message CARRIED are NOT content —
// they are Attachment rows (Message.Attachments); every surveyed platform
// (Slack files[], Graph attachments[], JMAP, WhatsApp media ids) separates
// body from files the same way. The closed set keeps consumers exhaustive
// per block type.
type ContentBlock struct {
	Type ContentBlockType `json:"type"`
	Text string           `json:"text,omitempty"`
	// Html is sanitized at ingest (allowlist: layout tables, inline styles,
	// https images, rewritten links; no scripts/forms/event handlers) and must
	// STILL be rendered in an isolated, script-less surface — sanitization is
	// defense one, the sandbox is defense two.
	Html string `json:"html,omitempty"`
}

// ContentBlockType discriminates ContentBlock variants.
type ContentBlockType string

const (
	ContentBlockTypeText ContentBlockType = "text"
	ContentBlockTypeHtml ContentBlockType = "html"
)

// NormalizedInboundMessage is the wire shape every ChannelIngestor produces and
// the domain consumes (IngestSink.IngestInbound) — it lives in the model package
// so connectors and domain share it without the connectors importing domain
// packages. ConnectionId is set when the ingestor already resolved the
// connection (e.g. from an OAuth-bound route); otherwise the domain resolves it
// via (ConnectorKey, ExternalAccountId).
type NormalizedInboundMessage struct {
	ConnectorKey           ConnectorKey `json:"connector_key"`
	Channel                Channel      `json:"channel"`
	ConnectionId           string       `json:"connection_id,omitempty"`
	ExternalAccountId      string       `json:"external_account_id"`
	ExternalConversationId string       `json:"external_conversation_id"`
	ExternalMessageId      string       `json:"external_message_id"`
	// ExternalParentConversationId is the external key of the MAIN conversation
	// this message's child conversation forked from (Slack: the channel id —
	// the message arrived inside a native Slack thread). Empty for top-level
	// messages and for channels without a second conversation level (email,
	// unipile). When set, the domain upserts the parent conversation and links
	// the child via parent_conversation_id.
	ExternalParentConversationId string `json:"external_parent_conversation_id,omitempty"`
	// ExternalStartedByMessageId is the external message id of the child
	// conversation's root message (Slack: "<channel>:<thread_ts>"); the domain
	// resolves it to Conversation.StartedByMessageId best-effort (empty when
	// the root predates ingest).
	ExternalStartedByMessageId string `json:"external_started_by_message_id,omitempty"`
	// Direction is the message's direction from the connected account's point of
	// view. Empty defaults to inbound — connectors that only ever ingest received
	// messages (Slack, echo) leave it unset. A connector that also captures mail
	// the account itself sent (Gmail SENT label, for full-thread visibility) sets
	// it to outbound; the domain then stores status=sent and the dispatcher skips
	// it (its loop-safety consumes inbound only).
	Direction MessageDirection `json:"direction,omitempty"`
	// ExternalRoomId is the venue the thread lives in (a Slack channel id, a
	// WhatsApp group's chat id) when the connector knows it; empty for direct
	// threads and venue-less channels (email). The domain resolves it against
	// the room directory — minting the row on first sight — to stamp
	// Conversation.RoomId.
	ExternalRoomId string `json:"external_room_id,omitempty"`
	// ExternalRoomName seeds an ingest-minted room's display name when the
	// webhook already carries it (a WhatsApp group's subject); empty when the
	// connector's directory lookup owns naming (Slack conversations.info).
	ExternalRoomName string `json:"external_room_name,omitempty"`
	// IsDirect marks a 1:1 thread (Slack DM, email person thread). The domain
	// builds the conversation roster from the enriched sender plus Recipients,
	// marking the connected account's own entry IsSelf.
	IsDirect bool `json:"is_direct,omitempty"`
	// Recipients is who the message was addressed to (To/Cc parsed from the
	// provider), each tagged with its role. The domain enriches these against the
	// participant directory and folds them into the conversation roster. Inbound
	// Bcc is not present (SMTP strips it).
	Recipients []MessageRecipient `json:"recipients,omitempty"`
	Sender     ContactRef         `json:"sender"`
	Content    []ContentBlock     `json:"content"`
	Subject    string             `json:"subject,omitempty"`
	// OccurredAt is when the message was said/sent on the external side; nil means
	// "now" (the domain stamps ingestion time).
	OccurredAt *time.Time     `json:"occurred_at,omitempty"`
	Metadata   map[string]any `json:"metadata,omitempty"`
	// Headers is a small fixed allow-list of filter-relevant email headers
	// (lowercased keys: auto-submitted, precedence, list-id, list-unsubscribe,
	// x-auto-response-suppress, return-path), populated by email ingestors only
	// — the input to the `automated` conversation-filter type. Nil on non-email
	// channels (those filters are then inert). Consumed at ingest, never
	// persisted; deliberately NOT a full header dump (PII + payload bloat).
	Headers map[string]string `json:"headers,omitempty"`
	// Attachments carries the message's files as decoded bytes; the domain
	// uploads them to storage-service and stores Attachment rows. Connectors
	// bound the sizes (a connector never hands over unbounded payloads).
	Attachments []NormalizedAttachment `json:"attachments,omitempty"`
}

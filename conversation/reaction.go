package conversationmodel

import (
	"time"

	"go.proteos.ai/model/common"
)

// MessageReaction is the stored reaction edge: one row per (message, emoji,
// participant). reaction_added = insert, reaction_removed = delete — rows are
// never updated, so there is no updated_* audit. Emoji is the CONNECTOR-NATIVE
// token (Slack shortcode "thumbsup", unicode elsewhere) so it round-trips per
// channel; cross-channel normalization is deliberately a read-time display
// concern, never the stored truth.
//
// Participant is the rich snapshot of WHO reacted on the channel, resolved
// from the participant directory at ingest; ParticipantExternalId is the same
// identity denormalized as a scalar for the dedupe unique index. CreatedBy is
// the PLATFORM audit stamp — a distinct subject: on ingest it is the system,
// on an API-set reaction it is the calling user while the participant is the
// connection's bot (who the provider records).
type MessageReaction struct {
	Id                    string     `json:"id"`
	OrgId                 string     `json:"org_id"`
	MessageId             string     `json:"message_id"`
	ConnectionId          string     `json:"connection_id"`
	Emoji                 string     `json:"emoji"`
	Participant           ContactRef `json:"participant"`
	ParticipantExternalId string     `json:"participant_external_id"`
	// ExternalReactionId is the provider-side identity when one exists; Slack
	// reactions have none (identity IS message+emoji+user), so it may be empty.
	ExternalReactionId string         `json:"external_reaction_id,omitempty"`
	OccurredAt         time.Time      `json:"occurred_at"`
	CreatedAt          time.Time      `json:"created_at"`
	CreatedBy          common.UserRef `json:"created_by"`
}

// Reaction is the aggregated read VO riding on Message.Reactions: one entry
// per emoji with the count and the reacting participants.
type Reaction struct {
	Emoji        string       `json:"emoji"`
	Count        int          `json:"count"`
	Participants []ContactRef `json:"participants"`
}

// ReactionAction discriminates a NormalizedReaction: the participant added or
// removed the emoji. Ingest-only — not persisted (add = insert, remove =
// delete).
type ReactionAction string

const (
	ReactionActionAdded   ReactionAction = "added"
	ReactionActionRemoved ReactionAction = "removed"
)

// NormalizedReaction is the wire shape a ChannelIngestor produces for a
// reaction event (sibling of NormalizedInboundMessage). It carries only the
// reactor's external id — reaction events have no name/email; the domain
// resolves the rich ContactRef snapshot from the participant directory
// before persisting.
type NormalizedReaction struct {
	ConnectorKey          ConnectorKey   `json:"connector_key"`
	Channel               Channel        `json:"channel"`
	ConnectionId          string         `json:"connection_id,omitempty"`
	ExternalAccountId     string         `json:"external_account_id"`
	ExternalMessageId     string         `json:"external_message_id"`
	Emoji                 string         `json:"emoji"`
	ParticipantExternalId string         `json:"participant_external_id"`
	Action                ReactionAction `json:"action"`
	// OccurredAt is when the reaction happened on the external side; nil means
	// "now" (the domain stamps ingestion time).
	OccurredAt *time.Time `json:"occurred_at,omitempty"`
}

// ReactionSetKind says how a channel's reaction vocabulary is shaped: open
// (any emoji token the provider accepts — Slack, including custom workspace
// emoji) or fixed (a closed set — iMessage tapbacks, LinkedIn).
type ReactionSetKind string

const (
	ReactionSetOpen  ReactionSetKind = "open"
	ReactionSetFixed ReactionSetKind = "fixed"
)

// ReactionOption is one entry of a fixed reaction vocabulary. Token is the
// connector-native value to SEND; Unicode is a best-effort display glyph
// (empty for custom emoji with no unicode equivalent) — display only, never
// the stored identity.
type ReactionOption struct {
	Token   string `json:"token"`
	Unicode string `json:"unicode,omitempty"`
	Label   string `json:"label,omitempty"`
}

// ReactionCapability is a connector's declarative answer to "what can I do
// with reactions here" — surfaced on the connection read so UIs adapt per
// connection (open → emoji picker; fixed → exactly the Allowed buttons;
// absent → no affordance). The connector is the sole authority: it is the
// only code that knows the protocol.
type ReactionCapability struct {
	Kind ReactionSetKind `json:"kind"`
	// Allowed is the closed vocabulary, populated only when Kind == fixed.
	Allowed []ReactionOption `json:"allowed,omitempty"`
	// MaxPerActor bounds how many reactions one participant may hold on one
	// message: 0 = unbounded (Slack), 1 = single-slot (a new reaction replaces
	// the previous one — WhatsApp, iMessage).
	MaxPerActor int `json:"max_per_actor"`
}

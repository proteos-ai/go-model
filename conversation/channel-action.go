package conversationmodel

import (
	"time"

	"go.proteos.ai/model/common"
)

// ChannelAction is ONE act performed through a channel connection that is
// neither a message nor a reaction: a LinkedIn connection invitation, a
// profile visit, an InMail (which ALSO mints a Message — see MessageId), later
// endorsements, post engagement, lead saving. The three outbound things a
// connection does, and how they differ:
//
//   - message  — what was said; the Message row, sent through Send.
//   - reaction — an emoji TOGGLED on a stored message; a set edge (present or
//     absent, idempotent, never counted), mirrored from the provider's event
//     stream — see MessageReaction.
//   - channel_action — an act with its own outcome and lifecycle; a LEDGER row
//     (append, transition, never toggle), counted against sending limits and
//     frequency caps, targeting a person or an external object rather than a
//     stored message. Our row is the source of truth; provider reconciliation
//     (webhooks + sweeps) moves its status.
//
// Rule of thumb: toggling a fact on something we already store (reaction,
// pin, read, archive) is a resource verb or edge; performing an outward act
// with a result is a channel_action.
//
// Direction: outbound = we acted (the connection's account invited someone);
// inbound = someone acted on us (a received invitation) — an inbound row is
// what Respond answers.
type ChannelAction struct {
	Id           string       `json:"id"`
	OrgId        string       `json:"org_id"`
	ConnectionId string       `json:"connection_id" sortable:""`
	ConnectorKey ConnectorKey `json:"connector_key"`
	// Channel is denormalized from the connection (like Message.Channel) so the
	// ledger filters without a registry lookup.
	Channel    Channel             `json:"channel" sortable:""`
	ActionType ChannelActionType   `json:"action_type" sortable:""`
	Direction  MessageDirection    `json:"direction" sortable:""`
	Status     ChannelActionStatus `json:"status" sortable:""`
	// Target — flat optional columns, never a nested union. Contact is the
	// person acted on (resolved from the address; empty when unresolvable),
	// ContactAddressId the address row on this channel, TargetExternalId the
	// provider-side identity the connector acted on (a LinkedIn member id; for
	// future external targets a post's social id).
	ContactId        string `json:"contact_id,omitempty"`
	ContactAddressId string `json:"contact_address_id,omitempty"`
	TargetExternalId string `json:"target_external_id"`
	// Contact is the person snapshot (name, resolved platform user) for
	// display — the same shape every message carries as Sender.
	Contact ContactRef `json:"contact"`
	// MessageId / ConversationId link the act to what it produced: an InMail's
	// minted message (and its conversation); an invitation's note-chat once the
	// invitee accepts and the note arrives as a message.
	MessageId      string `json:"message_id,omitempty"`
	ConversationId string `json:"conversation_id,omitempty"`
	// Params is the typed per-type input (tagged union keyed by ActionType —
	// see channel-action-params.go). Serializes to the bare variant.
	Params ChannelActionParams `json:"params,omitempty"`
	// ExternalActionId is the provider handle of the act (Unipile invitation
	// id): the cancel/respond key and the sweep dedupe key. Empty for acts the
	// provider does not identify (a profile visit).
	ExternalActionId string `json:"external_action_id,omitempty"`
	// Error holds the connector failure detail when status=failed.
	Error string `json:"error,omitempty"`
	// OccurredAt is when the act happened provider-side (inbound) or when we
	// performed it; ResolvedAt when a terminal outcome landed (accepted /
	// declined / withdrawn / expired).
	OccurredAt time.Time  `json:"occurred_at" sortable:""`
	ResolvedAt *time.Time `json:"resolved_at,omitempty"`
	// Metadata is provider enrichment (Unipile invitation usage %, the visited
	// profile's network_distance, InMail credit balance).
	Metadata map[string]any `json:"metadata"`
	// CreatedBy IS the performer on outbound rows (the calling user or agent);
	// SystemUserRef on ingested inbound rows, whose performer is the contact.
	// UpdatedBy is who accepted / declined / withdrew. Same split as
	// MessageReaction.Participant + CreatedBy — no separate performer field.
	CreatedAt time.Time      `json:"created_at" sortable:""`
	CreatedBy common.UserRef `json:"created_by"`
	UpdatedAt time.Time      `json:"updated_at" sortable:""`
	UpdatedBy common.UserRef `json:"updated_by"`
}

// ChannelActionCapability declares ONE action type a connector can perform,
// projected onto Connection.Actions (computed on read, never stored) so UIs
// and agents discover what a connection allows — the channel_action sibling of
// ReactionCapability (which stays its own field: reactions are not actions).
type ChannelActionCapability struct {
	ActionType ChannelActionType       `json:"action_type"`
	TargetKind ChannelActionTargetKind `json:"target_kind"`
	// IsCancelable: the act can be withdrawn after performing (an invitation).
	IsCancelable bool `json:"is_cancelable"`
	// IsRespondable: the connector can answer an INBOUND act of this type
	// (accept / decline a received invitation).
	IsRespondable bool `json:"is_respondable"`
	// IsMessageMinting: performing the act also sends a message that lands as a
	// Message + Conversation (InMail).
	IsMessageMinting bool `json:"is_message_minting"`
	// MaxNoteLength bounds a free-text note the act carries (LinkedIn invitation
	// note: 300); 0 = no note.
	MaxNoteLength int `json:"max_note_length,omitempty"`
}

// NormalizedChannelAction is the wire shape connectors hand the domain for
// acts observed provider-side — a received invitation (inbound, pending), a
// sent invitation seen on a sweep (outbound, performed), an accepted relation
// (outbound, accepted). Lives in the model package so connectors and domain
// share it without the connectors importing domain packages.
type NormalizedChannelAction struct {
	ConnectorKey      ConnectorKey        `json:"connector_key"`
	Channel           Channel             `json:"channel"`
	ConnectionId      string              `json:"connection_id,omitempty"`
	ExternalAccountId string              `json:"external_account_id"`
	ActionType        ChannelActionType   `json:"action_type"`
	Direction         MessageDirection    `json:"direction"`
	Status            ChannelActionStatus `json:"status"`
	ExternalActionId  string              `json:"external_action_id,omitempty"`
	// Participant is the counterpart as the provider names them: the inviter of
	// a received invitation, the invitee of a sent one, the new relation.
	Participant ContactRef `json:"participant"`
	// Note is the free text riding the act (an invitation's message).
	Note string `json:"note,omitempty"`
	// OccurredAt is the provider-side time; nil means "now".
	OccurredAt *time.Time     `json:"occurred_at,omitempty"`
	Metadata   map[string]any `json:"metadata,omitempty"`
}

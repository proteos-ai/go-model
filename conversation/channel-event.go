package conversationmodel

import (
	"time"

	"go.proteos.ai/model/common"
)

// ChannelEvent is ONE provider-observed occurrence about one outbound
// message on a sending platform: it was delivered, it bounced, the recipient
// opened it, clicked a link, reported it as spam, unsubscribed. The ledger
// sibling of channel_action (acts WE perform) and message_reaction (edges
// toggled on a message) — an event is neither: it is something the provider
// or the recipient did that we merely record. Append-only, deduped on the
// provider's event id (webhooks replay), never updated.
//
// Two projections derive from the ledger and are never written directly:
// Message.Delivery (the per-message read-side denorm the inbox renders) and
// the contact permission events a hard bounce / spam report / unsubscribe
// implies (source=provider), appended through the contact service.
type ChannelEvent struct {
	Id           string       `json:"id"`
	OrgId        string       `json:"org_id"`
	ConnectionId string       `json:"connection_id" sortable:""`
	ConnectorKey ConnectorKey `json:"connector_key"`
	// Channel is denormalized from the connection (like Message.Channel) so
	// the ledger filters without a registry lookup.
	Channel Channel `json:"channel" sortable:""`
	// MessageId / ConversationId link the event to the message it is about —
	// resolved from the correlation ids the send stamped onto the provider
	// message (custom args), else from the provider message id. Empty when
	// the provider reports on a message we never sent through this platform
	// (kept for account-level counting, never shown on a message).
	MessageId      string `json:"message_id,omitempty" sortable:""`
	ConversationId string `json:"conversation_id,omitempty"`
	// ContactId / ContactAddressId: the recipient the event is about, resolved
	// from RecipientAddress against the address directory (empty when
	// unresolvable — the address is still recorded).
	ContactId        string `json:"contact_id,omitempty"`
	ContactAddressId string `json:"contact_address_id,omitempty"`
	// RecipientAddress is the provider's recipient (lowercased email).
	RecipientAddress string           `json:"recipient_address"`
	EventType        ChannelEventType `json:"event_type" sortable:""`
	// ExternalEventId is the provider's event identity — UNIQUE per connection
	// (webhook replay dedupe). ExternalMessageId is the provider's message
	// identity the event carries (== Message.ExternalMessageId of the sent
	// row).
	ExternalEventId   string `json:"external_event_id"`
	ExternalMessageId string `json:"external_message_id,omitempty"`
	// OccurredAt is the provider-side time of the event (the list ordering).
	OccurredAt time.Time `json:"occurred_at" sortable:""`
	// Reason / Response: the provider's or the receiving server's text for a
	// bounce, deferral or drop (the SMTP response, the drop reason).
	Reason   string `json:"reason,omitempty"`
	Response string `json:"response,omitempty"`
	// BounceKind / BounceClassification: bounced events only.
	BounceKind           BounceKind           `json:"bounce_kind,omitempty"`
	BounceClassification BounceClassification `json:"bounce_classification,omitempty"`
	// Url is the clicked link (clicked events only).
	Url string `json:"url,omitempty"`
	// IsMachineOpen marks an opened event a mail client's privacy proxy
	// generated (Apple Mail Privacy Protection) rather than the person.
	IsMachineOpen bool `json:"is_machine_open,omitempty"`
	// Metadata is provider enrichment that has no column: user agent, ip,
	// tls, deferral attempt, unsubscribe group, categories.
	Metadata  map[string]any `json:"metadata"`
	CreatedAt time.Time      `json:"created_at" sortable:""`
	CreatedBy common.UserRef `json:"created_by"`
}

// NormalizedChannelEvent is the wire shape a connector's ingestor hands the
// domain for one provider event — the channel_event sibling of
// NormalizedInboundMessage. The connector fills what the provider payload
// carries (its ids, the recipient, the classification); the domain resolves
// the connection, the message row, and the contact. MessageId is the
// correlation id the send stamped onto the provider message when the
// payload echoes it back, else empty and the domain falls back to
// ExternalMessageId.
type NormalizedChannelEvent struct {
	ConnectorKey         ConnectorKey         `json:"connector_key"`
	Channel              Channel              `json:"channel"`
	ConnectionId         string               `json:"connection_id,omitempty"`
	ExternalAccountId    string               `json:"external_account_id,omitempty"`
	EventType            ChannelEventType     `json:"event_type"`
	ExternalEventId      string               `json:"external_event_id"`
	ExternalMessageId    string               `json:"external_message_id,omitempty"`
	MessageId            string               `json:"message_id,omitempty"`
	ConversationId       string               `json:"conversation_id,omitempty"`
	RecipientAddress     string               `json:"recipient_address"`
	OccurredAt           time.Time            `json:"occurred_at"`
	Reason               string               `json:"reason,omitempty"`
	Response             string               `json:"response,omitempty"`
	BounceKind           BounceKind           `json:"bounce_kind,omitempty"`
	BounceClassification BounceClassification `json:"bounce_classification,omitempty"`
	Url                  string               `json:"url,omitempty"`
	IsMachineOpen        bool                 `json:"is_machine_open,omitempty"`
	Metadata             map[string]any       `json:"metadata,omitempty"`
}

// MessageDelivery is the read-side projection of a message's channel events
// — one JSONB column on the message so the inbox renders the delivery state
// without joining the ledger. Status is the provider-side outcome; the
// first-time stamps + totals summarize engagement. Nil on channels that
// report nothing (mailbox connectors, chat).
type MessageDelivery struct {
	Status MessageDeliveryStatus `json:"status"`
	// DeliveredAt / BouncedAt / DroppedAt: the terminal delivery instants.
	DeliveredAt *time.Time `json:"delivered_at,omitempty"`
	BouncedAt   *time.Time `json:"bounced_at,omitempty"`
	DroppedAt   *time.Time `json:"dropped_at,omitempty"`
	// OpenedAt / ClickedAt are the FIRST human open / click; the totals count
	// every event (machine opens excluded from both).
	OpenedAt    *time.Time `json:"opened_at,omitempty"`
	OpensTotal  int        `json:"opens_total,omitempty"`
	ClickedAt   *time.Time `json:"clicked_at,omitempty"`
	ClicksTotal int        `json:"clicks_total,omitempty"`
	// SpamReportedAt / UnsubscribedAt: the recipient's negative signals.
	SpamReportedAt *time.Time `json:"spam_reported_at,omitempty"`
	UnsubscribedAt *time.Time `json:"unsubscribed_at,omitempty"`
	// Reason is the latest bounce / drop / deferral text for display.
	Reason string `json:"reason,omitempty"`
	// LastEventAt is the newest event folded in (idempotent replays and
	// out-of-order batches compare against it).
	LastEventAt *time.Time `json:"last_event_at,omitempty"`
}

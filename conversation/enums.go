package conversationmodel

// Channel is the platform medium a conversation lives on. It is a closed enum
// keyed on platform identity (the granularity users recognize and filter by —
// "Slack", "Email"), NOT on semantic category; deliberate decision 2026-07-02,
// matching the prevailing industry shape (Chatwoot channel_type, Intercom
// source.type). Multiple connectors may serve one channel (gmail and outlook
// both feed email). meeting/adhoc are connector-less time-bounded spoken media.
type Channel string

const (
	ChannelSlack     Channel = "slack"
	ChannelEmail     Channel = "email"
	ChannelLinkedin  Channel = "linkedin"
	ChannelSms       Channel = "sms"
	ChannelTeams     Channel = "teams"
	ChannelTelegram  Channel = "telegram"
	ChannelWhatsapp  Channel = "whatsapp"
	ChannelMeeting   Channel = "meeting"
	ChannelAdhoc     Channel = "adhoc"
	ChannelInstagram Channel = "instagram"
	ChannelMessenger Channel = "messenger"
	// ChannelX is the platform formerly known as Twitter. One name across
	// Go/SDK/UI/DB per the ubiquitous-naming rule; Unipile's provider constant
	// TWITTER stays adapter-internal (documented deviation).
	ChannelX Channel = "x"
)

// ConnectorKey identifies a concrete integration that produces and/or sends
// messages on a channel. A typed enum (not an open string — every connector is
// hand-coded, so the set is exactly the compiled-in connector packages; decision
// 2026-07-02). Runtime availability is decided by REGISTRY membership in the
// service (a key may exist here while the connector is not configured in an
// environment). Grows one value per shipped integration.
type ConnectorKey string

const (
	ConnectorKeySlack ConnectorKey = "slack"
	ConnectorKeyGmail ConnectorKey = "gmail"
	// ConnectorKeyEcho is the in-repo stub connector used by tests/e2e to prove
	// the registry + ingestor seam without an external dependency.
	ConnectorKeyEcho ConnectorKey = "echo"
	// The unipile-* connectors are the six messaging systems served through the
	// Unipile aggregator — one key per messaging system (not one "unipile" key),
	// so each shows up as its own connector in the catalog and a future native
	// connector on the same channel (e.g. twilio-whatsapp) coexists cleanly.
	ConnectorKeyUnipileWhatsapp  ConnectorKey = "unipile-whatsapp"
	ConnectorKeyUnipileLinkedin  ConnectorKey = "unipile-linkedin"
	ConnectorKeyUnipileTelegram  ConnectorKey = "unipile-telegram"
	ConnectorKeyUnipileInstagram ConnectorKey = "unipile-instagram"
	ConnectorKeyUnipileMessenger ConnectorKey = "unipile-messenger"
	ConnectorKeyUnipileX         ConnectorKey = "unipile-x"
)

// ConnectorProvider says who operates the integration mechanics behind a
// connector: Proteos' own hand-coded integration against the provider's API
// (native) or an account aggregated through Unipile's platform tenancy
// (unipile). Computed on Connection reads from the connector, never stored;
// the UI groups the connector catalog by it.
type ConnectorProvider string

const (
	ConnectorProviderNative  ConnectorProvider = "native"
	ConnectorProviderUnipile ConnectorProvider = "unipile"
)

// MessageDirection distinguishes what a participant sent to us (inbound) from
// what the platform sent out through a connector (outbound). The dispatcher's
// loop-safety keys on this: only inbound messages are bridged to agents.
type MessageDirection string

const (
	MessageDirectionInbound  MessageDirection = "inbound"
	MessageDirectionOutbound MessageDirection = "outbound"
)

// MessageStatus is the delivery state of a Message. Inbound messages are
// `received` on arrival; outbound messages start `pending` and flip to `sent`
// or `failed` after the connector Send (published as message.updated).
type MessageStatus string

const (
	MessageStatusReceived MessageStatus = "received"
	MessageStatusPending  MessageStatus = "pending"
	MessageStatusSent     MessageStatus = "sent"
	MessageStatusFailed   MessageStatus = "failed"
)

// ConnectionScope says who a connection belongs to: an org-wide integration
// (a Slack workspace, a shared mailbox) or a single user's grant (their personal
// Gmail). scope=user connections carry owner_user_id.
type ConnectionScope string

const (
	ConnectionScopeOrg  ConnectionScope = "org"
	ConnectionScopeUser ConnectionScope = "user"
)

// ConnectionStatus is the lifecycle of an integration instance: created but not
// yet installed (pending), healthy (active), failing (error), or de-authorized
// by the external side (revoked).
type ConnectionStatus string

const (
	ConnectionStatusPending ConnectionStatus = "pending"
	ConnectionStatusActive  ConnectionStatus = "active"
	ConnectionStatusError   ConnectionStatus = "error"
	ConnectionStatusRevoked ConnectionStatus = "revoked"
)

// ConversationStatus is the lifecycle of a conversation. `ended` carries real
// semantics for time-bounded media (a meeting that finished); `archived` is a
// user-facing tidy-up state.
type ConversationStatus string

const (
	ConversationStatusActive   ConversationStatus = "active"
	ConversationStatusEnded    ConversationStatus = "ended"
	ConversationStatusArchived ConversationStatus = "archived"
)

// TranscriptionStatus is the batch-transcription lifecycle. v1 transcribes
// synchronously inside the request, so rows move pending→processing→completed/
// failed within one call; the async (callback) path reuses the same states.
type TranscriptionStatus string

const (
	TranscriptionStatusPending    TranscriptionStatus = "pending"
	TranscriptionStatusProcessing TranscriptionStatus = "processing"
	TranscriptionStatusCompleted  TranscriptionStatus = "completed"
	TranscriptionStatusFailed     TranscriptionStatus = "failed"
)

// AgentListenerTriggerType says when a listener fires for an inbound message:
// on every message (always), only when the connection's bot user is mentioned
// (mention), only in a specific external channel (channel), or when the text
// matches configured phrases (keyword).
type AgentListenerTriggerType string

const (
	TriggerTypeAlways  AgentListenerTriggerType = "always"
	TriggerTypeMention AgentListenerTriggerType = "mention"
	TriggerTypeChannel AgentListenerTriggerType = "channel"
	TriggerTypeKeyword AgentListenerTriggerType = "keyword"
)

// RecipientKind discriminates the send-time Recipient VO: the target is a
// person (participant) or a venue (room).
type RecipientKind string

const (
	RecipientKindParticipant RecipientKind = "participant"
	RecipientKindRoom        RecipientKind = "room"
)

// RecipientRole is how a MessageRecipient was addressed on a message: a primary
// recipient (to), a carbon copy (cc), or a blind carbon copy (bcc). Only email
// distinguishes the three; other channels use `to`. Bcc is only known for
// outbound messages the platform itself sent.
type RecipientRole string

const (
	RecipientRoleTo  RecipientRole = "to"
	RecipientRoleCc  RecipientRole = "cc"
	RecipientRoleBcc RecipientRole = "bcc"
)

// ParticipantSource records how a participant directory row was populated: a
// full provider directory sweep (sync — prunable when the provider stops
// returning it) or derived passively from ingested messages (ingest — a
// "recent correspondent").
type ParticipantSource string

const (
	ParticipantSourceSync   ParticipantSource = "sync"
	ParticipantSourceIngest ParticipantSource = "ingest"
)

// RoomSource records how a room directory row was populated — the deliberate
// mirror of ParticipantSource for the venue half of the directory (separate
// entity-scoped types, same values, same prune semantics): a provider
// directory sweep (sync — prunable when the provider stops returning it) or
// minted at ingest for a channel the sweep hasn't seen yet (ingest — survives
// pruning until the sweep confirms it, then flips to sync).
type RoomSource string

const (
	RoomSourceSync   RoomSource = "sync"
	RoomSourceIngest RoomSource = "ingest"
)

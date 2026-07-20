package conversationmodel

// Channel is the platform medium a conversation lives on. It is a closed enum
// keyed on platform identity (the granularity users recognize and filter by —
// "Slack", "Email"), NOT on semantic category; deliberate decision 2026-07-02,
// matching the prevailing industry shape (Chatwoot channel_type, Intercom
// source.type). Multiple connectors may serve one channel (gmail and outlook
// both feed email). meeting/adhoc are connector-less time-bounded spoken media.
type Channel string

const (
	ChannelSlack    Channel = "slack"
	ChannelEmail    Channel = "email"
	ChannelLinkedin Channel = "linkedin"
	ChannelSms      Channel = "sms"
	// ChannelTeamsChat is Microsoft Teams CHAT (messaging). Where a provider
	// brand spans both a chat product and a meeting platform, the channel is
	// suffixed -chat / -meeting so the two media never collide (renamed from
	// "teams" 2026-07-12, before any Teams connector shipped).
	ChannelTeamsChat Channel = "teams-chat"
	ChannelTelegram  Channel = "telegram"
	ChannelWhatsapp  Channel = "whatsapp"
	// ChannelMeeting is the generic meeting fallback (connector-less recordings,
	// platforms classifyPlatform doesn't recognize). The *-meeting channels below
	// are the connector-backed meeting platforms the Ava meeting-bot family
	// serves; each meeting conversation is tagged with its platform, derived
	// from the meeting URL.
	ChannelMeeting      Channel = "meeting"
	ChannelZoomMeeting  Channel = "zoom-meeting"
	ChannelGoogleMeet   Channel = "google-meet"
	ChannelTeamsMeeting Channel = "teams-meeting"
	ChannelWebexMeeting Channel = "webex-meeting"
	ChannelAdhoc        Channel = "adhoc"
	ChannelInstagram    Channel = "instagram"
	ChannelMessenger    Channel = "messenger"
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
	// The Ava meeting-bot family (Recall.ai-backed; the vendor name stays
	// adapter-internal). adhoc-meeting is the org-level connection bots are
	// dispatched from by meeting URL (user, agent or workflow); the calendar
	// keys are per-user OAuth grants Ava auto-joins scheduled meetings through.
	// Keys are brand-neutral (no "ava"/"recall") so the product brand and the
	// vendor can each change without touching identifiers.
	ConnectorKeyAdhocMeeting             ConnectorKey = "adhoc-meeting"
	ConnectorKeyGoogleCalendarMeeting    ConnectorKey = "google-calendar-meeting"
	ConnectorKeyMicrosoftCalendarMeeting ConnectorKey = "microsoft-calendar-meeting"
)

// ConnectorProvider says who operates the integration mechanics behind a
// connector: Proteos' own hand-coded integration against the provider's API
// (native), an account aggregated through Unipile's platform tenancy
// (unipile), or the Ava meeting-bot family (ava — Recall.ai-backed; the
// vendor stays adapter-internal). Computed on Connection reads from the
// connector, never stored; the UI groups the connector catalog by it.
type ConnectorProvider string

const (
	ConnectorProviderNative  ConnectorProvider = "native"
	ConnectorProviderUnipile ConnectorProvider = "unipile"
	ConnectorProviderAva     ConnectorProvider = "ava"
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
// person's contact address or a venue (room).
type RecipientKind string

const (
	RecipientKindContactAddress RecipientKind = "contact-address"
	RecipientKindRoom           RecipientKind = "room"
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

// ContactAddressKind is the identifier NAMESPACE of a contact address — not a
// channel: email and phone span channels (an email address is reachable from
// any email connection; a phone number serves sms AND whatsapp), while
// slack/messenger are provider namespaces whose ids only mean something inside
// a provider tenant (see ContactAddress.Scope). whatsapp exists only for
// non-phone JIDs (e.g. @lid) — standard WhatsApp JIDs canonicalize to phone.
type ContactAddressKind string

const (
	ContactAddressKindEmail     ContactAddressKind = "email"
	ContactAddressKindPhone     ContactAddressKind = "phone"
	ContactAddressKindSlack     ContactAddressKind = "slack"
	ContactAddressKindLinkedin  ContactAddressKind = "linkedin"
	ContactAddressKindTelegram  ContactAddressKind = "telegram"
	ContactAddressKindInstagram ContactAddressKind = "instagram"
	ContactAddressKindMessenger ContactAddressKind = "messenger"
	ContactAddressKindX         ContactAddressKind = "x"
	ContactAddressKindWhatsapp  ContactAddressKind = "whatsapp"
)

// ContactAddressSource records how a contact address row was established:
// provider directory sweep (sync), derived from ingested messages (ingest),
// attached by a user via the API (manual), or repointed during a contact merge
// (merge). Provenance only — it drives no prune or lifecycle (addresses are
// permanent identity).
type ContactAddressSource string

const (
	ContactAddressSourceSync   ContactAddressSource = "sync"
	ContactAddressSourceIngest ContactAddressSource = "ingest"
	ContactAddressSourceManual ContactAddressSource = "manual"
	ContactAddressSourceMerge  ContactAddressSource = "merge"
)

// ContactSource records how a contact row was minted — the same provenance
// axis (and value set) as ContactAddressSource, so the mint path stays
// visible: a provider directory sweep (sync), derived from corresponded
// messages (ingest — we actually talked), created by a user via the API
// (manual), or as the surviving side of a merge (merge). Provenance only;
// sync/ingest + no manual edits = the "thin" contact that deterministic
// auto-merge may fold away.
type ContactSource string

const (
	ContactSourceSync   ContactSource = "sync"
	ContactSourceIngest ContactSource = "ingest"
	ContactSourceManual ContactSource = "manual"
	ContactSourceMerge  ContactSource = "merge"
)

// ContactStatus is the lifecycle of a contact. merged rows are tombstones
// redirecting to the winner via merged_into_contact_id; erased rows are GDPR
// tombstones (PII blanked, id retained). Blocking is NOT a status — it is the
// orthogonal IsBlocked flag (a blocked contact stays active in the directory).
type ContactStatus string

const (
	ContactStatusActive   ContactStatus = "active"
	ContactStatusMerged   ContactStatus = "merged"
	ContactStatusArchived ContactStatus = "archived"
	ContactStatusErased   ContactStatus = "erased"
)

// ConsentStatus is the denormalized subject-consent state derived from the
// opt_in/opt_out permission events (the subject's axis only — our own blocks
// live on the orthogonal IsBlocked flag). unknown = never captured; consent is
// an opt-out gate, so unknown addresses stay contactable.
type ConsentStatus string

const (
	ConsentStatusUnknown  ConsentStatus = "unknown"
	ConsentStatusOptedIn  ConsentStatus = "opted_in"
	ConsentStatusOptedOut ConsentStatus = "opted_out"
)

// PermissionEventType is what a permission-ledger event did: the subject's
// consent axis (opt_in/opt_out — projects onto ConsentStatus) or our own
// suppression axis (block/unblock — projects onto IsBlocked). The axes are
// deliberately orthogonal: an opt_in never clears a block.
type PermissionEventType string

const (
	PermissionEventOptIn   PermissionEventType = "opt_in"
	PermissionEventOptOut  PermissionEventType = "opt_out"
	PermissionEventBlock   PermissionEventType = "block"
	PermissionEventUnblock PermissionEventType = "unblock"
)

// PermissionEventSource records where a permission event came from. import is
// special-cased in derivation: an imported opt_in never overrides an existing
// opt_out (imports cannot prove recency).
type PermissionEventSource string

const (
	PermissionEventSourceManual    PermissionEventSource = "manual"
	PermissionEventSourceImport    PermissionEventSource = "import"
	PermissionEventSourceLinkClick PermissionEventSource = "link_click"
	PermissionEventSourceReply     PermissionEventSource = "reply"
	PermissionEventSourceApi       PermissionEventSource = "api"
	PermissionEventSourceSystem    PermissionEventSource = "system"
)

// MergeProposalStatus is the review-queue lifecycle of a duplicate-contact
// proposal. All non-proposed states are terminal; superseded means one side
// was merged or erased through another path before review.
type MergeProposalStatus string

const (
	MergeProposalStatusProposed   MergeProposalStatus = "proposed"
	MergeProposalStatusApproved   MergeProposalStatus = "approved"
	MergeProposalStatusRejected   MergeProposalStatus = "rejected"
	MergeProposalStatusSuperseded MergeProposalStatus = "superseded"
)

// MergeInitiator says HOW a merge fired — deterministic business logic
// (system), an LLM adjudication (agent), or a human decision (user). WHO wrote
// the row is the standard created_by audit field, a separate axis.
type MergeInitiator string

const (
	MergeInitiatorSystem MergeInitiator = "system"
	MergeInitiatorAgent  MergeInitiator = "agent"
	MergeInitiatorUser   MergeInitiator = "user"
)

// ErasureRequestStatus is the GDPR erasure-request lifecycle: captured and
// awaiting the scrub runner (requested) or fully scrubbed (completed).
type ErasureRequestStatus string

const (
	ErasureRequestStatusRequested ErasureRequestStatus = "requested"
	ErasureRequestStatusCompleted ErasureRequestStatus = "completed"
)

// RoomSource records how a room directory row was populated — the venue half
// of the directory (separate entity-scoped type, same prune semantics as the
// old participant sweep): a provider directory sweep (sync — prunable when the
// provider stops returning it) or minted at ingest for a channel the sweep
// hasn't seen yet (ingest — survives pruning until the sweep confirms it, then
// flips to sync).
type RoomSource string

const (
	RoomSourceSync   RoomSource = "sync"
	RoomSourceIngest RoomSource = "ingest"
)

// ConversationFilterType discriminates a conversation filter's matching rule —
// the tagged-union key for ConversationFilterConfig (see
// conversation-filter-config.go). Ordered here by evaluation specificity:
// address (one exact person-address) beats domain (a whole email domain) beats
// role_based/automated (email heuristics) beats internal_conversations (the
// all-participants-internal classification) beats all (unconditional).
type ConversationFilterType string

const (
	FilterTypeAddress               ConversationFilterType = "address"
	FilterTypeDomain                ConversationFilterType = "domain"
	FilterTypeRoleBased             ConversationFilterType = "role_based"
	FilterTypeAutomated             ConversationFilterType = "automated"
	FilterTypeInternalConversations ConversationFilterType = "internal_conversations"
	FilterTypeAll                   ConversationFilterType = "all"
)

// ConversationFilterAction says what a matching filter does to the inbound
// message: block drops it at ingest (never persisted, no contact minted),
// allow exempts it — within a specificity class and scope, allow beats block,
// so an allow rule punches a hole through a broader block (block domain X,
// allow alice@X). internal_conversations rows are always block (an "allow
// internal" rule is meaningless — internal exclusion is already the most
// general class an address/domain allow overrides).
type ConversationFilterAction string

const (
	FilterActionBlock ConversationFilterAction = "block"
	FilterActionAllow ConversationFilterAction = "allow"
)

// FilterMatchOn says which side of the message an address/domain filter tests:
// the sender only (the default — HubSpot-style ingest suppression), or any
// participant (sender OR any To/Cc recipient — the EAC/Attio matching model,
// for "never ingest anything this address is on").
type FilterMatchOn string

const (
	FilterMatchOnSender         FilterMatchOn = "sender"
	FilterMatchOnAnyParticipant FilterMatchOn = "any_participant"
)

// AutomatedSignal names one detectable class of machine-generated email for
// the automated filter type — all are deterministic header signals (RFC 3834
// and industry practice) read from NormalizedInboundMessage.Headers:
// auto_submitted (Auto-Submitted present and not "no": auto-replies,
// out-of-office), bulk (Precedence: bulk/junk/list), mailing_list (List-Id or
// List-Unsubscribe present), bounce (empty Return-Path <> or a
// MAILER-DAEMON/postmaster sender: DSNs), auto_response_suppress (Microsoft's
// X-Auto-Response-Suppress: OOO/read receipts/delivery reports).
type AutomatedSignal string

const (
	AutomatedSignalAutoSubmitted        AutomatedSignal = "auto_submitted"
	AutomatedSignalBulk                 AutomatedSignal = "bulk"
	AutomatedSignalMailingList          AutomatedSignal = "mailing_list"
	AutomatedSignalBounce               AutomatedSignal = "bounce"
	AutomatedSignalAutoResponseSuppress AutomatedSignal = "auto_response_suppress"
)

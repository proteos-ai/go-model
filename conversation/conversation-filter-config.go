package conversationmodel

import (
	"encoding/json"
	"fmt"
)

// ConversationFilterConfig is the typed, per-type configuration of a
// ConversationFilter — a tagged union discriminated by the sibling FilterType
// (mirrors AgentListenerTriggerConfig). all and self_originated carry no
// parameters; the other types do.
type ConversationFilterConfig interface {
	isConversationFilterConfig()
	FilterType() ConversationFilterType
}

// AddressFilterConfig — match one exact canonical address, any channel kind.
// Value must be canonical (lowercased email, E.164 phone, raw provider id) —
// the write path canonicalizes via logic.BuildConversationFilterConfig, and
// the evaluator compares against the sender's DeriveContactAddressKeys output,
// so a {kind: phone, value: "+49..."} rule covers WhatsApp AND SMS (standard
// WhatsApp JIDs bridge to kind=phone). ContactAddressKey.Scope is ignored: a
// Slack-id rule applies across all the org's workspaces.
type AddressFilterConfig struct {
	Kind  ContactAddressKind `json:"kind"`
	Value string             `json:"value"`
	// MatchOn defaults to sender; any_participant also fires when the address
	// appears as a To/Cc recipient.
	MatchOn FilterMatchOn `json:"match_on"`
}

func (AddressFilterConfig) isConversationFilterConfig()        {}
func (AddressFilterConfig) FilterType() ConversationFilterType { return FilterTypeAddress }

// DomainFilterConfig — match a whole email domain (equals on the parsed
// domain part, lowercase, no leading '@'). Deliberately equals, not
// substring: contains-matching is a false-positive footgun ("x.com" matches
// "max.company.de") and no vendor ships it for ingest suppression.
type DomainFilterConfig struct {
	Domain string `json:"domain"`
	// MatchOn defaults to sender; any_participant also tests recipient domains.
	MatchOn FilterMatchOn `json:"match_on"`
}

func (DomainFilterConfig) isConversationFilterConfig()        {}
func (DomainFilterConfig) FilterType() ConversationFilterType { return FilterTypeDomain }

// TitleKeywordFilterConfig — match when the conversation's title/subject
// contains any of these keywords (case-insensitive substring; keywords are
// stored lowercased by the write path). At ingest it tests the message
// Subject (meetings carry the meeting title as Subject on every utterance);
// at the meeting pre-join gate it tests the calendar event title — the
// Fireflies-style "don't record 1:1 / Personal meetings" rule.
type TitleKeywordFilterConfig struct {
	Keywords []string `json:"keywords"`
}

func (TitleKeywordFilterConfig) isConversationFilterConfig()        {}
func (TitleKeywordFilterConfig) FilterType() ConversationFilterType { return FilterTypeTitleKeyword }

// RoleBasedFilterConfig — match when the sender's email local-part is a role
// mailbox (info@, noreply@, postmaster@ …). Empty Prefixes = the built-in
// logic.DefaultRolePrefixes set; non-empty overrides it. Email-only; the
// connection's own identity is exempt (an info@ mailbox's outbound-ingested
// mail must not self-drop).
type RoleBasedFilterConfig struct {
	Prefixes []string `json:"prefixes,omitempty"`
}

func (RoleBasedFilterConfig) isConversationFilterConfig()        {}
func (RoleBasedFilterConfig) FilterType() ConversationFilterType { return FilterTypeRoleBased }

// AutomatedFilterConfig — match machine-generated email by deterministic
// header signals (see AutomatedSignal). Empty Signals = all signals. Only
// meaningful on email connections whose ingestor populates
// NormalizedInboundMessage.Headers; inert elsewhere. The connection's own
// identity is exempt, like role_based.
type AutomatedFilterConfig struct {
	Signals []AutomatedSignal `json:"signals,omitempty"`
}

func (AutomatedFilterConfig) isConversationFilterConfig()        {}
func (AutomatedFilterConfig) FilterType() ConversationFilterType { return FilterTypeAutomated }

// SelfOriginatedFilterConfig — match when the conversation originates from
// the connection's own identity. No parameters. At ingest: the sender is the
// connection self (outbound-ingested mail, own Slack bot, …). At the meeting
// pre-join gate: the calendar event's organizer is the connected calendar's
// own account. Its allow form is the "organized by me" auto-join composition:
// a connection-scoped all-block plus a self_originated allow joins only
// meetings the calendar owner organizes.
type SelfOriginatedFilterConfig struct{}

func (SelfOriginatedFilterConfig) isConversationFilterConfig() {}
func (SelfOriginatedFilterConfig) FilterType() ConversationFilterType {
	return FilterTypeSelfOriginated
}

// InternalParticipantFilterConfig — match when at least ONE participant other
// than the connection self has an email domain in Domains (any-internal, the
// counterpart to InternalConversationsFilterConfig's all-internal). The org
// declares what "internal" means via the explicit domain list, exactly as it
// does for internal_conversations. Its allow form is the "with teammates"
// auto-join composition: a connection-scoped all-block plus an
// internal_participant allow joins only meetings with at least one teammate.
type InternalParticipantFilterConfig struct {
	Domains []string `json:"domains"`
}

func (InternalParticipantFilterConfig) isConversationFilterConfig() {}
func (InternalParticipantFilterConfig) FilterType() ConversationFilterType {
	return FilterTypeInternalParticipant
}

// InternalConversationsFilterConfig — drop a message only when the sender AND
// every To/Cc recipient are on one of these email domains (all-participants-
// internal, the Salesforce EAC/Copper classification). A recipient with no
// derivable email is NOT internal, so customer threads that cc a teammate
// survive — this is the safe alternative to blocklisting your own domain.
// Action is always block.
type InternalConversationsFilterConfig struct {
	Domains []string `json:"domains"`
}

func (InternalConversationsFilterConfig) isConversationFilterConfig() {}
func (InternalConversationsFilterConfig) FilterType() ConversationFilterType {
	return FilterTypeInternalConversations
}

// AllFilterConfig — match every message unconditionally. No parameters. The
// scope-control primitive: connection all-allow bypasses global filters,
// connection all-block pauses that connection's ingest.
type AllFilterConfig struct{}

func (AllFilterConfig) isConversationFilterConfig()        {}
func (AllFilterConfig) FilterType() ConversationFilterType { return FilterTypeAll }

// MarshalFilterConfig encodes a variant to its stored (JSONB) bare shape. The
// discriminator lives in the sibling filter_type column, so no {kind,data}
// envelope is needed. nil → '{}'.
func MarshalFilterConfig(config ConversationFilterConfig) (json.RawMessage, error) {
	if config == nil {
		return json.RawMessage("{}"), nil
	}
	return json.Marshal(config)
}

// DecodeFilterConfig rebuilds the typed variant from the stored bare JSON plus
// the discriminating filter type. Empty/absent config decodes to the empty
// variant for that type (all carries none).
func DecodeFilterConfig(filterType ConversationFilterType, raw []byte) (ConversationFilterConfig, error) {
	switch filterType {
	case FilterTypeAddress:
		config := AddressFilterConfig{}
		if err := unmarshalConfig(raw, &config); err != nil {
			return nil, err
		}
		return config, nil
	case FilterTypeDomain:
		config := DomainFilterConfig{}
		if err := unmarshalConfig(raw, &config); err != nil {
			return nil, err
		}
		return config, nil
	case FilterTypeTitleKeyword:
		config := TitleKeywordFilterConfig{}
		if err := unmarshalConfig(raw, &config); err != nil {
			return nil, err
		}
		return config, nil
	case FilterTypeRoleBased:
		config := RoleBasedFilterConfig{}
		if err := unmarshalConfig(raw, &config); err != nil {
			return nil, err
		}
		return config, nil
	case FilterTypeSelfOriginated:
		return SelfOriginatedFilterConfig{}, nil
	case FilterTypeInternalParticipant:
		config := InternalParticipantFilterConfig{}
		if err := unmarshalConfig(raw, &config); err != nil {
			return nil, err
		}
		return config, nil
	case FilterTypeAutomated:
		config := AutomatedFilterConfig{}
		if err := unmarshalConfig(raw, &config); err != nil {
			return nil, err
		}
		return config, nil
	case FilterTypeInternalConversations:
		config := InternalConversationsFilterConfig{}
		if err := unmarshalConfig(raw, &config); err != nil {
			return nil, err
		}
		return config, nil
	case FilterTypeAll:
		return AllFilterConfig{}, nil
	default:
		return nil, fmt.Errorf("unknown filter type %q", filterType)
	}
}

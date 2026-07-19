package conversationmodel

import (
	"time"

	"go.proteos.ai/model/common"
)

// ContactAddress is ONE canonical DIGITAL endpoint to reach a contact — an
// email address, a phone number, a Slack user id, a LinkedIn URN. It replaces
// the old connection-scoped participant directory row with an org-canonical
// one: rows are unique per (org, kind, scope, value), and that unique key IS
// the deterministic dedup backbone — resolution is "canonicalize, point-probe,
// hit = same person".
//
// Every row is an atomic canonical scalar; structured real-world locations
// (postal/billing/shipping) are a future contact_physical_address sibling,
// never a kind here — so Value stays a scalar forever.
//
// Addresses are permanent identity: there is no last_seen and no prune. A
// person departing a provider directory does not delete the address (we
// did/could reach them; dedup and history depend on the row).
type ContactAddress struct {
	Id        string             `json:"id"`
	OrgId     string             `json:"org_id"`
	ContactId string             `json:"contact_id"`
	Kind      ContactAddressKind `json:"kind" sortable:""`
	// Scope qualifies identifiers that are only unique within a provider
	// tenant — the Slack workspace id, the Messenger page id. Empty for global
	// kinds (email, phone, linkedin, …). Part of the dedup key.
	Scope string `json:"scope,omitempty"`
	// Value is the CANONICAL scalar form (lowercased email, E.164 phone, the
	// opaque provider id) produced by logic.DeriveContactAddressKeys. It is the
	// resolution lookup key and is never structured.
	Value string `json:"value" sortable:""`
	// RawValue is the identifier as observed on the wire (the WhatsApp JID, the
	// original casing) — display/debug only, never the lookup key.
	RawValue string `json:"raw_value,omitempty"`
	// Name is the per-channel display name the provider showed for this
	// address (was participant.name). Merges never-downgrade.
	Name string `json:"name" sortable:""`
	// ConsentStatus is the address-level denormalization of the subject's
	// opt_in/opt_out permission events ("this email opted out").
	ConsentStatus    ConsentStatus `json:"consent_status"`
	ConsentUpdatedAt *time.Time    `json:"consent_updated_at,omitempty"`
	// IsBlocked is the address-level suppression flag ("block this number" —
	// the person is otherwise contactable). Orthogonal to consent.
	IsBlocked bool                 `json:"is_blocked"`
	Source    ContactAddressSource `json:"source"`
	// Metadata is per-channel provider enrichment (the Slack handle today;
	// avatar/title/profile-url as connectors expose them). Cold-written — only
	// on sync sweeps and first sight, never per message. Heterogeneous by kind,
	// which is exactly what makes a bag the right shape here.
	Metadata  map[string]any `json:"metadata"`
	CreatedAt time.Time      `json:"created_at" sortable:""`
	CreatedBy common.UserRef `json:"created_by"`
	UpdatedAt time.Time      `json:"updated_at" sortable:""`
	UpdatedBy common.UserRef `json:"updated_by"`
}

// Ref projects the address (plus its resolved contact identity) onto the
// inline snapshot shape carried by messages, reactions, and rosters. The
// caller supplies the contact-level facts (platform user) since they live on
// Contact, not here.
func (address ContactAddress) Ref(platformUser *common.UserRef) ContactRef {
	ref := ContactRef{
		ExternalId:   address.RawValueOrValue(),
		Name:         address.Name,
		ContactId:    address.ContactId,
		PlatformUser: platformUser,
	}
	if address.Kind == ContactAddressKindEmail {
		ref.Email = address.Value
	}
	return ref
}

// RawValueOrValue returns the wire-facing identifier: the raw observed form
// when we have it (connector send paths expect the provider's own id — the
// JID, not the derived E.164), the canonical value otherwise.
func (address ContactAddress) RawValueOrValue() string {
	if address.RawValue != "" {
		return address.RawValue
	}
	return address.Value
}

// ContactAddressKey is one canonical identity key — the (kind, scope, value)
// triple the unique index dedups on, plus the raw wire form for display/send.
// Produced by the canonicalization rules (logic.DeriveContactAddressKeys) and
// consumed by the resolution point-probe.
type ContactAddressKey struct {
	Kind     ContactAddressKind `json:"kind"`
	Scope    string             `json:"scope,omitempty"`
	Value    string             `json:"value"`
	RawValue string             `json:"raw_value,omitempty"`
}

// Key projects an address row back onto its identity key.
func (address ContactAddress) Key() ContactAddressKey {
	return ContactAddressKey{
		Kind:     address.Kind,
		Scope:    address.Scope,
		Value:    address.Value,
		RawValue: address.RawValue,
	}
}

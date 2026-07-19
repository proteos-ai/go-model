package conversationmodel

import (
	"time"

	"go.proteos.ai/model/common"
)

// ContactPermissionEvent is one entry in the append-only ledger of everything
// that changes our permission to contact a person, from BOTH authorities: the
// subject's consent (opt_in/opt_out) and our own suppression decisions
// (block/unblock). Named "permission" — the superset — so GDPR's legal term
// "consent" stays precise for the one axis it names.
//
// The ledger is the legal source of truth (GDPR Art. 7(1) demonstrability);
// the ConsentStatus and IsBlocked fields on Contact/ContactAddress are derived
// projections of it by event type. Append-only: rows are never updated, so
// there is no updated_* audit.
type ContactPermissionEvent struct {
	Id        string `json:"id"`
	OrgId     string `json:"org_id"`
	ContactId string `json:"contact_id"`
	// ContactAddressId scopes the event to one address ("this email opted
	// out"); empty = contact-level (all channels).
	ContactAddressId string              `json:"contact_address_id,omitempty"`
	Event            PermissionEventType `json:"event"`
	// Basis is the legal basis the event asserts (consent, legitimate_interest,
	// contract, erasure, …). Free-form string, snake_case by convention.
	Basis  string                `json:"basis,omitempty"`
	Source PermissionEventSource `json:"source"`
	// OccurredAt is when the consent/suppression act actually happened — the
	// derivation ordering key. Distinct from CreatedAt (when we recorded it):
	// an imported historical opt-out occurred long before it was recorded.
	OccurredAt time.Time `json:"occurred_at"`
	Note       string    `json:"note,omitempty"`
	// Evidence is the proof trail: message id, unsubscribe-link token, import
	// batch reference — whatever demonstrates the event really happened.
	Evidence  map[string]any `json:"evidence,omitempty"`
	CreatedAt time.Time      `json:"created_at"`
	CreatedBy common.UserRef `json:"created_by"`
}

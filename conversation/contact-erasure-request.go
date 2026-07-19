package conversationmodel

import (
	"time"

	"go.proteos.ai/model/common"
)

// ContactErasureRequest is a GDPR Art. 17 right-to-erasure request — a
// STATEFUL job (requested → completed), deliberately not a permission event:
// it has a lifecycle, an SLA, suppression hashes, and a legal-hold
// interaction. v1 captures the obligation and immediately blocks the person
// for contact (a system opt_out permission event is emitted on create); the
// PII-scrub runner that drains open requests is a follow-up.
//
// Append-only shape: created_at/created_by are the request time + requester;
// completed_at (not updated_*) marks the scrub.
type ContactErasureRequest struct {
	Id        string `json:"id"`
	OrgId     string `json:"org_id"`
	ContactId string `json:"contact_id"`
	// IdentifierHashes are per-org-salted HMACs of the contact's canonical
	// address values at request time ([{kind, hmac}]) — the suppression keys
	// that survive the PII delete so background syncs don't silently resurrect
	// the person.
	IdentifierHashes []ErasureIdentifierHash `json:"identifier_hashes"`
	Status           ErasureRequestStatus    `json:"status" sortable:""`
	CompletedAt      *time.Time              `json:"completed_at,omitempty"`
	// Stats is the scrub outcome (rows touched per table); empty until the
	// runner completes the request.
	Stats     map[string]any `json:"stats,omitempty"`
	CreatedAt time.Time      `json:"created_at" sortable:""`
	CreatedBy common.UserRef `json:"created_by"`
}

// ErasureIdentifierHash is one suppression key: the address kind plus the
// per-org-salted HMAC of its canonical value.
type ErasureIdentifierHash struct {
	Kind ContactAddressKind `json:"kind"`
	Hmac string             `json:"hmac"`
}

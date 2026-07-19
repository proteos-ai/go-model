package conversationmodel

import (
	"time"

	"go.proteos.ai/model/common"
)

// Contact is the actual PERSON an org communicates with — org-scoped and
// channel-agnostic. It aggregates every way to reach them (ContactAddress
// rows) and is the anchor for consent, blocking, and GDPR lifecycle. It is
// what the participant directory rows deduplicate into: one contact per human,
// however many channels they show up on.
//
// Distinct from the sales module's CRM `contact` entity (org module data in
// the record universe) — this is the platform's communication-identity layer.
//
// Deliberately lean: no metadata bag (named columns only — provider extras
// live on ContactAddress.Metadata, CRM custom fields on the sales-module
// contact) and no last_seen_at (recency is a message-layer fact,
// conversation.last_message_at — never a per-message write here).
type Contact struct {
	Id    string `json:"id"`
	OrgId string `json:"org_id"`
	// Name is the person's display name — filled from observations only while
	// empty; a manual edit sticks (HasManualEdits) and is never downgraded by
	// the resolution pipeline.
	Name   string        `json:"name" sortable:""`
	Status ContactStatus `json:"status" sortable:""`
	// IsBlocked is the person-level do-not-contact flag — OUR suppression
	// decision, orthogonal to Status and to the subject's consent axis. Only a
	// block/unblock permission event moves it; an opt_in never clears it.
	IsBlocked bool `json:"is_blocked"`
	// HasLegalHold parks GDPR erasure: RequestErasure is rejected while set.
	HasLegalHold bool `json:"has_legal_hold"`
	// ConsentStatus is the person-level denormalization of the subject's
	// opt_in/opt_out permission events (the ledger is the source of truth).
	ConsentStatus    ConsentStatus `json:"consent_status"`
	ConsentUpdatedAt *time.Time    `json:"consent_updated_at,omitempty"`
	// PlatformUser is the resolved platform identity when this person is a
	// colleague (email match against the account directory); nil otherwise.
	// A person-fact, so it lives here — not on the per-channel address.
	PlatformUser *common.UserRef `json:"platform_user,omitempty"`
	// MergedIntoContactId is the tombstone redirect: set when this contact lost
	// a merge (Status=merged). One hop — merges flatten existing chains so a
	// stale snapshot id resolves in a single redirect. Empty = not merged.
	MergedIntoContactId string        `json:"merged_into_contact_id,omitempty"`
	Source              ContactSource `json:"source"`
	// HasManualEdits flips true on any user-driven update and is the thinness
	// signal for auto-merge safety: only never-edited, pipeline-minted
	// (sync/ingest) contacts may be folded away without human review.
	HasManualEdits bool `json:"has_manual_edits"`
	// Addresses embeds the contact's reachable endpoints on reads (the
	// aggregate is small); persisted in contact_address rows, never as JSONB.
	Addresses []ContactAddress `json:"addresses,omitempty"`
	CreatedAt time.Time        `json:"created_at" sortable:""`
	CreatedBy common.UserRef   `json:"created_by"`
	UpdatedAt time.Time        `json:"updated_at" sortable:""`
	UpdatedBy common.UserRef   `json:"updated_by"`
}

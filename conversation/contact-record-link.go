package conversationmodel

import (
	"time"

	"go.proteos.ai/model/common"
)

// ContactRecordLink is a directed edge from a Contact (the platform's
// communication-identity layer) to a business record in the data-service
// (EntitySlug + RecordId). It is how one person attaches to the domain
// entities they play a part in — the contact on a company, the candidate in
// recruiting. Deliberately a BARE edge: the relationship's semantics come
// from the target entity itself, so there is no role/label column.
//
// Two provenances share the table (Source):
//
//   - manual: a user's explicit link. Immutable once created, removed via
//     unlink — the KnowledgeRecordLink shape.
//   - record: THE binding conversation-service maintains from the record's
//     contact-address attributes. Exactly one per record (partial unique
//     index); it REPOINTS when the record's identity moves to another
//     contact and REFRESHES AddressKeys on every observation, so it carries
//     updated_* audit. AddressKeys are the canonical keys the record
//     contributed on its last observation — the detach ledger: a key the
//     record dropped is detached from the contact when nothing else still
//     carries it. A binding is never unlinked by hand (409
//     contact_record_link_bound); clearing the record's attributes unbinds it.
//
// The target record is NOT validated at creation — a dangling link surfaces
// lazily (the record fetch 404s). Contacts merge: ExecuteMerge repoints the
// loser's links to the winner (a manual edge the winner already had is
// upgraded in place when the loser's edge was the binding).
type ContactRecordLink struct {
	Id         string                  `json:"id" sortable:""`
	OrgId      string                  `json:"org_id"`
	ContactId  string                  `json:"contact_id" sortable:""`
	EntitySlug string                  `json:"entity_slug" sortable:""`
	RecordId   string                  `json:"record_id" sortable:""`
	Source     ContactRecordLinkSource `json:"source" sortable:""`
	// AddressKeys are empty for manual links.
	AddressKeys []ContactAddressKey `json:"address_keys"`
	CreatedAt   time.Time           `json:"created_at" sortable:""`
	CreatedBy   common.UserRef      `json:"created_by"`
	UpdatedAt   time.Time           `json:"updated_at" sortable:""`
	UpdatedBy   common.UserRef      `json:"updated_by"`
}

package conversationapi

import (
	"time"

	"go.proteos.ai/model/common"
	conversationmodel "go.proteos.ai/model/conversation"
)

// SearchContactsQuery filters the org's contact directory. Q matches
// (case-insensitive) against the contact name AND its address values/names —
// handled in the repo across the join. Status defaults to active (merged and
// erased tombstones never appear unless asked for explicitly).
type SearchContactsQuery struct {
	Q      *string                          `json:"q" form:"q"`
	Status *conversationmodel.ContactStatus `json:"status" form:"status"`
	common.Pagination
	common.Sorting
}

type SearchContactsResponse struct {
	Meta common.ResponseMeta         `json:"meta"`
	Data []conversationmodel.Contact `json:"data"`
}

// UpdateContactRequest is the PATCH surface — the user-editable subset only.
// Any accepted change flips HasManualEdits (the auto-merge thinness gate).
// Status accepts active|archived (merged/erased are lifecycle outcomes, never
// set directly); blocking goes through the block/unblock endpoints so it is
// ledgered, never patched.
type UpdateContactRequest struct {
	Name         *string                          `json:"name"`
	Status       *conversationmodel.ContactStatus `json:"status"`
	HasLegalHold *bool                            `json:"has_legal_hold"`
}

// MergeContactsRequest folds source_contact_id INTO the path contact (the path
// contact is the winner).
type MergeContactsRequest struct {
	SourceContactId string `json:"source_contact_id" binding:"required"`
}

// RecordPermissionEventRequest appends one entry to the contact's permission
// ledger. ContactAddressId scopes it to one address; empty = contact-level.
// OccurredAt defaults to now; Source defaults to api. block/unblock events are
// accepted here too (they flip IsBlocked), though the dedicated block/unblock
// endpoints are the ergonomic path.
type RecordPermissionEventRequest struct {
	ContactAddressId string                                  `json:"contact_address_id"`
	Event            conversationmodel.PermissionEventType   `json:"event" binding:"required"`
	Basis            string                                  `json:"basis"`
	Source           conversationmodel.PermissionEventSource `json:"source"`
	OccurredAt       *time.Time                              `json:"occurred_at"`
	Note             string                                  `json:"note"`
	Evidence         map[string]any                          `json:"evidence"`
}

// BlockContactRequest optionally scopes a block/unblock to one address; empty
// = person-level. Note rides into the ledger event.
type BlockContactRequest struct {
	ContactAddressId string `json:"contact_address_id"`
	Note             string `json:"note"`
}

// SearchContactMergeProposalsQuery pages the duplicate review queue.
type SearchContactMergeProposalsQuery struct {
	Status *conversationmodel.MergeProposalStatus `json:"status" form:"status"`
	common.Pagination
	common.Sorting
}

type SearchContactMergeProposalsResponse struct {
	Meta common.ResponseMeta                      `json:"meta"`
	Data []conversationmodel.ContactMergeProposal `json:"data"`
}

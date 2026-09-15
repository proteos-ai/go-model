package conversationapi

import (
	"go.proteos.ai/model/common"
	conversationmodel "go.proteos.ai/model/conversation"
)

// CreateContactRecordLinkRequest links one contact to one business record.
// The contact must exist (and not be a merged/erased tombstone); the record
// is deliberately not validated — a dangling link surfaces lazily.
type CreateContactRecordLinkRequest struct {
	ContactId  string `json:"contact_id" binding:"required"`
	EntitySlug string `json:"entity_slug" binding:"required"`
	RecordId   string `json:"record_id" binding:"required"`
}

// GetManyContactRecordLinksQuery filters the org's contact-record edges:
// by contact (a person's records), by entity_slug + record_id (a record's
// people), and/or by source (manual links vs record bindings).
type GetManyContactRecordLinksQuery struct {
	ContactId  *string                                    `json:"contact_id" form:"contact_id"`
	EntitySlug *string                                    `json:"entity_slug" form:"entity_slug"`
	RecordId   *string                                    `json:"record_id" form:"record_id"`
	Source     *conversationmodel.ContactRecordLinkSource `json:"source" form:"source"`
	common.Pagination
	common.Sorting
}

type GetManyContactRecordLinksResponse struct {
	Meta common.ResponseMeta                   `json:"meta"`
	Data []conversationmodel.ContactRecordLink `json:"data"`
}

// ResolveContactRecordLinksRequest — POST /contact-record-links/resolve — binds
// each record of ONE entity to exactly one contact by its canonical
// contact-address values (see ContactService.ResolveRecordLinks). data-service
// calls it inline for duplicate_policy=reject writes; the async consumer
// builds the same request (one record) from a record_contact_observation.
// DuplicatePolicy "" defaults to flag. Records are independent: one record's
// failure never aborts the others.
type ResolveContactRecordLinksRequest struct {
	EntitySlug      string                                  `json:"entity_slug" binding:"required"`
	DuplicatePolicy conversationmodel.RecordDuplicatePolicy `json:"duplicate_policy"`
	Records         []ResolveContactRecordLinkItem          `json:"records" binding:"required,dive"`
}

// ResolveContactRecordLinkItem is one record of the batch. Empty Addresses =
// the record carries no address any more → unbind.
type ResolveContactRecordLinkItem struct {
	RecordId  string                                   `json:"record_id" binding:"required"`
	Name      string                                   `json:"name"`
	Addresses []conversationmodel.RecordContactAddress `json:"addresses"`
}

// ContactRecordLinkResolution is the per-record outcome. Duplicate is set on
// a rejected resolution; ErrorCode/ErrorMessage on a failed one.
type ContactRecordLinkResolution struct {
	RecordId     string                                     `json:"record_id"`
	Outcome      conversationmodel.ContactRecordLinkOutcome `json:"outcome"`
	ContactId    string                                     `json:"contact_id,omitempty"`
	Link         *conversationmodel.ContactRecordLink       `json:"link,omitempty"`
	Duplicate    *conversationmodel.RecordDuplicateConflict `json:"duplicate,omitempty"`
	ErrorCode    string                                     `json:"error_code,omitempty"`
	ErrorMessage string                                     `json:"error_message,omitempty"`
}

type ResolveContactRecordLinksResponse struct {
	Data []ContactRecordLinkResolution `json:"data"`
}

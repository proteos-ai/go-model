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
// by contact (a person's records) or by entity_slug + record_id (a record's
// people).
type GetManyContactRecordLinksQuery struct {
	ContactId  *string `json:"contact_id" form:"contact_id"`
	EntitySlug *string `json:"entity_slug" form:"entity_slug"`
	RecordId   *string `json:"record_id" form:"record_id"`
	common.Pagination
	common.Sorting
}

type GetManyContactRecordLinksResponse struct {
	Meta common.ResponseMeta                   `json:"meta"`
	Data []conversationmodel.ContactRecordLink `json:"data"`
}

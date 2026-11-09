package dataapi

import (
	"go.proteos.ai/model/common"
	datamodel "go.proteos.ai/model/data"
)

// CreateRecordDraftRequest proposes one or many writes of one intent to one
// entity. Items missing a transaction_id get the 1-based index as a string.
// Source is free-form provenance the caller knows (agent_key, session_id,
// execution_id); the acting user is stamped as created_by regardless.
type CreateRecordDraftRequest struct {
	EntitySlug string                      `json:"entity_slug"`
	Intent     datamodel.RecordDraftIntent `json:"intent"`
	Items      []datamodel.RecordDraftItem `json:"items"`
	Summary    string                      `json:"summary"`
	Source     map[string]any              `json:"source,omitempty"`
}

// ApproveRecordDraftRequest: Items absent = apply as proposed; present = the
// reviewer's edited proposal (values changed, fields removed or added, rows
// dropped). The server derives approved vs approved_with_changes from the
// diff — the caller never names the outcome.
type ApproveRecordDraftRequest struct {
	Items []datamodel.RecordDraftItem `json:"items,omitempty"`
}

// RejectRecordDraftRequest: an optional reason, kept on the draft for the
// proposer to learn from.
type RejectRecordDraftRequest struct {
	Reason string `json:"reason,omitempty"`
}

// GetManyRecordDraftsQuery lists the org's drafts. Status defaults to
// proposed — the review queue. RecordId narrows to drafts whose items name
// that record (what a record page would show).
type GetManyRecordDraftsQuery struct {
	EntitySlug *string                      `json:"entity_slug" form:"entity_slug"`
	Status     *datamodel.RecordDraftStatus `json:"status" form:"status"`
	Intent     *datamodel.RecordDraftIntent `json:"intent" form:"intent"`
	RecordId   *string                      `json:"record_id" form:"record_id"`
	common.Pagination
	common.Sorting
}

type GetManyRecordDraftsResponse struct {
	Meta common.ResponseMeta     `json:"meta"`
	Data []datamodel.RecordDraft `json:"data"`
}

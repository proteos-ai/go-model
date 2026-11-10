package dataapi

import (
	"time"

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

// ApproveRecordDraftRequest: Items absent = apply the draft's items as they
// stand (the proposal plus any saved edits); present = save them first as the
// reviewer's edit (values changed, fields removed or added, rows dropped or
// added), then apply them. The server derives approved vs
// approved_with_changes from the diff against the proposer's latest version —
// the caller never names the outcome. ExpectedUpdatedAt is the updated_at the
// reviewer last saw: a draft that moved on since fails with 409
// record_draft_conflict instead of applying what they never saw.
type ApproveRecordDraftRequest struct {
	Items             []datamodel.RecordDraftItem `json:"items,omitempty"`
	ExpectedUpdatedAt *time.Time                  `json:"expected_updated_at,omitempty"`
}

// UpdateRecordDraftRequest replaces a PROPOSED draft's items — the complete
// set, like a PUT: an item left out is dropped, an item without a
// transaction_id is added (the server gives it a fresh one). Kind says who is
// speaking: `revised` (the proposer changing its proposal) or `edited` (a
// reviewer). Summary is optional (absent keeps it). ExpectedUpdatedAt is the
// updated_at the editor last saw: a draft that has moved on since fails with
// 409 record_draft_conflict instead of being overwritten; omit for last write
// wins.
type UpdateRecordDraftRequest struct {
	Items             []datamodel.RecordDraftItem      `json:"items"`
	Summary           *string                          `json:"summary,omitempty"`
	Kind              datamodel.RecordDraftVersionKind `json:"kind"`
	ExpectedUpdatedAt *time.Time                       `json:"expected_updated_at,omitempty"`
}

// GetRecordDraftVersionsResponse is a draft's history, oldest first, and
// what the reviewers changed so far: the current items diffed against the
// proposer's latest version (the same diff an approve would store).
type GetRecordDraftVersionsResponse struct {
	Data    []datamodel.RecordDraftVersion `json:"data"`
	Changes []datamodel.RecordDraftChange  `json:"changes"`
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

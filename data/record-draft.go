package datamodel

import (
	"time"

	"go.proteos.ai/model/common"
)

// RecordDraftIntent is WHAT a draft proposes to do to the entity's records.
// Record-only (the entity rides the draft), hence the short verbs — the
// page-button intent vocabulary (`create_record`, …) spans actions and
// workflows too and is mapped at the UI boundary.
type RecordDraftIntent string

const (
	RecordDraftIntentCreate RecordDraftIntent = "create"
	RecordDraftIntentUpdate RecordDraftIntent = "update"
	RecordDraftIntentDelete RecordDraftIntent = "delete"
	RecordDraftIntentUpsert RecordDraftIntent = "upsert"
)

// IsValid reports whether the intent is a known value.
func (intent RecordDraftIntent) IsValid() bool {
	switch intent {
	case RecordDraftIntentCreate, RecordDraftIntentUpdate, RecordDraftIntentDelete, RecordDraftIntentUpsert:
		return true
	}
	return false
}

// RecordDraftStatus is the review lifecycle of a draft: proposed (awaiting a
// verdict), approved (applied exactly as proposed), approved_with_changes
// (the reviewer edited the proposal, then it was applied — Changes says what
// they altered), rejected (never applied). The last three are terminal.
//
// approved vs approved_with_changes is DERIVED by the server from Changes —
// a resolver never names the outcome, they only send what to apply.
type RecordDraftStatus string

const (
	RecordDraftStatusProposed            RecordDraftStatus = "proposed"
	RecordDraftStatusApproved            RecordDraftStatus = "approved"
	RecordDraftStatusApprovedWithChanges RecordDraftStatus = "approved_with_changes"
	RecordDraftStatusRejected            RecordDraftStatus = "rejected"
)

// IsValid reports whether the status is a known value.
func (status RecordDraftStatus) IsValid() bool {
	switch status {
	case RecordDraftStatusProposed, RecordDraftStatusApproved, RecordDraftStatusApprovedWithChanges, RecordDraftStatusRejected:
		return true
	}
	return false
}

// IsTerminal reports whether the draft can still be approved or rejected.
func (status RecordDraftStatus) IsTerminal() bool {
	return status.IsValid() && status != RecordDraftStatusProposed
}

// RecordDraftItem is ONE proposed write, keyed like a batch transaction
// (transaction_id is the batch endpoints' name for "this item's key").
//
// Shape per intent — create: Data, no Id · update: Id + Data · delete: Id, no
// Data · upsert: Data, Id optional (present ⇒ update, absent ⇒ create, the
// RecordService.Upsert routing). `id` never rides inside Data.
type RecordDraftItem struct {
	TransactionId string         `json:"transaction_id"`
	Id            string         `json:"id,omitempty"`
	Data          map[string]any `json:"data,omitempty"`
}

// RecordDraftItemResultStatus is one item's apply outcome.
type RecordDraftItemResultStatus string

const (
	RecordDraftItemResultStatusSuccess RecordDraftItemResultStatus = "success"
	RecordDraftItemResultStatusError   RecordDraftItemResultStatus = "error"
)

// RecordDraftItemError carries the SAME facts the single-record endpoint
// would have answered with (code, message, status, details) — a draft row's
// failure is as actionable as a standalone write's.
type RecordDraftItemError struct {
	Code    string         `json:"code"`
	Message string         `json:"message"`
	Status  int            `json:"status,omitempty"`
	Details map[string]any `json:"details,omitempty"`
}

// RecordDraftItemResult is the outcome of applying one item — the batch
// result minus the record itself (the jsonb stays small; a client re-reads
// the record by id).
type RecordDraftItemResult struct {
	TransactionId string                      `json:"transaction_id"`
	Status        RecordDraftItemResultStatus `json:"status"`
	RecordId      string                      `json:"record_id,omitempty"`
	Error         *RecordDraftItemError       `json:"error,omitempty"`
}

// RecordDraftChangeKind names ONE thing the reviewer altered before applying.
type RecordDraftChangeKind string

const (
	// RecordDraftChangeItemDropped: a whole item (a row of a mass draft) was
	// left out — "do not do this one".
	RecordDraftChangeItemDropped RecordDraftChangeKind = "item_dropped"
	// RecordDraftChangeAttributeChanged: a proposed value was edited.
	RecordDraftChangeAttributeChanged RecordDraftChangeKind = "attribute_changed"
	// RecordDraftChangeAttributeRemoved: a proposed field was removed from the
	// item (the record keeps what it had).
	RecordDraftChangeAttributeRemoved RecordDraftChangeKind = "attribute_removed"
	// RecordDraftChangeAttributeAdded: the reviewer brought in a field the
	// proposer had not named.
	RecordDraftChangeAttributeAdded RecordDraftChangeKind = "attribute_added"
	// RecordDraftChangeItemAdded: the reviewer added a whole item (a row) the
	// proposer had not proposed. Applied carries its data; no Attribute.
	RecordDraftChangeItemAdded RecordDraftChangeKind = "item_added"
)

// RecordDraftVersionKind says how a version of a draft's items came to be.
// Who the actor was cannot tell it — an agent acts AS its user — so the
// caller names it: the proposer revising its proposal, or a reviewer editing.
type RecordDraftVersionKind string

const (
	// RecordDraftVersionKindProposed is version 1: what was proposed.
	RecordDraftVersionKindProposed RecordDraftVersionKind = "proposed"
	// RecordDraftVersionKindRevised is the proposer changing its proposal
	// (an agent's update_record_draft). It moves the baseline the reviewer's
	// changes are measured against.
	RecordDraftVersionKindRevised RecordDraftVersionKind = "revised"
	// RecordDraftVersionKindEdited is a reviewer's edit. Saves by the same
	// person update one version, so editing makes one version, not hundreds.
	RecordDraftVersionKindEdited RecordDraftVersionKind = "edited"
)

// IsValid reports whether the kind is a known value.
func (kind RecordDraftVersionKind) IsValid() bool {
	switch kind {
	case RecordDraftVersionKindProposed, RecordDraftVersionKindRevised, RecordDraftVersionKindEdited:
		return true
	}
	return false
}

// IsProposer reports whether the version speaks for the proposer — the
// baseline a reviewer's changes are diffed against.
func (kind RecordDraftVersionKind) IsProposer() bool {
	return kind == RecordDraftVersionKindProposed || kind == RecordDraftVersionKindRevised
}

// RecordDraftVersion is one step in a draft's history: the complete items
// (and summary) as they stood after that step.
type RecordDraftVersion struct {
	Id            string                 `json:"id"`
	OrgId         string                 `json:"org_id"`
	RecordDraftId string                 `json:"record_draft_id"`
	Version       int                    `json:"version"`
	Kind          RecordDraftVersionKind `json:"kind"`
	Summary       string                 `json:"summary"`
	Items         []RecordDraftItem      `json:"items"`
	CreatedAt     time.Time              `json:"created_at"`
	CreatedBy     common.UserRef         `json:"created_by"`
	UpdatedAt     time.Time              `json:"updated_at"`
	UpdatedBy     common.UserRef         `json:"updated_by"`
}

// RecordDraftChange is one entry of the reviewer's delta, diffed against the
// proposer's LATEST version — the proposal, or the proposer's last revision
// of it — never against an earlier attempt or a reviewer's own edit.
//
// Current is the record's value BEFORE the write — what the reviewer saw —
// captured for existing targets only (a create, an upsert without id and an
// item_dropped have none). HasCurrent says whether it was captured: a
// captured null is a real value ("the field was empty"), an uncaptured one
// is unknown. Current, Proposed and Applied are never omitted from the JSON.
//
// Current is STORED for the proposer's benefit (an agent learning from its
// reviews) but not yet exposed by the API: it is a record value, and a
// draft reader holding record-drafts:read may not be allowed to read the
// record or that attribute. Responses carry it redacted (HasCurrent false,
// Current null) until reads apply the record's own access rules.
type RecordDraftChange struct {
	TransactionId string                `json:"transaction_id"`
	Kind          RecordDraftChangeKind `json:"kind"`
	Attribute     string                `json:"attribute,omitempty"`
	HasCurrent    bool                  `json:"has_current"`
	Current       any                   `json:"current"`
	Proposed      any                   `json:"proposed"`
	Applied       any                   `json:"applied"`
}

// RecordDraft is a PROPOSED record write awaiting a human (or reviewing
// agent) verdict — the record sibling of a draft message. An agent, or any
// caller holding record-drafts:write, describes the write it would make; a
// resolver approves (applies as proposed), approves with changes (edits, then
// applies) or rejects. Approve executes through the normal record write path
// AS THE RESOLVER, so entity permissions, hooks, validation and events apply
// exactly as for a direct write.
//
// Items is the draft as it stands NOW: the proposal, as revised by the
// proposer and edited by reviewers while it is proposed (each save is a
// RecordDraftVersion; Version is the newest one's number). AppliedItems is
// what the last approve attempt sent (nil until one). Results is the latest
// outcome per transaction_id across attempts — a re-approve after a partial
// failure only re-sends the items that have not succeeded. Changes is the
// reviewer's delta against the proposer's LATEST version (empty ⇔ approved,
// non-empty ⇔ approved_with_changes).
// Source is the caller-supplied provenance (agent_key, session_id, …);
// CreatedBy is always the acting user.
type RecordDraft struct {
	Id           string                  `json:"id" sortable:""`
	OrgId        string                  `json:"org_id"`
	EntitySlug   string                  `json:"entity_slug" sortable:""`
	Intent       RecordDraftIntent       `json:"intent" sortable:""`
	Status       RecordDraftStatus       `json:"status" sortable:""`
	Summary      string                  `json:"summary"`
	Items        []RecordDraftItem       `json:"items"`
	Version      int                     `json:"version"`
	AppliedItems []RecordDraftItem       `json:"applied_items,omitempty"`
	Results      []RecordDraftItemResult `json:"results,omitempty"`
	Changes      []RecordDraftChange     `json:"changes,omitempty"`
	Source       map[string]any          `json:"source,omitempty"`
	Reason       *string                 `json:"reason,omitempty"`
	ResolvedBy   *common.UserRef         `json:"resolved_by,omitempty"`
	ResolvedAt   *time.Time              `json:"resolved_at,omitempty"`
	CreatedAt    time.Time               `json:"created_at" sortable:""`
	CreatedBy    common.UserRef          `json:"created_by"`
	UpdatedAt    time.Time               `json:"updated_at" sortable:""`
	UpdatedBy    common.UserRef          `json:"updated_by"`
	// ApplyingSince is set while an approve is writing the draft's records:
	// saves, rejects and other approves are refused (409
	// record_draft_applying) until it has stored its outcome.
	ApplyingSince *time.Time `json:"applying_since,omitempty"`
	// LastTransactionId is the highest numeric transaction_id the draft has
	// ever handed out, so an id dropped by a later edit is never reused.
	LastTransactionId int `json:"-"`
}

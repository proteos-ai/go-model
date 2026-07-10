package dataapi

import (
	"encoding/json"

	"go.proteos.ai/model/common"
)

// RecordEventPayload is the Message.Payload for a record lifecycle event on a
// RecordEventsTopic. It is intentionally self-contained — full records, not
// deltas — so consumers (after-hooks today; search/analytics/webhooks later)
// have every field without a re-fetch, and a deleted row remains fully available.
//
// The envelope carries the rest and it is not duplicated here: org →
// Message.OrgID, record id → Message.Key, entity slug → Message.Headers, verb →
// Message.Headers / the Message.Type suffix.
type RecordEventPayload struct {
	// Record is the full row. For create and delete it is the created/deleted
	// row; for update it is the full post-update row.
	Record json.RawMessage `json:"record"`
	// PreviousRecord is the full pre-update row, set only for update events so an
	// after_update hook can diff old vs new (mirrors before_update's
	// currentRecord). Empty for create and delete.
	PreviousRecord json.RawMessage `json:"previous_record,omitempty"`
	// Actor is who caused the write — ctx.GetUserRef() captured at publish time.
	// Carried explicitly (not derived from the record's updated_by) so token
	// identity stays independent of record contents and a delete's deleter is
	// available even though it is not stamped on the row.
	Actor common.UserRef `json:"actor"`
}

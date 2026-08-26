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
// Mirrors knowledge's KnowledgeRecordLink: immutable once created (created_*
// audit only, no update path), and the target record is NOT validated at
// creation — a dangling link surfaces lazily (the record fetch 404s) and is
// removed via unlink. Unlike knowledge nodes, contacts merge: ExecuteMerge
// repoints the loser's links to the winner.
type ContactRecordLink struct {
	Id         string         `json:"id" sortable:""`
	OrgId      string         `json:"org_id"`
	ContactId  string         `json:"contact_id" sortable:""`
	EntitySlug string         `json:"entity_slug" sortable:""`
	RecordId   string         `json:"record_id" sortable:""`
	CreatedAt  time.Time      `json:"created_at" sortable:""`
	CreatedBy  common.UserRef `json:"created_by"`
}

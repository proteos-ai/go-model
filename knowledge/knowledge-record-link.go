package knowledgemodel

import (
	"time"

	"go.proteos.ai/model/common"
)

// KnowledgeRecordLink is a directed edge from a KnowledgeNode to a business
// record in the data-service (identified by EntitySlug + RecordId). It is how a
// knowledge node is attached to the thing it is *about* — e.g. a node holding a
// meeting transcript linked to the `meeting-transcript` record, or research
// about a contact linked to the `contact` record.
//
// Unlike KnowledgeLink (node→node) this edge is immutable once created: it
// carries only created_* audit fields, no update path. These links are surfaced
// only on the knowledge node detail panel (a "Linked records" section), not in
// the whole-org graph layout; the panel fetches each record live for its title.
//
// Record existence is NOT validated at creation (v1): a dangling link surfaces
// lazily in the UI (the record fetch 404s) and is removed via DeleteOne.
type KnowledgeRecordLink struct {
	Id         string         `json:"id" sortable:""`
	OrgId      string         `json:"org_id"`
	NodeId     string         `json:"node_id" sortable:""`     // references KnowledgeNode.id
	EntitySlug string         `json:"entity_slug" sortable:""` // data-service entity slug
	RecordId   string         `json:"record_id" sortable:""`   // data-service record id
	CreatedAt  time.Time      `json:"created_at" sortable:""`
	CreatedBy  common.UserRef `json:"created_by"`
}

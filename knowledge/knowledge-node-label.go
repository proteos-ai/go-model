package knowledgemodel

import (
	"time"

	"go.proteos.ai/model/common"
)

// KnowledgeNodeLabel is the many-to-many association between a node and a
// label. Keyed by (org_id, node_id, label_id); immutable once created (detach
// removes the row), so it carries no UpdatedAt/UpdatedBy.
type KnowledgeNodeLabel struct {
	OrgId     string         `json:"org_id"`
	NodeId    string         `json:"node_id"`
	LabelId   string         `json:"label_id"`
	CreatedAt time.Time      `json:"created_at"`
	CreatedBy common.UserRef `json:"created_by"`
}

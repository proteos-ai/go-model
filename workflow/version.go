package workflowmodel

import (
	"time"

	"go.proteos.ai/model/common"
)

// WorkflowVersion is one immutable snapshot of a workflow's graph. Every save
// bumps the workflow's version int and inserts a row here; executions pin the
// version int (never a graph copy), so an execution's graph stays stable while
// the workflow is edited, and history/diff/rollback build on the same table.
type WorkflowVersion struct {
	OrgId       string         `json:"org_id"`
	WorkflowKey string         `json:"workflow_key"`
	Version     int            `json:"version" sortable:""`
	Graph       WorkflowGraph  `json:"graph"`
	CreatedAt   time.Time      `json:"created_at" sortable:""`
	CreatedBy   common.UserRef `json:"created_by"`
}

// Package workflowapi holds the wire request/response shapes for workflow-service's
// HTTP v1 surface.
package workflowapi

import (
	"time"

	workflowmodel "go.proteos.ai/model/workflow"

	"go.proteos.ai/model/common"
)

// CreateWorkflowRequest creates a workflow definition. Key is the immutable
// kebab-case identifier (composite-PK'd with org_id). Status starts active; the
// graph is validated server-side (exactly one trigger node → action nodes).
type CreateWorkflowRequest struct {
	Key         string                      `json:"key" validate:"required"`
	Name        string                      `json:"name" validate:"required"`
	Description string                      `json:"description"`
	Graph       workflowmodel.WorkflowGraph `json:"graph"`
}

// UpdateWorkflowRequest fully replaces the workflow's metadata + graph (a partial
// graph patch is ambiguous); version is bumped server-side. Status changes go
// through the dedicated pause/unpause endpoints, not here.
type UpdateWorkflowRequest struct {
	Name        string                      `json:"name" validate:"required"`
	Description string                      `json:"description"`
	Graph       workflowmodel.WorkflowGraph `json:"graph"`
}

// RunWorkflowRequest is the body of a manual "run now" call. DestinationNodeId
// turns the run into a partial "run until here" execution: only nodes on a
// path from the trigger to the destination (inclusive) execute.
type RunWorkflowRequest struct {
	DestinationNodeId string `json:"destination_node_id"`
}

type GetManyWorkflowsQuery struct {
	Name         *string `json:"name" form:"name" db:"name"`
	NameContains *string `json:"name[contains]" form:"name[contains]" db:"name" op:"contains"`
	Status       *string `json:"status" form:"status" db:"status"`
	common.Pagination
	common.Sorting
}

type GetManyWorkflowsResponse struct {
	Meta common.ResponseMeta      `json:"meta"`
	Data []workflowmodel.Workflow `json:"data"`
}

// ── Version history ─────────────────────────────────────────────────────────

// WorkflowVersionSummary is one row of the version-history list: everything
// the history panel needs without shipping each version's full graph.
type WorkflowVersionSummary struct {
	Version    int            `json:"version"`
	NodesTotal int            `json:"nodes_total"`
	CreatedAt  time.Time      `json:"created_at"`
	CreatedBy  common.UserRef `json:"created_by"`
}

// GetManyVersionsQuery pages through a workflow's version history (newest
// first).
type GetManyVersionsQuery struct {
	common.Pagination
}

type GetManyVersionsResponse struct {
	Meta common.ResponseMeta      `json:"meta"`
	Data []WorkflowVersionSummary `json:"data"`
}

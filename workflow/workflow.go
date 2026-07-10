// Package workflowmodel is the shared domain model for workflow-service — the
// first, deliberately-minimal cut of an n8n-compatible, Temporal-backed workflow
// engine.
//
// A Workflow is a persistent definition stored n8n-shaped (a node graph as
// JSONB); a WorkflowExecution is one firing of it. v1 supports a single trigger
// node (cron | manual | webhook | event) wired to a single agent action node, but
// the graph schema is already general so new node types never force a migration.
package workflowmodel

import (
	"time"

	"go.proteos.ai/model/common"
)

// WorkflowStatus is the lifecycle state of a Workflow definition. A paused
// workflow keeps its definition but suppresses scheduled/triggered firings;
// archived is terminal (the workflow becomes read-only).
type WorkflowStatus string

const (
	WorkflowStatusActive   WorkflowStatus = "active"
	WorkflowStatusPaused   WorkflowStatus = "paused"
	WorkflowStatusArchived WorkflowStatus = "archived"
)

// WorkflowStatuses is the canonical, ordered set of workflow statuses.
var WorkflowStatuses = []WorkflowStatus{
	WorkflowStatusActive,
	WorkflowStatusPaused,
	WorkflowStatusArchived,
}

// Workflow is a stored automation definition. It is keyed by an immutable
// kebab-case key, composite-PK'd by (org_id, key) like every other config
// resource (Agent/Tool/Prompt/…) and referenced by key (workflow_key) — a
// WorkflowExecution carries the key, and the key is immutable so that reference
// is stable. Graph holds the n8n-shaped node graph (persisted as JSONB).
type Workflow struct {
	OrgId       string         `json:"org_id"`
	Key         string         `json:"key" sortable:""`
	Name        string         `json:"name" sortable:""`
	Description string         `json:"description"`
	Status      WorkflowStatus `json:"status" sortable:""`
	Graph       WorkflowGraph  `json:"graph"`
	Version     int            `json:"version"`
	CreatedAt   time.Time      `json:"created_at" sortable:""`
	CreatedBy   common.UserRef `json:"created_by"`
	UpdatedAt   time.Time      `json:"updated_at" sortable:""`
	UpdatedBy   common.UserRef `json:"updated_by"`
}

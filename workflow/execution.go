package workflowmodel

import (
	"encoding/json"
	"time"

	"go.proteos.ai/model/common"
)

// ExecutionStatus is the coarse, platform-side status of one workflow firing.
// Per-node progress lives in the append-only NodeExecution rows; this status
// only tracks whether the run started and reached a terminal outcome.
type ExecutionStatus string

const (
	ExecutionStatusPending   ExecutionStatus = "pending"
	ExecutionStatusRunning   ExecutionStatus = "running"
	ExecutionStatusSucceeded ExecutionStatus = "succeeded"
	ExecutionStatusFailed    ExecutionStatus = "failed"
	ExecutionStatusCancelled ExecutionStatus = "cancelled"
	// ExecutionStatusSkipped is a node-execution-only status: the node never ran
	// because no items reached it (dead branch / empty input).
	ExecutionStatusSkipped ExecutionStatus = "skipped"
)

// ExecutionStatuses is the canonical, ordered set of execution statuses.
var ExecutionStatuses = []ExecutionStatus{
	ExecutionStatusPending,
	ExecutionStatusRunning,
	ExecutionStatusSucceeded,
	ExecutionStatusFailed,
	ExecutionStatusCancelled,
}

// IsExecutionTerminal reports whether a status is final (no further transitions).
func IsExecutionTerminal(status ExecutionStatus) bool {
	switch status {
	case ExecutionStatusSucceeded, ExecutionStatusFailed, ExecutionStatusCancelled, ExecutionStatusSkipped:
		return true
	default:
		return false
	}
}

// TriggerKind names what caused an execution to fire.
type TriggerKind string

const (
	TriggerKindSchedule TriggerKind = "schedule"
	TriggerKindManual   TriggerKind = "manual"
	TriggerKindWebhook  TriggerKind = "webhook"
	TriggerKindEvent    TriggerKind = "event"
	TriggerKindMessage  TriggerKind = "message"
	// TriggerKindConnector fires on a connector-service sync.item_changed event.
	TriggerKindConnector TriggerKind = "connector"
	// TriggerKindWorkflow marks a child execution started by a parent workflow's
	// execute-workflow intrinsic.
	TriggerKindWorkflow TriggerKind = "workflow"
)

// ExecutionTriggerContext records why and how an execution fired. It is a flat,
// kind-discriminated struct (persisted as JSONB) — Kind selects which optional
// fields are meaningful. Payload carries the trigger's data payload (the record
// for event triggers, the message for message triggers, the request body for
// webhooks) and seeds the workflow's first item.
type ExecutionTriggerContext struct {
	Kind TriggerKind `json:"kind"`
	// schedule
	ScheduledAt *time.Time `json:"scheduled_at,omitempty"`
	// manual — the user/actor who triggered the run
	Actor *common.UserRef `json:"actor,omitempty"`
	// webhook
	ReceivedAt *time.Time `json:"received_at,omitempty"`
	// event
	Topic     string `json:"topic,omitempty"`
	EventType string `json:"event_type,omitempty"`
	RecordId  string `json:"record_id,omitempty"`
	// message
	ConversationId string `json:"conversation_id,omitempty"`
	Channel        string `json:"channel,omitempty"`
	// connector
	ConnectorKey string `json:"connector_key,omitempty"`
	ConnectionId string `json:"connection_id,omitempty"`
	// workflow (child executions)
	ParentExecutionId string `json:"parent_execution_id,omitempty"`
	// The trigger's data payload (event record / message / webhook body).
	Payload json.RawMessage `json:"payload,omitempty"`
}

// WorkflowExecution is one firing of a Workflow — the platform-side projection
// of a Temporal workflow execution. It pins the immutable WorkflowVersion it
// ran against (never a graph copy); per-node results live in the append-only
// NodeExecution rows keyed by (execution_id, node_id, run_index).
type WorkflowExecution struct {
	OrgId              string                  `json:"org_id"`
	Id                 string                  `json:"id" sortable:""`
	WorkflowKey        string                  `json:"workflow_key"`
	WorkflowVersion    int                     `json:"workflow_version"`
	Status             ExecutionStatus         `json:"status" sortable:""`
	TriggerContext     ExecutionTriggerContext `json:"trigger_context"`
	TemporalWorkflowId string                  `json:"temporal_workflow_id,omitempty"`
	TemporalRunId      string                  `json:"temporal_run_id,omitempty"`
	Error              *ExecutionError         `json:"error,omitempty"`
	StartedAt          *time.Time              `json:"started_at,omitempty"`
	FinishedAt         *time.Time              `json:"finished_at,omitempty"`
	CreatedAt          time.Time               `json:"created_at" sortable:""`
	CreatedBy          common.UserRef          `json:"created_by"`
}

// NodeExecution is one run of one node within an execution — append-only, a new
// row per (node_id, run_index) so loops and retries never overwrite history.
// ItemCounts counts output items per port (always present even when the items
// themselves live in the payload store); Metadata carries node-specific extras
// (e.g. the agent session id).
type NodeExecution struct {
	OrgId       string          `json:"org_id"`
	ExecutionId string          `json:"execution_id"`
	NodeId      string          `json:"node_id"`
	RunIndex    int             `json:"run_index"`
	NodeName    string          `json:"node_name"`
	NodeType    NodeType        `json:"node_type"`
	Status      ExecutionStatus `json:"status"`
	InputCounts map[string]int  `json:"input_counts,omitempty"`
	ItemCounts  map[string]int  `json:"item_counts,omitempty"`
	Error       *ExecutionError `json:"error,omitempty"`
	Metadata    map[string]any  `json:"metadata,omitempty"`
	StartedAt   *time.Time      `json:"started_at,omitempty"`
	FinishedAt  *time.Time      `json:"finished_at,omitempty"`
	CreatedAt   time.Time       `json:"created_at"`
}

// ExecutionError is a structured failure reason (snake_case code + message),
// mirroring the platform CustomError shape. IsUserError marks business failures
// the user can fix (bad input, 4xx) as opposed to infrastructure faults.
type ExecutionError struct {
	Code        string `json:"code"`
	Message     string `json:"message"`
	IsUserError bool   `json:"is_user_error,omitempty"`
}

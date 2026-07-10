package workflowapi

import (
	"encoding/json"

	workflowmodel "go.proteos.ai/model/workflow"
)

// TestNodeInputSource selects where a standalone node test's input comes from:
// the mirrored input of the workflow's most recent execution (default), or the
// pinned output of the node's upstream neighbors (workflow.graph.pin_data).
const (
	TestNodeInputSourceLastExecution = "last_execution"
	TestNodeInputSourcePinned        = "pinned"
)

// TestNodeRequest runs ONE node ephemerally — no execution rows, real side
// effects. Node carries the editor's IN-MEMORY candidate definition (possibly
// unsaved), never reloaded from the stored graph, so edits are testable before
// saving.
type TestNodeRequest struct {
	Node        TestNodeCandidate `json:"node" validate:"required"`
	InputSource string            `json:"input_source" validate:"omitempty,oneof=last_execution pinned"`
}

// TestNodeCandidate is the candidate node definition under test. Parameters
// stay raw — Liquid resolves on the node host per item, exactly like a real
// run.
type TestNodeCandidate struct {
	Type           workflowmodel.NodeType `json:"type" validate:"required"`
	TypeVersion    int                    `json:"type_version"`
	Name           string                 `json:"name"`
	Parameters     json.RawMessage        `json:"parameters"`
	CredentialRefs map[string]string      `json:"credential_refs"`
}

// TestNodeResponse carries the full test round trip: the input that was fed in
// (so the panel can show it even when the node fails) and the produced output,
// items inline per port. Both sides are capped (IsTruncated) — this is an
// interactive editor call, not a data export.
type TestNodeResponse struct {
	Status            string                          `json:"status"` // succeeded | failed
	InputSource       string                          `json:"input_source"`
	HasPriorExecution bool                            `json:"has_prior_execution"`
	Input             map[string][]workflowmodel.Item `json:"input"`
	Output            map[string][]workflowmodel.Item `json:"output"`
	ItemCounts        map[string]int                  `json:"item_counts"`
	Error             *workflowmodel.ExecutionError   `json:"error,omitempty"`
	Metadata          map[string]any                  `json:"metadata,omitempty"`
	IsTruncated       bool                            `json:"is_truncated"`
}

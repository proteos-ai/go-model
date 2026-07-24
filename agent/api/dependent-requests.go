package agentapi

import "time"

// DependentRefKind names the resource type a dependents query is scoped to — the
// agent field the reference lives in.
type DependentRefKind string

const (
	DependentRefMcpServer DependentRefKind = "mcp_server"
	DependentRefPrompt    DependentRefKind = "prompt"
	DependentRefTool      DependentRefKind = "tool"
	DependentRefSkill     DependentRefKind = "skill"
)

// DependentAgent is one agent referencing a resource, with its provider-sync
// staleness relative to that resource: is_stale means the resource changed after
// the agent's last successful sync (or the agent was never synced), so the runtime
// still runs the agent against the resource's previous configuration.
type DependentAgent struct {
	AgentKey     string     `json:"agent_key"`
	Name         string     `json:"name"`
	IsStale      bool       `json:"is_stale"`
	LastSyncedAt *time.Time `json:"last_synced_at"`
}

// DependentSyncStatus values for DependentSyncResult.
const (
	DependentSyncStatusSynced = "synced"
	DependentSyncStatusFailed = "failed"
)

// DependentSyncResult reports the outcome of re-syncing one stale dependent agent.
type DependentSyncResult struct {
	AgentKey string `json:"agent_key"`
	Status   string `json:"status"` // "synced" | "failed"
	Error    string `json:"error,omitempty"`
}

// ListDependentsResponse / SyncDependentsResponse wrap the collections in the
// `data` envelope the other subresource list endpoints use.
type ListDependentsResponse struct {
	Data []DependentAgent `json:"data"`
}

type SyncDependentsResponse struct {
	Data []DependentSyncResult `json:"data"`
}

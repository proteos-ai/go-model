package agentmodel

import (
	"time"

	"go.proteos.ai/model/common"
)

// Session is one conversation with a configured Agent. It is a thin mutable
// parent over the append-only agent_session_event log: the events are the source
// of truth, the session carries only the identity, the agent it runs, a title,
// and the latest lifecycle status. The session ↔ external provider-session id
// mapping is NOT here — it lives in the provider adapter's own store. Keyed by id
// (a session is a graph node referenced by every event, so it is surrogate-keyed,
// not (org_id, slug)).
type Session struct {
	OrgId    string        `json:"org_id"`
	Id       string        `json:"id"`
	AgentKey string        `json:"agent_key"`
	Title    string        `json:"title" sortable:""`
	Status   SessionStatus `json:"status"`
	// Tools is the session's resolved tool manifest (the agent's tools + the surface's
	// client tools), snapshotted at create. It maps a custom tool's name to its kind so
	// the runtime can route an inbound custom_tool_use and tag the ledger tool_kind.
	// Empty for sessions opened with no client tools. Persisted as the `tools` jsonb.
	Tools []SessionTool `json:"tools"`
	// Usage is the running token total across the whole session, projected from the
	// span.model_request_end events (recomputed from the log on each). Persisted as
	// the `usage` jsonb; zero-valued until the first model request closes.
	Usage     SessionUsage   `json:"usage"`
	CreatedAt time.Time      `json:"created_at" sortable:""`
	CreatedBy common.UserRef `json:"created_by"`
	UpdatedAt time.Time      `json:"updated_at" sortable:""`
	UpdatedBy common.UserRef `json:"updated_by"`
}

// SessionUsage is the rolled-up token accounting for a session: the sum of every
// model request's usage plus the number of requests. It is a projection of the
// span.model_request_end event log (the events remain the source of truth), kept
// on the session row so listings and headers can show cost without replaying the
// log. TotalTokens is the derived input+output sum.
//
// Model is the model the usage is attributed to. Managed Agents does not report a
// model id on its span events, so it is seeded once at session create from the
// session's configured agent (Agent.ModelConfig.ModelId) and preserved across
// recomputes — it is not derived from the usage events.
type SessionUsage struct {
	Model                    string `json:"model"`
	InputTokens              int64  `json:"input_tokens"`
	OutputTokens             int64  `json:"output_tokens"`
	CacheReadInputTokens     int64  `json:"cache_read_input_tokens"`
	CacheCreationInputTokens int64  `json:"cache_creation_input_tokens"`
	TotalTokens              int64  `json:"total_tokens"`
	RequestCount             int64  `json:"request_count"`
}

// ModelUsageTotals is the raw aggregate of a session's model-request usage —
// the SUMs over the span.model_request_end payloads plus the request COUNT,
// before any derived field is computed. It is the plain data the repository's
// aggregate read returns and logic.BuildSessionUsage consumes; defining it here
// (not in spi-ports) keeps the pure logic layer free of any port dependency.
type ModelUsageTotals struct {
	InputTokens              int64
	OutputTokens             int64
	CacheReadInputTokens     int64
	CacheCreationInputTokens int64
	RequestCount             int64
}

// ToolKind resolves a custom tool's kind from the session manifest, or "" when the
// name is not in the manifest (an Anthropic builtin-toolset tool, or an unknown tool).
func (session Session) ToolKind(name string) ToolKind {
	for _, tool := range session.Tools {
		if tool.Key == name {
			return tool.Kind
		}
	}
	return ""
}

// SessionStatus is the latest lifecycle state of a session, projected from the
// session.status_* events. idle = ready for the next user message; running = a
// turn is in flight; terminated = ended; error = the last turn failed.
type SessionStatus string

const (
	SessionStatusIdle       SessionStatus = "idle"
	SessionStatusRunning    SessionStatus = "running"
	SessionStatusTerminated SessionStatus = "terminated"
	SessionStatusError      SessionStatus = "error"
)

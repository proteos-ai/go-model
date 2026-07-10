package agentapi

import (
	"encoding/json"

	agentmodel "go.proteos.ai/model/agent"
	"go.proteos.ai/model/common"
)

// CreateSessionRequest opens a session against a configured agent. The agent must
// already be synced to the provider; Title is optional. ClientTools is the launching
// surface's client-tool catalog (custom tools the browser renders/executes) — attached
// to this session only, so a surface that can't service them simply omits them.
type CreateSessionRequest struct {
	AgentKey    string                      `json:"agent_key" validate:"required"`
	Title       *string                     `json:"title,omitempty"`
	ClientTools []agentmodel.ClientToolSpec `json:"client_tools,omitempty"`
}

// AppendEventRequest posts a client event onto a session. Payload is raw JSON: it
// can only be decoded once Type is known, so the service decodes it via
// agentmodel.DecodeEventPayload (validating the shape per type) rather than the
// JSON binder. v1 accepts only client-postable user.* types (see
// logic.ValidateClientPostable).
type AppendEventRequest struct {
	Type agentmodel.EventType `json:"type" validate:"required"`
	// TurnId optionally lets the client supply the turn's id (a UUID) so it can
	// render the message optimistically under a stable identity that matches the
	// event echoed back over the stream. Omitted, the service mints one.
	TurnId  *string         `json:"turn_id,omitempty"`
	Payload json.RawMessage `json:"payload,omitempty"`
}

type GetManySessionsQuery struct {
	AgentKey *string `json:"agent_key" form:"agent_key" db:"agent_key"`
	common.Pagination
	common.Sorting
}

type GetManySessionsResponse struct {
	Meta common.ResponseMeta  `json:"meta"`
	Data []agentmodel.Session `json:"data"`
}

// GetSessionEventsQuery is a forward cursor over the event log. AfterSeq is the
// exclusive lower bound (the client's last-seen seq); Order defaults to asc
// (chronological). It reads from Postgres (the durable channel); the live channel
// is the SSE endpoint.
type GetSessionEventsQuery struct {
	AfterSeq int64  `json:"after_seq" form:"after_seq"`
	Limit    int    `json:"limit" form:"limit"`
	Order    string `json:"order" form:"order"` // asc | desc; default asc
	// Types optionally restricts the result to specific event types. Accepts a
	// comma-separated list (e.g. "agent.message,agent.tool_use"), translated to a
	// `type IN (...)` filter by the shared url-to-db helper; omitted, every type
	// is returned.
	Types string `json:"types" form:"types" db:"type"`
}

type GetSessionEventsResponse struct {
	Meta common.ResponseMeta       `json:"meta"`
	Data []agentmodel.SessionEvent `json:"data"`
}

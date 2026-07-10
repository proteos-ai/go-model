package agentmodel

import "encoding/json"

// EventPayload is the type-discriminated body of a SessionEvent. The concrete
// variant is selected by the event's Type (see DecodeEventPayload) — the same
// tagged-union pattern as ToolBinding. Single-content payloads are FLAT (no inner
// `type`); only message/result content slices self-tag (see session-content.go).
// On the way out the concrete value marshals itself; nothing carries the union
// discriminator because the sibling EventType already is it.
type EventPayload interface{ isEventPayload() }

// UserMessagePayload is a user message body (text + file blocks).
type UserMessagePayload struct {
	Content []ContentBlock `json:"content"`
}

func (UserMessagePayload) isEventPayload() {}

// AgentMessagePayload is an agent message body (text blocks only).
type AgentMessagePayload struct {
	Content []TextBlock `json:"content"`
}

func (AgentMessagePayload) isEventPayload() {}

// AgentArtifactsPayload carries the file(s) the agent produced in its sandbox this
// turn (reconciled from the session outputs), surfaced as downloadable attachments.
// Files are storage-service FileBlocks (the same value object user attachments use);
// a single turn can produce several.
type AgentArtifactsPayload struct {
	Files []FileBlock `json:"files"`
}

func (AgentArtifactsPayload) isEventPayload() {}

// OutcomeRubric is the grading rubric of a define_outcome kickoff: inline text or
// a reference to an uploaded file. Type is "text" | "file".
type OutcomeRubric struct {
	Type    string `json:"type"`
	Content string `json:"content,omitempty"`
	FileId  string `json:"file_id,omitempty"`
}

// DefineOutcomePayload is a user.define_outcome body — the task description, the
// grading rubric, and an optional iteration cap. Anthropic Managed Agents runs the
// grade/iterate loop server-side; this payload is forwarded to the provider.
type DefineOutcomePayload struct {
	Description   string        `json:"description"`
	Rubric        OutcomeRubric `json:"rubric"`
	MaxIterations *int          `json:"max_iterations,omitempty"`
}

func (DefineOutcomePayload) isEventPayload() {}

// OutcomeEvaluationPayload carries one grader-iteration span (start/ongoing/end).
// Result and Explanation are populated on the end event ("satisfied" |
// "needs_revision" | "max_iterations_reached" | "failed" | "interrupted").
type OutcomeEvaluationPayload struct {
	OutcomeId   string `json:"outcome_id,omitempty"`
	Iteration   int    `json:"iteration,omitempty"`
	Result      string `json:"result,omitempty"`
	Explanation string `json:"explanation,omitempty"`
}

func (OutcomeEvaluationPayload) isEventPayload() {}

// ReasoningPayload is an agent reasoning ("thinking") block. v1 whole-events
// carry no text/signature from the provider (those only arrive on the deferred
// delta channel), so both are typically empty until then.
type ReasoningPayload struct {
	Text      string `json:"text"`
	Signature string `json:"signature,omitempty"`
}

func (ReasoningPayload) isEventPayload() {}

// ToolUsePayload is a tool invocation. Id is the provider's tool_use event id —
// the correlation key a matching ToolResultPayload.ToolUseId points back to.
// ToolKind is the execution location (action | mcp | client), reusing ToolKind.
type ToolUsePayload struct {
	Id       string          `json:"id"`
	Name     string          `json:"name"`
	Input    json.RawMessage `json:"input"`
	ToolKind ToolKind        `json:"tool_kind"`
}

func (ToolUsePayload) isEventPayload() {}

// ToolResultPayload is the result of a tool invocation, linked to its use by
// ToolUseId. Outcome is "success" | "error".
type ToolResultPayload struct {
	ToolUseId string        `json:"tool_use_id"`
	Result    []ResultBlock `json:"result"`
	Outcome   string        `json:"outcome"`
}

func (ToolResultPayload) isEventPayload() {}

// ToolConfirmationPayload is a user's decision on a tool that requested approval.
type ToolConfirmationPayload struct {
	ToolUseId string `json:"tool_use_id"`
	Decision  string `json:"decision"`
}

func (ToolConfirmationPayload) isEventPayload() {}

// ContextCompactedPayload announces that earlier context was summarized.
type ContextCompactedPayload struct {
	Summary string `json:"summary"`
}

func (ContextCompactedPayload) isEventPayload() {}

// ErrorPayload is a session-level error.
type ErrorPayload struct {
	Message string `json:"message"`
	Code    string `json:"code,omitempty"`
}

func (ErrorPayload) isEventPayload() {}

// StopReason is why a turn went idle. Type is an open string ("turn_ended" |
// "user_action_required" | …) so new provider reasons don't break decoding.
type StopReason struct {
	Type string `json:"type"`
}

// SessionIdlePayload is emitted when a turn ends. EventIds carries the tool_use
// event ids awaiting a client response when Type is user_action_required.
type SessionIdlePayload struct {
	StopReason StopReason `json:"stop_reason"`
	EventIds   []string   `json:"event_ids,omitempty"`
}

func (SessionIdlePayload) isEventPayload() {}

// SessionUpdatedPayload announces a server-side session metadata change. v1
// carries the (often auto-generated) Title; the service projects it onto
// Session.Title. Emitted only when the title actually changed.
type SessionUpdatedPayload struct {
	Title string `json:"title,omitempty"`
}

func (SessionUpdatedPayload) isEventPayload() {}

// ModelRequestStartPayload opens a model-request span.
type ModelRequestStartPayload struct {
	ModelId string `json:"model_id"`
}

func (ModelRequestStartPayload) isEventPayload() {}

// ModelUsage is per-request token accounting.
type ModelUsage struct {
	InputTokens              int64  `json:"input_tokens"`
	OutputTokens             int64  `json:"output_tokens"`
	CacheReadInputTokens     *int64 `json:"cache_read_input_tokens,omitempty"`
	CacheCreationInputTokens *int64 `json:"cache_creation_input_tokens,omitempty"`
}

// ModelRequestEndPayload closes a model-request span with usage.
type ModelRequestEndPayload struct {
	ModelUsage ModelUsage `json:"model_usage"`
}

func (ModelRequestEndPayload) isEventPayload() {}

// EmptyPayload is the body of events that carry no fields (user.interrupt,
// session.status_running, session.status_terminated). Marshals to `{}`.
type EmptyPayload struct{}

func (EmptyPayload) isEventPayload() {}

// RawPayload preserves an unrecognized payload verbatim, so the append-forever
// log forward-decodes any future event type without data loss.
type RawPayload struct {
	Raw json.RawMessage
}

func (RawPayload) isEventPayload() {}

// MarshalJSON emits the preserved bytes (or `{}` when absent) so a RawPayload
// round-trips unchanged.
func (payload RawPayload) MarshalJSON() ([]byte, error) {
	if len(payload.Raw) == 0 {
		return []byte("{}"), nil
	}
	return payload.Raw, nil
}

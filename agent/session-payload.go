package agentmodel

import "encoding/json"

// EventPayload is the type-discriminated body of a SessionEvent. The concrete
// variant is selected by the event's Type (see DecodeEventPayload) — the same
// tagged-union pattern as ToolBinding. Single-content payloads are FLAT (no inner
// `type`); only message/result content slices self-tag (see session-content.go).
// On the way out the concrete value marshals itself; nothing carries the union
// discriminator because the sibling EventType already is it.
type EventPayload interface{ isEventPayload() }

// UserMessagePayload is a user message body (text + file blocks). Context is
// what the launching surface hands the model alongside the message (a call's
// summary and linked records, say): stored with the event like every payload
// field, sent to the model as an extra text block, never rendered as part of
// the person's message. On the first message of a session whose agent has
// memory, the service appends the owner's memory block to it. Not a place for
// secrets: everyone who can read the session reads it back through the events API.
type UserMessagePayload struct {
	Content []ContentBlock `json:"content"`
	Context string         `json:"context,omitempty"`
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

// ReasoningPayload is an agent reasoning ("thinking") block. Managed Agents sends
// neither text nor signature on the whole event, so both are empty for its
// sessions. The embedded runtime stores what the model provider needs back to
// continue a conversation that used reasoning: the text with its signature, or
// RedactedData for reasoning the provider returned encrypted in place of text.
// Signature and RedactedData are opaque and go back unchanged, and only to the
// Provider that wrote them (ModelProvider*): another one cannot read them.
type ReasoningPayload struct {
	Text         string `json:"text"`
	Signature    string `json:"signature,omitempty"`
	RedactedData string `json:"redacted_data,omitempty"`
	Provider     string `json:"provider,omitempty"`
}

func (ReasoningPayload) isEventPayload() {}

// DeltaBoundaryPayload is the body of a .start or .end chunk frame of an
// agent.message or agent.reasoning. PartIndex names the content part of the
// whole event the frames belong to.
type DeltaBoundaryPayload struct {
	PartIndex int `json:"part_index"`
}

func (DeltaBoundaryPayload) isEventPayload() {}

// MessageDeltaPayload is one agent.message.delta: the next piece of Text of the
// part. A client appends the pieces in order.
type MessageDeltaPayload struct {
	PartIndex int    `json:"part_index"`
	Text      string `json:"text"`
}

func (MessageDeltaPayload) isEventPayload() {}

// ReasoningDeltaPayload is one agent.reasoning.delta: the next piece of the
// reasoning Text, or the Signature that closes it.
type ReasoningDeltaPayload struct {
	PartIndex int    `json:"part_index"`
	Text      string `json:"text,omitempty"`
	Signature string `json:"signature,omitempty"`
}

func (ReasoningDeltaPayload) isEventPayload() {}

// ToolUseStartPayload is an agent.tool_use.start: the call is known by Id, Name
// and ToolKind before its input has been written.
type ToolUseStartPayload struct {
	PartIndex int      `json:"part_index"`
	Id        string   `json:"id"`
	Name      string   `json:"name"`
	ToolKind  ToolKind `json:"tool_kind"`
}

func (ToolUseStartPayload) isEventPayload() {}

// ToolUseDeltaPayload is one agent.tool_use.delta: the next piece of the call's
// input as JSON text. The pieces concatenate to the input of the whole event.
type ToolUseDeltaPayload struct {
	PartIndex   int    `json:"part_index"`
	PartialJson string `json:"partial_json"`
}

func (ToolUseDeltaPayload) isEventPayload() {}

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

// ContextCompactedPayload announces that earlier context was summarized. On the
// embedded runtime the summary replaces the conversation up to and including
// the event with seq ThroughSeq whenever a model is sent the conversation. Managed
// Agents compacts on its own and announces it with neither field set.
type ContextCompactedPayload struct {
	Summary    string `json:"summary"`
	ThroughSeq int64  `json:"through_seq,omitempty"`
}

func (ContextCompactedPayload) isEventPayload() {}

// ErrorPayload is a session-level error.
type ErrorPayload struct {
	Message string `json:"message"`
	Code    string `json:"code,omitempty"`
}

func (ErrorPayload) isEventPayload() {}

// StopReason is why a turn went idle. Type is an open string so new provider
// reasons don't break decoding; the ones we emit are the constants below.
type StopReason struct {
	Type string `json:"type"`
}

const (
	// StopReasonTurnEnded means the agent finished its turn and the session is
	// ready for the next message.
	StopReasonTurnEnded = "turn_ended"
	// StopReasonUserActionRequired means the turn is PARKED: something outside
	// the server has to act before it can continue. EventIds names what.
	StopReasonUserActionRequired = "user_action_required"
	// StopReasonBudgetReached means the session used up the spend ceiling it was
	// created with. Its provider takes no further message for it: one sent to it
	// is refused (session_budget_reached).
	StopReasonBudgetReached = "budget_reached"
)

// SessionIdlePayload ends a turn. It is published by the DOMAIN, once per turn —
// it is not an echo of the provider's session status, which flips idle/running
// once per server-tool round inside a single turn.
//
// EventIds is set only with StopReasonUserActionRequired, and lists the tool_use
// events still OUTSTANDING — the ones a client must answer (a client tool).
// Server-executed tools never appear here: they are resolved before this event
// exists, so seeing this event at all means the turn cannot progress without
// help.
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

// DelegatedSessionUpdatedPayload reports where a delegated session stands at the
// end of one of its turns, in the log of the session that delegated to it.
//
// Status is the delegated session's status at that moment. StopReason says why
// the turn went idle and is set only with status idle: turn_ended, or
// user_action_required when the session waits for an answer, in which case
// AwaitingToolUses lists the calls it waits on (answer_delegated_session takes
// their ids). Summary is the last thing the delegated agent wrote, cut to a bound;
// the full conversation stays in the delegated session.
//
// IsStoppedByParent is true when the turn ended because the delegating agent
// interrupted it: that agent asked for the stop, so the update is no news to it.
//
// IsHandedBack is true when the person was talking to the delegated agent in its
// session and that agent handed them back (the hand_back_session client tool).
// Summary is then what the agent gave the delegating agent with that call, when
// it gave one.
type DelegatedSessionUpdatedPayload struct {
	SessionId         string            `json:"session_id"`
	AgentKey          string            `json:"agent_key"`
	Title             string            `json:"title"`
	Status            SessionStatus     `json:"status"`
	StopReason        *StopReason       `json:"stop_reason,omitempty"`
	AwaitingToolUses  []AwaitingToolUse `json:"awaiting_tool_uses,omitempty"`
	IsStoppedByParent bool              `json:"is_stopped_by_parent,omitempty"`
	IsHandedBack      bool              `json:"is_handed_back,omitempty"`
	Summary           string            `json:"summary"`
}

func (DelegatedSessionUpdatedPayload) isEventPayload() {}

// SpeechPayload is one utterance of a voice session, of the person (user.speech)
// or of the voice assistant (agent.speech): the text as the voice provider
// transcribed it, and when it was said, in milliseconds from the start of the
// voice session. The text is what the page reported, not audio this server heard.
type SpeechPayload struct {
	VoiceSessionId string `json:"voice_session_id"`
	Text           string `json:"text"`
	StartMs        int64  `json:"start_ms"`
	EndMs          int64  `json:"end_ms"`
}

func (SpeechPayload) isEventPayload() {}

// VoiceRequestPayload is a user.voice_request body: which voice session handed
// the request over. The request itself is the speech before the event.
// DelegationId is the voice provider's delegation the request is of, as the page
// reported it: the page speaks the turn's answer as that delegation's result, and
// a request is recorded once per delegation. Context is what the page tells the
// model with the request and nobody said (see UserMessagePayload.Context).
type VoiceRequestPayload struct {
	VoiceSessionId string `json:"voice_session_id"`
	DelegationId   string `json:"delegation_id,omitempty"`
	Context        string `json:"context,omitempty"`
}

func (VoiceRequestPayload) isEventPayload() {}

// VoiceStartedPayload is a session.voice_started body.
type VoiceStartedPayload struct {
	VoiceSessionId string `json:"voice_session_id"`
}

func (VoiceStartedPayload) isEventPayload() {}

// VoiceEndedPayload is a session.voice_ended body: why the voice session ended
// and how long the voice provider says it lasted (nil when nobody learned it).
type VoiceEndedPayload struct {
	VoiceSessionId  string         `json:"voice_session_id"`
	EndReason       VoiceEndReason `json:"end_reason"`
	DurationSeconds *int64         `json:"duration_seconds"`
}

func (VoiceEndedPayload) isEventPayload() {}

// AwaitingToolUse is a tool call a session is parked on: what it asks, for
// whoever answers it.
type AwaitingToolUse struct {
	ToolUseId string          `json:"tool_use_id"`
	Name      string          `json:"name"`
	Input     json.RawMessage `json:"input"`
}

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

// SessionRunningPayload opens a turn. It is published by the domain, once per
// turn, like the SessionIdlePayload that ends it.
//
// ActingUserId is the person the turn's tools run as. It is not always the
// author of the message that opened the turn: a message an agent posts to a
// session it delegated to is the agent's, and the turn acts for the person that
// agent acted for. Recorded so that a turn taken up after its replica stopped
// acts for the same person. Empty on a turn started before it was recorded.
type SessionRunningPayload struct {
	ActingUserId string `json:"acting_user_id,omitempty"`
}

func (SessionRunningPayload) isEventPayload() {}

// EmptyPayload is the body of events that carry no fields (user.interrupt,
// session.status_terminated, session.status_ended). Marshals to `{}`.
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

package agentmodel

import (
	"encoding/json"
	"fmt"

	"golang.org/x/exp/slices"
)

// EventType is the discriminant of a SessionEvent — it selects which payload
// shape the event carries (see DecodeEventPayload). The set mirrors the chat-v2
// model: user.*, agent.*, session.*, span.*. v1 carries WHOLE events only — the
// delta/chunk variants (agent.message.delta, …) are intentionally absent.
type EventType string

const (
	// User-originated events (posted via the API or echoed back by the provider).
	EventTypeUserMessage          EventType = "user.message"
	EventTypeUserDefineOutcome    EventType = "user.define_outcome"
	EventTypeUserInterrupt        EventType = "user.interrupt"
	EventTypeUserToolConfirmation EventType = "user.tool_confirmation"
	EventTypeUserToolResult       EventType = "user.tool_result"

	// Agent-originated events.
	EventTypeAgentMessage          EventType = "agent.message"
	EventTypeAgentReasoning        EventType = "agent.reasoning"
	EventTypeAgentToolUse          EventType = "agent.tool_use"
	EventTypeAgentToolResult       EventType = "agent.tool_result"
	EventTypeAgentContextCompacted EventType = "agent.context_compacted"
	// EventTypeAgentArtifacts carries one or more files the agent produced in its
	// sandbox (e.g. wrote to /mnt/session/outputs), surfaced as downloadable
	// attachments on the agent turn. Emitted at turn end, distinct from the
	// text-only agent.message.
	EventTypeAgentArtifacts EventType = "agent.artifacts"

	// Session lifecycle events.
	EventTypeSessionStatusRunning    EventType = "session.status_running"
	EventTypeSessionStatusIdle       EventType = "session.status_idle"
	EventTypeSessionStatusTerminated EventType = "session.status_terminated"
	EventTypeSessionError            EventType = "session.error"
	// EventTypeSessionUpdated carries a server-side session metadata change — v1 the
	// (often auto-generated) title. Projected onto Session.Title; not rendered.
	EventTypeSessionUpdated EventType = "session.updated"

	// Span (observability) events.
	EventTypeSpanModelRequestStart EventType = "span.model_request_start"
	EventTypeSpanModelRequestEnd   EventType = "span.model_request_end"

	// Outcome-evaluation span events (emitted during an outcome-driven turn — the
	// Anthropic grader's iterate/grade loop).
	EventTypeSpanOutcomeEvaluationStart   EventType = "span.outcome_evaluation_start"
	EventTypeSpanOutcomeEvaluationOngoing EventType = "span.outcome_evaluation_ongoing"
	EventTypeSpanOutcomeEvaluationEnd     EventType = "span.outcome_evaluation_end"
)

// EventTypes is the canonical, ordered set of v1 event types.
var EventTypes = []EventType{
	EventTypeUserMessage,
	EventTypeUserDefineOutcome,
	EventTypeUserInterrupt,
	EventTypeUserToolConfirmation,
	EventTypeUserToolResult,
	EventTypeAgentMessage,
	EventTypeAgentReasoning,
	EventTypeAgentToolUse,
	EventTypeAgentToolResult,
	EventTypeAgentContextCompacted,
	EventTypeAgentArtifacts,
	EventTypeSessionStatusRunning,
	EventTypeSessionStatusIdle,
	EventTypeSessionStatusTerminated,
	EventTypeSessionError,
	EventTypeSessionUpdated,
	EventTypeSpanModelRequestStart,
	EventTypeSpanModelRequestEnd,
	EventTypeSpanOutcomeEvaluationStart,
	EventTypeSpanOutcomeEvaluationOngoing,
	EventTypeSpanOutcomeEvaluationEnd,
}

func (EventType) Enum() []interface{} {
	enums := make([]interface{}, 0, len(EventTypes))
	for _, element := range EventTypes {
		enums = append(enums, element)
	}
	return enums
}

func (eventType *EventType) UnmarshalJSON(byteArray []byte) error {
	if string(byteArray) == "null" {
		*eventType = ""
		return nil
	}

	type _EventType EventType
	value := (*_EventType)(eventType)
	if err := json.Unmarshal(byteArray, value); err != nil {
		return err
	}

	if slices.Contains(EventTypes, *eventType) {
		return nil
	}

	return fmt.Errorf("invalid session event type: %s", *eventType)
}

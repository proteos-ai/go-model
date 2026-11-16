package agentmodel

import (
	"encoding/json"
	"fmt"

	"golang.org/x/exp/slices"
)

// EventType is the discriminant of a SessionEvent — it selects which payload
// shape the event carries (see DecodeEventPayload). The set mirrors the chat-v2
// model: user.*, agent.*, session.*, span.*. Whole events (EventTypes) make up the
// stored log; the chunk frames of an event still being written (DeltaEventTypes)
// exist on the live delta channel only.
type EventType string

const (
	// User-originated events (posted via the API or echoed back by the provider).
	EventTypeUserMessage          EventType = "user.message"
	EventTypeUserDefineOutcome    EventType = "user.define_outcome"
	EventTypeUserInterrupt        EventType = "user.interrupt"
	EventTypeUserToolConfirmation EventType = "user.tool_confirmation"
	EventTypeUserToolResult       EventType = "user.tool_result"
	// EventTypeUserSpeech is what the person said in a voice session, one event
	// per utterance, as the voice provider transcribed it. It starts no turn:
	// the model reads it with what opens the next one (see EventTypeUserVoiceRequest).
	EventTypeUserSpeech EventType = "user.speech"
	// EventTypeUserVoiceRequest is the voice assistant handing what the person
	// asked for to the session's agent. It starts a turn like a user.message and
	// has no content of its own: the request is the speech before it.
	EventTypeUserVoiceRequest EventType = "user.voice_request"

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
	// EventTypeAgentSpeech is what the voice assistant said aloud in a voice
	// session, one event per utterance. The voice assistant speaks for the
	// session's agent in its own words, so this is not the agent's message.
	EventTypeAgentSpeech EventType = "agent.speech"

	// Session lifecycle events.
	EventTypeSessionStatusRunning    EventType = "session.status_running"
	EventTypeSessionStatusIdle       EventType = "session.status_idle"
	EventTypeSessionStatusTerminated EventType = "session.status_terminated"
	EventTypeSessionError            EventType = "session.error"
	// EventTypeSessionStatusEnded records a user ending the session. Its created_by
	// and created_at are projected onto Session.EndedBy / EndedAt.
	EventTypeSessionStatusEnded EventType = "session.status_ended"
	// EventTypeSessionUpdated carries a server-side session metadata change — v1 the
	// (often auto-generated) title. Projected onto Session.Title; not rendered.
	EventTypeSessionUpdated EventType = "session.updated"
	// EventTypeSessionVoiceStarted and EventTypeSessionVoiceEnded bracket a voice
	// session: the time the person talked to the session's agent through a voice
	// assistant. Neither changes the session's status.
	EventTypeSessionVoiceStarted EventType = "session.voice_started"
	EventTypeSessionVoiceEnded   EventType = "session.voice_ended"

	// EventTypeDelegatedSessionUpdated is written to a session's log when a turn
	// of one of its delegated sessions ends: how that session stands and what its
	// agent last said. It is the only event of a session that reports on another.
	EventTypeDelegatedSessionUpdated EventType = "delegated_session.updated"

	// Span (observability) events.
	EventTypeSpanModelRequestStart EventType = "span.model_request_start"
	EventTypeSpanModelRequestEnd   EventType = "span.model_request_end"

	// Outcome-evaluation span events (emitted during an outcome-driven turn — the
	// Anthropic grader's iterate/grade loop).
	EventTypeSpanOutcomeEvaluationStart   EventType = "span.outcome_evaluation_start"
	EventTypeSpanOutcomeEvaluationOngoing EventType = "span.outcome_evaluation_ongoing"
	EventTypeSpanOutcomeEvaluationEnd     EventType = "span.outcome_evaluation_end"
)

// Chunk frames of the three long agent events, sent on the live delta channel while
// the model writes the event. A frame's type is the whole event's type plus .start,
// .delta or .end, and every frame of one event carries that event's id, so a client
// builds a draft from them and replaces it when the whole event arrives. Frames are
// never stored and never posted by a client.
const (
	EventTypeAgentMessageStart   EventType = "agent.message.start"
	EventTypeAgentMessageDelta   EventType = "agent.message.delta"
	EventTypeAgentMessageEnd     EventType = "agent.message.end"
	EventTypeAgentReasoningStart EventType = "agent.reasoning.start"
	EventTypeAgentReasoningDelta EventType = "agent.reasoning.delta"
	EventTypeAgentReasoningEnd   EventType = "agent.reasoning.end"
	EventTypeAgentToolUseStart   EventType = "agent.tool_use.start"
	EventTypeAgentToolUseDelta   EventType = "agent.tool_use.delta"
	EventTypeAgentToolUseEnd     EventType = "agent.tool_use.end"
)

// DeltaEventTypes is the set of chunk frame types of the live delta channel.
var DeltaEventTypes = []EventType{
	EventTypeAgentMessageStart,
	EventTypeAgentMessageDelta,
	EventTypeAgentMessageEnd,
	EventTypeAgentReasoningStart,
	EventTypeAgentReasoningDelta,
	EventTypeAgentReasoningEnd,
	EventTypeAgentToolUseStart,
	EventTypeAgentToolUseDelta,
	EventTypeAgentToolUseEnd,
}

// IsDeltaEventType reports whether the type is a chunk frame of the live delta
// channel rather than a whole event of the stored log.
func IsDeltaEventType(eventType EventType) bool {
	return slices.Contains(DeltaEventTypes, eventType)
}

// EventTypes is the canonical, ordered set of whole event types: the stored log.
var EventTypes = []EventType{
	EventTypeUserMessage,
	EventTypeUserDefineOutcome,
	EventTypeUserInterrupt,
	EventTypeUserToolConfirmation,
	EventTypeUserToolResult,
	EventTypeUserSpeech,
	EventTypeUserVoiceRequest,
	EventTypeAgentMessage,
	EventTypeAgentReasoning,
	EventTypeAgentToolUse,
	EventTypeAgentToolResult,
	EventTypeAgentContextCompacted,
	EventTypeAgentArtifacts,
	EventTypeAgentSpeech,
	EventTypeSessionStatusRunning,
	EventTypeSessionStatusIdle,
	EventTypeSessionStatusTerminated,
	EventTypeSessionError,
	EventTypeSessionStatusEnded,
	EventTypeSessionUpdated,
	EventTypeSessionVoiceStarted,
	EventTypeSessionVoiceEnded,
	EventTypeDelegatedSessionUpdated,
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

	if slices.Contains(EventTypes, *eventType) || IsDeltaEventType(*eventType) {
		return nil
	}

	return fmt.Errorf("invalid session event type: %s", *eventType)
}

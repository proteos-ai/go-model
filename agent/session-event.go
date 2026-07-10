package agentmodel

import (
	"encoding/json"
	"time"

	"go.proteos.ai/model/common"
)

// SessionEvent is one entry of the append-only session log — the envelope shared
// byte-for-byte across both channels (Redis live + Postgres durable). The same
// envelope is built once and written to both; Seq is the per-session monotonic
// cursor (the Redis entry id, <seq>-0) on both. Id is the dedupe key: our uuid
// for user-originated events, the provider's event id for agent events.
// ProcessedAt is meaningful only for user-sent events — null = queued,
// timestamp = consumed by the agent (stamped from the provider's echo).
type SessionEvent struct {
	Id          string         `json:"id"`
	SessionId   string         `json:"session_id"`
	OrgId       string         `json:"org_id"`
	TurnId      string         `json:"turn_id"`
	Seq         int64          `json:"seq"`
	Type        EventType      `json:"type"`
	Payload     EventPayload   `json:"payload"`
	CreatedBy   common.UserRef `json:"created_by"`
	ProcessedAt *time.Time     `json:"processed_at"`
	CreatedAt   time.Time      `json:"created_at"`
}

// DecodeEventPayload decodes a raw payload according to the event type (mirrors
// DecodeToolBinding). An unrecognized type yields a RawPayload so the
// append-forever log forward-decodes future event types without loss.
func DecodeEventPayload(eventType EventType, raw json.RawMessage) (EventPayload, error) {
	empty := len(raw) == 0 || string(raw) == "null"
	switch eventType {
	case EventTypeUserMessage:
		return decodePayload(raw, empty, &UserMessagePayload{})
	case EventTypeUserDefineOutcome:
		return decodePayload(raw, empty, &DefineOutcomePayload{})
	case EventTypeSpanOutcomeEvaluationStart, EventTypeSpanOutcomeEvaluationOngoing, EventTypeSpanOutcomeEvaluationEnd:
		return decodePayload(raw, empty, &OutcomeEvaluationPayload{})
	case EventTypeAgentMessage:
		return decodePayload(raw, empty, &AgentMessagePayload{})
	case EventTypeAgentReasoning:
		return decodePayload(raw, empty, &ReasoningPayload{})
	case EventTypeAgentToolUse:
		return decodePayload(raw, empty, &ToolUsePayload{})
	case EventTypeAgentToolResult, EventTypeUserToolResult:
		return decodePayload(raw, empty, &ToolResultPayload{})
	case EventTypeUserToolConfirmation:
		return decodePayload(raw, empty, &ToolConfirmationPayload{})
	case EventTypeAgentContextCompacted:
		return decodePayload(raw, empty, &ContextCompactedPayload{})
	case EventTypeAgentArtifacts:
		return decodePayload(raw, empty, &AgentArtifactsPayload{})
	case EventTypeSessionError:
		return decodePayload(raw, empty, &ErrorPayload{})
	case EventTypeSessionStatusIdle:
		return decodePayload(raw, empty, &SessionIdlePayload{})
	case EventTypeSessionUpdated:
		return decodePayload(raw, empty, &SessionUpdatedPayload{})
	case EventTypeSpanModelRequestStart:
		return decodePayload(raw, empty, &ModelRequestStartPayload{})
	case EventTypeSpanModelRequestEnd:
		return decodePayload(raw, empty, &ModelRequestEndPayload{})
	case EventTypeUserInterrupt, EventTypeSessionStatusRunning, EventTypeSessionStatusTerminated:
		return EmptyPayload{}, nil
	default:
		return RawPayload{Raw: raw}, nil
	}
}

// decodePayload unmarshals raw into dst (skipping when empty so a flat/empty body
// decodes to the zero payload) and returns it as the EventPayload interface.
func decodePayload[T EventPayload](raw json.RawMessage, empty bool, dst *T) (EventPayload, error) {
	if !empty {
		if err := json.Unmarshal(raw, dst); err != nil {
			return nil, err
		}
	}
	return *dst, nil
}

// UnmarshalJSON decodes the discriminated payload by the sibling Type. (Marshalling
// uses the default: the concrete payload value marshals itself.)
func (event *SessionEvent) UnmarshalJSON(data []byte) error {
	type alias SessionEvent
	aux := struct {
		*alias
		Payload json.RawMessage `json:"payload"`
	}{alias: (*alias)(event)}
	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}
	payload, err := DecodeEventPayload(event.Type, aux.Payload)
	if err != nil {
		return err
	}
	event.Payload = payload
	return nil
}

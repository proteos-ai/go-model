package agentmodel

import (
	"encoding/json"
	"testing"
)

func TestDecodeEventPayload_KnownTypes(t *testing.T) {
	userMsg, err := DecodeEventPayload(EventTypeUserMessage, json.RawMessage(`{"content":[{"type":"text","text":"hi"}]}`))
	if err != nil {
		t.Fatalf("user.message: %v", err)
	}
	if got := userMsg.(UserMessagePayload); len(got.Content) != 1 || got.Content[0].Text != "hi" {
		t.Errorf("user.message payload = %+v", got)
	}

	toolUse, err := DecodeEventPayload(EventTypeAgentToolUse, json.RawMessage(`{"id":"toolu_1","name":"get_weather","input":{"location":"SF"},"tool_kind":"action"}`))
	if err != nil {
		t.Fatalf("agent.tool_use: %v", err)
	}
	if got := toolUse.(ToolUsePayload); got.Id != "toolu_1" || got.ToolKind != ToolKindAction {
		t.Errorf("agent.tool_use payload = %+v", got)
	}

	idle, err := DecodeEventPayload(EventTypeSessionStatusIdle, json.RawMessage(`{"stop_reason":{"type":"turn_ended"}}`))
	if err != nil {
		t.Fatalf("session.status_idle: %v", err)
	}
	if got := idle.(SessionIdlePayload); got.StopReason.Type != "turn_ended" {
		t.Errorf("session.status_idle payload = %+v", got)
	}
}

func TestDecodeEventPayload_EmptyAndStatus(t *testing.T) {
	// A turn start written before it named its acting user decodes with none.
	for raw, want := range map[string]string{`{"acting_user_id":"user-1"}`: "user-1", `{}`: "", "": ""} {
		running, err := DecodeEventPayload(EventTypeSessionStatusRunning, json.RawMessage(raw))
		if err != nil {
			t.Fatalf("session.status_running %q: %v", raw, err)
		}
		if got, ok := running.(SessionRunningPayload); !ok || got.ActingUserId != want {
			t.Errorf("session.status_running %q payload = %#v, want acting user %q", raw, running, want)
		}
	}
	for _, eventType := range []EventType{EventTypeUserInterrupt, EventTypeSessionStatusTerminated, EventTypeSessionStatusEnded} {
		payload, err := DecodeEventPayload(eventType, nil)
		if err != nil {
			t.Fatalf("%s: %v", eventType, err)
		}
		if _, ok := payload.(EmptyPayload); !ok {
			t.Errorf("%s payload = %T, want EmptyPayload", eventType, payload)
		}
	}
}

func TestDecodeEventPayload_ChunkFrames(t *testing.T) {
	cases := []struct {
		eventType EventType
		raw       string
		want      EventPayload
	}{
		{EventTypeAgentMessageStart, `{"part_index":1}`, DeltaBoundaryPayload{PartIndex: 1}},
		{EventTypeAgentMessageDelta, `{"part_index":1,"text":"Hel"}`, MessageDeltaPayload{PartIndex: 1, Text: "Hel"}},
		{EventTypeAgentMessageEnd, `{"part_index":1}`, DeltaBoundaryPayload{PartIndex: 1}},
		{EventTypeAgentReasoningStart, `{"part_index":0}`, DeltaBoundaryPayload{}},
		{EventTypeAgentReasoningDelta, `{"part_index":0,"text":"hm"}`, ReasoningDeltaPayload{Text: "hm"}},
		{EventTypeAgentReasoningDelta, `{"part_index":0,"signature":"sig"}`, ReasoningDeltaPayload{Signature: "sig"}},
		{EventTypeAgentReasoningEnd, `{"part_index":0}`, DeltaBoundaryPayload{}},
		{EventTypeAgentToolUseStart, `{"part_index":2,"id":"toolu_1","name":"get_record","tool_kind":"platform"}`,
			ToolUseStartPayload{PartIndex: 2, Id: "toolu_1", Name: "get_record", ToolKind: ToolKindPlatform}},
		{EventTypeAgentToolUseDelta, `{"part_index":2,"partial_json":"{\"id\":"}`, ToolUseDeltaPayload{PartIndex: 2, PartialJson: `{"id":`}},
		{EventTypeAgentToolUseEnd, `{"part_index":2}`, DeltaBoundaryPayload{PartIndex: 2}},
	}
	covered := map[EventType]bool{}
	for _, testCase := range cases {
		covered[testCase.eventType] = true
		payload, err := DecodeEventPayload(testCase.eventType, json.RawMessage(testCase.raw))
		if err != nil {
			t.Fatalf("%s: %v", testCase.eventType, err)
		}
		if payload != testCase.want {
			t.Errorf("%s payload = %#v, want %#v", testCase.eventType, payload, testCase.want)
		}
		if !IsDeltaEventType(testCase.eventType) {
			t.Errorf("%s is not listed as a chunk frame type", testCase.eventType)
		}
	}
	for _, eventType := range DeltaEventTypes {
		if !covered[eventType] {
			t.Errorf("%s has no decode case", eventType)
		}
		if IsDeltaEventType(eventType) == false || containsEventType(EventTypes, eventType) {
			t.Errorf("%s must be a chunk frame type and not a whole event type", eventType)
		}
	}
}

func containsEventType(eventTypes []EventType, eventType EventType) bool {
	for _, candidate := range eventTypes {
		if candidate == eventType {
			return true
		}
	}
	return false
}

// A chunk frame round-trips through the envelope the delta channel stores.
func TestSessionEvent_ChunkFrameRoundTrip(t *testing.T) {
	frame := SessionEvent{
		Id: "evt_msg", SessionId: "sess_1", OrgId: "org_1", TurnId: "turn_1",
		Type: EventTypeAgentMessageDelta, Payload: MessageDeltaPayload{PartIndex: 0, Text: "Hel"},
	}
	raw, err := json.Marshal(frame)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	var decoded SessionEvent
	if err := json.Unmarshal(raw, &decoded); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if decoded.Type != EventTypeAgentMessageDelta || decoded.Payload != frame.Payload || decoded.Id != "evt_msg" {
		t.Errorf("round trip = %+v", decoded)
	}
}

func TestParseDeltaPosition(t *testing.T) {
	valid := map[string]DeltaPosition{
		"42":   {Seq: 42},
		"42-0": {Seq: 42},
		"42-3": {Seq: 42, Index: 3},
		"0":    {},
	}
	for raw, want := range valid {
		got, err := ParseDeltaPosition(raw)
		if err != nil || got != want {
			t.Errorf("ParseDeltaPosition(%q) = %+v, %v; want %+v", raw, got, err, want)
		}
	}
	if got := (DeltaPosition{Seq: 42, Index: 3}).String(); got != "42-3" {
		t.Errorf("String() = %q, want 42-3", got)
	}
	for _, raw := range []string{"", "abc", "-1", "42-", "42-x", "42--1", "4.2"} {
		if _, err := ParseDeltaPosition(raw); err == nil {
			t.Errorf("ParseDeltaPosition(%q) = nil error, want an error", raw)
		}
	}
}

func TestDecodeEventPayload_UnknownYieldsRaw(t *testing.T) {
	raw := json.RawMessage(`{"some":"future_field"}`)
	payload, err := DecodeEventPayload(EventType("agent.future_event"), raw)
	if err != nil {
		t.Fatalf("unknown: %v", err)
	}
	if got, ok := payload.(RawPayload); !ok || string(got.Raw) != string(raw) {
		t.Errorf("unknown payload = %#v, want RawPayload preserving bytes", payload)
	}
}

func TestSessionEvent_UnmarshalDiscriminatedPayload(t *testing.T) {
	raw := `{
		"id":"evt_1","session_id":"sess_1","org_id":"org_1","turn_id":"turn_1","seq":5,
		"type":"agent.tool_result",
		"payload":{"tool_use_id":"toolu_1","result":[{"type":"json","json":{"temp_f":61}}],"outcome":"success"},
		"created_at":"2026-06-02T09:00:00Z"
	}`
	var event SessionEvent
	if err := json.Unmarshal([]byte(raw), &event); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if event.Type != EventTypeAgentToolResult || event.Seq != 5 {
		t.Errorf("envelope = %+v", event)
	}
	result, ok := event.Payload.(ToolResultPayload)
	if !ok {
		t.Fatalf("payload type = %T, want ToolResultPayload", event.Payload)
	}
	if result.ToolUseId != "toolu_1" || result.Outcome != "success" || len(result.Result) != 1 {
		t.Errorf("payload = %+v", result)
	}
}

func TestSessionEvent_MarshalRoundTrip(t *testing.T) {
	original := SessionEvent{
		Id: "evt_1", SessionId: "sess_1", OrgId: "org_1", TurnId: "turn_1", Seq: 7,
		Type:    EventTypeAgentMessage,
		Payload: AgentMessagePayload{Content: []TextBlock{{Type: "text", Text: "It's foggy."}}},
	}

	encoded, err := json.Marshal(original)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}

	var decoded SessionEvent
	if err := json.Unmarshal(encoded, &decoded); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	message, ok := decoded.Payload.(AgentMessagePayload)
	if !ok {
		t.Fatalf("payload type = %T, want AgentMessagePayload", decoded.Payload)
	}
	if len(message.Content) != 1 || message.Content[0].Text != "It's foggy." {
		t.Errorf("round-tripped payload = %+v", message)
	}
}

func TestContentBlock_FileMarshalsInline(t *testing.T) {
	name := "q2.pdf"
	block := ContentBlock{Type: "file", FileBlock: FileBlock{FileId: "file_1", FileName: &name}}
	encoded, err := json.Marshal(block)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	// file_id/file_name promoted inline; text omitted.
	want := `{"type":"file","file_id":"file_1","file_name":"q2.pdf"}`
	if string(encoded) != want {
		t.Errorf("file content block = %s, want %s", encoded, want)
	}
}

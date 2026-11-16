package agentmodel

import (
	"encoding/json"
	"testing"
)

// The session tool manifest is persisted as a jsonb column; bun marshals/scans it with
// encoding/json. This verifies the kind-discriminated Binding survives a round trip
// (action/mcp decode to their concrete bindings, client stays nil) the same way it does
// through the database.
func TestSessionTool_JSONRoundTrip(t *testing.T) {
	manifest := []SessionTool{
		{Key: "recompute", Kind: ToolKindAction, Binding: ActionBinding{ActionKey: "recompute-totals"}},
		{Key: "search", Kind: ToolKindMcp, Binding: McpBinding{ServerKey: "kb", ToolName: "search"}},
		{Key: "navigate", Kind: ToolKindClient, Description: "Open a page",
			InputSchema: json.RawMessage(`{"type":"object","properties":{"to":{"type":"string"}}}`)},
	}

	data, err := json.Marshal(manifest)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}

	var decoded []SessionTool
	if err := json.Unmarshal(data, &decoded); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if len(decoded) != 3 {
		t.Fatalf("len = %d, want 3", len(decoded))
	}

	action, ok := decoded[0].Binding.(ActionBinding)
	if !ok || action.ActionKey != "recompute-totals" {
		t.Errorf("action binding = %#v", decoded[0].Binding)
	}
	mcp, ok := decoded[1].Binding.(McpBinding)
	if !ok || mcp.ServerKey != "kb" || mcp.ToolName != "search" {
		t.Errorf("mcp binding = %#v", decoded[1].Binding)
	}
	if decoded[2].Binding != nil {
		t.Errorf("client binding = %#v, want nil", decoded[2].Binding)
	}
	if decoded[2].Kind != ToolKindClient {
		t.Errorf("client kind = %q, want client", decoded[2].Kind)
	}
	// How a tool is described to the model survives with it; a tool stored
	// without one (every manifest written before these fields) decodes empty.
	if decoded[2].Description != "Open a page" || string(decoded[2].InputSchema) != `{"type":"object","properties":{"to":{"type":"string"}}}` {
		t.Errorf("client definition = %q %s", decoded[2].Description, decoded[2].InputSchema)
	}
	if decoded[0].Description != "" || decoded[0].InputSchema != nil {
		t.Errorf("a tool stored without a definition decoded with one: %q %s", decoded[0].Description, decoded[0].InputSchema)
	}
}

// An empty manifest (a session opened with no client tools) round-trips as an empty
// slice, matching the jsonb column default '[]'.
func TestSessionTool_EmptyManifest(t *testing.T) {
	data, err := json.Marshal([]SessionTool{})
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	if string(data) != "[]" {
		t.Errorf("marshal empty = %s, want []", data)
	}
	var decoded []SessionTool
	if err := json.Unmarshal([]byte("[]"), &decoded); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if len(decoded) != 0 {
		t.Errorf("len = %d, want 0", len(decoded))
	}
}

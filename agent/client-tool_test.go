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
		{Key: "navigate", Kind: ToolKindClient},
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

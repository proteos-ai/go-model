package agentmodel

import (
	"encoding/json"
	"testing"
)

func TestDecodeToolBinding(t *testing.T) {
	action, err := DecodeToolBinding(ToolKindAction, json.RawMessage(`{"action_key":"mark-as-signed"}`))
	if err != nil {
		t.Fatalf("action: %v", err)
	}
	if action.(ActionBinding).ActionKey != "mark-as-signed" {
		t.Errorf("action key = %q", action.(ActionBinding).ActionKey)
	}

	mcp, err := DecodeToolBinding(ToolKindMcp, json.RawMessage(`{"server_key":"github","tool_name":"create_issue"}`))
	if err != nil {
		t.Fatalf("mcp: %v", err)
	}
	if mcp.(McpBinding).ServerKey != "github" || mcp.(McpBinding).ToolName != "create_issue" {
		t.Errorf("mcp binding = %+v", mcp)
	}

	client, err := DecodeToolBinding(ToolKindClient, nil)
	if err != nil {
		t.Fatalf("client: %v", err)
	}
	if client != nil {
		t.Errorf("client binding = %v, want nil", client)
	}

	if _, err := DecodeToolBinding(ToolKindAction, nil); err == nil {
		t.Error("action with no binding: want error")
	}
}

func TestToolUnmarshalDiscriminatedBinding(t *testing.T) {
	raw := `{
		"org_id":"o1","key":"get-weather","name":"Get Weather","description":"d",
		"kind":"mcp","input_schema":[],
		"binding":{"server_key":"weather","tool_name":"forecast"},"version":1
	}`
	var tool Tool
	if err := json.Unmarshal([]byte(raw), &tool); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	binding, ok := tool.Binding.(McpBinding)
	if !ok {
		t.Fatalf("binding type = %T, want McpBinding", tool.Binding)
	}
	if binding.ServerKey != "weather" || binding.ToolName != "forecast" {
		t.Errorf("binding = %+v", binding)
	}

	// client tool: binding absent → nil
	var clientTool Tool
	if err := json.Unmarshal([]byte(`{"key":"k","kind":"client"}`), &clientTool); err != nil {
		t.Fatalf("client unmarshal: %v", err)
	}
	if clientTool.Binding != nil {
		t.Errorf("client binding = %v, want nil", clientTool.Binding)
	}
}

func TestToolKindUnmarshalRejectsUnknown(t *testing.T) {
	var kind ToolKind
	if err := kind.UnmarshalJSON([]byte(`"backend"`)); err == nil {
		t.Error(`ToolKind "backend" should be rejected (it's action in the registry)`)
	}
	if err := kind.UnmarshalJSON([]byte(`"action"`)); err != nil {
		t.Errorf(`ToolKind "action": %v`, err)
	}
	if err := kind.UnmarshalJSON([]byte(`null`)); err != nil {
		t.Errorf("ToolKind null should be tolerated: %v", err)
	}
}

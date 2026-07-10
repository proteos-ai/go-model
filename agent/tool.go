package agentmodel

import (
	"encoding/json"
	"fmt"
	"time"

	"go.proteos.ai/model/common"
)

// Tool is a thin registry entry over one of three binding sources (action / mcp /
// client). Key is the wire name the model calls (tool_use.name) and what
// Agent.tools lists. Binding is a tagged union decoded by Kind; client tools carry
// no binding — the host owns the implementation and the key is the identifier.
//
// InputSchema / OutputSchema are NOT persisted — they are read-only JSON Schema
// (draft-07) computed on read from the binding's source of truth (the action's
// params/returns, or the MCP server's tools/list). They are empty on writes and
// on plain list reads; the read path populates them. Keyed by (org_id, key).
type Tool struct {
	OrgId        string          `json:"org_id"`
	Key          string          `json:"key" sortable:""`
	Name         string          `json:"name" sortable:""`
	ModuleSlug   string          `json:"module_slug" sortable:""`
	Description  string          `json:"description"`
	Kind         ToolKind        `json:"kind"`
	InputSchema  json.RawMessage `json:"input_schema,omitempty"`
	OutputSchema json.RawMessage `json:"output_schema,omitempty"`
	Binding      ToolBinding     `json:"binding,omitempty"`
	Version      int             `json:"version"`
	CreatedAt    time.Time       `json:"created_at" sortable:""`
	CreatedBy    common.UserRef  `json:"created_by"`
	UpdatedAt    time.Time       `json:"updated_at" sortable:""`
	UpdatedBy    common.UserRef  `json:"updated_by"`
}

// ToolBinding is the kind-discriminated payload of a Tool. Concrete variants:
// ActionBinding (kind=action), McpBinding (kind=mcp). kind=client has no binding.
type ToolBinding interface{ isToolBinding() }

// ActionBinding binds to a function-service Action by its key.
type ActionBinding struct {
	ActionKey string `json:"action_key"`
}

func (ActionBinding) isToolBinding() {}

// McpBinding binds to one tool on a registered McpServer.
type McpBinding struct {
	ServerKey string `json:"server_key"`
	ToolName  string `json:"tool_name"`
}

func (McpBinding) isToolBinding() {}

// DecodeToolBinding decodes a raw binding payload according to kind (mirrors the
// v3 model's DecodePayload). Returns (nil, nil) for client — it has no binding.
func DecodeToolBinding(kind ToolKind, raw json.RawMessage) (ToolBinding, error) {
	empty := len(raw) == 0 || string(raw) == "null"
	if kind == ToolKindClient {
		return nil, nil
	}
	if empty {
		return nil, fmt.Errorf("missing binding for tool kind %q", kind)
	}
	switch kind {
	case ToolKindAction:
		var binding ActionBinding
		if err := json.Unmarshal(raw, &binding); err != nil {
			return nil, err
		}
		return binding, nil
	case ToolKindMcp:
		var binding McpBinding
		if err := json.Unmarshal(raw, &binding); err != nil {
			return nil, err
		}
		return binding, nil
	default:
		return nil, fmt.Errorf("unknown tool kind %q", kind)
	}
}

// UnmarshalJSON decodes the discriminated binding by kind. (Marshalling uses the
// default: the concrete binding value marshals itself; a nil binding is omitted.)
func (tool *Tool) UnmarshalJSON(data []byte) error {
	type alias Tool
	aux := struct {
		*alias
		Binding json.RawMessage `json:"binding"`
	}{alias: (*alias)(tool)}
	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}
	binding, err := DecodeToolBinding(tool.Kind, aux.Binding)
	if err != nil {
		return err
	}
	tool.Binding = binding
	return nil
}

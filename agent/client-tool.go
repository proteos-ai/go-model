package agentmodel

import "encoding/json"

// ClientToolSpec is one frontend-declared client tool a surface attaches to a session
// at create time: a custom tool the agent calls but the browser renders/executes,
// returning the result. Only these three fields cross the wire — the render slot and
// the React component are frontend-only. The host owns the implementation; Name is the
// wire identifier the model calls (tool_use.name).
type ClientToolSpec struct {
	Name        string          `json:"name"`
	Description string          `json:"description"`
	InputSchema json.RawMessage `json:"input_schema,omitempty"`
}

// SessionTool is one entry of a session's resolved tool manifest: the wire name the
// model calls (Key), its kind (action | mcp | client), and the binding for action/mcp
// (nil for client). The manifest is snapshotted at session create from the agent's
// resolved tools plus the surface's client tools, and persisted as the session's
// `tools` jsonb. It is the single source of truth the runtime uses to route an inbound
// custom_tool_use (action→execute, mcp→proxy, client→wait for the frontend) and to tag
// the ledger event's tool_kind. Binding is the same kind-discriminated union as Tool.
//
// Description and InputSchema are how the tool is described to the model, taken at
// the same moment. A runtime that sends the tools with every model request (the
// embedded runtime) reads them here; Managed Agents holds its own copy.
//
// IsReadOnly is Tool.IsReadOnly at the same moment. A manifest written before the
// field existed has it false, so its tools count as writing.
type SessionTool struct {
	Key         string          `json:"key"`
	Kind        ToolKind        `json:"kind"`
	Binding     ToolBinding     `json:"binding,omitempty"`
	Description string          `json:"description,omitempty"`
	InputSchema json.RawMessage `json:"input_schema,omitempty"`
	IsReadOnly  bool            `json:"is_read_only,omitempty"`
}

// UnmarshalJSON decodes the discriminated binding by kind (mirrors Tool.UnmarshalJSON),
// so a manifest round-trips through the session's jsonb column. Marshalling uses the
// default: the concrete binding marshals itself; a nil (client) binding is omitted.
func (sessionTool *SessionTool) UnmarshalJSON(data []byte) error {
	type alias SessionTool
	aux := struct {
		*alias
		Binding json.RawMessage `json:"binding"`
	}{alias: (*alias)(sessionTool)}
	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}
	binding, err := DecodeToolBinding(sessionTool.Kind, aux.Binding)
	if err != nil {
		return err
	}
	sessionTool.Binding = binding
	return nil
}

package agentmodel

import (
	"encoding/json"
	"fmt"

	"golang.org/x/exp/slices"
)

// ToolKind is the binding SOURCE of a Tool (registry-level): where the callable
// comes from. It is distinct from the runtime tool_kind in the agent message
// protocol (execution location); the worker maps action→backend, mcp→mcp,
// client→client.
type ToolKind string

const (
	// ToolKindAction binds to a function-service Action.
	ToolKindAction ToolKind = "action"
	// ToolKindMcp binds to one tool on a registered McpServer.
	ToolKindMcp ToolKind = "mcp"
	// ToolKindClient is a host-provided builtin — it carries no binding.
	ToolKindClient ToolKind = "client"
)

var ToolKinds = []ToolKind{ToolKindAction, ToolKindMcp, ToolKindClient}

func (ToolKind) Enum() []interface{} {
	enums := []interface{}{}
	for _, element := range ToolKinds {
		enums = append(enums, element)
	}
	return enums
}

func (toolKind *ToolKind) UnmarshalJSON(byteArray []byte) error {
	if string(byteArray) == "null" {
		*toolKind = ""
		return nil
	}

	type _ToolKind ToolKind
	value := (*_ToolKind)(toolKind)
	if err := json.Unmarshal(byteArray, value); err != nil {
		return err
	}

	if slices.Contains(ToolKinds, *toolKind) {
		return nil
	}

	return fmt.Errorf("invalid tool kind: %s", *toolKind)
}

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
	// ToolKindPlatform binds to one tool of the platform MCP server (mcp-service),
	// executed server-side by agent-service as the acting user.
	ToolKindPlatform ToolKind = "platform"
	// ToolKindQuery stores a SQL query (data-service dialect) with declared, typed
	// params; agent-service substitutes escaped literals and executes it via
	// data-service as the acting user.
	ToolKindQuery ToolKind = "query"
	// ToolKindDelegation is one of the tools agent-service itself gives a session
	// to hand work to a sub-agent in a session of its own and to manage those
	// delegated sessions. It carries no binding and is not a registry entry: a
	// session's manifest takes these tools at create, when its agent has sub-agents
	// and its runtime runs them.
	ToolKindDelegation ToolKind = "delegation"
	// ToolKindSystem is a tool agent-service gives every agent to tell the person
	// what it is doing, such as `new_phase`. It carries no binding, is not a
	// registry entry, and agent-service answers it itself without doing anything.
	ToolKindSystem ToolKind = "system"
)

var ToolKinds = []ToolKind{ToolKindAction, ToolKindMcp, ToolKindClient, ToolKindPlatform, ToolKindQuery, ToolKindDelegation, ToolKindSystem}

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

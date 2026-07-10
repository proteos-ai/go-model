package connectormodel

import metamodel "go.proteos.ai/model/meta"

// MethodDeclaration describes one callable operation a connector exposes
// (get_event, update_event, …). Params and Returns use the platform's
// attribute language — the same schema shape actions use — so codegen,
// validation and pickers (the workflow connector node) work identically for
// pre-built and custom connectors.
type MethodDeclaration struct {
	// Key is the invocation name: snake_case verb_noun (list_events).
	Key         string `json:"key"`
	Title       string `json:"title"`
	Description string `json:"description,omitempty"`
	// Access gates the invoke permission: read methods require connections
	// read, write methods require connections write. Declarations omitting it
	// are treated as write (conservative default).
	Access MethodAccess `json:"access"`
	// ActionSlug points at the function-service action (scope
	// connector_method) that implements this method — set ONLY for custom
	// (wasm) connectors. Pre-built connectors execute in-process and leave it
	// empty. This is the connector↔action link: the method references the
	// action, never the reverse.
	ActionSlug string                `json:"action_slug,omitempty"`
	Params     []metamodel.Attribute `json:"params"`
	Returns    []metamodel.Attribute `json:"returns"`
}

// MethodAccess classifies a method's effect on the remote system.
type MethodAccess string

const (
	MethodAccessRead  MethodAccess = "read"
	MethodAccessWrite MethodAccess = "write"
)

// EffectiveAccess resolves the conservative default: anything not explicitly
// read is write.
func (declaration MethodDeclaration) EffectiveAccess() MethodAccess {
	if declaration.Access == MethodAccessRead {
		return MethodAccessRead
	}
	return MethodAccessWrite
}

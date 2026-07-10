package agentapi

import (
	"encoding/json"

	agentmodel "go.proteos.ai/model/agent"
	"go.proteos.ai/model/common"
)

// Binding is carried on the wire as raw JSON: it can only be decoded once `kind`
// is known, so the service decodes it (via agentmodel.DecodeToolBinding) rather
// than the JSON binder. client tools omit it.
//
// Input/output schema are not authored — they are resolved on read from the
// binding's source — so neither request carries a schema field.
type CreateToolRequest struct {
	Key         string              `json:"key" validate:"required"`
	Name        string              `json:"name" validate:"required"`
	ModuleSlug  string              `json:"module_slug"`
	Description string              `json:"description"`
	Kind        agentmodel.ToolKind `json:"kind" validate:"required"`
	Binding     json.RawMessage     `json:"binding,omitempty"`
}

// UpdateToolRequest fully replaces the tool's definition (the kind↔binding
// coupling makes a partial patch ambiguous); version is bumped server-side.
// ModuleSlug is preserved when empty (see tool-service.UpdateOne) so a UI edit
// that omits it never orphans the tool from its module.
type UpdateToolRequest struct {
	Name        string              `json:"name" validate:"required"`
	ModuleSlug  string              `json:"module_slug"`
	Description string              `json:"description"`
	Kind        agentmodel.ToolKind `json:"kind" validate:"required"`
	Binding     json.RawMessage     `json:"binding,omitempty"`
}

type GetManyToolsQuery struct {
	Key          *string `json:"key" form:"key" db:"key"`
	Name         *string `json:"name" form:"name" db:"name"`
	ModuleSlug   *string `json:"module_slug" form:"module_slug" db:"module_slug"`
	NameContains *string `json:"name[contains]" form:"name[contains]" db:"name" op:"contains"`
	Kind         *string `json:"kind" form:"kind" db:"kind"`
	// Expand opts into read-time schema resolution per tool (expand=schema).
	// Off by default so a plain list avoids N function-service / MCP round trips.
	Expand string `json:"expand" form:"expand"`
	common.Pagination
	common.Sorting
}

type GetManyToolsResponse struct {
	Meta common.ResponseMeta `json:"meta"`
	Data []agentmodel.Tool   `json:"data"`
}

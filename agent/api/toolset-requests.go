package agentapi

import (
	agentmodel "go.proteos.ai/model/agent"
	"go.proteos.ai/model/common"
)

// Toolset writes create/replace CUSTOM toolsets only — platform toolsets are
// hardcoded server-side and read-only, so no request carries a kind. Tools are
// the member Tool keys; membership existence is validated on write.
type CreateToolsetRequest struct {
	Key         string   `json:"key" validate:"required"`
	Name        string   `json:"name" validate:"required"`
	ModuleSlug  string   `json:"module_slug"`
	Description string   `json:"description"`
	Tools       []string `json:"tools"`
}

// UpdateToolsetRequest fully replaces the custom toolset's definition (membership
// is a set — a partial patch is ambiguous); version is bumped server-side.
// ModuleSlug is preserved when empty, mirroring tool updates.
type UpdateToolsetRequest struct {
	Name        string   `json:"name" validate:"required"`
	ModuleSlug  string   `json:"module_slug"`
	Description string   `json:"description"`
	Tools       []string `json:"tools"`
}

type GetManyToolsetsQuery struct {
	Key          *string `json:"key" form:"key" db:"key"`
	Name         *string `json:"name" form:"name" db:"name"`
	ModuleSlug   *string `json:"module_slug" form:"module_slug" db:"module_slug"`
	NameContains *string `json:"name[contains]" form:"name[contains]" db:"name" op:"contains"`
	// Kind filters the merged listing (platform | custom). Applied in the service —
	// platform toolsets are not rows, so it is not a repository column filter.
	Kind *string `json:"kind" form:"kind"`
	common.Pagination
	common.Sorting
}

type GetManyToolsetsResponse struct {
	Meta common.ResponseMeta `json:"meta"`
	Data []agentmodel.Toolset `json:"data"`
}

// ListToolsetToolsResponse lists the tools inside one toolset — platform members
// proxied from mcp-service's tools/list, custom members summarized from the
// org's Tool rows.
type ListToolsetToolsResponse struct {
	Data []agentmodel.ToolsetToolSummary `json:"data"`
}

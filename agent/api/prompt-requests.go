package agentapi

import (
	agentmodel "go.proteos.ai/model/agent"
	"go.proteos.ai/model/common"
	metamodel "go.proteos.ai/model/meta"
)

type CreatePromptRequest struct {
	Key         string                `json:"key" validate:"required"`
	Name        string                `json:"name" validate:"required"`
	ModuleSlug  string                `json:"module_slug"`
	Description string                `json:"description"`
	Body        string                `json:"body" validate:"required"`
	Inputs      []metamodel.Attribute `json:"inputs"`
}

// UpdatePromptRequest is a patch: Name/Description edit metadata in place; supplying
// Body forks a new immutable version (Inputs carries forward from the current version
// unless also supplied). A metadata-only patch creates no version.
type UpdatePromptRequest struct {
	Name        *string                `json:"name,omitempty"`
	ModuleSlug  *string                `json:"module_slug,omitempty"`
	Description *string                `json:"description,omitempty"`
	Body        *string                `json:"body,omitempty"`
	Inputs      *[]metamodel.Attribute `json:"inputs,omitempty"`
}

type GetManyPromptsQuery struct {
	Key          *string `json:"key" form:"key" db:"key"`
	Name         *string `json:"name" form:"name" db:"name"`
	ModuleSlug   *string `json:"module_slug" form:"module_slug" db:"module_slug"`
	NameContains *string `json:"name[contains]" form:"name[contains]" db:"name" op:"contains"`
	common.Pagination
	common.Sorting
}

type GetManyPromptsResponse struct {
	Meta common.ResponseMeta `json:"meta"`
	Data []agentmodel.Prompt `json:"data"`
}

type GetManyPromptVersionsResponse struct {
	Data []agentmodel.PromptVersion `json:"data"`
}

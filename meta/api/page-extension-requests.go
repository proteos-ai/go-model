package metaapi

import (
	"go.proteos.ai/model/common"
	"go.proteos.ai/model/meta"
)

type CreatePageExtensionRequest struct {
	Key         string `json:"key" validate:"required"`
	Name        string `json:"name"`
	Description string `json:"description"`
	// PageSlug is the host page. Immutable after create —
	// UpdatePageExtensionRequest carries no page_slug; delete and recreate to
	// re-host.
	PageSlug   string `json:"page_slug" validate:"required"`
	ModuleSlug string `json:"module_slug"`
	// Placements are the contributed blocks: at least one, every element and
	// tab with an authored id, no extension_key (the server stamps it on the
	// host).
	Placements []metamodel.PageExtensionPlacement `json:"placements" validate:"required"`
}

// UpdatePageExtensionRequest is a partial update. page_slug is the row's
// binding and not updatable. placements, when present, replaces the whole
// contributed list.
type UpdatePageExtensionRequest struct {
	Name        *string                             `json:"name,omitempty"`
	Description *string                             `json:"description,omitempty"`
	ModuleSlug  *string                             `json:"module_slug,omitempty"`
	Placements  *[]metamodel.PageExtensionPlacement `json:"placements,omitempty"`
}

type GetManyPageExtensionsQuery struct {
	Key          *string `json:"key" db:"key"`
	Name         *string `json:"name" db:"name"`
	NameContains *string `json:"name[contains]" db:"name" op:"contains"`
	PageSlug     *string `json:"page_slug" db:"page_slug"`
	ModuleSlug   *string `json:"module_slug" db:"module_slug"`
	CreatedBy    *string `json:"created_by" db:"created_by->>'id'"`
	UpdatedBy    *string `json:"updated_by" db:"updated_by->>'id'"`
	common.Pagination
	common.Sorting
}

type GetManyPageExtensionsResponse struct {
	Meta common.ResponseMeta       `json:"meta"`
	Data []metamodel.PageExtension `json:"data"`
}

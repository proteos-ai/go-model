package metaapi

import (
	"go.proteos.ai/model/common"
	"go.proteos.ai/model/meta"
)

type CreateMenuConfigurationExtensionRequest struct {
	Key         string `json:"key" validate:"required"`
	Name        string `json:"name"`
	Description string `json:"description"`
	// MenuSlug is the host menu configuration. Immutable after create —
	// UpdateMenuConfigurationExtensionRequest carries no menu_slug; delete and
	// recreate to re-host.
	MenuSlug   string `json:"menu_slug" validate:"required"`
	ModuleSlug string `json:"module_slug"`
	// Placements are the contributed blocks: at least one, every item with an
	// authored id, no extension_key (the server stamps it on the host).
	Placements []metamodel.MenuConfigurationExtensionPlacement `json:"placements" validate:"required"`
}

// UpdateMenuConfigurationExtensionRequest is a partial update. menu_slug is
// the row's binding and not updatable. placements, when present, replaces the
// whole contributed list.
type UpdateMenuConfigurationExtensionRequest struct {
	Name        *string                                          `json:"name,omitempty"`
	Description *string                                          `json:"description,omitempty"`
	ModuleSlug  *string                                          `json:"module_slug,omitempty"`
	Placements  *[]metamodel.MenuConfigurationExtensionPlacement `json:"placements,omitempty"`
}

type GetManyMenuConfigurationExtensionsQuery struct {
	Key          *string `json:"key" db:"key"`
	Name         *string `json:"name" db:"name"`
	NameContains *string `json:"name[contains]" db:"name" op:"contains"`
	MenuSlug     *string `json:"menu_slug" db:"menu_slug"`
	ModuleSlug   *string `json:"module_slug" db:"module_slug"`
	CreatedBy    *string `json:"created_by" db:"created_by->>'id'"`
	UpdatedBy    *string `json:"updated_by" db:"updated_by->>'id'"`
	common.Pagination
	common.Sorting
}

type GetManyMenuConfigurationExtensionsResponse struct {
	Meta common.ResponseMeta                    `json:"meta"`
	Data []metamodel.MenuConfigurationExtension `json:"data"`
}

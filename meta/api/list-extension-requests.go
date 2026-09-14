package metaapi

import (
	"go.proteos.ai/model/common"
	"go.proteos.ai/model/meta"
)

type CreateListExtensionRequest struct {
	Key         string `json:"key" validate:"required"`
	Name        string `json:"name"`
	Description string `json:"description"`
	// ListSlug is the host list. Immutable after create —
	// UpdateListExtensionRequest carries no list_slug; delete and recreate to
	// re-host.
	ListSlug   string `json:"list_slug" validate:"required"`
	ModuleSlug string `json:"module_slug"`
	// Columns are the contributed column placements; Actions the appended
	// toolbar buttons. At least one of the two must be non-empty; every
	// contributed column and action is unstamped (the server stamps the host).
	Columns []metamodel.ListExtensionColumnPlacement `json:"columns"`
	Actions []metamodel.PageAction                   `json:"actions"`
}

// UpdateListExtensionRequest is a partial update. list_slug is the row's
// binding and not updatable. columns / actions, when present, each replace
// their whole list.
type UpdateListExtensionRequest struct {
	Name        *string                                   `json:"name,omitempty"`
	Description *string                                   `json:"description,omitempty"`
	ModuleSlug  *string                                   `json:"module_slug,omitempty"`
	Columns     *[]metamodel.ListExtensionColumnPlacement `json:"columns,omitempty"`
	Actions     *[]metamodel.PageAction                   `json:"actions,omitempty"`
}

type GetManyListExtensionsQuery struct {
	Key          *string `json:"key" db:"key"`
	Name         *string `json:"name" db:"name"`
	NameContains *string `json:"name[contains]" db:"name" op:"contains"`
	ListSlug     *string `json:"list_slug" db:"list_slug"`
	ModuleSlug   *string `json:"module_slug" db:"module_slug"`
	CreatedBy    *string `json:"created_by" db:"created_by->>'id'"`
	UpdatedBy    *string `json:"updated_by" db:"updated_by->>'id'"`
	common.Pagination
	common.Sorting
}

type GetManyListExtensionsResponse struct {
	Meta common.ResponseMeta       `json:"meta"`
	Data []metamodel.ListExtension `json:"data"`
}

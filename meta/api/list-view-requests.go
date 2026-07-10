package metaapi

import (
	"go.proteos.ai/model/common"
	"go.proteos.ai/model/meta"
)

type CreateListViewRequest struct {
	Slug       string                 `json:"slug" validate:"required"`
	ModuleSlug string                 `json:"module_slug"`
	ListSlug   string                 `json:"list_slug" validate:"required"`
	Name       string                 `json:"name" validate:"required"`
	Columns    []metamodel.Column     `json:"columns" validate:"required"`
	Sorting    []metamodel.SortConfig `json:"sorting"`
	Filters    []common.FilterGroup   `json:"filters"`
}

type UpdateListViewRequest struct {
	Name       *string                 `json:"name,omitempty"`
	ModuleSlug *string                 `json:"module_slug,omitempty"`
	Columns    *[]metamodel.Column     `json:"columns,omitempty"`
	Sorting    *[]metamodel.SortConfig `json:"sorting,omitempty"`
	Filters    *[]common.FilterGroup   `json:"filters,omitempty"`
}

type GetManyListViewsQuery struct {
	Slug       *string `json:"slug" db:"slug"`
	Name       *string `json:"name" db:"name"`
	ModuleSlug *string `json:"module_slug" db:"module_slug"`
	ListSlug   *string `json:"list_slug" db:"list_slug"`
	CreatedBy  *string `json:"created_by" db:"created_by->>'id'"`
	UpdatedBy  *string `json:"updated_by" db:"updated_by->>'id'"`
	common.Pagination
	common.Sorting
}

type GetManyListViewsResponse struct {
	Meta common.ResponseMeta  `json:"meta"`
	Data []metamodel.ListView `json:"data"`
}

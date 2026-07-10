package metaapi

import (
	"go.proteos.ai/model/common"
	"go.proteos.ai/model/meta"
)

type CreateListRequest struct {
	Slug       string                 `json:"slug" validate:"required"`
	ModuleSlug string                 `json:"module_slug"`
	EntitySlug string                 `json:"entity_slug" validate:"required"`
	Name       string                 `json:"name" validate:"required"`
	Columns    []metamodel.Column     `json:"columns" validate:"required"`
	Sorting    []metamodel.SortConfig `json:"sorting"`
	Filters    []common.FilterGroup   `json:"filters"`
}

type UpdateListRequest struct {
	Name       *string                 `json:"name,omitempty"`
	ModuleSlug *string                 `json:"module_slug,omitempty"`
	Columns    *[]metamodel.Column     `json:"columns,omitempty"`
	Sorting    *[]metamodel.SortConfig `json:"sorting,omitempty"`
	Filters    *[]common.FilterGroup   `json:"filters,omitempty"`
}

type GetManyListsQuery struct {
	Slug         *string `json:"slug" db:"slug"`
	Name         *string `json:"name" db:"name"`
	ModuleSlug   *string `json:"module_slug" db:"module_slug"`
	NameContains *string `json:"name[contains]" db:"name" op:"contains"`
	EntitySlug   *string `json:"entity_slug" db:"entity_slug"`
	CreatedBy    *string `json:"created_by" db:"created_by->>'id'"`
	UpdatedBy    *string `json:"updated_by" db:"updated_by->>'id'"`
	common.Pagination
	common.Sorting
}

type GetManyListsResponse struct {
	Meta common.ResponseMeta `json:"meta"`
	Data []metamodel.List    `json:"data"`
}

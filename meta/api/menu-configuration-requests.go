package metaapi

import (
	"go.proteos.ai/model/common"
	"go.proteos.ai/model/meta"
)

type CreateMenuConfigurationRequest struct {
	Slug       string               `json:"slug" validate:"required"`
	ModuleSlug string               `json:"module_slug"`
	Name       string               `json:"name" validate:"required"`
	AppSlug    string               `json:"app_slug" validate:"required"`
	Items      []metamodel.MenuItem `json:"items"`
	IsDefault  bool                 `json:"is_default"`
}

type UpdateMenuConfigurationRequest struct {
	Name       *string               `json:"name,omitempty"`
	ModuleSlug *string               `json:"module_slug,omitempty"`
	AppSlug    *string               `json:"app_slug,omitempty"`
	Items      *[]metamodel.MenuItem `json:"items,omitempty"`
	IsDefault  *bool                 `json:"is_default,omitempty"`
}

type GetManyMenuConfigurationsQuery struct {
	Id         *string `json:"id" db:"id"`
	Slug       *string `json:"slug" db:"slug"`
	ModuleSlug *string `json:"module_slug" db:"module_slug"`
	AppSlug    *string `json:"app_slug" db:"app_slug"`
	Name       *string `json:"name" db:"name"`
	IsDefault  *bool   `json:"is_default" db:"is_default"`
	CreatedBy  *string `json:"created_by" db:"created_by->>'id'"`
	UpdatedBy  *string `json:"updated_by" db:"updated_by->>'id'"`
	common.Pagination
	common.Sorting
}

type GetManyMenuConfigurationsResponse struct {
	Meta common.ResponseMeta           `json:"meta"`
	Data []metamodel.MenuConfiguration `json:"data"`
}

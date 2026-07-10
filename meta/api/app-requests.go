package metaapi

import (
	"go.proteos.ai/model/common"
	"go.proteos.ai/model/meta"
)

type CreateAppRequest struct {
	Slug        string `json:"slug" validate:"required"`
	Name        string `json:"name" validate:"required"`
	ModuleSlug  string `json:"module_slug"`
	Description string `json:"description"`
	IconSlug    string `json:"icon_slug" validate:"required"`
}

type UpdateAppRequest struct {
	Name        *string `json:"name,omitempty"`
	ModuleSlug  *string `json:"module_slug,omitempty"`
	Description *string `json:"description,omitempty"`
	IconSlug    *string `json:"icon_slug,omitempty"`
}

type GetManyAppsQuery struct {
	Slug       *string `json:"slug" db:"slug"`
	Name       *string `json:"name" db:"name"`
	ModuleSlug *string `json:"module_slug" db:"module_slug"`
	IconSlug   *string `json:"icon_slug" db:"icon_slug"`
	CreatedBy  *string `json:"created_by" db:"created_by->>'id'"`
	UpdatedBy  *string `json:"updated_by" db:"updated_by->>'id'"`
	common.Pagination
	common.Sorting
}

type GetManyAppsResponse struct {
	Meta common.ResponseMeta `json:"meta"`
	Data []metamodel.App     `json:"data"`
}

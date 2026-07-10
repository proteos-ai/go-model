package metaapi

import (
	"go.proteos.ai/model/common"
	"go.proteos.ai/model/meta"
)

type CreateEntityRequest struct {
	Slug          string                `json:"slug" validate:"required"`
	Name          string                `json:"name"`
	IsRemote      bool                  `json:"is_remote"`
	ModuleSlug    string                `json:"module_slug"`
	Description   string                `json:"description"`
	TitleTemplate string                `json:"title_template"`
	Attributes    []metamodel.Attribute `json:"attributes"`
}

type UpdateEntityRequest struct {
	Name          *string                `json:"name,omitempty"`
	IsRemote      *bool                  `json:"is_remote,omitempty"`
	ModuleSlug    *string                `json:"module_slug,omitempty"`
	Description   *string                `json:"description,omitempty"`
	TitleTemplate *string                `json:"title_template,omitempty"`
	Attributes    *[]metamodel.Attribute `json:"attributes,omitempty"`
}

// GetOneEntityQuery contains query parameters for GetOne endpoint
type GetOneEntityQuery struct {
	WithSchema bool `json:"with_schema" form:"with_schema"`
}

type GetManyEntitiesQuery struct {
	Slug         *string `json:"slug" db:"slug"`
	Name         *string `json:"name" db:"name"`
	NameContains *string `json:"name[contains]" db:"name" op:"contains"`
	IsRemote     *bool   `json:"is_remote" db:"is_remote"`
	ModuleSlug   *string `json:"module_slug" db:"module_slug"`
	CreatedBy    *string `json:"created_by" db:"created_by->>'id'"`
	UpdatedBy    *string `json:"updated_by" db:"updated_by->>'id'"`
	WithSchema   *bool   `json:"with_schema"`
	common.Pagination
	common.Sorting
}

type GetManyEntitiesResponse struct {
	Meta common.ResponseMeta `json:"meta"`
	Data []metamodel.Entity  `json:"data"`
}

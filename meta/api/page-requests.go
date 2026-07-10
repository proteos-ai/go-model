package metaapi

import (
	"go.proteos.ai/model/common"
	"go.proteos.ai/model/meta"
)

type CreatePageRequest struct {
	Slug string `json:"slug" validate:"required"`
	Name string `json:"name" validate:"required"`
	// Type defaults to "record" when empty. The type↔entity_slug invariant
	// (record requires entity_slug, platform forbids it) is enforced in the
	// service's domain logic, not via struct tags.
	ModuleSlug string                 `json:"module_slug"`
	Type       metamodel.PageType     `json:"type"`
	EntitySlug string                 `json:"entity_slug"`
	Actions    []metamodel.PageAction `json:"actions"`
	Layout     metamodel.PageLayout   `json:"layout" validate:"required"`
}

type UpdatePageRequest struct {
	Name       *string                 `json:"name,omitempty"`
	ModuleSlug *string                 `json:"module_slug,omitempty"`
	Actions    *[]metamodel.PageAction `json:"actions,omitempty"`
	Layout     *metamodel.PageLayout   `json:"layout,omitempty"`
}

type GetManyPagesQuery struct {
	Id         *string `json:"id" db:"id"`
	Slug       *string `json:"slug" db:"slug"`
	ModuleSlug *string `json:"module_slug" db:"module_slug"`
	Type       *string `json:"type" db:"type"`
	EntitySlug *string `json:"entity_slug" db:"entity_slug"`
	Name       *string `json:"name" db:"name"`
	CreatedBy  *string `json:"created_by" db:"created_by->>'id'"`
	UpdatedBy  *string `json:"updated_by" db:"updated_by->>'id'"`
	common.Pagination
	common.Sorting
}

type GetManyPagesResponse struct {
	Meta common.ResponseMeta `json:"meta"`
	Data []metamodel.Page    `json:"data"`
}

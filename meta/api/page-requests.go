package metaapi

import (
	"go.proteos.ai/model/common"
	"go.proteos.ai/model/meta"
)

type CreatePageRequest struct {
	Slug string `json:"slug" validate:"required"`
	Name string `json:"name" validate:"required"`
	// Type defaults to "record" when empty. The type↔entity_slug invariant
	// (record requires entity_slug, standalone types forbid it) is enforced in
	// the service's domain logic, not via struct tags.
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

// PublicPageComponent is the slice of component metadata a public page's
// renderer needs — just the slug and props schema. Deliberately NOT the full
// Component (no file ids, no audit fields): this rides on an unauthenticated
// response.
type PublicPageComponent struct {
	Slug        string         `json:"slug"`
	PropsSchema map[string]any `json:"props_schema,omitempty"`
}

// PublicPageResponse is the payload of the unauthenticated
// GET /meta/v1/public/orgs/{orgId}/pages/{slug}: the page plus the
// props_schema of every component its layout references, so the renderer
// needs no follow-up authenticated calls.
type PublicPageResponse struct {
	Page       metamodel.Page        `json:"page"`
	Components []PublicPageComponent `json:"components"`
}

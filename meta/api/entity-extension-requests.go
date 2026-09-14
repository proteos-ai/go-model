package metaapi

import (
	"go.proteos.ai/model/common"
	"go.proteos.ai/model/meta"
)

type CreateEntityExtensionRequest struct {
	Key         string `json:"key" validate:"required"`
	Name        string `json:"name"`
	Description string `json:"description"`
	// EntitySlug is the host entity. Immutable after create —
	// UpdateEntityExtensionRequest carries no entity_slug; delete and recreate
	// to re-host.
	EntitySlug string `json:"entity_slug" validate:"required"`
	ModuleSlug string `json:"module_slug"`
	// Attributes are the contributed definitions: at least one, no platform
	// names, no extension_key (the server stamps it on the host entity).
	Attributes []metamodel.Attribute `json:"attributes" validate:"required"`
}

// UpdateEntityExtensionRequest is a partial update. entity_slug is the row's
// binding and not updatable. attributes, when present, replaces the whole
// contributed list (same full-replacement semantics as an entity's attributes).
type UpdateEntityExtensionRequest struct {
	Name        *string                `json:"name,omitempty"`
	Description *string                `json:"description,omitempty"`
	ModuleSlug  *string                `json:"module_slug,omitempty"`
	Attributes  *[]metamodel.Attribute `json:"attributes,omitempty"`
}

type GetManyEntityExtensionsQuery struct {
	Key          *string `json:"key" db:"key"`
	Name         *string `json:"name" db:"name"`
	NameContains *string `json:"name[contains]" db:"name" op:"contains"`
	EntitySlug   *string `json:"entity_slug" db:"entity_slug"`
	ModuleSlug   *string `json:"module_slug" db:"module_slug"`
	CreatedBy    *string `json:"created_by" db:"created_by->>'id'"`
	UpdatedBy    *string `json:"updated_by" db:"updated_by->>'id'"`
	common.Pagination
	common.Sorting
}

type GetManyEntityExtensionsResponse struct {
	Meta common.ResponseMeta         `json:"meta"`
	Data []metamodel.EntityExtension `json:"data"`
}

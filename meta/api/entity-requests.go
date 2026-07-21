package metaapi

import (
	"go.proteos.ai/model/common"
	"go.proteos.ai/model/meta"
)

type CreateEntityRequest struct {
	Slug     string `json:"slug" validate:"required"`
	Name     string `json:"name"`
	IsRemote bool   `json:"is_remote"`
	// PublicRecordAccess opts all records of the entity into unauthenticated
	// access (see metamodel.Entity.PublicRecordAccess; only ["read"] honored
	// today). Manifest-driven full-replacement on upsert: an upsert without
	// the field resets it to private.
	PublicRecordAccess common.PublicAccess   `json:"public_record_access"`
	ModuleSlug         string                `json:"module_slug"`
	Description        string                `json:"description"`
	TitleTemplate      string                `json:"title_template"`
	Attributes         []metamodel.Attribute `json:"attributes"`
}

type UpdateEntityRequest struct {
	Name               *string                `json:"name,omitempty"`
	IsRemote           *bool                  `json:"is_remote,omitempty"`
	PublicRecordAccess *common.PublicAccess   `json:"public_record_access,omitempty"`
	ModuleSlug         *string                `json:"module_slug,omitempty"`
	Description        *string                `json:"description,omitempty"`
	TitleTemplate      *string                `json:"title_template,omitempty"`
	Attributes         *[]metamodel.Attribute `json:"attributes,omitempty"`
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

// InboundRelation is one relation attribute on EntitySlug that references a
// target entity, together with the delete policy to apply when a target record
// is removed. Returned by GET /meta/v1/entities/:slug/inbound-relations and
// consumed by data-service's record-delete cascade.
type InboundRelation struct {
	// EntitySlug is the referencing (child) entity.
	EntitySlug string `json:"entity_slug"`
	// AttributeName is the relation attribute on EntitySlug holding the reference.
	AttributeName string `json:"attribute_name"`
	// RelatedAttribute is the target-entity attribute the reference points at
	// (normally "id").
	RelatedAttribute string `json:"related_attribute"`
	// OnDelete is the policy: cascade | restrict | set-null.
	OnDelete metamodel.OnDeleteAction `json:"on_delete"`
}

type GetInboundRelationsResponse struct {
	Data []InboundRelation `json:"data"`
}

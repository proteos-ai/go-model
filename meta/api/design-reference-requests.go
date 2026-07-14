package metaapi

import (
	"go.proteos.ai/model/common"
	"go.proteos.ai/model/meta"
)

type CreateDesignReferenceRequest struct {
	Slug        string `json:"slug" validate:"required"`
	Name        string `json:"name" validate:"required"`
	Description string `json:"description"`
	// Content optionally seeds the DESIGN.md body at create time. It can be edited
	// afterwards only via PUT /design-references/:id/content.
	Content string `json:"content"`
}

type UpdateDesignReferenceRequest struct {
	Slug        *string `json:"slug,omitempty"`
	Name        *string `json:"name,omitempty"`
	Description *string `json:"description,omitempty"`
	// Content is deliberately NOT patchable here — it has its own endpoint
	// (PUT /design-references/:id/content) so a metadata edit can never blank the
	// document.
}

// SetDesignReferenceContentRequest is the body of PUT /design-references/:id/content.
type SetDesignReferenceContentRequest struct {
	Content string `json:"content"`
}

// DesignReferenceContentResponse is the body of GET /design-references/:id/content.
type DesignReferenceContentResponse struct {
	Content string `json:"content"`
}

type GetManyDesignReferencesQuery struct {
	Id          *string `json:"id" db:"id"`
	Slug        *string `json:"slug" db:"slug"`
	Name        *string `json:"name" db:"name"`
	Description *string `json:"description" db:"description"`
	CreatedBy   *string `json:"created_by" db:"created_by->>'id'"`
	UpdatedBy   *string `json:"updated_by" db:"updated_by->>'id'"`
	common.Pagination
	common.Sorting
}

type GetManyDesignReferencesResponse struct {
	Meta common.ResponseMeta         `json:"meta"`
	Data []metamodel.DesignReference `json:"data"`
}

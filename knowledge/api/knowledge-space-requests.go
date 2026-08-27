package knowledgeapi

import (
	"go.proteos.ai/model/common"
	knowledgemodel "go.proteos.ai/model/knowledge"
)

type CreateKnowledgeSpaceRequest struct {
	Slug        string  `json:"slug" validate:"required"`
	Name        string  `json:"name" validate:"required"`
	Description *string `json:"description,omitempty"`
	Color       *string `json:"color,omitempty"`
	Icon        *string `json:"icon,omitempty"`
}

// UpdateKnowledgeSpaceRequest is a partial update. `slug` is deliberately absent:
// it is half the primary key, node rows reference it, and access grants name it —
// renaming would orphan both, so a rename is creating a different space.
type UpdateKnowledgeSpaceRequest struct {
	Name        *string `json:"name,omitempty"`
	Description *string `json:"description,omitempty"`
	Color       *string `json:"color,omitempty"`
	Icon        *string `json:"icon,omitempty"`
}

type GetManyKnowledgeSpacesQuery struct {
	Slug         *string `json:"slug" db:"slug"`
	Name         *string `json:"name" db:"name"`
	NameContains *string `json:"name[contains]" db:"name" op:"contains"`
	common.Pagination
	common.Sorting
}

type GetManyKnowledgeSpacesResponse struct {
	Meta common.ResponseMeta             `json:"meta"`
	Data []knowledgemodel.KnowledgeSpace `json:"data"`
}

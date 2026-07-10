package knowledgeapi

import (
	"go.proteos.ai/model/common"
	knowledgemodel "go.proteos.ai/model/knowledge"
)

type CreateKnowledgeLabelRequest struct {
	Name        string  `json:"name" validate:"required"`
	Slug        string  `json:"slug" validate:"required"`
	Description *string `json:"description,omitempty"`
	Color       *string `json:"color,omitempty"`
	Icon        *string `json:"icon,omitempty"`
}

type UpdateKnowledgeLabelRequest struct {
	Name        *string `json:"name,omitempty"`
	Description *string `json:"description,omitempty"`
	Color       *string `json:"color,omitempty"`
	Icon        *string `json:"icon,omitempty"`
}

type GetManyKnowledgeLabelsQuery struct {
	Slug         *string `json:"slug" db:"slug"`
	Name         *string `json:"name" db:"name"`
	NameContains *string `json:"name[contains]" db:"name" op:"contains"`
	common.Pagination
	common.Sorting
}

type GetManyKnowledgeLabelsResponse struct {
	Meta common.ResponseMeta             `json:"meta"`
	Data []knowledgemodel.KnowledgeLabel `json:"data"`
}

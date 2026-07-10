package knowledgeapi

import (
	"go.proteos.ai/model/common"
	knowledgemodel "go.proteos.ai/model/knowledge"
)

type CreateKnowledgeLinkRequest struct {
	FromId      string  `json:"from_id" validate:"required"`
	ToId        string  `json:"to_id" validate:"required"`
	Type        string  `json:"type" validate:"required,oneof=references relates_to depends_on part_of derived_from contradicts supports duplicates superseded_by"`
	Description *string `json:"description,omitempty"`
}

type UpdateKnowledgeLinkRequest struct {
	Type        *string `json:"type,omitempty" validate:"omitempty,oneof=references relates_to depends_on part_of derived_from contradicts supports duplicates superseded_by"`
	Description *string `json:"description,omitempty"`
}

type GetManyKnowledgeLinksQuery struct {
	FromId *string `json:"from_id" db:"from_id"`
	ToId   *string `json:"to_id" db:"to_id"`
	Type   *string `json:"type" db:"type"`
	// Date-range filters (RFC3339 timestamps). The `op` tag maps each field to a
	// comparison on its `db` column, so a pair sharing a column forms a closed
	// range (created_after + created_before).
	CreatedAfter  *string `json:"created_after" db:"created_at" op:"gte"`
	CreatedBefore *string `json:"created_before" db:"created_at" op:"lte"`
	UpdatedAfter  *string `json:"updated_after" db:"updated_at" op:"gte"`
	UpdatedBefore *string `json:"updated_before" db:"updated_at" op:"lte"`
	common.Pagination
	common.Sorting
}

type GetManyKnowledgeLinksResponse struct {
	Meta common.ResponseMeta            `json:"meta"`
	Data []knowledgemodel.KnowledgeLink `json:"data"`
}

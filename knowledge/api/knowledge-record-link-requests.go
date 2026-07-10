package knowledgeapi

import (
	"go.proteos.ai/model/common"
	knowledgemodel "go.proteos.ai/model/knowledge"
)

type CreateKnowledgeRecordLinkRequest struct {
	NodeId     string `json:"node_id" validate:"required"`
	EntitySlug string `json:"entity_slug" validate:"required"`
	RecordId   string `json:"record_id" validate:"required"`
}

type GetManyKnowledgeRecordLinksQuery struct {
	NodeId     *string `json:"node_id" db:"node_id"`
	EntitySlug *string `json:"entity_slug" db:"entity_slug"`
	RecordId   *string `json:"record_id" db:"record_id"`
	common.Pagination
	common.Sorting
}

type GetManyKnowledgeRecordLinksResponse struct {
	Meta common.ResponseMeta                  `json:"meta"`
	Data []knowledgemodel.KnowledgeRecordLink `json:"data"`
}

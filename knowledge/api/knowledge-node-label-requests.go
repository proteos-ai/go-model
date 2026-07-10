package knowledgeapi

import (
	knowledgemodel "go.proteos.ai/model/knowledge"
)

// AttachLabelRequest attaches an existing label to a node.
type AttachLabelRequest struct {
	LabelId string `json:"label_id" validate:"required"`
}

// ListNodeLabelsResponse returns the labels currently attached to a node.
type ListNodeLabelsResponse struct {
	Data []knowledgemodel.KnowledgeLabel `json:"data"`
}

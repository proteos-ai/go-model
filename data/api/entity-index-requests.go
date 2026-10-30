package dataapi

import datamodel "go.proteos.ai/model/data"

// === Entity indexes (/data/v1/entities/:slug/indexes) ===

// CreateEntityIndexRequest adds one APPLIED index to an entity's record table.
// Only lookup and search may be applied: uniqueness is declared on the entity
// (is_unique / unique_constraints) and automatic indexes are the platform's.
type CreateEntityIndexRequest struct {
	// Name is optional; derived from the attributes when omitted
	// (by_<a>_<b>, search_<a>). snake_case, unique within the entity.
	Name       string                    `json:"name,omitempty"`
	Type       datamodel.EntityIndexType `json:"type" binding:"required"`
	Attributes []string                  `json:"attributes" binding:"required"`
	// Replaces names APPLIED lookup indexes this one supersedes — each must
	// be covered by the new index (its attributes are a leading prefix). They
	// are removed as part of the create, so the table never carries both. An
	// advisor `add` recommendation fills this in.
	Replaces []string `json:"replaces,omitempty"`
}

type GetManyEntityIndexesResponse struct {
	Data []datamodel.EntityIndex `json:"data"`
}

type GetManyIndexRecommendationsResponse struct {
	Data []datamodel.IndexRecommendation `json:"data"`
}

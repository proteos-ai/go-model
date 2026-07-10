package knowledgeapi

import (
	"go.proteos.ai/model/common"
	knowledgemodel "go.proteos.ai/model/knowledge"
)

const (
	// MatchModeHybrid fuses lexical (tsvector) and semantic (pgvector) ranking.
	MatchModeHybrid = "hybrid"
	// MatchModeSemantic ranks purely by embedding similarity.
	MatchModeSemantic = "semantic"
	// MatchModeKeyword ranks purely by lexical ts_rank. The only mode wired today
	// — semantic/hybrid are reserved until embedding ingestion lands.
	MatchModeKeyword = "keyword"
)

// SearchNodesRequest is the body of POST /v1/nodes/actions/search — hybrid
// retrieval with optional metadata/graph filters. `match_mode` defaults to
// keyword; semantic and hybrid are accepted in the contract but not yet wired.
type SearchNodesRequest struct {
	Query     string `json:"query" validate:"required"`
	MatchMode string `json:"match_mode" validate:"omitempty,oneof=hybrid semantic keyword"`
	// Filters. LabelIds matches nodes carrying ANY of the given labels; LinkedToId
	// matches nodes connected (in either direction) to that node.
	Type       *string  `json:"type,omitempty" validate:"omitempty,oneof=markdown file url"`
	Status     *string  `json:"status,omitempty" validate:"omitempty,oneof=draft published archived"`
	LabelIds   []string `json:"label_ids,omitempty"`
	LinkedToId *string  `json:"linked_to_id,omitempty"`
	// Date-range filters (RFC3339). These only narrow the candidate set; ordering
	// stays relevance-based (score), so a date filter never reorders results.
	CreatedAfter  *string `json:"created_after,omitempty" validate:"omitempty,datetime=2006-01-02T15:04:05Z07:00"`
	CreatedBefore *string `json:"created_before,omitempty" validate:"omitempty,datetime=2006-01-02T15:04:05Z07:00"`
	UpdatedAfter  *string `json:"updated_after,omitempty" validate:"omitempty,datetime=2006-01-02T15:04:05Z07:00"`
	UpdatedBefore *string `json:"updated_before,omitempty" validate:"omitempty,datetime=2006-01-02T15:04:05Z07:00"`
	// Point-in-time validity filter: keep only hits valid AT IsValidAt (half-open,
	// null-aware), or INVALID at it when IsValid is false (default true). Narrows
	// the candidate set like the other filters; ordering stays relevance-based.
	IsValidAt *string `json:"is_valid_at,omitempty" validate:"omitempty,datetime=2006-01-02T15:04:05Z07:00"`
	IsValid   *bool   `json:"is_valid,omitempty"`
	common.Pagination
}

type SearchNodesResponse struct {
	Meta common.ResponseMeta                        `json:"meta"`
	Data []knowledgemodel.KnowledgeNodeSearchResult `json:"data"`
}

package knowledgeapi

import (
	"time"

	"go.proteos.ai/model/common"
	knowledgemodel "go.proteos.ai/model/knowledge"
)

type CreateKnowledgeNodeRequest struct {
	Title   string  `json:"title" validate:"required"`
	Type    string  `json:"type" validate:"required,oneof=markdown file url"`
	Status  string  `json:"status" validate:"required,oneof=draft published archived"`
	Content string  `json:"content"`
	FileId  *string `json:"file_id,omitempty"`
	Url     *string `json:"url,omitempty"`
	Summary *string `json:"summary,omitempty"`
	// ValidFrom / ValidUntil set the node's temporal validity window at creation.
	// Both optional and unbounded when omitted (null). Plain pointers — create has
	// no need to distinguish "clear" from "omit" (omitted = NULL = unbounded).
	ValidFrom  *time.Time `json:"valid_from,omitempty"`
	ValidUntil *time.Time `json:"valid_until,omitempty"`
	// LabelIds and Links let a caller create a node and connect it to the graph
	// in one call (create-and-connect). LabelIds attaches existing labels; Links
	// creates typed edges FROM the new node TO existing nodes. Both are applied
	// atomically with the node insert — if any target is missing the whole create
	// is rolled back.
	LabelIds []string            `json:"label_ids,omitempty"`
	Links    []InlineLinkRequest `json:"links,omitempty"`
}

// InlineLinkRequest is a typed edge from a node being created to an existing
// node, supplied inline on create. The `from` side is the new node, so only the
// target and the edge type are given here.
type InlineLinkRequest struct {
	ToId        string  `json:"to_id" validate:"required"`
	Type        string  `json:"type" validate:"required,oneof=references relates_to depends_on part_of derived_from contradicts supports duplicates superseded_by"`
	Description *string `json:"description,omitempty"`
}

// UpdateKnowledgeNodeRequest is a partial metadata update. `type` and the source
// pointers (file_id / url) are fixed at create time in v1, so they cannot be
// changed here — that keeps the type↔field invariant trivially intact across
// updates. The `content` body is NOT updatable here: it is owned by the content
// sub-resource (`PUT /nodes/:id/content` and `…/content/actions/edit`) so there
// is a single write path for the (potentially huge) body.
type UpdateKnowledgeNodeRequest struct {
	Title   *string `json:"title,omitempty"`
	Status  *string `json:"status,omitempty" validate:"omitempty,oneof=draft published archived"`
	Summary *string `json:"summary,omitempty"`
	// ValidFrom / ValidUntil are tri-state so the window can be set, moved, OR
	// cleared back to NULL (re-opening a fact's validity): absent = leave
	// unchanged, JSON null = clear, value = set. See common.Optional.
	ValidFrom  common.Optional[time.Time] `json:"valid_from"`
	ValidUntil common.Optional[time.Time] `json:"valid_until"`
}

type GetManyKnowledgeNodesQuery struct {
	Type          *string `json:"type" db:"type"`
	Status        *string `json:"status" db:"status"`
	Title         *string `json:"title" db:"title"`
	TitleContains *string `json:"title[contains]" db:"title" op:"contains"`
	// LabelIds keeps only nodes carrying ANY of the given labels (union). Sent as a
	// single comma-separated query param (the URL binder reads only the first value
	// of a repeated param, so the union list is passed CSV-joined). No `db` tag — the
	// repository applies the EXISTS predicate manually, like IsValidAt.
	LabelIds string `json:"label_ids"`
	// Date-range filters (RFC3339 timestamps). The `op` tag maps each field to a
	// comparison on its `db` column, so a pair sharing a column forms a closed
	// range (created_after + created_before).
	CreatedAfter  *string `json:"created_after" db:"created_at" op:"gte"`
	CreatedBefore *string `json:"created_before" db:"created_at" op:"lte"`
	UpdatedAfter  *string `json:"updated_after" db:"updated_at" op:"gte"`
	UpdatedBefore *string `json:"updated_before" db:"updated_at" op:"lte"`
	// Validity boundary range filters — query where a node's window start/end
	// falls (e.g. "facts expiring soon": valid_until_after=now & valid_until_before=now+30d).
	// Same db+op tag mechanism as the created/updated ranges.
	ValidFromAfter   *string `json:"valid_from_after" db:"valid_from" op:"gte"`
	ValidFromBefore  *string `json:"valid_from_before" db:"valid_from" op:"lte"`
	ValidUntilAfter  *string `json:"valid_until_after" db:"valid_until" op:"gte"`
	ValidUntilBefore *string `json:"valid_until_before" db:"valid_until" op:"lte"`
	// IsValidAt is a point-in-time membership filter: keep only nodes valid AT this
	// instant (half-open, null-aware). IsValid flips it — false returns nodes
	// INVALID at the instant. No `db` tag: the repository applies the compound
	// predicate manually (UrlToDbQuery skips it).
	IsValidAt *string `json:"is_valid_at" validate:"omitempty,datetime=2006-01-02T15:04:05Z07:00"`
	IsValid   *bool   `json:"is_valid"`
	common.Pagination
	common.Sorting
}

type GetManyKnowledgeNodesResponse struct {
	Meta common.ResponseMeta                    `json:"meta"`
	Data []knowledgemodel.KnowledgeNodeMetadata `json:"data"`
}

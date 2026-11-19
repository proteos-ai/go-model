package datamodel

import (
	"time"

	"go.proteos.ai/model/common"
)

// QueryBuilder is the visual editor's state for a Query: the entity it reads,
// the selected column paths (relation fields are prefixed with the foreign-key
// attribute, e.g. `accountId.name`) and the filter. It only reopens the
// editor — it is never what runs.
type QueryBuilder struct {
	EntitySlug string              `json:"entity_slug"`
	Columns    []string            `json:"columns"`
	Filter     *common.FilterGroup `json:"filter,omitempty"`
}

// Query is a stored SQL query — what the Reports tab saves, reopens and
// shares. SQL is what runs, through POST /query/execute as the VIEWER, so a
// query never carries data access with it. Builder is nil for a query that
// was never built visually.
//
// CreatedBy is the owner. Everyone else reaches a query through a share
// (platform shares on the `queries` slug) or the unscoped entity permission.
// CanWrite says whether the CALLER may update this query; it is computed per
// request and never stored.
type Query struct {
	Id        string         `json:"id"`
	OrgId     string         `json:"org_id"`
	Name      string         `json:"name"`
	SQL       string         `json:"sql"`
	Builder   *QueryBuilder  `json:"builder,omitempty"`
	CanWrite  bool           `json:"can_write"`
	CreatedAt time.Time      `json:"created_at"`
	CreatedBy common.UserRef `json:"created_by"`
	UpdatedAt time.Time      `json:"updated_at"`
	UpdatedBy common.UserRef `json:"updated_by"`
}

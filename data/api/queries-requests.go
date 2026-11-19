package dataapi

import (
	"go.proteos.ai/model/common"
	datamodel "go.proteos.ai/model/data"
)

// CreateQueryRequest stores a query. The caller becomes its owner. The SQL is
// stored as text — nothing is executed or validated on save; every run goes
// through POST /query/execute as the viewer.
type CreateQueryRequest struct {
	Name    string                  `json:"name"`
	SQL     string                  `json:"sql"`
	Builder *datamodel.QueryBuilder `json:"builder,omitempty"`
}

// UpdateQueryRequest changes the named fields of a query; an absent field
// keeps its value.
type UpdateQueryRequest struct {
	Name    *string                 `json:"name,omitempty"`
	SQL     *string                 `json:"sql,omitempty"`
	Builder *datamodel.QueryBuilder `json:"builder,omitempty"`
}

// GetManyQueriesQuery pages the queries the caller can read: their own and
// those shared with them, their teams or the organization. Ordered by name.
type GetManyQueriesQuery struct {
	common.Pagination
}

type GetManyQueriesResponse struct {
	Meta common.ResponseMeta `json:"meta"`
	Data []datamodel.Query   `json:"data"`
}

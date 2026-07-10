package dataapi

import (
	"go.proteos.ai/model/data"
)

// QueryRequest represents a request to validate or execute a SQL query
type QueryRequest struct {
	SQL string `json:"sql" binding:"required" validate:"required"`
}

// QueryValidationMeta contains metadata about query validation
type QueryValidationMeta struct {
	LimitApplied           int64 `json:"limit_applied"`                       // The LIMIT value that would be applied
	WasDefaultLimitApplied bool  `json:"was_default_limit_applied,omitempty"` // Whether a default LIMIT was added
}

// QueryValidationResponse represents the response from validating a SQL query
type QueryValidationResponse struct {
	Valid        bool                 `json:"valid"`
	RewrittenSQL string               `json:"rewritten_sql,omitempty"` // The SQL that would be executed (with schema)
	Tables       []string             `json:"tables,omitempty"`        // Tables referenced in the query
	Meta         *QueryValidationMeta `json:"meta,omitempty"`
}

// QueryResponseMeta contains metadata about query execution
type QueryResponseMeta struct {
	Columns                []string `json:"columns"`
	Items                  int      `json:"items"`
	LimitApplied           int64    `json:"limit_applied"`
	WasDefaultLimitApplied bool     `json:"was_default_limit_applied,omitempty"` // Only present when true
	ExecutionTimeMs        int64    `json:"execution_time_ms"`
}

// QueryResponse represents the response from a SQL query
type QueryResponse struct {
	Data []datamodel.Record `json:"data"`
	Meta *QueryResponseMeta `json:"meta,omitempty"`
}

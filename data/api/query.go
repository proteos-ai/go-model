// Package dataapi carries the wire shapes for data-service's HTTP API.
// Lives in model/ (not sdk/) so that wasip1 guests in
// `go.proteos.ai/functions-sdk-go/fn` can share the exact same JSON
// shape with the host-side SDK without dragging in net/http.
package dataapi

// QueryRow is one row from a SQL query result. Keys are column names;
// values are JSON-decoded primitives, nested objects, or arrays.
type QueryRow map[string]any

// QueryExecuteMeta describes the result-set returned by /query/execute.
type QueryExecuteMeta struct {
	Columns                []string `json:"columns"`
	Items                  int      `json:"items"`
	LimitApplied           int      `json:"limit_applied"`
	WasDefaultLimitApplied bool     `json:"was_default_limit_applied,omitempty"`
	ExecutionTimeMs        int      `json:"execution_time_ms"`
}

// QueryExecuteResponse is the response shape of /query/execute.
type QueryExecuteResponse struct {
	Data []QueryRow        `json:"data"`
	Meta *QueryExecuteMeta `json:"meta,omitempty"`
}

// QueryValidateMeta is the meta returned from /query/validate.
type QueryValidateMeta struct {
	LimitApplied           int  `json:"limit_applied"`
	WasDefaultLimitApplied bool `json:"was_default_limit_applied,omitempty"`
}

// QueryValidateResponse is the response shape of /query/validate.
type QueryValidateResponse struct {
	Valid        bool               `json:"valid"`
	RewrittenSQL string             `json:"rewritten_sql,omitempty"`
	Tables       []string           `json:"tables,omitempty"`
	Meta         *QueryValidateMeta `json:"meta,omitempty"`
}

package common

// Pagination carries the page window of a list request/response. The `json`
// tags are the wire names — both for response meta and for query binding via
// httputils.BindQueryAsJson; the `form` tags cover controllers that still bind
// with gin's ShouldBindQuery (gin ignores `json` tags on query binds).
type Pagination struct {
	Page     int `json:"page" form:"page"`
	PageSize int `json:"page_size" form:"page_size"`
}

package knowledgeapi

// GetNeighborsQuery is the query string of GET /v1/nodes/:id/neighbors. All
// fields are optional: direction defaults to "both", depth to 1 (clamped to a
// max server-side), and an empty link_types means "any link type". link_types is
// a comma-separated list (e.g. "part_of,depends_on").
type GetNeighborsQuery struct {
	Direction string `json:"direction"`
	LinkTypes string `json:"link_types"`
	Depth     int    `json:"depth"`
}

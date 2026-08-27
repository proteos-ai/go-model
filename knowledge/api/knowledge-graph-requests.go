package knowledgeapi

// GetKnowledgeGraphQuery is the query string of GET /v1/graph. LabelIds is an
// optional comma-separated list of label ids; when present the graph is pruned
// to nodes carrying ANY of those labels (union) plus the links among them. Empty
// means the whole org graph. (Comma-separated, mirroring neighbors' link_types,
// because the query binder only supports primitive query params.)
type GetKnowledgeGraphQuery struct {
	LabelIds string `json:"label_ids"`
	// SpaceSlugs prunes the graph to nodes in ANY of these spaces (union),
	// comma-separated like LabelIds. Empty means every space. Unassigned nodes
	// (space_slug IS NULL) are excluded by a non-empty filter — ask for them
	// with the reserved sentinel, see UnassignedSpaceSlug.
	SpaceSlugs string `json:"space_slugs"`
}

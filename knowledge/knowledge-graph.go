package knowledgemodel

// KnowledgeGraphNode is a lean node projection for the whole-graph (radiant
// graph) view: only what the renderer needs — identity, type/status for
// styling, the ordered LabelIds (for color and label filtering) and Degree (for
// node sizing). It deliberately omits the content body, embedding provenance,
// source pointers and audit fields so a payload of tens of thousands of nodes
// stays small on the wire.
type KnowledgeGraphNode struct {
	Id     string `json:"id"`
	Title  string `json:"title"`
	Type   string `json:"type"`   // markdown | file | url
	Status string `json:"status"` // draft | published | archived
	// LabelIds are this node's attached labels, ordered by attachment time
	// (knowledge_node_labels.created_at ASC) — so LabelIds[0] is the "first"
	// label, which the UI uses to colour the node. Empty (never null) when the
	// node carries no labels.
	LabelIds []string `json:"label_ids"`
	// Degree is the count of links incident to this node (in + out) within the
	// returned graph. The renderer scales node size by degree.
	Degree int `json:"degree"`
}

// KnowledgeGraphLink is a lean directed-edge projection: just the endpoints and
// the typed relation, no audit fields.
type KnowledgeGraphLink struct {
	Id     string `json:"id"`
	FromId string `json:"from_id"`
	ToId   string `json:"to_id"`
	Type   string `json:"type"`
}

// KnowledgeGraph is the whole org graph in a single payload: lean nodes + links,
// the full org label set (sent once so the client can render every label — and
// colour every node — even when the graph is pruned by a label filter), and
// Total: the unfiltered org node count. Total lets the client decide whether to
// load the whole graph or switch to progressive/level-of-detail loading.
//
// Node→record links (KnowledgeRecordLink) are intentionally NOT part of this
// payload: they are surfaced only on the node detail panel, not the graph.
type KnowledgeGraph struct {
	Nodes  []KnowledgeGraphNode `json:"nodes"`
	Links  []KnowledgeGraphLink `json:"links"`
	Labels []KnowledgeLabel     `json:"labels"`
	Total  int                  `json:"total"`
}

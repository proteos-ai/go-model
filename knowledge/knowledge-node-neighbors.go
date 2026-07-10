package knowledgemodel

// NeighborNode is a node reached by a graph walk, with its hop distance from the
// origin (1 = directly linked). Bodies are excluded — a traversal returns shape,
// not content.
type NeighborNode struct {
	Id       string `json:"id"`
	Title    string `json:"title"`
	Type     string `json:"type"`
	Status   string `json:"status"`
	Distance int    `json:"distance"`
}

// NeighborEdge is a typed edge within the explored neighborhood (both endpoints
// are in the reachable set).
type NeighborEdge struct {
	Id     string `json:"id"`
	FromId string `json:"from_id"`
	ToId   string `json:"to_id"`
	Type   string `json:"type"`
}

// NodeNeighborhood is the result of GET /v1/nodes/:id/neighbors: the nodes
// reachable from the origin within the requested depth/direction, plus the edges
// connecting that neighborhood. The origin node itself is not included in Nodes.
type NodeNeighborhood struct {
	Nodes []NeighborNode `json:"nodes"`
	Edges []NeighborEdge `json:"edges"`
}

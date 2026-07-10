package common

// KnowledgeNodeRef is the stored value of a `knowledge-text` attribute — a
// reference to the knowledge node (knowledge-service) that owns the text body.
// The record JSONB persists ONLY { id }; the data-service materializes
// client-sent text into the node before persisting.
//
// Content is transient: it is accepted on writes (the text the client wants
// materialized) and filled on single-record reads (data-service fetches the
// node body), but it is never persisted on the record — omitempty keeps the
// stored shape minimal.
type KnowledgeNodeRef struct {
	Id      string `json:"id"`
	Content string `json:"content,omitempty"`
}

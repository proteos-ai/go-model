package metamodel

// KnowledgeTextAttributeMeta holds the metadata for an attribute of type
// `knowledge-text`.
//
// A knowledge-text attribute stores long-form text (usually markdown) whose
// body lives in a knowledge node owned by the knowledge-service — the record
// itself persists only the composite common.KnowledgeNodeRef `{ id }`. Clients
// write the text (bare string or `{ content }`); the data-service materializes
// it into the node before persisting and fills `{ id, content }` on
// single-record reads. The value is not filterable (the body is not in the
// record JSONB).
//
// The attribute carries no required configuration today; this struct exists
// for parity with the other AttributeMeta variants and reserves room for
// future options (e.g. node labels, status overrides).
type KnowledgeTextAttributeMeta struct {
	Description string `json:"description,omitempty"`
}

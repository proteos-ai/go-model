package knowledgemodel

import (
	"time"

	"go.proteos.ai/model/common"
)

// KnowledgeNodeMetadata is everything about a node except its (potentially huge)
// `Content` body and the derived `ContentEmbedding` vector. It's what reads
// return by default — list and `GET /nodes/:id` hand back metadata so a caller
// browsing the graph never pays to move bodies it didn't ask for (the body is
// fetched explicitly via the content sub-resource). Reads never select the
// vector column either, to avoid detoasting it on every row.
type KnowledgeNodeMetadata struct {
	Id     string `json:"id" sortable:""`
	OrgId  string `json:"org_id"`
	Title  string `json:"title" sortable:""`
	Type   string `json:"type" sortable:""`   // markdown | file | url
	Status string `json:"status" sortable:""` // draft | published | archived
	// EmbeddingModel / EmbeddedAt are the embedding's provenance: which model
	// produced ContentEmbedding and when. They make re-embedding deterministic
	// and staleness detectable (null model / EmbeddedAt < UpdatedAt ⇒ re-embed).
	EmbeddingModel *string    `json:"embedding_model,omitempty"`
	EmbeddedAt     *time.Time `json:"embedded_at,omitempty"`
	FileId         *string    `json:"file_id,omitempty"` // source pointer (type=file → storage-service)
	Url            *string    `json:"url,omitempty"`     // source pointer (type=url)
	Summary        *string    `json:"summary,omitempty"`
	// ValidFrom / ValidUntil are the node's temporal validity window — when the
	// fact it asserts starts and stops being true, decoupled from CreatedAt (a
	// node ingested today can be valid_from years ago). Both nullable and
	// unbounded when null: null ValidFrom = true since forever, null ValidUntil =
	// still true. Validity is half-open [ValidFrom, ValidUntil).
	ValidFrom  *time.Time     `json:"valid_from,omitempty" sortable:""`
	ValidUntil *time.Time     `json:"valid_until,omitempty" sortable:""`
	CreatedAt  time.Time      `json:"created_at" sortable:""`
	CreatedBy  common.UserRef `json:"created_by"`
	UpdatedAt  time.Time      `json:"updated_at" sortable:""`
	UpdatedBy  common.UserRef `json:"updated_by"`
}

// KnowledgeNode is the central node of the knowledge graph: its metadata plus
// the `Content` body. Every node carries a markdown `Content` that is always the
// searchable representation — for a `file` node it's the extracted markdown, for
// a `url` node the fetched markdown, for a `markdown` node the authored body.
// Nodes are id-keyed (surrogate id); links reference that stable id.
type KnowledgeNode struct {
	KnowledgeNodeMetadata
	Content string `json:"content"` // markdown representation — always the searchable body
	// ContentEmbedding is the pgvector embedding of Content — the vector twin of
	// the content FTS index, backing semantic/hybrid search_nodes. Derived data:
	// populated (and refreshed on content change) by the ingestion path (deferred
	// follow-up). Never serialized and never selected on normal reads.
	ContentEmbedding []float32 `json:"-"`
}

const (
	NodeTypeMarkdown = "markdown"
	NodeTypeFile     = "file"
	NodeTypeUrl      = "url"
)

// NodeTypes is the closed set of KnowledgeNode.Type values.
var NodeTypes = []string{NodeTypeMarkdown, NodeTypeFile, NodeTypeUrl}

const (
	NodeStatusDraft     = "draft"
	NodeStatusPublished = "published"
	NodeStatusArchived  = "archived"
)

// NodeStatuses is the closed set of KnowledgeNode.Status values.
var NodeStatuses = []string{NodeStatusDraft, NodeStatusPublished, NodeStatusArchived}

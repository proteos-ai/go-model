package knowledgemodel

import (
	"time"

	"go.proteos.ai/model/common"
)

// KnowledgeLink is a typed, directed edge between two KnowledgeNodes. Links are
// explicit and are the source of truth for the graph (no wikilink parsing in
// v1) — edges are created directly via the API.
type KnowledgeLink struct {
	Id          string         `json:"id" sortable:""`
	OrgId       string         `json:"org_id"`
	FromId      string         `json:"from_id" sortable:""` // references KnowledgeNode.id
	ToId        string         `json:"to_id" sortable:""`   // references KnowledgeNode.id
	Type        string         `json:"type" sortable:""`    // see LinkTypes
	Description *string        `json:"description,omitempty"`
	CreatedAt   time.Time      `json:"created_at" sortable:""`
	CreatedBy   common.UserRef `json:"created_by"`
	UpdatedAt   time.Time      `json:"updated_at" sortable:""`
	UpdatedBy   common.UserRef `json:"updated_by"`
}

const (
	LinkTypeReferences  = "references"
	LinkTypeRelatesTo   = "relates_to"
	LinkTypeDependsOn   = "depends_on"
	LinkTypePartOf      = "part_of"
	LinkTypeDerivedFrom = "derived_from"
	LinkTypeContradicts = "contradicts"
	LinkTypeSupports    = "supports"
	LinkTypeDuplicates  = "duplicates"
	// LinkTypeSupersededBy points from an outdated node to the node that replaces
	// it (from = superseded, to = successor). Pairs with the temporal validity
	// window: when B supersedes A you both create A --superseded_by--> B and end
	// A's validity (set A.valid_until to B.valid_from).
	LinkTypeSupersededBy = "superseded_by"
)

// LinkTypes is the closed set of KnowledgeLink.Type values.
var LinkTypes = []string{
	LinkTypeReferences,
	LinkTypeRelatesTo,
	LinkTypeDependsOn,
	LinkTypePartOf,
	LinkTypeDerivedFrom,
	LinkTypeContradicts,
	LinkTypeSupports,
	LinkTypeDuplicates,
	LinkTypeSupersededBy,
}

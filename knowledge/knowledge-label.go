package knowledgemodel

import (
	"time"

	"go.proteos.ai/model/common"
)

// KnowledgeLabel is a flexible categorization handle. Organization comes from
// labels + `part_of` graph edges (there are no folders). A label's slug is its
// handle, unique within an org.
type KnowledgeLabel struct {
	Id          string         `json:"id" sortable:""`
	OrgId       string         `json:"org_id"`
	Name        string         `json:"name" sortable:""`
	Slug        string         `json:"slug" sortable:""` // unique (org_id, slug)
	Description *string        `json:"description,omitempty"`
	Color       *string        `json:"color,omitempty"`
	Icon        *string        `json:"icon,omitempty"`
	CreatedAt   time.Time      `json:"created_at" sortable:""`
	CreatedBy   common.UserRef `json:"created_by"`
	UpdatedAt   time.Time      `json:"updated_at" sortable:""`
	UpdatedBy   common.UserRef `json:"updated_by"`
}

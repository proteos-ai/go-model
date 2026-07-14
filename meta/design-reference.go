package metamodel

import (
	"go.proteos.ai/model/common"
	"time"
)

// DesignReference is a stored DESIGN.md document — a named design reference
// (design-system spec + prose rules) an org authors and that design agents read
// as the source of truth for a surface. Name + Description are the selector
// ("which reference, and when to use it"); Content is the full markdown body.
//
// Content is intentionally NOT returned by list/get (the repository never selects
// the column there, and json:omitempty drops the empty value) — it is read and
// written through the dedicated /design-references/:id/content endpoint. This
// keeps the metadata endpoints cheap: list many references by their selector,
// then pull the full body only for the one you picked. Mirrors knowledge-service
// get_node_meta (metadata) vs read_node_content (body).
type DesignReference struct {
	Id    string `json:"id" sortable:""`
	OrgId string `json:"org_id" sortable:""`
	// Slug is the human/CLI/agent-facing reference, unique per org.
	Slug string `json:"slug" sortable:""`
	// Name is the display name of the design reference (e.g. "Proteos · Indigo").
	Name string `json:"name" sortable:""`
	// Description is the short "when to use this" hint an agent reads to pick a
	// reference for a surface.
	Description string `json:"description" sortable:""`
	// Content is the full DESIGN.md markdown body, stored opaque (not parsed).
	// Omitted from list/get; served via GET /design-references/:id/content.
	Content   string         `json:"content,omitempty"`
	CreatedAt time.Time      `json:"created_at" sortable:""`
	CreatedBy common.UserRef `json:"created_by" sortable:""`
	UpdatedAt time.Time      `json:"updated_at" sortable:""`
	UpdatedBy common.UserRef `json:"updated_by" sortable:""`
}

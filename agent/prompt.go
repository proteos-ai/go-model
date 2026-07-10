package agentmodel

import (
	"time"

	"go.proteos.ai/model/common"
	metamodel "go.proteos.ai/model/meta"
)

// Prompt is a reusable, Liquid-templated instruction text. It is versioned like a
// storage File: this parent row holds identity + mutable metadata + a pointer to
// the current immutable PromptVersion, while the content (body + inputs) lives on
// the version rows. Keyed by (org_id, key); key is immutable.
type Prompt struct {
	OrgId          string         `json:"org_id"`
	Key            string         `json:"key" sortable:""`
	Name           string         `json:"name" sortable:""`
	ModuleSlug     string         `json:"module_slug" sortable:""`
	Description    string         `json:"description"`
	CurrentVersion *PromptVersion `json:"current_version"`
	CreatedAt      time.Time      `json:"created_at" sortable:""`
	CreatedBy      common.UserRef `json:"created_by"`
	UpdatedAt      time.Time      `json:"updated_at" sortable:""`
	UpdatedBy      common.UserRef `json:"updated_by"`
}

// PromptVersion is an immutable snapshot of a Prompt's content. Editing body or
// inputs forks a new version (Number increments); metadata-only edits do not.
// Inputs declares the Liquid placeholders the template expects (same shape as an
// Action's params) so the host knows what context to supply at render.
type PromptVersion struct {
	Id        string                `json:"id"`
	Number    uint32                `json:"number"`
	OrgId     string                `json:"org_id"`
	PromptKey string                `json:"prompt_key"`
	Body      string                `json:"body"`
	Inputs    []metamodel.Attribute `json:"inputs"`
	Hash      string                `json:"hash"`
	CreatedAt time.Time             `json:"created_at"`
	CreatedBy common.UserRef        `json:"created_by"`
}

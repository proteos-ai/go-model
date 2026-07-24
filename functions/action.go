package functionsmodel

import (
	"go.proteos.ai/model/common"
	"time"

	metamodel "go.proteos.ai/model/meta"
)

// Action is a deployable user-invokable operation. Scope discriminates
// between record-scoped (invoked against a specific record via
// `POST /api/v1/entities/:entity/records/:recordId/actions/:slug/invoke`)
// and global (invoked via `POST /api/v1/actions/:slug/invoke`).
// EntitySlug is non-nil only when Scope == ActionScopeEntity.
//
// ParamsSchema and ReturnsSchema use the same `meta.Attribute` schema
// language as entity attributes — the UI form renderer consumes them
// directly.
//
// The author's `main.go` is, by convention, always at `./main.go`
// next to the `action.json` manifest; no `source` field is sent over
// the wire or persisted on the row.
type Action struct {
	Slug       string      `json:"slug" sortable:""`
	OrgId      string      `json:"org_id" sortable:""`
	ModuleSlug string      `json:"module_slug" sortable:""`
	Scope      ActionScope `json:"scope" sortable:""`
	EntitySlug *string     `json:"entity,omitempty" sortable:""`
	Name       string      `json:"name" sortable:""`
	IsActive   bool        `json:"is_active" sortable:""`
	// IsPublic exposes a global action on the unauthenticated public
	// dispatch endpoint (`POST /functions/v1/public/orgs/:orgId/actions/:slug/invoke`).
	// Opt-in; defaults false. Only honoured for ActionScopeGlobal — the
	// public dispatch query also gates on `scope = 'global'`.
	IsPublic      bool                  `json:"is_public" sortable:""`
	FileId string `json:"file_id" sortable:""`
	// Checksum is the content address of the deployed wasm blob in canonical
	// "<algo>:<hex>" form (see common.FormatChecksum), stamped at deploy time.
	// Lets the CLI's plan/deploy diff skip re-uploading an unchanged artifact.
	// Empty on rows deployed before checksum stamping.
	Checksum     string                `json:"checksum,omitempty"`
	ParamsSchema []metamodel.Attribute `json:"params"`
	ReturnsSchema []metamodel.Attribute `json:"returns"`
	CreatedAt     time.Time             `json:"created_at" sortable:""`
	CreatedBy     common.UserRef        `json:"created_by"`
	UpdatedAt     time.Time             `json:"updated_at" sortable:""`
	UpdatedBy     common.UserRef        `json:"updated_by"`
}

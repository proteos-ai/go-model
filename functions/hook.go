package functionsmodel

import (
	"go.proteos.ai/model/common"
	"time"
)

// Hook is a deployed record-lifecycle interceptor. One row per
// (orgId, slug). FileId points at the wasm blob in storage-service —
// the blob's content hash + size live on storage-service's `version`
// row, not on the hook itself. Versioning is deferred to LUM-63;
// deploys upsert the row and update FileId in place.
//
// The author's `main.go` is, by convention, always at `./main.go`
// next to the `hook.json` manifest; no `source` field is sent over the
// wire or persisted on the row.
type Hook struct {
	Slug       string         `json:"slug" sortable:""`
	OrgId      string         `json:"org_id" sortable:""`
	ModuleSlug string         `json:"module_slug" sortable:""`
	EntitySlug string         `json:"entity" sortable:""`
	Event      HookEvent      `json:"event" sortable:""`
	IsActive   bool           `json:"is_active" sortable:""`
	FileId     string         `json:"file_id" sortable:""`
	// Checksum is the content address of the deployed wasm blob in canonical
	// "<algo>:<hex>" form (see common.FormatChecksum), stamped at deploy time.
	// Lets the CLI's plan/deploy diff skip re-uploading an unchanged artifact.
	// Empty on rows deployed before checksum stamping.
	Checksum  string    `json:"checksum,omitempty"`
	CreatedAt time.Time `json:"created_at" sortable:""`
	CreatedBy  common.UserRef `json:"created_by"`
	UpdatedAt  time.Time      `json:"updated_at" sortable:""`
	UpdatedBy  common.UserRef `json:"updated_by"`
}

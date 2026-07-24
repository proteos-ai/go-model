package metamodel

import (
	"go.proteos.ai/model/common"
	"time"
)

type Component struct {
	Slug        string `json:"slug" sortable:""`
	OrgId       string `json:"org_id" sortable:""`
	Name        string `json:"name" sortable:""`
	ModuleSlug  string `json:"module_slug" sortable:""`
	Description string `json:"description" sortable:""`
	// BundleFileId points at the compiled, single-file ESM bundle in
	// storage-service (logical reference; storage lives in another service's DB,
	// so no enforced FK). Served back via GET /api/v1/components/:slug/bundle.
	BundleFileId string `json:"bundle_file_id" sortable:""`
	// SourceFileId points at a tar.gz of the component's source directory in
	// storage-service, kept for provenance / rebuild. Not served to the runtime.
	SourceFileId string `json:"source_file_id" sortable:""`
	// BundleChecksum is the content address of the compiled ESM bundle in
	// canonical "<algo>:<hex>" form (see common.FormatChecksum), stamped
	// server-side at deploy time (never client-supplied). Lets the CLI's
	// plan/deploy diff skip re-uploading an unchanged bundle. Empty on rows
	// deployed before checksum stamping.
	BundleChecksum string `json:"bundle_checksum,omitempty"`
	// PropsSchema is the component's JSON Schema (jsonb). Drives the
	// page-designer's props editor and runtime prop validation.
	PropsSchema map[string]any `json:"props_schema"`
	// IsPublic opts the component into UNAUTHENTICATED serving: its compiled
	// bundle becomes world-downloadable via
	// `GET /meta/v1/public/orgs/{orgId}/components/{slug}/bundle`, and public
	// (type='public') pages may only reference public components — enforced at
	// page save. It is also the author's declaration that the component is
	// written for the restricted public sdk (only
	// `functions.actions.invokePublic` reaches the platform on a public page).
	// Defaults to false.
	IsPublic  bool           `json:"is_public"`
	CreatedAt time.Time      `json:"created_at" sortable:""`
	CreatedBy common.UserRef `json:"created_by" sortable:""`
	UpdatedAt time.Time      `json:"updated_at" sortable:""`
	UpdatedBy common.UserRef `json:"updated_by" sortable:""`
}

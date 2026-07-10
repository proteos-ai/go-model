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
	// PropsSchema is the component's JSON Schema (jsonb). Drives the
	// page-designer's props editor and runtime prop validation.
	PropsSchema map[string]any `json:"props_schema"`
	CreatedAt   time.Time      `json:"created_at" sortable:""`
	CreatedBy   common.UserRef `json:"created_by" sortable:""`
	UpdatedAt   time.Time      `json:"updated_at" sortable:""`
	UpdatedBy   common.UserRef `json:"updated_by" sortable:""`
}

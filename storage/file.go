package storagemodel

import (
	"time"

	"go.proteos.ai/model/common"
)

type File struct {
	Id             string   `json:"id" example:"0027.8f5b88e0-b875-4aa7-a73e-08995671047e"`
	OrgId          *string  `json:"org_id,omitempty" sortable:""`
	Name           string   `json:"name" example:"pdr-max-mustermann.pdf"`
	ContentType    string   `json:"content_type" example:"application/pdf"`
	CurrentVersion *Version `json:"current_version"`
	IsDeleted      bool     `json:"is_deleted" example:"false"`
	IsPersisted    bool     `json:"is_persisted" example:"false"`
	IsLocked       bool     `json:"is_locked" example:"false"`
	// PublicAccess is the set of operations this file is exposed for on the
	// unauthenticated public surface. Only ["read"] is honored (download via
	// GET /storage/v1/public/orgs/{orgId}/files/{id}/download); a file has no
	// public write/delete. Org-transcending (NULL-org) files are never
	// publicly served regardless. Empty = private (default).
	PublicAccess common.PublicAccess `json:"public_access"`
	CreatedAt    time.Time           `json:"created_at"`
	UpdatedAt    time.Time           `json:"updated_at"`
}

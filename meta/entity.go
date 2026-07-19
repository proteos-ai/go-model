package metamodel

import (
	"go.proteos.ai/model/common"
	"time"
)

// Entity represents a metadata entity definition with its attributes
type Entity struct {
	Slug     string `json:"slug" sortable:""`
	OrgId    string `json:"org_id" sortable:""`
	Name     string `json:"name" sortable:""`
	IsRemote bool   `json:"is_remote" sortable:""`
	// PublicRecordAccess is the set of operations ALL records of this entity
	// are exposed for on the UNAUTHENTICATED public surface. Today only
	// ["read"] is honored: records become world-readable via
	// GET /data/v1/public/orgs/:orgId/records/:slug… (the entity definition is
	// implicitly readable too, via GET /meta/v1/public/orgs/:orgId/entities/:slug,
	// so the records can be interpreted). `write` and `delete` are reserved
	// (rejected until their backends land). Empty = fully private (default).
	PublicRecordAccess common.PublicAccess `json:"public_record_access"`
	ModuleSlug         string              `json:"module_slug" sortable:""`
	Description        string              `json:"description" sortable:""`
	TitleTemplate      string              `json:"title_template" sortable:""`
	Attributes         []Attribute         `json:"attributes"`
	CreatedAt          time.Time           `json:"created_at" sortable:""`
	CreatedBy          common.UserRef      `json:"created_by" sortable:""`
	UpdatedAt          time.Time           `json:"updated_at" sortable:""`
	UpdatedBy          common.UserRef      `json:"updated_by" sortable:""`
}

type EntityWithSchema struct {
	Entity
	Schema any `json:"schema,omitempty"`
}

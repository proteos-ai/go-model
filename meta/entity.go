package metamodel

import (
	"go.proteos.ai/model/common"
	"time"
)

// Entity represents a metadata entity definition with its attributes
type Entity struct {
	Slug          string         `json:"slug" sortable:""`
	OrgId         string         `json:"org_id" sortable:""`
	Name          string         `json:"name" sortable:""`
	IsRemote      bool           `json:"is_remote" sortable:""`
	ModuleSlug    string         `json:"module_slug" sortable:""`
	Description   string         `json:"description" sortable:""`
	TitleTemplate string         `json:"title_template" sortable:""`
	Attributes    []Attribute    `json:"attributes"`
	CreatedAt     time.Time      `json:"created_at" sortable:""`
	CreatedBy     common.UserRef `json:"created_by" sortable:""`
	UpdatedAt     time.Time      `json:"updated_at" sortable:""`
	UpdatedBy     common.UserRef `json:"updated_by" sortable:""`
}

type EntityWithSchema struct {
	Entity
	Schema any `json:"schema,omitempty"`
}

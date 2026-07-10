package metamodel

import (
	"go.proteos.ai/model/common"
	"time"
)

type App struct {
	Slug        string         `json:"slug" sortable:""`
	OrgId       string         `json:"org_id" sortable:""`
	Name        string         `json:"name" sortable:""`
	ModuleSlug  string         `json:"module_slug" sortable:""`
	Description string         `json:"description" sortable:""`
	IconSlug    string         `json:"icon_slug" sortable:""`
	CreatedAt   time.Time      `json:"created_at" sortable:""`
	CreatedBy   common.UserRef `json:"created_by" sortable:""`
	UpdatedAt   time.Time      `json:"updated_at" sortable:""`
	UpdatedBy   common.UserRef `json:"updated_by" sortable:""`
}

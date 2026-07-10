package metamodel

import (
	"go.proteos.ai/model/common"
	"time"
)

type Module struct {
	Slug          string         `json:"slug" sortable:""`
	OrgId         string         `json:"org_id" sortable:""`
	Name          string         `json:"name" sortable:""`
	Description   string         `json:"description" sortable:""`
	FileId        string         `json:"file_id" sortable:""`
	Version       string         `json:"version" sortable:""`
	IsDeactivated bool           `json:"is_deactivated" sortable:""`
	Status        string         `json:"status" sortable:""`
	StatusDetails string         `json:"status_details" sortable:""`
	CreatedAt     time.Time      `json:"created_at" sortable:""`
	CreatedBy     common.UserRef `json:"created_by" sortable:""`
	UpdatedAt     time.Time      `json:"updated_at" sortable:""`
	UpdatedBy     common.UserRef `json:"updated_by" sortable:""`
}

package accountmodel

import (
	"go.proteos.ai/model/common"
	"time"
)

type Role struct {
	Slug        string         `json:"slug" sortable:""`
	Name        string         `json:"name" sortable:""`
	OrgId       string         `json:"org_id" sortable:""`
	Description string         `json:"description" sortable:""`
	CreatedAt   time.Time      `json:"created_at" sortable:""`
	CreatedBy   common.UserRef `json:"created_by" sortable:""`
	UpdatedAt   time.Time      `json:"updated_at" sortable:""`
	UpdatedBy   common.UserRef `json:"updated_by" sortable:""`
}

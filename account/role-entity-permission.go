package accountmodel

import (
	"go.proteos.ai/model/common"
	"time"
)

type RoleEntityPermission struct {
	Id         string         `json:"id" sortable:""`
	OrgId      string         `json:"org_id" sortable:""`
	RoleSlug   string         `json:"role_slug" sortable:""`
	EntitySlug string         `json:"entity_slug" sortable:""`
	Permission Permission     `json:"permission" sortable:""`
	CreatedAt  time.Time      `json:"created_at" sortable:""`
	CreatedBy  common.UserRef `json:"created_by" sortable:""`
	UpdatedAt  time.Time      `json:"updated_at" sortable:""`
	UpdatedBy  common.UserRef `json:"updated_by" sortable:""`
}

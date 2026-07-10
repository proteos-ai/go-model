package accountmodel

import (
	"go.proteos.ai/model/common"
	"time"
)

type UserRoleAssignment struct {
	Id        string         `json:"id" sortable:""`
	UserId    string         `json:"user_id" sortable:""`
	RoleSlug  string         `json:"role_slug" sortable:""`
	OrgId     string         `json:"org_id" sortable:""`
	CreatedAt time.Time      `json:"created_at" sortable:""`
	CreatedBy common.UserRef `json:"created_by" sortable:""`
	UpdatedAt time.Time      `json:"updated_at" sortable:""`
	UpdatedBy common.UserRef `json:"updated_by" sortable:""`
}

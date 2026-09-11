package accountmodel

import (
	"time"

	"go.proteos.ai/model/common"
)

// UserProfileAssignment binds a user to their ONE profile in an org (unique on
// (org_id, user_id)). Setting a profile replaces the previous one; clearing it
// leaves the user on the app's base configuration.
type UserProfileAssignment struct {
	Id          string         `json:"id" sortable:""`
	UserId      string         `json:"user_id" sortable:""`
	ProfileSlug string         `json:"profile_slug" sortable:""`
	OrgId       string         `json:"org_id" sortable:""`
	CreatedAt   time.Time      `json:"created_at" sortable:""`
	CreatedBy   common.UserRef `json:"created_by" sortable:""`
	UpdatedAt   time.Time      `json:"updated_at" sortable:""`
	UpdatedBy   common.UserRef `json:"updated_by" sortable:""`
}

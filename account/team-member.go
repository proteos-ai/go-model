package accountmodel

import (
	"time"

	"go.proteos.ai/model/common"
)

// TeamMember is one user's direct membership of one team.
//
// Only DIRECT memberships are stored. Membership of ancestor teams is derived —
// it comes from the FGA closure rather than being materialised here — so moving
// a team under a new parent changes every descendant's effective membership
// with zero rows rewritten.
type TeamMember struct {
	Id        string         `json:"id" sortable:""`
	OrgId     string         `json:"org_id" sortable:""`
	TeamSlug  string         `json:"team_slug" sortable:""`
	UserId    string         `json:"user_id" sortable:""`
	CreatedAt time.Time      `json:"created_at" sortable:""`
	CreatedBy common.UserRef `json:"created_by" sortable:""`
	UpdatedAt time.Time      `json:"updated_at" sortable:""`
	UpdatedBy common.UserRef `json:"updated_by" sortable:""`
}

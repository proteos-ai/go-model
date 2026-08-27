package accountmodel

import (
	"time"

	"go.proteos.ai/model/common"
)

// Team is a node in the org's structure, and a principal that can hold access
// anywhere a user can — on a record, a conversation, a knowledge space.
//
// Hierarchy rolls UP: a member of a child team is a member of every ancestor
// FOR GRANTS MADE TO THAT ANCESTOR. Granting `sales` reaches the members of
// `sales-emea`; granting `sales-emea` never reaches someone who is only in
// `sales`. This is the direction people get backwards.
//
// The slug is immutable — it is half the primary key, and grants reference
// `team:<orgId>/<slug>`, so a rename would orphan every grant naming it.
type Team struct {
	Slug        string `json:"slug" sortable:""`
	Name        string `json:"name" sortable:""`
	OrgId       string `json:"org_id" sortable:""`
	Description string `json:"description" sortable:""`
	// ParentTeamSlug is this team's parent within the same org; empty for a root.
	ParentTeamSlug string         `json:"parent_team_slug" sortable:""`
	CreatedAt      time.Time      `json:"created_at" sortable:""`
	CreatedBy      common.UserRef `json:"created_by" sortable:""`
	UpdatedAt      time.Time      `json:"updated_at" sortable:""`
	UpdatedBy      common.UserRef `json:"updated_by" sortable:""`
}

// PrincipalRef renders the team as the principal it is granted as.
func (team Team) PrincipalRef() common.PrincipalRef {
	return common.PrincipalRef{Type: common.PrincipalTypeTeam, Id: team.Slug}
}

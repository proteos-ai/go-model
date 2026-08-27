package accountapi

import (
	"go.proteos.ai/model/account"
	"go.proteos.ai/model/common"
)

type CreateTeamRequest struct {
	OrgId string `json:"org_id" form:"org_id" validate:"required"`
	Slug  string `json:"slug" form:"slug" validate:"required"`
	Name  string `json:"name" form:"name" validate:"required"`
	// ParentTeamSlug nests this team under an existing one in the same org.
	// Empty makes it a root team.
	ParentTeamSlug string `json:"parent_team_slug,omitempty" form:"parent_team_slug,omitempty"`
	Description    string `json:"description" form:"description"`
}

// UpdateTeamRequest omits Slug deliberately: it is half the primary key and is
// referenced by every grant as `team:<orgId>/<slug>`, so renaming it would
// orphan those grants. Re-parenting is allowed and is re-validated against the
// cycle rules.
type UpdateTeamRequest struct {
	Name        *string `json:"name,omitempty" form:"name,omitempty"`
	Description *string `json:"description,omitempty" form:"description,omitempty"`
	// ParentTeamSlug re-parents the team; an explicit empty string promotes it
	// to a root team.
	ParentTeamSlug *string `json:"parent_team_slug,omitempty" form:"parent_team_slug,omitempty"`
}

type GetManyTeamsQuery struct {
	OrgId          string `json:"org_id" db:"org_id"`
	Slug           string `json:"slug" db:"slug"`
	Name           string `json:"name" db:"name"`
	Description    string `json:"description" db:"description"`
	ParentTeamSlug string `json:"parent_team_slug" db:"parent_team_slug"`
	CreatedBy      string `json:"created_by" db:"created_by->>'id'"`
	UpdatedBy      string `json:"updated_by" db:"updated_by->>'id'"`
	common.Pagination
	common.Sorting
}

type GetManyTeamsResponse struct {
	Meta common.ResponseMeta `json:"meta"`
	Data []accountmodel.Team `json:"data"`
}

type CreateTeamMemberRequest struct {
	UserId string `json:"user_id" form:"user_id" validate:"required"`
}

type GetManyTeamMembersQuery struct {
	OrgId    string `json:"org_id" db:"org_id"`
	TeamSlug string `json:"team_slug" db:"team_slug"`
	UserId   string `json:"user_id" db:"user_id"`
	common.Pagination
	common.Sorting
}

type GetManyTeamMembersResponse struct {
	Meta common.ResponseMeta       `json:"meta"`
	Data []accountmodel.TeamMember `json:"data"`
}

package accountapi

import (
	"go.proteos.ai/model/account"
	"go.proteos.ai/model/common"
)

type CreateUserRoleAssignmentRequest struct {
	UserId   string `json:"user_id" form:"user_id" validate:"required"`
	RoleSlug string `json:"role_slug" form:"role_slug" validate:"required"`
	// OrgId optionally targets the org the assignment lands in. Honored only for
	// privileged callers (platform admin / system); a regular caller's value is
	// ignored and pinned to their token org. nil/empty = use the token org.
	OrgId *string `json:"org_id,omitempty" form:"org_id,omitempty"`
}

type GetManyUserRoleAssignmentsQuery struct {
	Id       *string `json:"id" db:"id"`
	UserId   *string `json:"user_id" db:"user_id"`
	RoleSlug *string `json:"role_slug" db:"role_slug"`
	common.Pagination
	common.Sorting
}

type GetManyUserRoleAssignmentsResponse struct {
	Meta common.ResponseMeta               `json:"meta"`
	Data []accountmodel.UserRoleAssignment `json:"data"`
}

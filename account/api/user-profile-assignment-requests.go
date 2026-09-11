package accountapi

import (
	"go.proteos.ai/model/account"
	"go.proteos.ai/model/common"
)

// SetUserProfileRequest is the body of PUT /users/:id/profile — it REPLACES the
// user's profile in the token org (one profile per user per org).
type SetUserProfileRequest struct {
	UserId      string `json:"user_id" form:"user_id" validate:"required"`
	ProfileSlug string `json:"profile_slug" form:"profile_slug" validate:"required"`
}

type GetManyUserProfileAssignmentsQuery struct {
	Id          *string `json:"id" db:"id"`
	UserId      *string `json:"user_id" db:"user_id"`
	ProfileSlug *string `json:"profile_slug" db:"profile_slug"`
	common.Pagination
	common.Sorting
}

type GetManyUserProfileAssignmentsResponse struct {
	Meta common.ResponseMeta                  `json:"meta"`
	Data []accountmodel.UserProfileAssignment `json:"data"`
}

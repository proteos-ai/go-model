package accountapi

import (
	"go.proteos.ai/model/account"
	"go.proteos.ai/model/common"
)

type CreateUserRequest struct {
	GivenName    string  `json:"given_name" form:"given_name" validate:"required"`
	FamilyName   string  `json:"family_name" form:"family_name" validate:"required"`
	Email        string  `json:"email" form:"email" validate:"required"`
	DefaultOrgId *string `json:"default_org_id,omitempty" form:"default_org_id,omitempty"`
}

type UpdateUserRequest struct {
	GivenName    *string `json:"given_name,omitempty" form:"given_name,omitempty"`
	FamilyName   *string `json:"family_name,omitempty" form:"family_name,omitempty"`
	DefaultOrgId *string `json:"default_org_id,omitempty" form:"default_org_id,omitempty"`
}

type GetManyUsersQuery struct {
	Id           string `json:"id" db:"id"`
	GivenName    string `json:"given_name" db:"given_name"`
	FamilyName   string `json:"family_name" db:"family_name"`
	Email        string `json:"email" db:"email"`
	DefaultOrgId string `json:"default_org_id" db:"default_org_id"`
	CreatedBy    string `json:"created_by" db:"created_by->>'id'"`
	UpdatedBy    string `json:"updated_by" db:"updated_by->>'id'"`
	common.Pagination
	common.Sorting
}

type GetManyUsersResponse struct {
	Meta common.ResponseMeta `json:"meta"`
	Data []accountmodel.User `json:"data"`
}

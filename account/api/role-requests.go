package accountapi

import (
	"go.proteos.ai/model/account"
	"go.proteos.ai/model/common"
)

type CreateRoleRequest struct {
	OrgId       string `json:"org_id" form:"org_id" validate:"required"`
	Slug        string `json:"slug" form:"slug" validate:"required"`
	Name        string `json:"name" form:"name" validate:"required"`
	Description string `json:"description" form:"description"`
}

type UpdateRoleRequest struct {
	OrgId       *string `json:"org_id,omitempty" form:"org_id,omitempty"`
	Name        *string `json:"name,omitempty" form:"name,omitempty"`
	Description *string `json:"description,omitempty" form:"description,omitempty"`
}

type GetManyRolesQuery struct {
	OrgId       string `json:"org_id" db:"org_id"`
	Slug        string `json:"slug" db:"slug"`
	Name        string `json:"name" db:"name"`
	Description string `json:"description" db:"description"`
	CreatedBy   string `json:"created_by" db:"created_by->>'id'"`
	UpdatedBy   string `json:"updated_by" db:"updated_by->>'id'"`
	common.Pagination
	common.Sorting
}

type GetManyRolesResponse struct {
	Meta common.ResponseMeta `json:"meta"`
	Data []accountmodel.Role `json:"data"`
}

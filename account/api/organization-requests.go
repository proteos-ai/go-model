package accountapi

import (
	"go.proteos.ai/model/account"
	"go.proteos.ai/model/common"
)

type CreateOrganizationRequest struct {
	Name        string `json:"name" form:"name" validate:"required"`
	Description string `json:"description" form:"description"`
}

type UpdateOrganizationRequest struct {
	Name        *string `json:"name,omitempty" form:"name,omitempty"`
	Description *string `json:"description,omitempty" form:"description,omitempty"`
}

type GetManyOrganizationsQuery struct {
	Id          string `json:"id" db:"id"`
	Name        string `json:"name" db:"name"`
	Description string `json:"description" db:"description"`
	CreatedBy   string `json:"created_by" db:"created_by->>'id'"`
	UpdatedBy   string `json:"updated_by" db:"updated_by->>'id'"`
	common.Pagination
	common.Sorting
}

type GetManyOrganizationsResponse struct {
	Meta common.ResponseMeta         `json:"meta"`
	Data []accountmodel.Organization `json:"data"`
}

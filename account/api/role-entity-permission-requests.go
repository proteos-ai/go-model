package accountapi

import (
	"go.proteos.ai/model/account"
	"go.proteos.ai/model/common"
)

type CreateRoleEntityPermissionRequest struct {
	RoleSlug   string                  `json:"role_slug" form:"role_slug" validate:"required"`
	EntitySlug string                  `json:"entity_slug" form:"entity_slug" validate:"required"`
	Permission accountmodel.Permission `json:"permission" form:"permission" validate:"required"`
}

type GetManyRoleEntityPermissionsQuery struct {
	Id         *string                  `json:"id" db:"id"`
	RoleSlug   *string                  `json:"role_slug" db:"role_slug"`
	EntitySlug *string                  `json:"entity_slug" db:"entity_slug"`
	Permission *accountmodel.Permission `json:"permission" db:"permission"`
	common.Pagination
	common.Sorting
}

type GetManyRoleEntityPermissionsResponse struct {
	Meta common.ResponseMeta                 `json:"meta"`
	Data []accountmodel.RoleEntityPermission `json:"data"`
}

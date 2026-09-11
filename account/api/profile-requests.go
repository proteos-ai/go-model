package accountapi

import (
	"go.proteos.ai/model/account"
	"go.proteos.ai/model/common"
)

type CreateProfileRequest struct {
	OrgId          string   `json:"org_id" form:"org_id" validate:"required"`
	Slug           string   `json:"slug" form:"slug" validate:"required"`
	Name           string   `json:"name" form:"name" validate:"required"`
	Description    string   `json:"description" form:"description"`
	ModuleSlug     string   `json:"module_slug" form:"module_slug"`
	DefaultAppSlug string   `json:"default_app_slug" form:"default_app_slug"`
	AppSlugs       []string `json:"app_slugs" form:"app_slugs"`
}

type UpdateProfileRequest struct {
	OrgId          *string   `json:"org_id,omitempty" form:"org_id,omitempty"`
	Name           *string   `json:"name,omitempty" form:"name,omitempty"`
	Description    *string   `json:"description,omitempty" form:"description,omitempty"`
	ModuleSlug     *string   `json:"module_slug,omitempty" form:"module_slug,omitempty"`
	DefaultAppSlug *string   `json:"default_app_slug,omitempty" form:"default_app_slug,omitempty"`
	AppSlugs       *[]string `json:"app_slugs,omitempty" form:"app_slugs,omitempty"`
}

type GetManyProfilesQuery struct {
	OrgId          string `json:"org_id" db:"org_id"`
	Slug           string `json:"slug" db:"slug"`
	Name           string `json:"name" db:"name"`
	ModuleSlug     string `json:"module_slug" db:"module_slug"`
	DefaultAppSlug string `json:"default_app_slug" db:"default_app_slug"`
	CreatedBy      string `json:"created_by" db:"created_by->>'id'"`
	UpdatedBy      string `json:"updated_by" db:"updated_by->>'id'"`
	common.Pagination
	common.Sorting
}

type GetManyProfilesResponse struct {
	Meta common.ResponseMeta    `json:"meta"`
	Data []accountmodel.Profile `json:"data"`
}

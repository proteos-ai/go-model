package metaapi

import (
	"go.proteos.ai/model/common"
	"go.proteos.ai/model/meta"
)

type DeployModuleRequest struct {
	Slug        string `json:"slug" validate:"required"`
	Version     string `json:"version" validate:"required"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

type UpdateModuleRequest struct {
	Name        *string `json:"name,omitempty"`
	Description *string `json:"description,omitempty"`
}

type GetManyModulesQuery struct {
	Slug          *string                 `json:"slug" db:"slug"`
	Name          *string                 `json:"name" db:"name"`
	IsDeactivated *bool                   `json:"is_deactivated" db:"is_deactivated"`
	FileId        *string                 `json:"file_id" db:"file_id"`
	Status        *metamodel.ModuleStatus `json:"status" db:"status"`
	Version       *string                 `json:"version" db:"version"`
	CreatedBy     *string                 `json:"created_by" db:"created_by->>'id'"`
	UpdatedBy     *string                 `json:"updated_by" db:"updated_by->>'id'"`
	common.Pagination
	common.Sorting
}

type GetManyModulesResponse struct {
	Meta common.ResponseMeta `json:"meta"`
	Data []metamodel.Module  `json:"data"`
}

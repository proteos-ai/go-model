package metaapi

import (
	"go.proteos.ai/model/common"
	"go.proteos.ai/model/meta"
)

type CreateVariableRequest struct {
	Key      string `json:"key" validate:"required"`
	Value    string `json:"value"`
	IsSecret bool   `json:"is_secret"`
	Module   string `json:"module"`
}

type UpdateVariableRequest struct {
	Value *string `json:"value,omitempty"`
}

type GetManyVariablesQuery struct {
	Id        *string `json:"id" db:"id"`
	Key       *string `json:"key" db:"key"`
	IsSecret  *bool   `json:"is_secret" db:"is_secret"`
	Module    *string `json:"module" db:"module"`
	CreatedBy *string `json:"created_by" db:"created_by->>'id'"`
	UpdatedBy *string `json:"updated_by" db:"updated_by->>'id'"`
	common.Pagination
	common.Sorting
}

type GetManyVariablesResponse struct {
	Meta common.ResponseMeta  `json:"meta"`
	Data []metamodel.Variable `json:"data"`
}

package metamodel

import (
	"go.proteos.ai/model/common"
	"time"
)

type Variable struct {
	Id        string         `json:"id" query:"id"`
	OrgId     string         `json:"org_id" query:"org_id"`
	Module    string         `json:"module" query:"module"`
	Key       string         `json:"key" query:"key"`
	Value     string         `json:"value" query:"value"`
	IsSecret  bool           `json:"is_secret" query:"is_secret"`
	CreatedAt time.Time      `json:"created_at" query:"created_at"`
	UpdatedAt time.Time      `json:"updated_at" query:"updated_at"`
	CreatedBy common.UserRef `json:"created_by" query:"created_by"`
	UpdatedBy common.UserRef `json:"updated_by" query:"updated_by"`
}

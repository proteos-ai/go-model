package metamodel

import (
	"go.proteos.ai/model/common"
	"time"
)

type List struct {
	Slug       string               `json:"slug"`
	OrgId      string               `json:"org_id" sortable:""`
	ModuleSlug string               `json:"module_slug" sortable:""`
	EntitySlug string               `json:"entity_slug"`
	Name       string               `json:"name"`
	Columns    []Column             `json:"columns"`
	Sorting    []SortConfig         `json:"sorting"`
	Filters    []common.FilterGroup `json:"filters"`
	CreatedAt  time.Time            `json:"created_at"`
	UpdatedAt  time.Time            `json:"updated_at"`
	CreatedBy  common.UserRef       `json:"created_by"`
	UpdatedBy  common.UserRef       `json:"updated_by"`
}

type Column struct {
	Attribute string `json:"attribute"`
	Label     string `json:"label"`
	Width     int    `json:"width"`
}

type SortConfig struct {
	Attribute string               `json:"attribute"`
	Direction common.SortDirection `json:"direction"`
}

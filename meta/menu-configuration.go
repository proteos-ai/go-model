package metamodel

import (
	"go.proteos.ai/model/common"
	"time"
)

type MenuItemType string

const (
	MenuItemTypeLink   MenuItemType = "link"
	MenuItemTypeGroup  MenuItemType = "group"
	MenuItemTypeEntity MenuItemType = "entity"
	MenuItemTypePage   MenuItemType = "page"
	MenuItemList       MenuItemType = "list"
)

type Icon string

const (
	IconDashboard     Icon = "dashboard"
	IconPerson        Icon = "person"
	IconPersonCircled Icon = "person-circled"
	IconHome          Icon = "home"
	IconSettings      Icon = "settings"
	IconTrendingUp    Icon = "trending-up"
)

type MenuItem struct {
	Id        string       `json:"id"`
	Order     int          `json:"order"`
	Label     string       `json:"label"`
	Type      MenuItemType `json:"type"`
	Icon      Icon         `json:"icon"`
	Reference string       `json:"reference"`
	Children  []MenuItem   `json:"children"`
}

type MenuConfiguration struct {
	Slug       string         `json:"slug" sortable:""`
	OrgId      string         `json:"org_id" sortable:""`
	ModuleSlug string         `json:"module_slug" sortable:""`
	Name       string         `json:"name" sortable:""`
	AppSlug    string         `json:"app_slug" sortable:""`
	Items      []MenuItem     `json:"items"`
	IsDefault  bool           `json:"is_default" sortable:""`
	CreatedAt  time.Time      `json:"created_at" sortable:""`
	CreatedBy  common.UserRef `json:"created_by" sortable:""`
	UpdatedAt  time.Time      `json:"updated_at" sortable:""`
	UpdatedBy  common.UserRef `json:"updated_by" sortable:""`
}

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
	MenuItemTypeList   MenuItemType = "list"
)

var SupportedMenuItemTypes = []MenuItemType{
	MenuItemTypeLink,
	MenuItemTypeGroup,
	MenuItemTypeEntity,
	MenuItemTypePage,
	MenuItemTypeList,
}

type Icon string

const (
	IconDashboard     Icon = "dashboard"
	IconPerson        Icon = "person"
	IconPersonCircled Icon = "person-circled"
	IconHome          Icon = "home"
	IconSettings      Icon = "settings"
	IconTrendingUp    Icon = "trending-up"
)

// MenuItem is one node of a menu's tree. Id is the item's identity: it is
// required, unique across the whole tree, and it is the menu's extension API —
// menu-configuration extensions anchor on it (see MenuConfigurationExtension).
//
// The stored tree of a host menu is a MATERIALIZATION of the menu's own items
// plus every extension's placements; the three trailing fields are the
// server-derived stamps of that merge and are never authored:
//   - ExtensionKey: the item was contributed by that extension.
//   - IsHidden: a `remove` tombstone — consumers skip the item.
//   - ReplacedItem: the host's own item this stamped item displaced (`replace`
//     / `remove`), kept so stripping the stamps restores the host's own tree.
type MenuItem struct {
	Id           string       `json:"id"`
	Order        int          `json:"order"`
	Label        string       `json:"label"`
	Type         MenuItemType `json:"type"`
	Icon         Icon         `json:"icon"`
	Reference    string       `json:"reference"`
	Children     []MenuItem   `json:"children"`
	ExtensionKey string       `json:"extension_key,omitempty"`
	IsHidden     bool         `json:"is_hidden,omitempty"`
	ReplacedItem *MenuItem    `json:"replaced_item,omitempty"`
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

package metamodel

import (
	"go.proteos.ai/model/common"
	"time"
)

// PageType distinguishes record-scoped pages (rendered against a single record
// of EntitySlug) from platform pages (standalone, no record context — e.g. a
// dashboard launched from a menu item). Empty normalizes to PageTypeRecord for
// backward compatibility.
type PageType string

const (
	PageTypeRecord   PageType = "record"
	PageTypePlatform PageType = "platform"
	// PageTypeExternal is reserved for a future chromeless / no-app-shell page
	// variant. It is a known value but not yet accepted by validation.
	PageTypeExternal PageType = "external"
)

// SupportedPageTypes is the set of types this version accepts. PageTypeExternal
// is intentionally excluded — it is reserved-but-unsupported.
var SupportedPageTypes = []PageType{PageTypeRecord, PageTypePlatform}

// Normalized returns the type with the empty value defaulting to
// PageTypeRecord, so existing record pages (persisted before `type` existed)
// keep their behavior.
func (pageType PageType) Normalized() PageType {
	if pageType == "" {
		return PageTypeRecord
	}
	return pageType
}

type PageAction struct {
	Label  string `json:"label"`
	Icon   Icon   `json:"icon"`
	Action string `json:"action"`
}

type Page struct {
	Slug       string         `json:"slug" sortable:""`
	OrgId      string         `json:"org_id" sortable:""`
	ModuleSlug string         `json:"module_slug" sortable:""`
	Type       PageType       `json:"type" sortable:""`
	EntitySlug string         `json:"entity_slug,omitempty" sortable:""`
	Name       string         `json:"name" sortable:""`
	Actions    []PageAction   `json:"actions"`
	Layout     PageLayout     `json:"layout"`
	CreatedAt  time.Time      `json:"created_at" sortable:""`
	CreatedBy  common.UserRef `json:"created_by" sortable:""`
	UpdatedAt  time.Time      `json:"updated_at" sortable:""`
	UpdatedBy  common.UserRef `json:"updated_by" sortable:""`
}

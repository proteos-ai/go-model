package metamodel

import (
	"go.proteos.ai/model/common"
	"time"
)

// PageType encodes what a page binds to and how it is served — both the
// chrome (app shell vs bare) and the auth posture follow from it:
//
//   - record:   rendered against a single record of EntitySlug; app chrome;
//     authenticated.
//   - platform: standalone, no record context (e.g. a dashboard launched from
//     a menu item); app chrome; authenticated.
//   - kiosk:    standalone, NO app chrome (bare page at /k/…); authenticated.
//   - public:   standalone, NO app chrome (bare page at /p/…);
//     UNAUTHENTICATED — the page layout is world-readable and its components
//     may only call is_public global actions. Never bind org secrets into a
//     public page's layout or props.
//
// Empty normalizes to PageTypeRecord for backward compatibility.
type PageType string

const (
	PageTypeRecord   PageType = "record"
	PageTypePlatform PageType = "platform"
	PageTypeKiosk    PageType = "kiosk"
	PageTypePublic   PageType = "public"
)

// SupportedPageTypes is the set of types this version accepts. (The formerly
// reserved "external" placeholder was retired in favour of kiosk + public.)
var SupportedPageTypes = []PageType{PageTypeRecord, PageTypePlatform, PageTypeKiosk, PageTypePublic}

// IsStandalone reports whether the page renders without a record context —
// everything except record pages. Standalone pages forbid entity_slug and
// record-bound layout elements (field, related_list).
func (pageType PageType) IsStandalone() bool {
	return pageType.Normalized() != PageTypeRecord
}

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

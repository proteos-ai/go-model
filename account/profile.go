package accountmodel

import (
	"time"

	"go.proteos.ai/model/common"
)

// Profile is the org's user-type: what a person sees first. ONE per user per
// org (see UserProfileAssignment), unlike roles which are many and only decide
// access. The split mirrors the Salesforce Profile / Permission-Set lesson: UI
// defaults keyed by a single subject stay unambiguous; permissions compose.
//
// The profile owns the defaults that are one value PER SUBJECT (default app,
// app visibility). Defaults that vary per app × profile (home, menu, agents,
// record pages) live in metadata-service's app_configuration rows, keyed by
// this profile's slug.
//
// The slug is immutable — it is half the primary key, and app configurations
// reference it, so a rename would orphan every override naming it.
type Profile struct {
	Slug        string `json:"slug" sortable:""`
	Name        string `json:"name" sortable:""`
	OrgId       string `json:"org_id" sortable:""`
	Description string `json:"description" sortable:""`
	// ModuleSlug attributes the profile to the module that shipped it (empty
	// for org-authored profiles) — the same attribution apps and menus carry,
	// so `pro module plan/pull` can list, diff and orphan-detect by module.
	ModuleSlug string `json:"module_slug" sortable:""`
	// DefaultAppSlug is the app opened after login. A metadata-service slug,
	// not validated here: clients fail open to the first visible app when it
	// names no app (or one hidden by AppSlugs).
	DefaultAppSlug string `json:"default_app_slug" sortable:""`
	// AppSlugs is the app-switcher allowlist. Empty = every app in the org.
	// Not sortable — JSONB.
	AppSlugs  []string       `json:"app_slugs"`
	CreatedAt time.Time      `json:"created_at" sortable:""`
	CreatedBy common.UserRef `json:"created_by" sortable:""`
	UpdatedAt time.Time      `json:"updated_at" sortable:""`
	UpdatedBy common.UserRef `json:"updated_by" sortable:""`
}

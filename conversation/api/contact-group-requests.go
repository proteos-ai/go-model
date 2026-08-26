package conversationapi

import (
	"go.proteos.ai/model/common"
	conversationmodel "go.proteos.ai/model/conversation"
)

// CreateContactGroupRequest creates one contact group. Key is the immutable
// identity (lowercase kebab, like every platform key); Description tells both
// humans and the tone-synthesis model who belongs in the group. ModuleSlug is
// stamped by `pro module deploy` — interactive callers leave it empty.
// is_auto_created is server-owned and deliberately absent.
type CreateContactGroupRequest struct {
	Key         string `json:"key" validate:"required,max=255"`
	Name        string `json:"name" validate:"required,max=255"`
	Description string `json:"description" validate:"max=4096"`
	ModuleSlug  string `json:"module_slug" validate:"max=255"`
}

// UpdateContactGroupRequest is a partial update — nil leaves the stored value
// untouched. Key is immutable; is_auto_created and module_slug are
// server-owned (attribution moves only via the upsert path).
type UpdateContactGroupRequest struct {
	Name        *string `json:"name,omitempty" validate:"omitempty,max=255"`
	Description *string `json:"description,omitempty" validate:"omitempty,max=4096"`
}

type GetManyContactGroupsQuery struct {
	// ModuleSlug filters to groups deployed by one module — the `pro module`
	// remote-state discovery door.
	ModuleSlug *string `json:"module_slug" form:"module_slug"`
	// Search filters by a case-insensitive substring of key or name.
	Search *string `json:"search" form:"search"`
	common.Pagination
	common.Sorting
}

type GetManyContactGroupsResponse struct {
	Meta common.ResponseMeta              `json:"meta"`
	Data []conversationmodel.ContactGroup `json:"data"`
}

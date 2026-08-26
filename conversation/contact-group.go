package conversationmodel

import (
	"time"

	"go.proteos.ai/model/common"
)

// ContactGroup is a generic org-shared audience taxonomy entry ("external
// client", "internal colleague", …) — NOT tone-specific. Membership lives on
// the contact (Contact.GroupKey, one group per contact) and is assignable by
// any flow: manual UI/API assignment, tone-synthesis auto-assignment, and
// whatever comes later. Keyed (org_id, key) and module-deployable
// (contact-groups/<key>.json), like conversation types.
type ContactGroup struct {
	OrgId string `json:"org_id"`
	// Key is the immutable identity within the org — the value models answer
	// with and the value stored on Contact.GroupKey.
	Key  string `json:"key" sortable:""`
	Name string `json:"name" sortable:""`
	// Description tells both humans and the tone-synthesis model who belongs in
	// this group — it is injected verbatim into the synthesis prompt.
	Description string `json:"description,omitempty"`
	// IsAutoCreated is true when tone synthesis proposed the group. Server-owned:
	// never settable from a module manifest or the API.
	IsAutoCreated bool `json:"is_auto_created"`
	// ModuleSlug attributes the group to the module that deployed it; empty =
	// not module-owned.
	ModuleSlug string         `json:"module_slug,omitempty"`
	CreatedAt  time.Time      `json:"created_at" sortable:""`
	CreatedBy  common.UserRef `json:"created_by"`
	UpdatedAt  time.Time      `json:"updated_at" sortable:""`
	UpdatedBy  common.UserRef `json:"updated_by"`
}

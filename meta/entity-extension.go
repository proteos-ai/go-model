package metamodel

import (
	"time"

	"go.proteos.ai/model/common"
)

// EntityExtension contributes attributes to ONE existing host entity without
// owning it — the way a second module adds fields to an entity shipped by
// another. The host's stored attribute list is a materialization:
//
//	platform attributes + the entity's own attributes + every extension's
//	attributes (each stamped with Attribute.ExtensionKey)
//
// so every reader of entities.attributes sees the merged schema with no extra
// lookup. Attribute names must be unique across the host and all of its
// extensions; a collision is refused, never resolved silently.
//
// The row stores the author's attribute list UNSTAMPED — the stamp is derived
// when the host is (re)materialized — so an extension's own GET / pull / plan
// round-trips cleanly. Identity is Key (PK org_id + key); EntitySlug is
// immutable after create — delete and recreate to re-host.
type EntityExtension struct {
	Key         string `json:"key" sortable:""`
	OrgId       string `json:"org_id" sortable:""`
	Name        string `json:"name" sortable:""`
	Description string `json:"description" sortable:""`
	// EntitySlug names the host entity the attributes are merged into.
	EntitySlug string `json:"entity_slug" sortable:""`
	ModuleSlug string `json:"module_slug" sortable:""`
	// Attributes are the contributed definitions exactly as authored: no
	// extension_key stamp, no platform attributes.
	Attributes []Attribute    `json:"attributes"`
	CreatedAt  time.Time      `json:"created_at" sortable:""`
	CreatedBy  common.UserRef `json:"created_by" sortable:""`
	UpdatedAt  time.Time      `json:"updated_at" sortable:""`
	UpdatedBy  common.UserRef `json:"updated_by" sortable:""`
}

// IsEntityExtensionAttribute reports whether an attribute in a host entity's
// stored list was contributed by an entity extension (it carries the stamp).
func IsEntityExtensionAttribute(attr Attribute) bool {
	return attr.ExtensionKey != ""
}

// StripEntityExtensionAttributes returns the list without the
// extension-stamped attributes, order otherwise preserved — i.e. the entity's
// own schema (platform attributes included). Shared by the server, which never
// trusts a client-supplied stamp and re-merges from the extension rows on every
// write, and by the CLI, whose pull / plan compare a host entity against its
// own attributes only. Always returns a fresh, non-nil slice.
func StripEntityExtensionAttributes(attributes []Attribute) []Attribute {
	own := make([]Attribute, 0, len(attributes))
	for _, attr := range attributes {
		if IsEntityExtensionAttribute(attr) {
			continue
		}
		own = append(own, attr)
	}
	return own
}

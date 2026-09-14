package metamodel

import (
	"time"

	"go.proteos.ai/model/common"
)

// ListExtensionColumnPlacement is one anchored run of columns. AnchorAttribute
// names a column of the host's OWN column list by its attribute path (a
// column's identity — a list never shows the same attribute twice); it is
// required for `before`/`after` and must be empty for `first`/`last`, which
// address the ends of the list. The columns land as one contiguous run in
// author order.
type ListExtensionColumnPlacement struct {
	AnchorAttribute string            `json:"anchor_attribute,omitempty"`
	Position        ExtensionPosition `json:"position"`
	Columns         []Column          `json:"columns"`
}

// ListExtension contributes columns and toolbar actions to ONE existing host
// list without owning it — the way the CRM add-on shows the `crm_score` it
// added to `contact` as a column of the owner's `active-customers` list. The
// host's stored columns and actions are a materialization:
//
//	the list's own columns + every extension's column placements
//	the list's own actions + every extension's actions (appended)
//
// each contributed column and action stamped with extension_key, so every
// reader of lists.columns / lists.actions sees the merged list with no extra
// lookup. A column's attribute must be unique across the host and all of its
// extensions; a collision is refused, never resolved silently. An anchor the
// host no longer shows is refused too — on the extension write and on the
// host push that would remove it. Sorting, filters, selection mode and the
// rest of the list stay the owner's.
//
// The row stores the author's columns and actions UNSTAMPED. Identity is Key
// (PK org_id + key); ListSlug is immutable after create.
type ListExtension struct {
	Key         string `json:"key" sortable:""`
	OrgId       string `json:"org_id" sortable:""`
	Name        string `json:"name" sortable:""`
	Description string `json:"description" sortable:""`
	// ListSlug names the host list the columns and actions are merged into.
	ListSlug   string `json:"list_slug" sortable:""`
	ModuleSlug string `json:"module_slug" sortable:""`
	// Columns are the contributed column placements exactly as authored.
	Columns []ListExtensionColumnPlacement `json:"columns"`
	// Actions are appended to the host's toolbar, exactly as authored. A
	// PageAction has no identity, so actions carry no anchor.
	Actions   []PageAction   `json:"actions"`
	CreatedAt time.Time      `json:"created_at" sortable:""`
	CreatedBy common.UserRef `json:"created_by" sortable:""`
	UpdatedAt time.Time      `json:"updated_at" sortable:""`
	UpdatedBy common.UserRef `json:"updated_by" sortable:""`
}

// IsListExtensionColumn reports whether a column in a host list's stored list
// was contributed by a list extension (it carries the stamp).
func IsListExtensionColumn(column Column) bool {
	return column.ExtensionKey != ""
}

// IsListExtensionAction reports whether a toolbar action was contributed by a
// list extension.
func IsListExtensionAction(action PageAction) bool {
	return action.ExtensionKey != ""
}

// StripListExtensionColumns returns the columns without the stamped ones,
// order otherwise preserved — the list's own columns. Shared by the server
// (never trusts a client-supplied stamp; re-merges from the rows on every
// write) and by the CLI (pull / plan compare a host against its own columns).
// Always returns a fresh, non-nil slice.
func StripListExtensionColumns(columns []Column) []Column {
	own := make([]Column, 0, len(columns))
	for _, column := range columns {
		if IsListExtensionColumn(column) {
			continue
		}
		own = append(own, column)
	}
	return own
}

// StripListExtensionActions returns the actions without the stamped ones —
// the list's own toolbar. Always returns a fresh, non-nil slice.
func StripListExtensionActions(actions []PageAction) []PageAction {
	own := make([]PageAction, 0, len(actions))
	for _, action := range actions {
		if IsListExtensionAction(action) {
			continue
		}
		own = append(own, action)
	}
	return own
}

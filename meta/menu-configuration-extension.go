package metamodel

import (
	"slices"
	"strings"
	"time"

	"go.proteos.ai/model/common"
)

// The virtual anchor names the menu's root item list, so an extension can
// prepend/append top-level items without naming one of the host's ids. The
// dot keeps it out of the kebab-case authored-id namespace; `first`/`last`
// only.
const (
	MenuConfigurationExtensionAnchorItems = "menu.items"
	// MenuConfigurationExtensionReservedIdPrefix is the id prefix of the virtual
	// anchor — no authored menu item may use it.
	MenuConfigurationExtensionReservedIdPrefix = "menu."
)

// IsMenuConfigurationExtensionVirtualAnchor reports whether an anchor id names
// the root list rather than an authored item.
func IsMenuConfigurationExtensionVirtualAnchor(anchorId string) bool {
	return anchorId == MenuConfigurationExtensionAnchorItems
}

// MenuConfigurationExtensionPosition is where a placement's block lands
// relative to its anchor. The four insert positions share their wire values
// with ExtensionPosition; `replace` and `remove` are menu-only and displace
// the anchor itself — which is why this is its own type and page/list
// validators never accept them.
type MenuConfigurationExtensionPosition string

const (
	MenuConfigurationExtensionPositionBefore  MenuConfigurationExtensionPosition = "before"
	MenuConfigurationExtensionPositionAfter   MenuConfigurationExtensionPosition = "after"
	MenuConfigurationExtensionPositionFirst   MenuConfigurationExtensionPosition = "first"
	MenuConfigurationExtensionPositionLast    MenuConfigurationExtensionPosition = "last"
	MenuConfigurationExtensionPositionReplace MenuConfigurationExtensionPosition = "replace"
	MenuConfigurationExtensionPositionRemove  MenuConfigurationExtensionPosition = "remove"
)

var SupportedMenuConfigurationExtensionPositions = []MenuConfigurationExtensionPosition{
	MenuConfigurationExtensionPositionBefore,
	MenuConfigurationExtensionPositionAfter,
	MenuConfigurationExtensionPositionFirst,
	MenuConfigurationExtensionPositionLast,
	MenuConfigurationExtensionPositionReplace,
	MenuConfigurationExtensionPositionRemove,
}

// IsInside reports whether the block goes INTO the anchor (its children, or
// the root list) rather than next to it.
func (position MenuConfigurationExtensionPosition) IsInside() bool {
	return position == MenuConfigurationExtensionPositionFirst || position == MenuConfigurationExtensionPositionLast
}

// IsDisplacing reports whether the placement takes the anchor's own slot
// (`replace`) or hides it (`remove`) — the host's item is displaced either way.
func (position MenuConfigurationExtensionPosition) IsDisplacing() bool {
	return position == MenuConfigurationExtensionPositionReplace || position == MenuConfigurationExtensionPositionRemove
}

// MenuConfigurationExtensionPlacement is one anchored block: an anchor (an
// item id of the host's OWN tree, or the virtual root), a position, and the
// contributed items. Items are injected as one contiguous run in author
// order. `remove` carries no items; every other position needs at least one.
type MenuConfigurationExtensionPlacement struct {
	AnchorId string                             `json:"anchor_id"`
	Position MenuConfigurationExtensionPosition `json:"position"`
	Items    []MenuItem                         `json:"items,omitempty"`
}

// MenuConfigurationExtension contributes items to ONE existing host menu
// without owning it — the way a second module adds an entry to a menu shipped
// by another. The host's stored items are a materialization:
//
//	the menu's own items + every extension's placements (each contributed
//	item stamped with extension_key; a replaced / removed host item kept
//	under the stamped item's replaced_item)
//
// so every reader of menu_configurations.items sees the merged tree with no
// extra lookup. Item ids must be unique across the host and all of its
// extensions; a collision is refused, never resolved silently. An anchor the
// host no longer has is refused too — on the extension write and on the host
// push that would remove it — so a host refactor can never drop or misplace a
// contribution.
//
// The row stores the author's placements UNSTAMPED — the stamp is derived
// when the host is (re)materialized — so an extension's own GET / pull / plan
// round-trips cleanly. Identity is Key (PK org_id + key); MenuSlug is
// immutable after create — delete and recreate to re-host.
type MenuConfigurationExtension struct {
	Key         string `json:"key" sortable:""`
	OrgId       string `json:"org_id" sortable:""`
	Name        string `json:"name" sortable:""`
	Description string `json:"description" sortable:""`
	// MenuSlug names the host menu configuration the placements are merged into.
	MenuSlug   string `json:"menu_slug" sortable:""`
	ModuleSlug string `json:"module_slug" sortable:""`
	// Placements are the contributed blocks exactly as authored: no
	// extension_key stamp anywhere.
	Placements []MenuConfigurationExtensionPlacement `json:"placements"`
	CreatedAt  time.Time                             `json:"created_at" sortable:""`
	CreatedBy  common.UserRef                        `json:"created_by" sortable:""`
	UpdatedAt  time.Time                             `json:"updated_at" sortable:""`
	UpdatedBy  common.UserRef                        `json:"updated_by" sortable:""`
}

// IsMenuConfigurationExtensionItem reports whether an item in a host menu's
// stored tree was contributed by a menu-configuration extension (it carries
// the stamp).
func IsMenuConfigurationExtensionItem(item MenuItem) bool {
	return item.ExtensionKey != ""
}

// CloneMenuItems deep-copies a tree (children and replaced items included).
func CloneMenuItems(items []MenuItem) []MenuItem {
	out := make([]MenuItem, 0, len(items))
	for _, item := range items {
		out = append(out, cloneMenuItem(item))
	}
	return out
}

func cloneMenuItem(item MenuItem) MenuItem {
	cloned := item
	cloned.Children = CloneMenuItems(item.Children)
	if item.ReplacedItem != nil {
		replaced := cloneMenuItem(*item.ReplacedItem)
		cloned.ReplacedItem = &replaced
	}
	return cloned
}

// MenuItemVisit is one node handed to WalkMenuItems: the item, its parent (nil
// at the root) and its index in the sibling list.
type MenuItemVisit struct {
	Item   *MenuItem
	Parent *MenuItem
	Index  int
}

// WalkMenuItems visits every item of the VISIBLE tree depth-first, parents
// before children. Replaced items (the host's originals kept under a stamp)
// are not visited: they are not part of the rendered menu.
func WalkMenuItems(items []MenuItem, visit func(MenuItemVisit)) {
	walkMenuItems(items, nil, visit)
}

func walkMenuItems(items []MenuItem, parent *MenuItem, visit func(MenuItemVisit)) {
	for i := range items {
		visit(MenuItemVisit{Item: &items[i], Parent: parent, Index: i})
		walkMenuItems(items[i].Children, &items[i], visit)
	}
}

// CanonicalizeMenuItemOrder puts every sibling list into canonical form in
// place: stable-sorted by `order`, then `order` rewritten to the index. Array
// position and `order` then always agree, so a merge can splice by position
// and every consumer can sort by `order`, and re-canonicalizing is a no-op
// (the host's own orders never drift once extended). Applied to merge output,
// strip output and both sides of a CLI plan.
func CanonicalizeMenuItemOrder(items []MenuItem) {
	slices.SortStableFunc(items, func(a, b MenuItem) int { return a.Order - b.Order })
	for i := range items {
		CanonicalizeMenuItemOrder(items[i].Children)
	}
	RenumberMenuItemOrder(items)
}

// RenumberMenuItemOrder rewrites `order` to the array index on every sibling
// list, in place, without sorting — for a tree whose array positions are
// already authoritative (a merge that spliced blocks by position).
func RenumberMenuItemOrder(items []MenuItem) {
	for i := range items {
		items[i].Order = i
		RenumberMenuItemOrder(items[i].Children)
	}
}

// StripMenuConfigurationExtensionItems returns a fresh copy of the tree
// without the extension-stamped items — i.e. the menu's own items, in
// canonical order: a stamped item that displaced a host item (replace /
// remove) gives its slot back to the original under replaced_item, any other
// stamped item (whole subtree) is dropped. Shared by the server, which never
// trusts a client-supplied stamp and re-merges from the extension rows on
// every write, and by the CLI, whose pull / plan compare a host menu against
// its own items only. Always returns a non-nil slice; the input is never
// mutated.
func StripMenuConfigurationExtensionItems(items []MenuItem) []MenuItem {
	own := stripMenuItems(CloneMenuItems(items))
	CanonicalizeMenuItemOrder(own)
	return own
}

func stripMenuItems(items []MenuItem) []MenuItem {
	kept := make([]MenuItem, 0, len(items))
	for _, item := range items {
		if IsMenuConfigurationExtensionItem(item) {
			if item.ReplacedItem == nil {
				continue
			}
			// The displaced original is host-owned and unstamped itself, and
			// keeps the slot the stamped item took — but its subtree is
			// stripped recursively: a descendant an earlier placement (of this
			// or another extension) displaced was snapshotted as a stamped
			// tombstone, and only unwinding that too yields the host's own tree.
			restored := *item.ReplacedItem
			// The slot's materialized order, not the original's pre-merge one:
			// surviving siblings carry post-merge numbers (inserts above may
			// have shifted them), and the canonical sort below must keep the
			// host's relative order.
			restored.Order = item.Order
			restored.Children = stripMenuItems(restored.Children)
			kept = append(kept, restored)
			continue
		}
		item.Children = stripMenuItems(item.Children)
		kept = append(kept, item)
	}
	return kept
}

// HasMenuConfigurationExtensionReservedId reports whether an authored id uses
// the prefix reserved for virtual anchors.
func HasMenuConfigurationExtensionReservedId(id string) bool {
	return strings.HasPrefix(id, MenuConfigurationExtensionReservedIdPrefix)
}

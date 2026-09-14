package metamodel

import (
	"encoding/json"
	"fmt"
	"time"

	"go.proteos.ai/model/common"
)

// Virtual anchors name the layout roots, so an extension can contribute to a
// host that authored no ids of its own, and can add a side panel to a host
// that has none. Dots keep them out of the kebab-case authored-id namespace;
// `first`/`last` only.
const (
	PageExtensionAnchorMain      = "layout.main"
	PageExtensionAnchorSidePanel = "layout.side_panel"
	// PageExtensionReservedIdPrefix is the id prefix of the virtual anchors —
	// no authored element or tab may use it.
	PageExtensionReservedIdPrefix = "layout."
)

// IsPageExtensionVirtualAnchor reports whether an anchor id names a layout
// root rather than an authored element.
func IsPageExtensionVirtualAnchor(anchorId string) bool {
	return anchorId == PageExtensionAnchorMain || anchorId == PageExtensionAnchorSidePanel
}

// PageExtensionPlacement is one anchored block: the anchor (an element id, a
// tab id, or a virtual root anchor) in the host's OWN layout, a position, and
// the payload. Exactly one of Elements / Tabs is set, and which one is right
// follows from the anchor: a tabs element (`first`/`last`) or a tab id
// (`before`/`after`) takes Tabs, everything else takes Elements. The payload is
// injected as one contiguous run in author order, so several siblings land in
// the host's row/column without a wrapper.
type PageExtensionPlacement struct {
	AnchorId string            `json:"anchor_id"`
	Position ExtensionPosition `json:"position"`
	Elements []LayoutElement   `json:"elements,omitempty"`
	Tabs     []LayoutTab       `json:"tabs,omitempty"`
}

// UnmarshalJSON exists because Elements is interface-typed and needs the
// discriminated dispatch.
func (placement *PageExtensionPlacement) UnmarshalJSON(data []byte) error {
	var wire struct {
		AnchorId string            `json:"anchor_id"`
		Position ExtensionPosition `json:"position"`
		Elements []json.RawMessage `json:"elements"`
		Tabs     []LayoutTab       `json:"tabs"`
	}
	if err := json.Unmarshal(data, &wire); err != nil {
		return err
	}
	placement.AnchorId = wire.AnchorId
	placement.Position = wire.Position
	placement.Tabs = wire.Tabs
	placement.Elements = nil
	for i, raw := range wire.Elements {
		if len(raw) == 0 || string(raw) == "null" {
			return fmt.Errorf("elements[%d]: element is required", i)
		}
		element, err := unmarshalLayoutElement(raw)
		if err != nil {
			return fmt.Errorf("elements[%d]: %w", i, err)
		}
		placement.Elements = append(placement.Elements, element)
	}
	return nil
}

// PageExtension contributes layout to ONE existing host page without owning
// it — the way a second module adds a card or a tab to a page shipped by
// another. The host's stored layout is a materialization:
//
//	the page's own layout + every extension's placements (each contributed
//	element and tab stamped with extension_key)
//
// so every reader of pages.layout sees the merged tree with no extra lookup.
// Element and tab ids must be unique across the host and all of its
// extensions; a collision is refused, never resolved silently. An anchor that
// the host no longer has is refused too — on the extension write and on the
// host page push that would remove it — so a host refactor can never drop or
// misplace a contribution.
//
// The row stores the author's placements UNSTAMPED — the stamp is derived
// when the host is (re)materialized — so an extension's own GET / pull / plan
// round-trips cleanly. Identity is Key (PK org_id + key); PageSlug is
// immutable after create — delete and recreate to re-host.
type PageExtension struct {
	Key         string `json:"key" sortable:""`
	OrgId       string `json:"org_id" sortable:""`
	Name        string `json:"name" sortable:""`
	Description string `json:"description" sortable:""`
	// PageSlug names the host page the placements are merged into.
	PageSlug   string `json:"page_slug" sortable:""`
	ModuleSlug string `json:"module_slug" sortable:""`
	// Placements are the contributed blocks exactly as authored: no
	// extension_key stamp anywhere.
	Placements []PageExtensionPlacement `json:"placements"`
	CreatedAt  time.Time                `json:"created_at" sortable:""`
	CreatedBy  common.UserRef           `json:"created_by" sortable:""`
	UpdatedAt  time.Time                `json:"updated_at" sortable:""`
	UpdatedBy  common.UserRef           `json:"updated_by" sortable:""`
}

// IsPageExtensionElement reports whether an element in a host page's stored
// layout was contributed by a page extension (it carries the stamp).
func IsPageExtensionElement(element LayoutElement) bool {
	return LayoutElementCommonProps(element).ExtensionKey != ""
}

// IsPageExtensionTab reports whether a tab was contributed by a page extension.
func IsPageExtensionTab(tab LayoutTab) bool {
	return tab.ExtensionKey != ""
}

// StripPageExtensionElements returns a fresh copy of the layout without the
// extension-stamped elements and tabs (whole subtrees), order otherwise
// preserved — i.e. the page's own layout. A stamped side-panel content (the
// wrapper the server creates for the layout.side_panel virtual anchor) removes
// the side panel. Shared by the server, which never trusts a client-supplied
// stamp and re-merges from the extension rows on every write, and by the CLI,
// whose pull / plan compare a host page against its own layout only. The
// input is never mutated.
func StripPageExtensionElements(layout PageLayout) (PageLayout, error) {
	own, err := ClonePageLayout(layout)
	if err != nil {
		return PageLayout{}, err
	}
	if own.Main != nil {
		own.Main = stripElement(own.Main)
	}
	if own.SidePanel != nil {
		if own.SidePanel.Content == nil || IsPageExtensionElement(own.SidePanel.Content) {
			own.SidePanel = nil
		} else {
			own.SidePanel.Content = stripElement(own.SidePanel.Content)
		}
	}
	return own, nil
}

// stripElement removes stamped descendants in place and returns the element,
// or nil when the element itself is stamped.
func stripElement(element LayoutElement) LayoutElement {
	if element == nil || IsPageExtensionElement(element) {
		return nil
	}
	switch typed := element.(type) {
	case *RowElement:
		typed.Children = stripChildren(typed.Children)
	case *ColumnElement:
		typed.Children = stripChildren(typed.Children)
	case *SectionElement:
		typed.Content = stripElement(typed.Content)
	case *CardElement:
		typed.Content = stripElement(typed.Content)
	case *TabsElement:
		kept := make([]LayoutTab, 0, len(typed.Tabs))
		for _, tab := range typed.Tabs {
			if IsPageExtensionTab(tab) {
				continue
			}
			tab.Content = stripElement(tab.Content)
			kept = append(kept, tab)
		}
		typed.Tabs = kept
	}
	return element
}

func stripChildren(children []LayoutElement) []LayoutElement {
	kept := make([]LayoutElement, 0, len(children))
	for _, child := range children {
		if stripped := stripElement(child); stripped != nil {
			kept = append(kept, stripped)
		}
	}
	return kept
}

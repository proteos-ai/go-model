package metamodel

import (
	"encoding/json"
	"fmt"
)

// LayoutVisit is one element seen by WalkLayout: the element, the element it
// hangs off (nil for the roots — layout.main and side_panel.content) and, when
// the parent is a tabs element, the tab whose content it is.
type LayoutVisit struct {
	Element LayoutElement
	Parent  LayoutElement
	Tab     *LayoutTab
}

// WalkLayout visits every element of a layout depth-first, parents before
// children: layout.main first, then side_panel.content. Nil slots are skipped.
// The single walker behind stripping, anchor lookup and id collection — the
// validator keeps its own path-aware walk because it also reports positions.
func WalkLayout(layout PageLayout, visit func(LayoutVisit)) {
	if layout.Main != nil {
		WalkLayoutElement(layout.Main, visit)
	}
	if layout.SidePanel != nil && layout.SidePanel.Content != nil {
		WalkLayoutElement(layout.SidePanel.Content, visit)
	}
}

// WalkLayoutElement visits root and every descendant, depth-first.
func WalkLayoutElement(root LayoutElement, visit func(LayoutVisit)) {
	walkLayoutElement(root, nil, nil, visit)
}

func walkLayoutElement(element, parent LayoutElement, tab *LayoutTab, visit func(LayoutVisit)) {
	if element == nil {
		return
	}
	visit(LayoutVisit{Element: element, Parent: parent, Tab: tab})
	switch typed := element.(type) {
	case *RowElement:
		for _, child := range typed.Children {
			walkLayoutElement(child, element, nil, visit)
		}
	case *ColumnElement:
		for _, child := range typed.Children {
			walkLayoutElement(child, element, nil, visit)
		}
	case *SectionElement:
		walkLayoutElement(typed.Content, element, nil, visit)
	case *CardElement:
		walkLayoutElement(typed.Content, element, nil, visit)
	case *TabsElement:
		for i := range typed.Tabs {
			walkLayoutElement(typed.Tabs[i].Content, element, &typed.Tabs[i], visit)
		}
	}
}

// LayoutElementCommonProps returns the CommonProps embedded by any concrete
// element (a copy — set through the concrete type).
func LayoutElementCommonProps(element LayoutElement) CommonProps {
	switch typed := element.(type) {
	case *RowElement:
		return typed.CommonProps
	case *ColumnElement:
		return typed.CommonProps
	case *SectionElement:
		return typed.CommonProps
	case *CardElement:
		return typed.CommonProps
	case *TabsElement:
		return typed.CommonProps
	case *FieldElement:
		return typed.CommonProps
	case *RelatedListElement:
		return typed.CommonProps
	case *RelatedRecordElement:
		return typed.CommonProps
	case *ComponentElement:
		return typed.CommonProps
	case *RecordFilterElement:
		return typed.CommonProps
	case *ListElement:
		return typed.CommonProps
	case *DividerElement:
		return typed.CommonProps
	case *WorkflowTriggerElement:
		return typed.CommonProps
	case *TextElement:
		return typed.CommonProps
	}
	return CommonProps{}
}

// LayoutElementID is the element's id ("" when it has none).
func LayoutElementID(element LayoutElement) string {
	return LayoutElementCommonProps(element).ID
}

// SetLayoutElementExtensionKey stamps one element (not its descendants).
func SetLayoutElementExtensionKey(element LayoutElement, key string) {
	switch typed := element.(type) {
	case *RowElement:
		typed.ExtensionKey = key
	case *ColumnElement:
		typed.ExtensionKey = key
	case *SectionElement:
		typed.ExtensionKey = key
	case *CardElement:
		typed.ExtensionKey = key
	case *TabsElement:
		typed.ExtensionKey = key
	case *FieldElement:
		typed.ExtensionKey = key
	case *RelatedListElement:
		typed.ExtensionKey = key
	case *RelatedRecordElement:
		typed.ExtensionKey = key
	case *ComponentElement:
		typed.ExtensionKey = key
	case *RecordFilterElement:
		typed.ExtensionKey = key
	case *ListElement:
		typed.ExtensionKey = key
	case *DividerElement:
		typed.ExtensionKey = key
	case *WorkflowTriggerElement:
		typed.ExtensionKey = key
	case *TextElement:
		typed.ExtensionKey = key
	}
}

// CloneLayoutElement deep-copies an element tree. Every concrete element is a
// pointer, so anything that stamps or splices a tree it did not build must
// clone first. A nil element clones to nil.
func CloneLayoutElement(element LayoutElement) (LayoutElement, error) {
	if element == nil {
		return nil, nil
	}
	raw, err := json.Marshal(element)
	if err != nil {
		return nil, fmt.Errorf("clone layout element: %w", err)
	}
	return unmarshalLayoutElement(raw)
}

// CloneLayoutTab deep-copies a tab, content included.
func CloneLayoutTab(tab LayoutTab) (LayoutTab, error) {
	raw, err := json.Marshal(tab)
	if err != nil {
		return LayoutTab{}, fmt.Errorf("clone layout tab: %w", err)
	}
	var out LayoutTab
	if err := json.Unmarshal(raw, &out); err != nil {
		return LayoutTab{}, fmt.Errorf("clone layout tab: %w", err)
	}
	return out, nil
}

// ClonePageLayout deep-copies a layout.
func ClonePageLayout(layout PageLayout) (PageLayout, error) {
	raw, err := json.Marshal(layout)
	if err != nil {
		return PageLayout{}, fmt.Errorf("clone page layout: %w", err)
	}
	var out PageLayout
	if err := json.Unmarshal(raw, &out); err != nil {
		return PageLayout{}, fmt.Errorf("clone page layout: %w", err)
	}
	return out, nil
}

package metamodel

import (
	"encoding/json"
	"fmt"

	"go.proteos.ai/model/common"
	"golang.org/x/exp/slices"
)

// LayoutElementType is the discriminator for nodes in a PageLayout tree.
type LayoutElementType string

const (
	LayoutElementTypeRow         LayoutElementType = "row"
	LayoutElementTypeColumn      LayoutElementType = "column"
	LayoutElementTypeSection     LayoutElementType = "section"
	LayoutElementTypeTabs        LayoutElementType = "tabs"
	LayoutElementTypeField       LayoutElementType = "field"
	LayoutElementTypeRelatedList LayoutElementType = "related_list"
	LayoutElementTypeComponent   LayoutElementType = "component"
	LayoutElementTypeDivider     LayoutElementType = "divider"
	LayoutElementTypeText        LayoutElementType = "text"
)

// LayoutElementTypes enumerates every valid type discriminator.
var LayoutElementTypes = []LayoutElementType{
	LayoutElementTypeRow,
	LayoutElementTypeColumn,
	LayoutElementTypeSection,
	LayoutElementTypeTabs,
	LayoutElementTypeField,
	LayoutElementTypeRelatedList,
	LayoutElementTypeComponent,
	LayoutElementTypeDivider,
	LayoutElementTypeText,
}

// UnmarshalJSON validates the wire value against LayoutElementTypes. Mirrors
// the hardening pattern used by ComparisonOperator / LogicalOperator etc.
func (k *LayoutElementType) UnmarshalJSON(b []byte) error {
	return unmarshalEnum(b, k, func(v LayoutElementType) bool {
		return slices.Contains(LayoutElementTypes, v)
	}, "layoutElementType")
}

// LayoutAlign values for align-self / cross-axis alignment.
type LayoutAlign string

const (
	LayoutAlignStart   LayoutAlign = "start"
	LayoutAlignCenter  LayoutAlign = "center"
	LayoutAlignEnd     LayoutAlign = "end"
	LayoutAlignStretch LayoutAlign = "stretch"
)

var LayoutAligns = []LayoutAlign{
	LayoutAlignStart, LayoutAlignCenter, LayoutAlignEnd, LayoutAlignStretch,
}

func (a *LayoutAlign) UnmarshalJSON(b []byte) error {
	return unmarshalEnum(b, a, func(v LayoutAlign) bool {
		return slices.Contains(LayoutAligns, v)
	}, "layoutAlign")
}

// LayoutJustify values for main-axis arrangement on Row/Column.
type LayoutJustify string

const (
	LayoutJustifyStart   LayoutJustify = "start"
	LayoutJustifyCenter  LayoutJustify = "center"
	LayoutJustifyEnd     LayoutJustify = "end"
	LayoutJustifyBetween LayoutJustify = "between"
	LayoutJustifyAround  LayoutJustify = "around"
)

var LayoutJustifies = []LayoutJustify{
	LayoutJustifyStart, LayoutJustifyCenter, LayoutJustifyEnd,
	LayoutJustifyBetween, LayoutJustifyAround,
}

func (j *LayoutJustify) UnmarshalJSON(b []byte) error {
	return unmarshalEnum(b, j, func(v LayoutJustify) bool {
		return slices.Contains(LayoutJustifies, v)
	}, "layoutJustify")
}

// LayoutGap values for spacing between children of a Row/Column.
type LayoutGap string

const (
	LayoutGapXS LayoutGap = "xs"
	LayoutGapSM LayoutGap = "sm"
	LayoutGapMD LayoutGap = "md"
	LayoutGapLG LayoutGap = "lg"
)

var LayoutGaps = []LayoutGap{LayoutGapXS, LayoutGapSM, LayoutGapMD, LayoutGapLG}

func (g *LayoutGap) UnmarshalJSON(b []byte) error {
	return unmarshalEnum(b, g, func(v LayoutGap) bool {
		return slices.Contains(LayoutGaps, v)
	}, "layoutGap")
}

// TextVariant values for TextElement.
type TextVariant string

const (
	TextVariantHeading    TextVariant = "heading"
	TextVariantSubheading TextVariant = "subheading"
	TextVariantBody       TextVariant = "body"
	TextVariantCaption    TextVariant = "caption"
	TextVariantCallout    TextVariant = "callout"
)

var TextVariants = []TextVariant{
	TextVariantHeading, TextVariantSubheading, TextVariantBody,
	TextVariantCaption, TextVariantCallout,
}

func (v *TextVariant) UnmarshalJSON(b []byte) error {
	return unmarshalEnum(b, v, func(x TextVariant) bool {
		return slices.Contains(TextVariants, x)
	}, "textVariant")
}

// unmarshalEnum is the shared body of the per-enum UnmarshalJSON: decode the
// raw string, treat null as zero-value, validate against the allow-list, and
// surface a typed error otherwise. Mirrors the inline implementation on
// ComparisonOperator / LogicalOperator.
func unmarshalEnum[T ~string](b []byte, dst *T, ok func(T) bool, name string) error {
	if string(b) == `null` {
		*dst = ""
		return nil
	}
	var s string
	if err := json.Unmarshal(b, &s); err != nil {
		return err
	}
	v := T(s)
	if !ok(v) {
		return fmt.Errorf("invalid %s: %q", name, s)
	}
	*dst = v
	return nil
}

// SizeValue is a polymorphic sizing value. Accepted forms:
//   - JSON number in [0, 1] for a fraction (0.5 → 50%)
//   - "n/m" fraction string ("1/2", "2/3", ...)
//   - "Npx" pixel string
//   - "N%" percent string
//   - "auto"  — size to content
//   - "fill"  — grow to fill remaining space
//
// The raw JSON is preserved so consumers and the renderer keep full fidelity
// over the wire form.
type SizeValue struct {
	raw json.RawMessage
}

// Raw returns the underlying JSON bytes (number or string literal).
func (s SizeValue) Raw() json.RawMessage { return s.raw }

// MarshalJSON emits the original wire form. A zero-value SizeValue is null.
func (s SizeValue) MarshalJSON() ([]byte, error) {
	if len(s.raw) == 0 {
		return []byte("null"), nil
	}
	return s.raw, nil
}

// UnmarshalJSON stores the raw bytes for later validation.
func (s *SizeValue) UnmarshalJSON(data []byte) error {
	s.raw = append([]byte(nil), data...)
	return nil
}

// SizingProps is the partial set of sizing knobs used by responsive overrides.
type SizingProps struct {
	Width  *SizeValue  `json:"width,omitempty"`
	Height *SizeValue  `json:"height,omitempty"`
	Grow   *float64    `json:"grow,omitempty"`
	Shrink *float64    `json:"shrink,omitempty"`
	Align  LayoutAlign `json:"align,omitempty"`
}

// ResponsiveSizing carries breakpoint-specific sizing overrides.
type ResponsiveSizing struct {
	SM *SizingProps `json:"sm,omitempty"`
	MD *SizingProps `json:"md,omitempty"`
	LG *SizingProps `json:"lg,omitempty"`
}

// CommonProps is embedded by every concrete LayoutElement struct.
//
// `VisibleWhen` and `ReadOnlyWhen` reuse the `common.FilterGroup` predicate model
// already used by lists and the data-service URL query convention — same
// operators, same and/or composition, pipe-joined values for in/not_in.
//
// `Align` here is align-self (override of the parent's cross-axis arrangement
// for this child). On Row/Column elements it is shadowed by an outer Align
// field whose semantic is the cross-axis arrangement applied to children.
type CommonProps struct {
	ID           string              `json:"id,omitempty"`
	VisibleWhen  *common.FilterGroup `json:"visible_when,omitempty"`
	ReadOnlyWhen *common.FilterGroup `json:"read_only_when,omitempty"`
	Width        *SizeValue          `json:"width,omitempty"`
	Height       *SizeValue          `json:"height,omitempty"`
	Grow         *float64            `json:"grow,omitempty"`
	Shrink       *float64            `json:"shrink,omitempty"`
	Align        LayoutAlign         `json:"align,omitempty"`
	Responsive   *ResponsiveSizing   `json:"responsive,omitempty"`
}

// LayoutElement is the discriminated-union interface implemented by every
// element type. Use *RowElement, *ColumnElement, ... as concrete values.
type LayoutElement interface {
	isLayoutElement()
	LayoutType() LayoutElementType
}

// ──────────────────────────────────────────────────────────── Row ──

type RowElement struct {
	Type LayoutElementType `json:"type"`
	CommonProps
	Gap        LayoutGap       `json:"gap,omitempty"`
	AllowsWrap *bool           `json:"allows_wrap,omitempty"`
	Align      LayoutAlign     `json:"align,omitempty"`
	Justify    LayoutJustify   `json:"justify,omitempty"`
	Children   []LayoutElement `json:"children"`
}

func (RowElement) isLayoutElement()              {}
func (RowElement) LayoutType() LayoutElementType { return LayoutElementTypeRow }
func (e *RowElement) UnmarshalJSON(data []byte) error {
	type wireRow RowElement // alias to avoid recursion
	var wire struct {
		wireRow
		Children []json.RawMessage `json:"children"`
	}
	if err := json.Unmarshal(data, &wire); err != nil {
		return err
	}
	*e = RowElement(wire.wireRow)
	e.Children = make([]LayoutElement, 0, len(wire.Children))
	for i, raw := range wire.Children {
		child, err := unmarshalLayoutElement(raw)
		if err != nil {
			return fmt.Errorf("children[%d]: %w", i, err)
		}
		e.Children = append(e.Children, child)
	}
	return nil
}

// ───────────────────────────────────────────────────────── Column ──

type ColumnElement struct {
	Type LayoutElementType `json:"type"`
	CommonProps
	Gap      LayoutGap       `json:"gap,omitempty"`
	Align    LayoutAlign     `json:"align,omitempty"`
	Justify  LayoutJustify   `json:"justify,omitempty"`
	Children []LayoutElement `json:"children"`
}

func (ColumnElement) isLayoutElement()              {}
func (ColumnElement) LayoutType() LayoutElementType { return LayoutElementTypeColumn }
func (e *ColumnElement) UnmarshalJSON(data []byte) error {
	type wireCol ColumnElement
	var wire struct {
		wireCol
		Children []json.RawMessage `json:"children"`
	}
	if err := json.Unmarshal(data, &wire); err != nil {
		return err
	}
	*e = ColumnElement(wire.wireCol)
	e.Children = make([]LayoutElement, 0, len(wire.Children))
	for i, raw := range wire.Children {
		child, err := unmarshalLayoutElement(raw)
		if err != nil {
			return fmt.Errorf("children[%d]: %w", i, err)
		}
		e.Children = append(e.Children, child)
	}
	return nil
}

// ──────────────────────────────────────────────────────── Section ──

type SectionElement struct {
	Type LayoutElementType `json:"type"`
	CommonProps
	Title            string        `json:"title,omitempty"`
	Description      string        `json:"description,omitempty"`
	IsCollapsible    *bool         `json:"is_collapsible,omitempty"`
	DefaultCollapsed *bool         `json:"default_collapsed,omitempty"`
	Content          LayoutElement `json:"content"`
}

func (SectionElement) isLayoutElement()              {}
func (SectionElement) LayoutType() LayoutElementType { return LayoutElementTypeSection }
func (e *SectionElement) UnmarshalJSON(data []byte) error {
	type wireSec SectionElement
	var wire struct {
		wireSec
		Content json.RawMessage `json:"content"`
	}
	if err := json.Unmarshal(data, &wire); err != nil {
		return err
	}
	*e = SectionElement(wire.wireSec)
	if len(wire.Content) > 0 && string(wire.Content) != "null" {
		content, err := unmarshalLayoutElement(wire.Content)
		if err != nil {
			return fmt.Errorf("content: %w", err)
		}
		e.Content = content
	}
	return nil
}

// ─────────────────────────────────────────────────────────── Tabs ──

type LayoutTab struct {
	ID          string              `json:"id"`
	Label       string              `json:"label"`
	Icon        string              `json:"icon,omitempty"`
	VisibleWhen *common.FilterGroup `json:"visible_when,omitempty"`
	Content     LayoutElement       `json:"content"`
}

func (t *LayoutTab) UnmarshalJSON(data []byte) error {
	type wireTab LayoutTab
	var wire struct {
		wireTab
		Content json.RawMessage `json:"content"`
	}
	if err := json.Unmarshal(data, &wire); err != nil {
		return err
	}
	*t = LayoutTab(wire.wireTab)
	if len(wire.Content) > 0 && string(wire.Content) != "null" {
		content, err := unmarshalLayoutElement(wire.Content)
		if err != nil {
			return fmt.Errorf("content: %w", err)
		}
		t.Content = content
	}
	return nil
}

type TabsElement struct {
	Type LayoutElementType `json:"type"`
	CommonProps
	Tabs         []LayoutTab `json:"tabs"`
	DefaultTabID string      `json:"default_tab_id,omitempty"`
}

func (TabsElement) isLayoutElement()              {}
func (TabsElement) LayoutType() LayoutElementType { return LayoutElementTypeTabs }

// ────────────────────────────────────────────────────────── Field ──

type FieldElement struct {
	Type LayoutElementType `json:"type"`
	CommonProps
	Attribute    string         `json:"attribute"`
	Label        *string        `json:"label,omitempty"` // null hides the label; absent uses default
	Description  string         `json:"description,omitempty"`
	Placeholder  string         `json:"placeholder,omitempty"`
	IsReadOnly   *bool          `json:"is_read_only,omitempty"`
	IsRequired   *bool          `json:"is_required,omitempty"`
	EmptyDisplay string         `json:"empty_display,omitempty"`
	Control      string         `json:"control,omitempty"`
	ControlProps map[string]any `json:"control_props,omitempty"`
}

func (FieldElement) isLayoutElement()              {}
func (FieldElement) LayoutType() LayoutElementType { return LayoutElementTypeField }

// ──────────────────────────────────────────────────── RelatedList ──

// RelatedListElement renders a list of records related to the current page
// record. The element points at a relation attribute on a related entity:
// `RelatedEntitySlug` names the entity whose records to list (e.g. "order" on
// a Customer page), and `ViaAttribute` names the relation attribute on that
// entity whose `Meta.RelatedEntitySlug` points back at the current entity
// (e.g. "customerId" on Order). `ListSlug` optionally pins which list
// definition drives the column / sort / filter model; when omitted, the
// renderer falls back to the first list configured for the related entity.
// `FollowsParentEditMode` controls whether the list enters row-edit mode
// together with the host page's edit mode: nil/true = follows the page
// (default), false = independent — rows are only editable through the
// list's own toggle.
type RelatedListElement struct {
	Type LayoutElementType `json:"type"`
	CommonProps
	RelatedEntitySlug     string `json:"related_entity_slug"`
	ViaAttribute          string `json:"via_attribute"`
	ListSlug              string `json:"list_slug,omitempty"`
	FollowsParentEditMode *bool  `json:"follows_parent_edit_mode,omitempty"`
}

func (RelatedListElement) isLayoutElement()              {}
func (RelatedListElement) LayoutType() LayoutElementType { return LayoutElementTypeRelatedList }

// ────────────────────────────────────────────────────── Component ──

type ComponentElement struct {
	Type LayoutElementType `json:"type"`
	CommonProps
	ComponentSlug string         `json:"component_slug"`
	Props         map[string]any `json:"props,omitempty"`
}

func (ComponentElement) isLayoutElement()              {}
func (ComponentElement) LayoutType() LayoutElementType { return LayoutElementTypeComponent }

// ──────────────────────────────────────────────────────── Divider ──

type DividerElement struct {
	Type LayoutElementType `json:"type"`
	CommonProps
}

func (DividerElement) isLayoutElement()              {}
func (DividerElement) LayoutType() LayoutElementType { return LayoutElementTypeDivider }

// ─────────────────────────────────────────────────────────── Text ──

type TextElement struct {
	Type LayoutElementType `json:"type"`
	CommonProps
	Variant TextVariant `json:"variant"`
	Content string      `json:"content"`
}

func (TextElement) isLayoutElement()              {}
func (TextElement) LayoutType() LayoutElementType { return LayoutElementTypeText }

// ──────────────────────────────────────────────────── PageLayout ──

// PageLayoutSidePanel is the named right-rail slot.
type PageLayoutSidePanel struct {
	Width    *SizeValue    `json:"width,omitempty"`
	IsSticky *bool         `json:"is_sticky,omitempty"`
	Content  LayoutElement `json:"content"`
}

func (s *PageLayoutSidePanel) UnmarshalJSON(data []byte) error {
	type wireSP PageLayoutSidePanel
	var wire struct {
		wireSP
		Content json.RawMessage `json:"content"`
	}
	if err := json.Unmarshal(data, &wire); err != nil {
		return err
	}
	*s = PageLayoutSidePanel(wire.wireSP)
	if len(wire.Content) > 0 && string(wire.Content) != "null" {
		content, err := unmarshalLayoutElement(wire.Content)
		if err != nil {
			return fmt.Errorf("content: %w", err)
		}
		s.Content = content
	}
	return nil
}

// PageStyle carries per-page presentation of the standalone page shell:
// background color and the content container's max-width / side & top padding.
// Every field is optional; an omitted field falls back to the renderer default
// (the historical bg-bg / 1200px / 24px-32px look). Setting max_width to "fill"
// (or "auto") together with zero padding lets a component fill the page
// edge-to-edge.
//
// Currently honored only on `public` and `kiosk` pages (see the metadata-service
// layout validator's type gate).
type PageStyle struct {
	// Background is a design-token key ("bg", "bg-2", "accent", …) resolved to
	// var(--color-<key>) by the renderer, or a raw CSS color for branded pages.
	Background *string `json:"background,omitempty"`
	// MaxWidth caps the content container. "fill"/"auto" removes the cap.
	MaxWidth *SizeValue `json:"max_width,omitempty"`
	// PaddingX / PaddingY are the horizontal / vertical padding of the content
	// container (mirrors the Tailwind px/py idiom). "0px" makes it flush.
	PaddingX *SizeValue `json:"padding_x,omitempty"`
	PaddingY *SizeValue `json:"padding_y,omitempty"`
}

// PageLayout is the typed layout document persisted on Page.
type PageLayout struct {
	Version   int                  `json:"version"`
	Main      LayoutElement        `json:"main"`
	SidePanel *PageLayoutSidePanel `json:"side_panel,omitempty"`
	Style     *PageStyle           `json:"style,omitempty"`
}

func (l *PageLayout) UnmarshalJSON(data []byte) error {
	var wire struct {
		Version   int                  `json:"version"`
		Main      json.RawMessage      `json:"main"`
		SidePanel *PageLayoutSidePanel `json:"side_panel,omitempty"`
		Style     *PageStyle           `json:"style,omitempty"`
	}
	if err := json.Unmarshal(data, &wire); err != nil {
		return err
	}
	l.Version = wire.Version
	if len(wire.Main) > 0 && string(wire.Main) != "null" {
		main, err := unmarshalLayoutElement(wire.Main)
		if err != nil {
			return fmt.Errorf("main: %w", err)
		}
		l.Main = main
	}
	l.SidePanel = wire.SidePanel
	l.Style = wire.Style
	return nil
}

// ─────────────────────────────────────────────── Dispatch helper ──

// unmarshalLayoutElement peeks at the `type` discriminator and unmarshals
// `data` into the matching concrete element type, returned as LayoutElement.
func unmarshalLayoutElement(data json.RawMessage) (LayoutElement, error) {
	var disc struct {
		Type LayoutElementType `json:"type"`
	}
	if err := json.Unmarshal(data, &disc); err != nil {
		return nil, fmt.Errorf("layout element: missing or invalid type: %w", err)
	}
	if !slices.Contains(LayoutElementTypes, disc.Type) {
		return nil, fmt.Errorf("layout element: unknown type %q", disc.Type)
	}
	switch disc.Type {
	case LayoutElementTypeRow:
		var v RowElement
		if err := json.Unmarshal(data, &v); err != nil {
			return nil, err
		}
		return &v, nil
	case LayoutElementTypeColumn:
		var v ColumnElement
		if err := json.Unmarshal(data, &v); err != nil {
			return nil, err
		}
		return &v, nil
	case LayoutElementTypeSection:
		var v SectionElement
		if err := json.Unmarshal(data, &v); err != nil {
			return nil, err
		}
		return &v, nil
	case LayoutElementTypeTabs:
		var v TabsElement
		if err := json.Unmarshal(data, &v); err != nil {
			return nil, err
		}
		return &v, nil
	case LayoutElementTypeField:
		var v FieldElement
		if err := json.Unmarshal(data, &v); err != nil {
			return nil, err
		}
		return &v, nil
	case LayoutElementTypeRelatedList:
		var v RelatedListElement
		if err := json.Unmarshal(data, &v); err != nil {
			return nil, err
		}
		return &v, nil
	case LayoutElementTypeComponent:
		var v ComponentElement
		if err := json.Unmarshal(data, &v); err != nil {
			return nil, err
		}
		return &v, nil
	case LayoutElementTypeDivider:
		var v DividerElement
		if err := json.Unmarshal(data, &v); err != nil {
			return nil, err
		}
		return &v, nil
	case LayoutElementTypeText:
		var v TextElement
		if err := json.Unmarshal(data, &v); err != nil {
			return nil, err
		}
		return &v, nil
	default:
		return nil, fmt.Errorf("layout element: unhandled type %q", disc.Type)
	}
}

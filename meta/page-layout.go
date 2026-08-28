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
	LayoutElementTypeRow     LayoutElementType = "row"
	LayoutElementTypeColumn  LayoutElementType = "column"
	LayoutElementTypeSection LayoutElementType = "section"
	// LayoutElementTypeCard is Section's carded twin: same grouping contract,
	// drawn as a real card with the title in its header bar.
	LayoutElementTypeCard        LayoutElementType = "card"
	LayoutElementTypeTabs        LayoutElementType = "tabs"
	LayoutElementTypeField       LayoutElementType = "field"
	LayoutElementTypeRelatedList LayoutElementType = "related_list"
	// LayoutElementTypeRelatedRecord is the singular counterpart of
	// LayoutElementTypeRelatedList: one related record, rendered as a page.
	LayoutElementTypeRelatedRecord LayoutElementType = "related_record"
	LayoutElementTypeComponent     LayoutElementType = "component"
	LayoutElementTypeDivider       LayoutElementType = "divider"
	LayoutElementTypeText          LayoutElementType = "text"
	// LayoutElementTypeRecordFilter is a filter builder placed on the page; it
	// publishes its filter under its element id for list elements to consume.
	LayoutElementTypeRecordFilter LayoutElementType = "record_filter"
	// LayoutElementTypeList renders a configured List's records — the
	// record-agnostic sibling of related_list.
	LayoutElementTypeList LayoutElementType = "list"
)

// LayoutElementTypes enumerates every valid type discriminator.
var LayoutElementTypes = []LayoutElementType{
	LayoutElementTypeRow,
	LayoutElementTypeColumn,
	LayoutElementTypeSection,
	LayoutElementTypeCard,
	LayoutElementTypeTabs,
	LayoutElementTypeField,
	LayoutElementTypeRelatedList,
	LayoutElementTypeRelatedRecord,
	LayoutElementTypeComponent,
	LayoutElementTypeDivider,
	LayoutElementTypeText,
	LayoutElementTypeRecordFilter,
	LayoutElementTypeList,
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

// ─────────────────────────────────────────────────────────── Card ──

// CardElement groups content the way SectionElement does, but renders as an
// actual card — border, surface fill, rounded corners — with Title in the
// card's header bar. Collapsing needs a header to hang the toggle off, so a
// card with neither Title nor Description never collapses.
type CardElement struct {
	Type LayoutElementType `json:"type"`
	CommonProps
	Title            string        `json:"title,omitempty"`
	Description      string        `json:"description,omitempty"`
	IsCollapsible    *bool         `json:"is_collapsible,omitempty"`
	DefaultCollapsed *bool         `json:"default_collapsed,omitempty"`
	Content          LayoutElement `json:"content"`
}

func (CardElement) isLayoutElement()              {}
func (CardElement) LayoutType() LayoutElementType { return LayoutElementTypeCard }

// UnmarshalJSON exists because Content is interface-typed and needs the
// discriminator dispatch; same shape as SectionElement's.
func (e *CardElement) UnmarshalJSON(data []byte) error {
	type wireCard CardElement
	var wire struct {
		wireCard
		Content json.RawMessage `json:"content"`
	}
	if err := json.Unmarshal(data, &wire); err != nil {
		return err
	}
	*e = CardElement(wire.wireCard)
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

// ────────────────────────────────────────────────── RelatedRecord ──

// RelatedRecordElement renders the FIRST record related to the current page
// record, laid out with a record page of its own. The relation is addressed
// exactly like RelatedListElement — inbound: `RelatedEntitySlug` names the
// entity to pull from and `ViaAttribute` names the relation attribute on that
// entity pointing back at the current one — but only the first match is
// rendered (oldest first, so the choice is stable).
//
// `PageSlug` optionally pins which record page supplies the layout; when
// omitted (or dangling) the renderer falls back to the related entity's
// default record page. `FollowsParentEditMode` mirrors RelatedListElement:
// nil/true = the element enters edit mode together with the host page,
// false = independent, editable only through its own hover control. Either
// way the element saves the related record itself — the host page's Save
// never covers it.
//
// The element carries no chrome of its own: it renders the nested page bare.
// Wrap it in a SectionElement for a title or collapse affordance.
type RelatedRecordElement struct {
	Type LayoutElementType `json:"type"`
	CommonProps
	RelatedEntitySlug     string `json:"related_entity_slug"`
	ViaAttribute          string `json:"via_attribute"`
	PageSlug              string `json:"page_slug,omitempty"`
	FollowsParentEditMode *bool  `json:"follows_parent_edit_mode,omitempty"`
}

func (RelatedRecordElement) isLayoutElement()              {}
func (RelatedRecordElement) LayoutType() LayoutElementType { return LayoutElementTypeRelatedRecord }

// ────────────────────────────────────────────────────── Component ──

type ComponentElement struct {
	Type LayoutElementType `json:"type"`
	CommonProps
	ComponentSlug string         `json:"component_slug"`
	Props         map[string]any `json:"props,omitempty"`
	// ReservedHeight is the placeholder height (px) the host reserves while the
	// component bundle loads, used ONLY when Height is unset/auto. It does not
	// constrain the rendered component — after load it auto-sizes to its real
	// content. Reserving roughly the right space keeps the reveal from shoving
	// the content below it. Nil falls back to the host's default (80px).
	ReservedHeight *float64 `json:"reserved_height,omitempty"`
}

func (ComponentElement) isLayoutElement()              {}
func (ComponentElement) LayoutType() LayoutElementType { return LayoutElementTypeComponent }

// ─────────────────────────────────────────────────── RecordFilter ──

// RecordFilterVariant selects the filter element's chrome.
type RecordFilterVariant string

const (
	// RecordFilterVariantToolbar is the Filter button plus active-filter chips,
	// matching every list view. The default.
	RecordFilterVariantToolbar RecordFilterVariant = "toolbar"
	// RecordFilterVariantPanel is the always-open AND/OR editor, for a page
	// whose point IS the filter.
	RecordFilterVariantPanel RecordFilterVariant = "panel"
)

var RecordFilterVariants = []RecordFilterVariant{
	RecordFilterVariantToolbar, RecordFilterVariantPanel,
}

func (variant *RecordFilterVariant) UnmarshalJSON(b []byte) error {
	return unmarshalEnum(b, variant, func(v RecordFilterVariant) bool {
		return slices.Contains(RecordFilterVariants, v)
	}, "recordFilterVariant")
}

// Normalized returns the variant with the empty value defaulting to toolbar.
func (variant RecordFilterVariant) Normalized() RecordFilterVariant {
	if variant == "" {
		return RecordFilterVariantToolbar
	}
	return variant
}

// RecordFilterElement is a filter builder placed directly on a page. It binds
// to no data itself: it publishes the authored filter, and the chosen subject
// entity, under its element ID. Every ListElement naming that ID in
// FilterElementID renders the filtered rows.
//
// SubjectEntity pins the entity the filter is authored against. Left empty, the
// element renders a subject picker over the entities its bound lists offer —
// one entry per bound list, since a List carries exactly one entity.
//
// IsComplexEnabled (default true) allows nested AND/OR groups. The records
// query carries the whole tree, so the affordance costs nothing.
type RecordFilterElement struct {
	Type LayoutElementType `json:"type"`
	CommonProps
	SubjectEntity    string              `json:"subject_entity,omitempty"`
	Variant          RecordFilterVariant `json:"variant,omitempty"`
	IsComplexEnabled *bool               `json:"is_complex_enabled,omitempty"`
}

func (RecordFilterElement) isLayoutElement()              {}
func (RecordFilterElement) LayoutType() LayoutElementType { return LayoutElementTypeRecordFilter }

// ─────────────────────────────────────────────────────────── List ──

// ListElement renders the records of a configured List. It is the
// record-agnostic sibling of RelatedListElement, and the only way to show
// records on a page that has no record of its own.
//
// The List supplies everything about presentation and behaviour — columns,
// sorting, base filters, toolbar actions, selection mode, and the page a row
// opens. Nothing is restated on the element: a second column or action model
// here would be the same concept under a second name.
//
// ListSlugs may name more than one list, which makes the element switchable;
// the active one is chosen by the bound filter's subject picker, or by the
// element's own switcher when it is unbound. FilterElementID binds the element
// to a RecordFilterElement on the same page; the bound filter is ANDed with the
// list's own saved filters, narrowing the list rather than replacing what it
// declared.
type ListElement struct {
	Type LayoutElementType `json:"type"`
	CommonProps
	ListSlugs       []string `json:"list_slugs"`
	FilterElementID string   `json:"filter_element_id,omitempty"`
	// FilterAttribute is the record-page alternative to FilterElementID: the
	// attribute on THIS page's record holding a saved filter, as written by the
	// `record-filter` field control. The list then renders what the record's own
	// filter selects — a saved-segment record showing its matches.
	//
	// Mutually exclusive with FilterElementID (two filters driving one list has
	// no defined precedence) and meaningless without a record, so the validator
	// rejects both-at-once and rejects it on standalone pages.
	FilterAttribute string `json:"filter_attribute,omitempty"`
	// SubjectEntityAttribute names the attribute holding the subject entity
	// slug; the element renders whichever of ListSlugs targets that entity.
	// Without it the first configured list wins.
	SubjectEntityAttribute string `json:"subject_entity_attribute,omitempty"`
	PageSize               *int   `json:"page_size,omitempty"`
}

func (ListElement) isLayoutElement()              {}
func (ListElement) LayoutType() LayoutElementType { return LayoutElementTypeList }

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
// Honored on `platform`, `public` and `kiosk` pages; record pages reject it
// (see the metadata-service layout validator's type gate). On platform pages
// it styles the app shell's content area rather than a standalone shell.
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
	case LayoutElementTypeCard:
		var v CardElement
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
	case LayoutElementTypeRelatedRecord:
		var v RelatedRecordElement
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
	case LayoutElementTypeRecordFilter:
		var v RecordFilterElement
		if err := json.Unmarshal(data, &v); err != nil {
			return nil, err
		}
		return &v, nil
	case LayoutElementTypeList:
		var v ListElement
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

package metamodel

import (
	"encoding/json"
	"go.proteos.ai/model/common"
	"strings"
	"testing"
)

// Round-trips the design-doc §11 worked example through Unmarshal+Marshal and
// asserts the result re-decodes to a structurally identical layout. The
// `visibleWhen` clause uses the common.FilterGroup wire shape (logicalOperator +
// elements/groups with pipe-joined values) — same as list filters.
func TestPageLayout_RoundTrip(t *testing.T) {
	src := `{
  "version": 1,
  "main": {
    "type": "column", "gap": "lg",
    "children": [
      {
        "type": "section", "title": "Details",
        "content": {
          "type": "column", "gap": "md",
          "children": [
            { "type": "field", "attribute": "name", "is_required": true },
            { "type": "row", "gap": "md", "children": [
              { "type": "field", "attribute": "ownerId", "control": "user-picker" },
              { "type": "field", "attribute": "stage" }
            ]},
            { "type": "row", "gap": "sm", "children": [
              { "type": "field", "attribute": "amount",   "width": "fill",  "control": "currency" },
              { "type": "field", "attribute": "currency", "width": "120px" },
              { "type": "field", "attribute": "closeDate", "width": "auto",
                "visible_when": {
                  "logical_operator": "and",
                  "elements": [
                    { "field": "stage", "operator": "in", "value": "proposal|negotiation|won" }
                  ]
                }
              }
            ]}
          ]
        }
      },
      {
        "type": "tabs",
        "tabs": [
          { "id": "activity", "label": "Activity",
            "content": { "type": "component", "component_slug": "activity-feed" } },
          { "id": "contacts", "label": "Contacts",
            "content": { "type": "related_list", "related_entity_slug": "contact", "via_attribute": "accountId", "follows_parent_edit_mode": false } },
          { "id": "primary", "label": "Primary contact",
            "content": { "type": "related_record", "related_entity_slug": "contact", "via_attribute": "accountId", "page_slug": "contact-detail" } }
        ]
      }
    ]
  },
  "side_panel": {
    "width": "320px", "is_sticky": true,
    "content": {
      "type": "column", "gap": "md",
      "children": [
        { "type": "section", "title": "Owner",
          "content": { "type": "field", "attribute": "ownerId", "control": "user-card" } }
      ]
    }
  }
}`

	var layout PageLayout
	if err := json.Unmarshal([]byte(src), &layout); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}

	if layout.Version != 1 {
		t.Errorf("version: want 1, got %d", layout.Version)
	}
	mainCol, ok := layout.Main.(*ColumnElement)
	if !ok {
		t.Fatalf("main: want *ColumnElement, got %T", layout.Main)
	}
	if len(mainCol.Children) != 2 {
		t.Fatalf("main.children: want 2, got %d", len(mainCol.Children))
	}
	if _, ok := mainCol.Children[0].(*SectionElement); !ok {
		t.Errorf("main.children[0]: want *SectionElement, got %T", mainCol.Children[0])
	}
	if _, ok := mainCol.Children[1].(*TabsElement); !ok {
		t.Errorf("main.children[1]: want *TabsElement, got %T", mainCol.Children[1])
	}
	if tabs, ok := mainCol.Children[1].(*TabsElement); ok {
		related, ok := tabs.Tabs[1].Content.(*RelatedListElement)
		if !ok {
			t.Fatalf("tabs[1].content: want *RelatedListElement, got %T", tabs.Tabs[1].Content)
		}
		if related.FollowsParentEditMode == nil || *related.FollowsParentEditMode {
			t.Errorf("followsParentEditMode: want false, got %v", related.FollowsParentEditMode)
		}
		record, ok := tabs.Tabs[2].Content.(*RelatedRecordElement)
		if !ok {
			t.Fatalf("tabs[2].content: want *RelatedRecordElement, got %T", tabs.Tabs[2].Content)
		}
		if record.RelatedEntitySlug != "contact" || record.ViaAttribute != "accountId" {
			t.Errorf("relation: want contact.accountId, got %s.%s", record.RelatedEntitySlug, record.ViaAttribute)
		}
		if record.PageSlug != "contact-detail" {
			t.Errorf("pageSlug: want contact-detail, got %q", record.PageSlug)
		}
		if record.FollowsParentEditMode != nil {
			t.Errorf("followsParentEditMode: want nil (defaults to follow), got %v", *record.FollowsParentEditMode)
		}
	}

	if layout.SidePanel == nil {
		t.Fatal("sidePanel: missing")
	}
	if _, ok := layout.SidePanel.Content.(*ColumnElement); !ok {
		t.Errorf("sidePanel.content: want *ColumnElement, got %T", layout.SidePanel.Content)
	}

	out, err := json.Marshal(&layout)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	var roundTripped PageLayout
	if err := json.Unmarshal(out, &roundTripped); err != nil {
		t.Fatalf("re-unmarshal: %v", err)
	}
}

// An unpinned related_record leaves page_slug off the wire entirely, so the
// renderer falls back to the related entity's default record page.
func TestRelatedRecordElement_OmitsUnpinnedPageSlug(t *testing.T) {
	src := `{"version":1,"main":{"type":"related_record","related_entity_slug":"contact","via_attribute":"accountId"}}`
	var layout PageLayout
	if err := json.Unmarshal([]byte(src), &layout); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if _, ok := layout.Main.(*RelatedRecordElement); !ok {
		t.Fatalf("main: want *RelatedRecordElement, got %T", layout.Main)
	}
	out, err := json.Marshal(&layout)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	if strings.Contains(string(out), "page_slug") {
		t.Errorf("page_slug: want omitted, got %s", string(out))
	}
}

func TestCardElement_RoundTrip(t *testing.T) {
	src := `{"version":1,"main":{"type":"card","title":"Details","description":"All of it",` +
		`"is_collapsible":true,"default_collapsed":true,` +
		`"content":{"type":"column","children":[{"type":"field","attribute":"name"}]}}}`
	var layout PageLayout
	if err := json.Unmarshal([]byte(src), &layout); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	card, ok := layout.Main.(*CardElement)
	if !ok {
		t.Fatalf("main: want *CardElement, got %T", layout.Main)
	}
	if card.Title != "Details" || card.Description != "All of it" {
		t.Errorf("title/description: got %q / %q", card.Title, card.Description)
	}
	if card.IsCollapsible == nil || !*card.IsCollapsible {
		t.Errorf("is_collapsible: want true, got %v", card.IsCollapsible)
	}
	if card.DefaultCollapsed == nil || !*card.DefaultCollapsed {
		t.Errorf("default_collapsed: want true, got %v", card.DefaultCollapsed)
	}
	// The interface-typed Content must dispatch, not decode as nil.
	col, ok := card.Content.(*ColumnElement)
	if !ok {
		t.Fatalf("content: want *ColumnElement, got %T", card.Content)
	}
	if len(col.Children) != 1 {
		t.Fatalf("content.children: want 1, got %d", len(col.Children))
	}
	if _, err := json.Marshal(&layout); err != nil {
		t.Fatalf("marshal: %v", err)
	}
}

func TestCardElement_OmitsUnsetOptionals(t *testing.T) {
	src := `{"version":1,"main":{"type":"card","content":{"type":"column","children":[]}}}`
	var layout PageLayout
	if err := json.Unmarshal([]byte(src), &layout); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	out, err := json.Marshal(&layout)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	for _, key := range []string{"title", "description", "is_collapsible", "default_collapsed"} {
		if strings.Contains(string(out), key) {
			t.Errorf("%s: want omitted, got %s", key, string(out))
		}
	}
}

func TestPageLayout_UnknownType(t *testing.T) {
	src := `{"version":1,"main":{"type":"fild","attribute":"x"}}`
	var layout PageLayout
	err := json.Unmarshal([]byte(src), &layout)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if !strings.Contains(err.Error(), "layoutElementType") || !strings.Contains(err.Error(), "fild") {
		t.Errorf("expected layoutElementType error mentioning 'fild', got: %v", err)
	}
}

// Verifies that nested common.FilterGroup (groups of groups, mixed with elements)
// round-trips cleanly when stored inside a layout's visibleWhen.
func TestPageLayout_NestedFilterGroup(t *testing.T) {
	src := `{"version":1,"main":{
		"type":"field","attribute":"x",
		"visible_when":{
			"logical_operator":"and",
			"elements":[
				{"field":"a","operator":"eq","value":"1"}
			],
			"groups":[
				{
					"logical_operator":"or",
					"elements":[
						{"field":"b","operator":"empty","value":""},
						{"field":"c","operator":"in","value":"x|y"}
					]
				}
			]
		}
	}}`
	var layout PageLayout
	if err := json.Unmarshal([]byte(src), &layout); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	field := layout.Main.(*FieldElement)
	if field.VisibleWhen == nil {
		t.Fatal("visibleWhen: missing")
	}
	if field.VisibleWhen.LogicalOperator != common.LogicalOperatorAnd {
		t.Errorf("logicalOperator: want and, got %v", field.VisibleWhen.LogicalOperator)
	}
	if len(field.VisibleWhen.Elements) != 1 {
		t.Errorf("elements: want 1, got %d", len(field.VisibleWhen.Elements))
	}
	if len(field.VisibleWhen.Groups) != 1 {
		t.Fatalf("groups: want 1, got %d", len(field.VisibleWhen.Groups))
	}
	inner := field.VisibleWhen.Groups[0]
	if inner.LogicalOperator != common.LogicalOperatorOr {
		t.Errorf("inner.logicalOperator: want or, got %v", inner.LogicalOperator)
	}
	if len(inner.Elements) != 2 {
		t.Fatalf("inner.elements: want 2, got %d", len(inner.Elements))
	}
	if inner.Elements[1].Operator != common.ComparisonOperatorIn || inner.Elements[1].Value != "x|y" {
		t.Errorf("inner.elements[1]: %+v", inner.Elements[1])
	}
}

func TestSizeValue_PreservesWireForm(t *testing.T) {
	cases := []string{`"fill"`, `"auto"`, `"1/2"`, `"320px"`, `"50%"`, `0.5`, `1`}
	for _, c := range cases {
		var s SizeValue
		if err := json.Unmarshal([]byte(c), &s); err != nil {
			t.Errorf("unmarshal %s: %v", c, err)
			continue
		}
		out, err := json.Marshal(&s)
		if err != nil {
			t.Errorf("marshal %s: %v", c, err)
			continue
		}
		if string(out) != c {
			t.Errorf("round-trip: want %s, got %s", c, string(out))
		}
	}
}

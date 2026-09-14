package metamodel

import (
	"encoding/json"
	"testing"
)

func parseTestLayout(t *testing.T, src string) PageLayout {
	t.Helper()
	var layout PageLayout
	if err := json.Unmarshal([]byte(src), &layout); err != nil {
		t.Fatalf("parse layout: %v", err)
	}
	return layout
}

func layoutIDs(layout PageLayout) []string {
	ids := []string{}
	WalkLayout(layout, func(visit LayoutVisit) {
		ids = append(ids, LayoutElementID(visit.Element))
		if tabs, ok := visit.Element.(*TabsElement); ok {
			for _, tab := range tabs.Tabs {
				ids = append(ids, "tab:"+tab.ID)
			}
		}
	})
	return ids
}

func equalStrings(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

const stampedLayout = `{"version":1,
  "main":{"type":"column","id":"main","children":[
    {"type":"field","id":"f-name","attribute":"name"},
    {"type":"row","id":"crm-row","extension_key":"crm","children":[
      {"type":"field","id":"crm-score","extension_key":"crm","attribute":"crm_score"}
    ]},
    {"type":"section","id":"sec","content":{"type":"column","id":"sec-col","children":[
      {"type":"field","id":"f-a","attribute":"a"},
      {"type":"field","id":"x-b","extension_key":"x","attribute":"b"}
    ]}},
    {"type":"tabs","id":"tabs","tabs":[
      {"id":"tab-a","label":"A","content":{"type":"divider","id":"d-a"}},
      {"id":"tab-crm","label":"CRM","extension_key":"crm","content":{"type":"divider","id":"d-crm","extension_key":"crm"}}
    ]}
  ]},
  "side_panel":{"content":{"type":"column","id":"layout.side_panel","extension_key":"crm","children":[
    {"type":"field","id":"sp-x","extension_key":"crm","attribute":"x"}
  ]}}}`

func TestStripPageExtensionElements_RemovesStampedSubtreesTabsAndSidePanel(t *testing.T) {
	layout := parseTestLayout(t, stampedLayout)

	own, err := StripPageExtensionElements(layout)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	want := []string{"main", "f-name", "sec", "sec-col", "f-a", "tabs", "tab:tab-a", "d-a"}
	if got := layoutIDs(own); !equalStrings(got, want) {
		t.Errorf("expected own ids %v, got %v", want, got)
	}
	if own.SidePanel != nil {
		t.Error("a stamped side-panel wrapper must remove the side panel")
	}
	// The input is untouched.
	if len(layout.Main.(*ColumnElement).Children) != 4 || layout.SidePanel == nil {
		t.Fatal("StripPageExtensionElements must not mutate its input")
	}
}

func TestStripPageExtensionElements_KeepsOwnSidePanel(t *testing.T) {
	layout := parseTestLayout(t, `{"version":1,"main":{"type":"divider","id":"m"},
	  "side_panel":{"width":"320px","content":{"type":"column","id":"sp","children":[
	    {"type":"field","id":"sp-own","attribute":"a"},
	    {"type":"field","id":"sp-ext","extension_key":"k","attribute":"b"}]}}}`)

	own, err := StripPageExtensionElements(layout)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if own.SidePanel == nil || own.SidePanel.Width == nil {
		t.Fatal("an own side panel must survive with its sizing")
	}
	if got := layoutIDs(own); !equalStrings(got, []string{"m", "sp", "sp-own"}) {
		t.Errorf("unexpected ids %v", got)
	}
}

func TestPageExtensionPlacement_RoundTripsElementsAndTabs(t *testing.T) {
	src := `{"key":"crm","page_slug":"p","placements":[
	  {"anchor_id":"f-name","position":"after","elements":[{"type":"field","id":"s","attribute":"score"}]},
	  {"anchor_id":"tabs","position":"last","tabs":[{"id":"tab-crm","label":"CRM","content":{"type":"divider","id":"d"}}]}]}`
	var extension PageExtension
	if err := json.Unmarshal([]byte(src), &extension); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if len(extension.Placements) != 2 {
		t.Fatalf("expected 2 placements, got %d", len(extension.Placements))
	}
	if _, ok := extension.Placements[0].Elements[0].(*FieldElement); !ok {
		t.Errorf("elements must decode to concrete types, got %T", extension.Placements[0].Elements[0])
	}
	if extension.Placements[1].Tabs[0].Content.LayoutType() != LayoutElementTypeDivider {
		t.Error("tab content must decode")
	}
	raw, err := json.Marshal(extension)
	if err != nil {
		t.Fatal(err)
	}
	var again PageExtension
	if err := json.Unmarshal(raw, &again); err != nil {
		t.Fatalf("re-unmarshal: %v", err)
	}
	if again.Placements[0].Position != ExtensionPositionAfter || again.Placements[1].Tabs[0].ID != "tab-crm" {
		t.Errorf("round trip lost data: %+v", again.Placements)
	}
}

func TestCloneLayoutElement_IsIndependent(t *testing.T) {
	layout := parseTestLayout(t, `{"version":1,"main":{"type":"column","id":"c","children":[{"type":"divider","id":"d"}]}}`)
	clone, err := CloneLayoutElement(layout.Main)
	if err != nil {
		t.Fatal(err)
	}
	SetLayoutElementExtensionKey(clone.(*ColumnElement).Children[0], "k")
	if IsPageExtensionElement(layout.Main.(*ColumnElement).Children[0]) {
		t.Error("stamping the clone must not touch the original")
	}
	if nilClone, err := CloneLayoutElement(nil); err != nil || nilClone != nil {
		t.Error("nil clones to nil")
	}
}

func TestWalkLayout_ReportsParentAndTab(t *testing.T) {
	layout := parseTestLayout(t, `{"version":1,"main":{"type":"tabs","id":"t","tabs":[{"id":"a","label":"A","content":{"type":"divider","id":"d"}}]}}`)
	var seen []LayoutVisit
	WalkLayout(layout, func(visit LayoutVisit) { seen = append(seen, visit) })
	if len(seen) != 2 || seen[0].Parent != nil || seen[1].Parent != layout.Main || seen[1].Tab == nil || seen[1].Tab.ID != "a" {
		t.Errorf("unexpected walk: %+v", seen)
	}
}

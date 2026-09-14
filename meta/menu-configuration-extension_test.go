package metamodel

import (
	"encoding/json"
	"fmt"
	"testing"
)

func menuIds(items []MenuItem) []string {
	out := []string{}
	WalkMenuItems(items, func(visit MenuItemVisit) {
		token := visit.Item.Id
		if IsMenuConfigurationExtensionItem(*visit.Item) {
			token += "*"
		}
		if visit.Item.IsHidden {
			token += "~"
		}
		out = append(out, token)
	})
	return out
}

func TestStripMenuConfigurationExtensionItems(t *testing.T) {
	original := MenuItem{Id: "pricing", Order: 1, Label: "Pricing", Type: MenuItemTypeGroup, Children: []MenuItem{
		{Id: "sales-prices", Order: 0, Type: MenuItemTypeList, Reference: "sales-prices"},
	}}
	removed := MenuItem{Id: "manufacturing", Order: 2, Label: "Manufacturing", Type: MenuItemTypeGroup, Children: []MenuItem{
		{Id: "bom", Order: 0, Type: MenuItemTypeList, Reference: "bom"},
	}}
	materialized := []MenuItem{
		{Id: "catalog", Order: 0, Type: MenuItemTypeGroup, Children: []MenuItem{
			{Id: "products", Order: 0, Type: MenuItemTypeList, Reference: "products"},
			{Id: "variants", Order: 1, Type: MenuItemTypeList, Reference: "variants", ExtensionKey: "plus"},
			{Id: "categories", Order: 2, Type: MenuItemTypeList, Reference: "categories"},
		}},
		{Id: "pricing", Order: 1, Label: "Pricing+", Type: MenuItemTypeGroup, ExtensionKey: "plus", ReplacedItem: &original, Children: []MenuItem{
			{Id: "price-lists", Order: 0, Type: MenuItemTypeList, ExtensionKey: "plus"},
		}},
		{Id: "manufacturing", Order: 2, Type: MenuItemTypeGroup, ExtensionKey: "plus", IsHidden: true, ReplacedItem: &removed},
		{Id: "reports", Order: 3, Type: MenuItemTypeGroup, ExtensionKey: "plus", Children: []MenuItem{
			{Id: "margins", Order: 0, Type: MenuItemTypePage, ExtensionKey: "plus"},
		}},
	}
	before, _ := json.Marshal(materialized)

	own := StripMenuConfigurationExtensionItems(materialized)

	if got := fmt.Sprint(menuIds(own)); got != "[catalog products categories pricing sales-prices manufacturing bom]" {
		t.Errorf("unexpected own ids %s", got)
	}
	if own[1].Label != "Pricing" || own[1].ReplacedItem != nil || own[1].ExtensionKey != "" {
		t.Errorf("replace must restore the original item, got %+v", own[1])
	}
	if own[2].IsHidden || own[2].Children[0].Id != "bom" {
		t.Errorf("remove must restore the hidden original with its subtree, got %+v", own[2])
	}
	if own[0].Children[1].Order != 1 || own[2].Order != 2 {
		t.Errorf("strip must renumber siblings, got %+v", own)
	}
	if after, _ := json.Marshal(materialized); string(after) != string(before) {
		t.Fatal("strip must not mutate its input")
	}
	if StripMenuConfigurationExtensionItems(nil) == nil {
		t.Fatal("expected an empty, non-nil slice for nil input")
	}
}

func TestStripMenuConfigurationExtensionItems_NestedDisplacements(t *testing.T) {
	// A child removed by one placement, then its parent removed by a later one
	// (or by another extension): the parent's replaced_item carries the child's
	// stamped tombstone, which must unwind to the host's own child.
	childOriginal := MenuItem{Id: "bom", Order: 0, Label: "BOM", Type: MenuItemTypeList, Reference: "bom", Children: []MenuItem{}}
	childTombstone := MenuItem{Id: "bom", Order: 0, Label: "BOM", Type: MenuItemTypeList, Reference: "bom", Children: []MenuItem{}, ExtensionKey: "a", IsHidden: true, ReplacedItem: &childOriginal}
	parentSnapshot := MenuItem{Id: "manufacturing", Order: 0, Label: "Manufacturing", Type: MenuItemTypeGroup, Children: []MenuItem{childTombstone}}
	materialized := []MenuItem{
		{Id: "manufacturing", Order: 0, Label: "Manufacturing", Type: MenuItemTypeGroup, Children: []MenuItem{}, ExtensionKey: "b", IsHidden: true, ReplacedItem: &parentSnapshot},
	}
	own := StripMenuConfigurationExtensionItems(materialized)
	if got := fmt.Sprint(menuIds(own)); got != "[manufacturing bom]" {
		t.Fatalf("unexpected own ids %s", got)
	}
	child := own[0].Children[0]
	if child.ExtensionKey != "" || child.IsHidden || child.ReplacedItem != nil {
		t.Errorf("nested tombstone must unwind to the host's own child, got %+v", child)
	}
}

func TestStripMenuConfigurationExtensionItems_RestoresDisplacedItemInItsSlot(t *testing.T) {
	// Host A,B,C; an extension replaced B (snapshotted with its pre-merge
	// order 1) and prepended X,Y, so the materialized tree is X0 Y1 A2 R3 C4.
	// The restored B must land in R's slot, not sort by its stale order.
	original := MenuItem{Id: "b", Order: 1, Type: MenuItemTypeList}
	materialized := []MenuItem{
		{Id: "x", Order: 0, Type: MenuItemTypeList, ExtensionKey: "e"},
		{Id: "y", Order: 1, Type: MenuItemTypeList, ExtensionKey: "e"},
		{Id: "a", Order: 2, Type: MenuItemTypeList},
		{Id: "r", Order: 3, Type: MenuItemTypeList, ExtensionKey: "e", ReplacedItem: &original},
		{Id: "c", Order: 4, Type: MenuItemTypeList},
	}
	own := StripMenuConfigurationExtensionItems(materialized)
	if got := fmt.Sprint(menuIds(own)); got != "[a b c]" {
		t.Fatalf("expected the host's own order, got %s", got)
	}
}

func TestCanonicalizeMenuItemOrder(t *testing.T) {
	items := []MenuItem{
		{Id: "c", Order: 20},
		{Id: "a", Order: 0, Children: []MenuItem{{Id: "a2", Order: 5}, {Id: "a1", Order: 1}}},
		{Id: "b", Order: 10},
		{Id: "b-dup", Order: 10},
	}
	CanonicalizeMenuItemOrder(items)
	if got := fmt.Sprint(menuIds(items)); got != "[a a1 a2 b b-dup c]" {
		t.Errorf("unexpected order %s", got)
	}
	for i, item := range items {
		if item.Order != i {
			t.Errorf("order must equal index, got %+v", item)
		}
	}
	snapshot, _ := json.Marshal(items)
	CanonicalizeMenuItemOrder(items)
	if again, _ := json.Marshal(items); string(again) != string(snapshot) {
		t.Error("canonical form must be a fixpoint")
	}
}

func TestMenuConfigurationExtension_RoundTrips(t *testing.T) {
	src := `{"key":"plus","menu_slug":"products","placements":[
	  {"anchor_id":"menu.items","position":"last","items":[{"id":"reports","order":0,"label":"Reports","type":"group","icon":"BarChart3","reference":"","children":[]}]},
	  {"anchor_id":"manufacturing","position":"remove"}]}`
	var extension MenuConfigurationExtension
	if err := json.Unmarshal([]byte(src), &extension); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if len(extension.Placements) != 2 || extension.Placements[1].Position != MenuConfigurationExtensionPositionRemove || extension.Placements[1].Items != nil {
		t.Errorf("placements lost data: %+v", extension.Placements)
	}
	raw, err := json.Marshal(extension)
	if err != nil {
		t.Fatal(err)
	}
	var again MenuConfigurationExtension
	if err := json.Unmarshal(raw, &again); err != nil {
		t.Fatal(err)
	}
	if again.Placements[0].Items[0].Id != "reports" || again.Placements[0].Items[0].ExtensionKey != "" {
		t.Errorf("round trip lost data: %+v", again)
	}
}

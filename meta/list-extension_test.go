package metamodel

import (
	"encoding/json"
	"testing"
)

func TestStripListExtensionColumnsAndActions(t *testing.T) {
	columns := []Column{
		{Attribute: "name"},
		{Attribute: "crm_score", ExtensionKey: "crm"},
		{Attribute: "status"},
	}
	actions := []PageAction{
		{Label: "Archive"},
		{Label: "Recalculate", ExtensionKey: "crm"},
	}

	ownColumns := StripListExtensionColumns(columns)
	ownActions := StripListExtensionActions(actions)

	if len(ownColumns) != 2 || ownColumns[0].Attribute != "name" || ownColumns[1].Attribute != "status" {
		t.Errorf("unexpected own columns %+v", ownColumns)
	}
	if len(ownActions) != 1 || ownActions[0].Label != "Archive" {
		t.Errorf("unexpected own actions %+v", ownActions)
	}
	if len(columns) != 3 || len(actions) != 2 {
		t.Fatal("strip must not mutate its inputs")
	}
	if StripListExtensionColumns(nil) == nil || StripListExtensionActions(nil) == nil {
		t.Fatal("expected empty, non-nil slices for nil input")
	}
}

func TestListExtension_RoundTrips(t *testing.T) {
	src := `{"key":"crm","list_slug":"l","columns":[
	  {"anchor_attribute":"status","position":"after","columns":[{"attribute":"crm_score","label":"Score","width":120}]},
	  {"position":"last","columns":[{"attribute":"crm_owner.name","label":"Owner","width":160}]}],
	  "actions":[{"label":"Recalculate","icon":"Sparkles","kind":"action","action":"recalc"}]}`
	var extension ListExtension
	if err := json.Unmarshal([]byte(src), &extension); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if len(extension.Columns) != 2 || extension.Columns[1].AnchorAttribute != "" || extension.Columns[1].Position != ExtensionPositionLast {
		t.Errorf("placements lost data: %+v", extension.Columns)
	}
	raw, err := json.Marshal(extension)
	if err != nil {
		t.Fatal(err)
	}
	var again ListExtension
	if err := json.Unmarshal(raw, &again); err != nil {
		t.Fatal(err)
	}
	if again.Columns[0].Columns[0].Attribute != "crm_score" || again.Actions[0].Action != "recalc" {
		t.Errorf("round trip lost data: %+v", again)
	}
}

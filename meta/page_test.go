package metamodel

import (
	"encoding/json"
	"testing"
)

func TestPageType_AllowsFields(t *testing.T) {
	cases := map[PageType]bool{
		"":               true, // normalizes to record
		PageTypeRecord:   true,
		PageTypeTrigger:  true,
		PageTypePlatform: false,
		PageTypeKiosk:    false,
		PageTypePublic:   false,
	}
	for pageType, want := range cases {
		if got := pageType.AllowsFields(); got != want {
			t.Errorf("%q.AllowsFields() = %v, want %v", pageType, got, want)
		}
	}
	// Trigger pages are standalone (no record) AND allow fields — orthogonal.
	if !PageTypeTrigger.IsStandalone() {
		t.Error("trigger pages must be standalone")
	}
}

func TestTriggerSchema_Attributes(t *testing.T) {
	schema := TriggerSchema{
		Params:  []Attribute{{Name: "outcome"}, {Name: "summary", Type: AttributeTypeString}},
		Returns: []Attribute{{Name: "summary", Type: AttributeTypeObject}, {Name: "total"}},
	}
	got := schema.Attributes()
	names := make([]string, 0, len(got))
	for _, attribute := range got {
		names = append(names, attribute.Name)
	}
	// Params first; a return colliding with a param resolves to the param.
	if len(names) != 3 || names[0] != "outcome" || names[1] != "summary" || names[2] != "total" {
		t.Fatalf("expected [outcome summary total], got %v", names)
	}
	if got[1].Type != AttributeTypeString {
		t.Errorf("the param must win a name collision, got %+v", got[1])
	}
}

// The binding maps are strict objects on the wire: the pre-binding string
// shape is rejected rather than silently read as a view / input binding.
func TestAppConfiguration_BindingsRejectStringShorthand(t *testing.T) {
	var configuration AppConfiguration
	if err := json.Unmarshal([]byte(`{"slug":"sales","record_pages":{"company":"company-sales"}}`), &configuration); err == nil {
		t.Fatal("expected record_pages string shorthand to be rejected")
	}
	if err := json.Unmarshal([]byte(`{"slug":"sales","action_pages":{"complete-meeting":"form"}}`), &configuration); err == nil {
		t.Fatal("expected action_pages string shorthand to be rejected")
	}
	if err := json.Unmarshal([]byte(`{"slug":"sales","record_pages":{"company":{"view":"company-sales","create":"company-new"}},"workflow_pages":{"onboard":{"running":"onboard-started"}}}`), &configuration); err != nil {
		t.Fatalf("object bindings must decode: %v", err)
	}
	if configuration.RecordPages["company"].ForMode(PageModeCreate) != "company-new" || configuration.RecordPages["company"].ForMode(PageModeEdit) != "" {
		t.Errorf("record binding not decoded: %+v", configuration.RecordPages)
	}
	if configuration.WorkflowPages["onboard"].ForPhase(TriggerPhaseRunning) != "onboard-started" || !(TriggerPageBinding{}).IsEmpty() {
		t.Errorf("trigger binding not decoded: %+v", configuration.WorkflowPages)
	}
}

package metamodel

import "testing"

func TestStripEntityExtensionAttributes_RemovesOnlyStampedAttributes(t *testing.T) {
	attributes := append(PlatformAttributes(),
		Attribute{Name: "name", Type: AttributeTypeString},
		Attribute{Name: "b_score", Type: AttributeTypeNumber, ExtensionKey: "widget-b"},
		Attribute{Name: "color", Type: AttributeTypeString},
		Attribute{Name: "c_flag", Type: AttributeTypeBoolean, ExtensionKey: "widget-c"},
	)

	own := StripEntityExtensionAttributes(attributes)

	want := []string{"id", "created_at", "updated_at", "created_by", "updated_by", "name", "color"}
	if len(own) != len(want) {
		t.Fatalf("expected %d attributes, got %d", len(want), len(own))
	}
	for i, name := range want {
		if own[i].Name != name {
			t.Errorf("position %d: expected %q, got %q", i, name, own[i].Name)
		}
		if IsEntityExtensionAttribute(own[i]) {
			t.Errorf("attribute %q must not carry a stamp after stripping", own[i].Name)
		}
	}
	// The input is untouched (the stamped attributes are still there).
	if len(attributes) != len(want)+2 {
		t.Fatal("StripEntityExtensionAttributes must not mutate its input")
	}
}

func TestStripEntityExtensionAttributes_NeverReturnsNil(t *testing.T) {
	if StripEntityExtensionAttributes(nil) == nil {
		t.Fatal("expected an empty, non-nil slice for nil input")
	}
}

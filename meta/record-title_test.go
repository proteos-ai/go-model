package metamodel

import "testing"

func TestValidateReservedAttributeNames(t *testing.T) {
	ok := []Attribute{{Name: "name", Type: AttributeTypeString}, {Name: "title", Type: AttributeTypeString}}
	if err := ValidateReservedAttributeNames(ok); err != nil {
		t.Fatalf("plain names must pass, got %v", err)
	}
	reserved := []Attribute{{Name: RecordTitleField, Type: AttributeTypeString}}
	if err := ValidateReservedAttributeNames(reserved); err == nil {
		t.Fatalf("%q must be rejected", RecordTitleField)
	}
	if err := ValidateAttributeDefinitions(reserved); err == nil {
		t.Fatalf("definition validation must reject %q at the top level", RecordTitleField)
	}
	// A leaf inside an object cannot collide with a top-level virtual column.
	nested := []Attribute{{
		Name: "meta", Type: AttributeTypeObject,
		Meta: map[string]any{"attributes": []map[string]any{{"name": RecordTitleField, "type": "string"}}},
	}}
	if err := ValidateAttributeDefinitions(nested); err != nil {
		t.Fatalf("nested leaf must pass, got %v", err)
	}
}

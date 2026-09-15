package metamodel

import (
	"encoding/json"
	"testing"
)

// The meta reaches production code as a map (the entity crossed JSONB + HTTP),
// as a struct (in-process construction) or as a pointer — ContactAddressMetaOf
// must read all three, exactly the trap PrincipalMetaOf documents.
func TestContactAddressMetaOfShapes(t *testing.T) {
	shapes := map[string]any{
		"struct":  ContactAddressAttributeMeta{Kind: "phone"},
		"pointer": &ContactAddressAttributeMeta{Kind: "phone"},
		"map":     map[string]any{"kind": "phone", "description": "mobile"},
	}
	for name, meta := range shapes {
		t.Run(name, func(t *testing.T) {
			got, ok := ContactAddressMetaOf(Attribute{Name: "mobile", Type: AttributeTypeContactAddress, Meta: meta})
			if !ok || got.Kind != "phone" {
				t.Fatalf("shape %s: got (%+v, %v), want kind phone", name, got, ok)
			}
		})
	}
	if _, ok := ContactAddressMetaOf(Attribute{Name: "note", Type: AttributeTypeString, Meta: map[string]any{"kind": "email"}}); ok {
		t.Fatalf("a non-contact-address attribute must not yield meta")
	}
	if _, ok := ContactAddressMetaOf(Attribute{Name: "email", Type: AttributeTypeContactAddress}); ok {
		t.Fatalf("missing meta must not yield meta")
	}
}

func TestContactAddressAttributesOrderAndPresence(t *testing.T) {
	attributes := []Attribute{
		{Name: "id", Type: AttributeTypeString},
		{Name: "email", Type: AttributeTypeContactAddress, Meta: map[string]any{"kind": "email"}},
		{Name: "name", Type: AttributeTypeString},
		{Name: "mobile", Type: AttributeTypeContactAddress, Meta: map[string]any{"kind": "phone"}},
	}
	got := ContactAddressAttributes(attributes)
	if len(got) != 2 || got[0].Name != "email" || got[1].Name != "mobile" {
		t.Fatalf("expected [email mobile] in attribute order, got %+v", got)
	}
	if !HasContactAddressAttributes(attributes) {
		t.Fatalf("presence must be detected")
	}
	if HasContactAddressAttributes(attributes[:1]) {
		t.Fatalf("no contact-address → false")
	}
}

func TestContactAddressJSONSchema(t *testing.T) {
	entity := Entity{Slug: "candidate", Name: "Candidate", Attributes: []Attribute{
		{Name: "email", Type: AttributeTypeContactAddress, Meta: map[string]any{"kind": "email"}},
		{Name: "mobile", Type: AttributeTypeContactAddress, Meta: map[string]any{"kind": "phone"}},
		{Name: "linkedin", Type: AttributeTypeContactAddress, Meta: map[string]any{"kind": "linkedin_public_identifier"}},
	}}
	schema := EntityToJSONSchema(entity)
	raw, _ := json.Marshal(schema.Properties)
	var properties map[string]map[string]any
	if err := json.Unmarshal(raw, &properties); err != nil {
		t.Fatal(err)
	}
	if properties["email"]["type"] != "string" || properties["email"]["format"] != "email" {
		t.Fatalf("email kind must export string + format email, got %v", properties["email"])
	}
	if properties["mobile"]["pattern"] != E164Pattern {
		t.Fatalf("phone kind must export the E.164 pattern, got %v", properties["mobile"])
	}
	if _, has := properties["linkedin"]["pattern"]; has || properties["linkedin"]["type"] != "string" {
		t.Fatalf("linkedin kind must stay a plain string, got %v", properties["linkedin"])
	}
}

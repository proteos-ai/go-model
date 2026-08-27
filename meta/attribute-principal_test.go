package metamodel

import (
	"encoding/json"
	"testing"

	"go.proteos.ai/model/common"
)

// decodeAttribute puts an attribute through the trip every entity schema really
// makes: marshalled into JSONB, served over HTTP, decoded into Attribute.Meta,
// which is `any` — so the meta arrives as map[string]any, NOT as the struct the
// producer built.
//
// Any test that only ever builds the struct in memory tests a shape production
// never sees. That is exactly how PrincipalMetaOf shipped broken.
func decodeAttribute(t *testing.T, attribute Attribute) Attribute {
	t.Helper()
	raw, err := json.Marshal(attribute)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	var decoded Attribute
	if err := json.Unmarshal(raw, &decoded); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if _, isMap := decoded.Meta.(map[string]any); !isMap && decoded.Meta != nil {
		t.Fatalf("expected the decoded meta to be a map, got %T — this helper no longer models the real path", decoded.Meta)
	}
	return decoded
}

func grantingAttribute() Attribute {
	return Attribute{
		Name: "account_manager",
		Type: AttributeTypePrincipal,
		Meta: PrincipalAttributeMeta{
			Grants:       []string{"read", "write"},
			AllowedTypes: []common.PrincipalType{common.PrincipalTypePerson, common.PrincipalTypeTeam},
		},
	}
}

// THE REGRESSION THIS PINS: PrincipalMetaOf type-switched on the struct, so it
// answered false for every attribute that had crossed JSON — which is every
// attribute on a real request. Granting attributes projected nothing,
// ValidateGrantingAttributeWrites short-circuited before its escalation guard,
// and the allowed_types narrowing never ran. Phase 5 was inert.
func TestPrincipalMetaOfReadsEveryShape(t *testing.T) {
	value := grantingAttribute()
	pointerMeta := PrincipalAttributeMeta{Grants: []string{"read", "write"}}
	pointer := Attribute{Name: "account_manager", Type: AttributeTypePrincipal, Meta: &pointerMeta}

	shapes := map[string]Attribute{
		"struct":  value,
		"pointer": pointer,
		"decoded": decodeAttribute(t, value),
	}

	for name, attribute := range shapes {
		meta, ok := PrincipalMetaOf(attribute)
		if !ok {
			t.Fatalf("%s: PrincipalMetaOf returned not-ok", name)
		}
		if !meta.IsGranting() {
			t.Errorf("%s: grants lost, got %v", name, meta.Grants)
		}
		if !HasGrantingAttributes([]Attribute{attribute}) {
			t.Errorf("%s: HasGrantingAttributes = false", name)
		}
	}
}

// allowed_types must narrow after decoding too.
func TestAllowedTypesNarrowAfterDecoding(t *testing.T) {
	decoded := decodeAttribute(t, grantingAttribute())

	meta, ok := PrincipalMetaOf(decoded)
	if !ok {
		t.Fatal("PrincipalMetaOf returned not-ok on a decoded attribute")
	}
	if !meta.AllowsType(common.PrincipalTypeTeam) {
		t.Error("team should be allowed")
	}
	if meta.AllowsType(common.PrincipalTypeOrg) {
		t.Error("org must stay excluded after decoding")
	}
}

// An attribute that confers `share` — what "owner" means now — refuses `org`
// whatever allowed_types says, in every shape. An org-held `share` would hand
// management rights to every scoped member of the organization.
func TestShareGrantingAttributeRefusesOrgAfterDecoding(t *testing.T) {
	owner := Attribute{
		Name: "account_manager",
		Type: AttributeTypePrincipal,
		// No allowed_types at all: "every kind" — except org, because of share.
		Meta: &PrincipalAttributeMeta{Grants: []string{"read", "write", "delete", "share"}},
	}
	for name, attribute := range map[string]Attribute{"pointer": owner, "decoded": decodeAttribute(t, owner)} {
		meta, ok := PrincipalMetaOf(attribute)
		if !ok {
			t.Fatalf("%s: PrincipalMetaOf returned not-ok", name)
		}
		if !meta.ConfersShare() {
			t.Errorf("%s: share verb lost", name)
		}
		if meta.AllowsType(common.PrincipalTypeOrg) {
			t.Errorf("%s: an attribute conferring share must refuse org", name)
		}
		for _, allowed := range []common.PrincipalType{
			common.PrincipalTypePerson, common.PrincipalTypeAgent, common.PrincipalTypeApi, common.PrincipalTypeTeam,
		} {
			if !meta.AllowsType(allowed) {
				t.Errorf("%s: should accept %q", name, allowed)
			}
		}
	}

	// A plain read-granting attribute may still name the org: that is how
	// "visible to everyone" is expressed on a field.
	visible := PrincipalAttributeMeta{Grants: []string{"read"}}
	if !visible.AllowsType(common.PrincipalTypeOrg) {
		t.Error("a read-only granting attribute must accept org")
	}
	if visible.ConfersShare() {
		t.Error("read alone must not confer share")
	}
}

// An ordinary principal attribute confers nothing — adding `principal` to an
// entity must not start granting access by accident, in any shape.
func TestNonGrantingPrincipalAttributeConfersNothing(t *testing.T) {
	plain := Attribute{Name: "mentioned_by", Type: AttributeTypePrincipal, Meta: PrincipalAttributeMeta{}}

	for name, attribute := range map[string]Attribute{"struct": plain, "decoded": decodeAttribute(t, plain)} {
		if meta, ok := PrincipalMetaOf(attribute); ok && meta.IsGranting() {
			t.Errorf("%s: a meta with no grants must not be granting", name)
		}
		if HasGrantingAttributes([]Attribute{attribute}) {
			t.Errorf("%s: HasGrantingAttributes = true", name)
		}
	}
}

// A non-principal attribute is never principal meta, whatever it carries.
func TestPrincipalMetaOfRejectsOtherTypes(t *testing.T) {
	if _, ok := PrincipalMetaOf(Attribute{Name: "stage", Type: AttributeTypeString}); ok {
		t.Error("a string attribute must not yield principal meta")
	}
	// Absent meta on a principal attribute is not an error, just nothing to read.
	if _, ok := PrincipalMetaOf(Attribute{Name: "owner", Type: AttributeTypePrincipal}); ok {
		t.Error("a principal attribute with no meta must yield not-ok")
	}
}

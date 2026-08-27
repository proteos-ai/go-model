package metamodel

import (
	"testing"
)

func TestPlatformAttributes_CanonicalSet(t *testing.T) {
	attrs := PlatformAttributes()
	want := []string{"id", "created_at", "updated_at", "created_by", "updated_by"}
	if len(attrs) != len(want) {
		t.Fatalf("expected %d platform attributes, got %d", len(want), len(attrs))
	}
	for i, name := range want {
		if attrs[i].Name != name {
			t.Errorf("attr[%d]: want name %q, got %q", i, name, attrs[i].Name)
		}
		if !attrs[i].IsPlatformManaged {
			t.Errorf("attr %q: expected IsPlatformManaged=true", attrs[i].Name)
		}
	}
}

// Every platform attribute is server-managed and read-only. There is no
// writable platform attribute any more: ownership moved to the business field
// (a principal attribute granting `share`), so the one exception is gone.
func TestPlatformAttributes_AllReadOnly(t *testing.T) {
	for _, attr := range PlatformAttributes() {
		if !attr.IsReadOnly {
			t.Errorf("attr %q: expected IsReadOnly=true", attr.Name)
		}
	}
}

func TestEnsurePlatformAttributes_AddsWhenMissing(t *testing.T) {
	got := EnsurePlatformAttributes([]Attribute{
		{Name: "email", Type: AttributeTypeString},
	})

	if len(got) != 6 {
		t.Fatalf("expected 5 platform + 1 user = 6 attributes, got %d", len(got))
	}
	// Platform attributes come first, in canonical order.
	for i, name := range []string{"id", "created_at", "updated_at", "created_by", "updated_by"} {
		if got[i].Name != name {
			t.Errorf("position %d: want %q, got %q", i, name, got[i].Name)
		}
	}
	if got[5].Name != "email" {
		t.Errorf("user attribute should follow platform attributes, got %q", got[5].Name)
	}
}

// THE RETIREMENT PIN. `owned_by` was persisted into every stored schema as a
// platform attribute. Once it stops being one, an unstripped entry would come
// back as a USER attribute — principal-typed, no column behind it — which is
// the worst of both worlds. It must vanish, in whatever shape it arrives.
func TestEnsurePlatformAttributes_StripsRetiredOwnedBy(t *testing.T) {
	got := EnsurePlatformAttributes([]Attribute{
		{Name: "owned_by", Type: AttributeTypePrincipal, IsPlatformManaged: true},
		{Name: "email", Type: AttributeTypeString},
	})
	for _, attr := range got {
		if attr.Name == "owned_by" {
			t.Fatalf("retired owned_by survived: %+v", got)
		}
	}
	if len(got) != 6 || got[5].Name != "email" {
		t.Fatalf("expected 5 platform + email, got %+v", got)
	}
	if !IsRetiredPlatformAttributeName("owned_by") || IsRetiredPlatformAttributeName("email") {
		t.Error("IsRetiredPlatformAttributeName is wrong")
	}
}

func TestEnsurePlatformAttributes_OverridesClientRedefinition(t *testing.T) {
	// A client attempts to redefine `id` as a mutable number and drop the rest.
	got := EnsurePlatformAttributes([]Attribute{
		{Name: "id", Type: AttributeTypeNumber, IsReadOnly: false, IsPlatformManaged: false},
		{Name: "email", Type: AttributeTypeString},
	})

	if len(got) != 6 {
		t.Fatalf("expected the canonical 5 + email, got %d", len(got))
	}
	id := got[0]
	if id.Name != "id" || id.Type != AttributeTypeString || !id.IsReadOnly || !id.IsPlatformManaged {
		t.Errorf("client redefinition of `id` was not overridden: %+v", id)
	}
	// `email` survives.
	if got[5].Name != "email" {
		t.Errorf("user attribute `email` was dropped: %+v", got)
	}
}

func TestIsPlatformAttributeName(t *testing.T) {
	for _, name := range []string{"id", "created_at", "updated_at", "created_by", "updated_by"} {
		if !IsPlatformAttributeName(name) {
			t.Errorf("%q should be a platform attribute name", name)
		}
	}
	for _, name := range []string{"email", "name", "createdAt", "owned_by", ""} {
		if IsPlatformAttributeName(name) {
			t.Errorf("%q should NOT be a platform attribute name", name)
		}
	}
}

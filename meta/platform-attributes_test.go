package metamodel

import "testing"

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
		if !attrs[i].IsReadOnly {
			t.Errorf("attr %q: expected IsReadOnly=true", attrs[i].Name)
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
	for _, name := range []string{"email", "name", "createdAt", ""} {
		if IsPlatformAttributeName(name) {
			t.Errorf("%q should NOT be a platform attribute name", name)
		}
	}
}

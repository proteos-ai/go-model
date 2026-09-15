package conversationmodel

import "testing"

func TestCanonicalizeEmail(t *testing.T) {
	cases := []struct {
		name   string
		raw    string
		want   string
		wantOk bool
	}{
		{"lowercases and trims", "  Ada.Lovelace@Example.COM ", "ada.lovelace@example.com", true},
		{"keeps dots (no gmail magic)", "a.d.a@gmail.com", "a.d.a@gmail.com", true},
		{"keeps plus tags", "ada+news@example.com", "ada+news@example.com", true},
		{"rejects empty", "", "", false},
		{"rejects no at", "not-an-email", "", false},
		{"rejects leading at", "@example.com", "", false},
		{"rejects trailing at", "ada@", "", false},
	}
	for _, testCase := range cases {
		t.Run(testCase.name, func(t *testing.T) {
			got, ok := CanonicalizeEmail(testCase.raw)
			if ok != testCase.wantOk || got != testCase.want {
				t.Fatalf("CanonicalizeEmail(%q) = (%q, %v), want (%q, %v)",
					testCase.raw, got, ok, testCase.want, testCase.wantOk)
			}
		})
	}
}

func TestCanonicalizePhone(t *testing.T) {
	cases := []struct {
		name   string
		raw    string
		want   string
		wantOk bool
	}{
		{"passes e164", "+4917612345678", "+4917612345678", true},
		{"strips separators", "+49 176 / 123-45.678", "+4917612345678", true},
		{"rejects national format", "0176 12345678", "", false},
		{"rejects letters", "+49abc", "", false},
		{"rejects too short", "+1234", "", false},
		{"rejects too long", "+1234567890123456", "", false},
	}
	for _, testCase := range cases {
		t.Run(testCase.name, func(t *testing.T) {
			got, ok := CanonicalizePhone(testCase.raw)
			if ok != testCase.wantOk || got != testCase.want {
				t.Fatalf("CanonicalizePhone(%q) = (%q, %v), want (%q, %v)",
					testCase.raw, got, ok, testCase.want, testCase.wantOk)
			}
		})
	}
}

func TestCanonicalizeLinkedin(t *testing.T) {
	cases := []struct {
		name   string
		raw    string
		want   string
		wantOk bool
	}{
		{"bare slug", "ada-lovelace", "ada-lovelace", true},
		{"lowercases", "Ada-Lovelace", "ada-lovelace", true},
		{"full url with query and trailing slash", "https://www.linkedin.com/in/Ada-Lovelace/?trk=public", "ada-lovelace", true},
		{"url without scheme", "linkedin.com/in/ada-lovelace", "ada-lovelace", true},
		{"regional host", "https://de.linkedin.com/in/ada-lovelace", "ada-lovelace", true},
		{"path only", "/in/ada-lovelace/", "ada-lovelace", true},
		{"in prefix without slashes", "in/ada-lovelace", "ada-lovelace", true},
		{"fragment", "https://www.linkedin.com/in/ada-lovelace#about", "ada-lovelace", true},
		{"rejects classic profile id", "ACoAAB1234xyz", "", false},
		{"rejects sales navigator profile id", "ACwAAB1234xyz", "", false},
		{"rejects recruiter profile id", "AEwAAB1234xyz", "", false},
		{"rejects profile id inside a profile url", "https://www.linkedin.com/in/ACoAAAcDMMQBODyLwZrRcgYhrkCafURGqva0U4E/", "", false},
		{"lowercase ae slug stays a slug", "aerial-photo", "aerial-photo", true},
		{"lowercase ae slug in a url stays a slug", "linkedin.com/in/aerial-photo", "aerial-photo", true},
		{"refuses a slug typed with a profile-id prefix casing", "linkedin.com/in/AErial-photo", "", false},
		{"rejects company page", "https://www.linkedin.com/company/proteos", "", false},
		{"rejects foreign host", "https://example.com/in/ada", "", false},
		{"rejects host only", "https://www.linkedin.com/", "", false},
		{"rejects empty", "   ", "", false},
		{"rejects leading hyphen", "-ada", "", false},
		{"rejects too short", "ab", "", false},
	}
	for _, testCase := range cases {
		t.Run(testCase.name, func(t *testing.T) {
			got, ok := CanonicalizeLinkedin(testCase.raw)
			if ok != testCase.wantOk || got != testCase.want {
				t.Fatalf("CanonicalizeLinkedin(%q) = (%q, %v), want (%q, %v)",
					testCase.raw, got, ok, testCase.want, testCase.wantOk)
			}
		})
	}
}

func TestCanonicalizeContactAddress(t *testing.T) {
	cases := []struct {
		name   string
		kind   ContactAddressKind
		raw    string
		want   string
		wantOk bool
	}{
		{"email dispatches", ContactAddressKindEmail, " Ada@Example.com", "ada@example.com", true},
		{"phone dispatches", ContactAddressKindPhone, "+49 176 12345678", "+4917612345678", true},
		{"linkedin public identifier dispatches", ContactAddressKindLinkedinPublicIdentifier, "https://www.linkedin.com/in/ada", "ada", true},
		{"linkedin profile id is not a record kind", ContactAddressKindLinkedinId, "ACoAAB1234xyz", "", false},
		{"legacy linkedin alias is not a record kind", ContactAddressKindLinkedin, "ada", "", false},
		{"unsupported kind", ContactAddressKindSlack, "U0123ABC", "", false},
		{"empty kind", "", "ada@example.com", "", false},
	}
	for _, testCase := range cases {
		t.Run(testCase.name, func(t *testing.T) {
			got, ok := CanonicalizeContactAddress(testCase.kind, testCase.raw)
			if ok != testCase.wantOk || got != testCase.want {
				t.Fatalf("CanonicalizeContactAddress(%q, %q) = (%q, %v), want (%q, %v)",
					testCase.kind, testCase.raw, got, ok, testCase.want, testCase.wantOk)
			}
		})
	}
}

func TestIsRecordContactAddressKind(t *testing.T) {
	for _, kind := range RecordContactAddressKinds {
		if !IsRecordContactAddressKind(kind) {
			t.Fatalf("%q must be a record kind", kind)
		}
	}
	for _, kind := range []ContactAddressKind{ContactAddressKindSlack, ContactAddressKindMessenger, ContactAddressKindWhatsapp,
		ContactAddressKindLinkedinId, ContactAddressKindLinkedin, ""} {
		if IsRecordContactAddressKind(kind) {
			t.Fatalf("%q must not be a record kind", kind)
		}
	}
}

func TestIsLinkedinProviderId(t *testing.T) {
	cases := []struct {
		value string
		want  bool
	}{
		{"ACoAAAcDMMQBODyLwZrRcgYhrkCafURGqva0U4E", true},
		{"ACwAAB1234-_xyz", true},
		{"AEwAAB1234xyz", true},
		{"  ACoAAB1234xyz ", true},
		{"ACo", false},
		{"ada-lovelace", false},
		{"aerial-photo", false},
		{"https://www.linkedin.com/in/ACoAAB1234xyz", false},
		{"ACo AAB", false},
		{"", false},
	}
	for _, testCase := range cases {
		if got := IsLinkedinProviderId(testCase.value); got != testCase.want {
			t.Fatalf("IsLinkedinProviderId(%q) = %v, want %v", testCase.value, got, testCase.want)
		}
	}
}

func TestContactAddressKindCanonical(t *testing.T) {
	if got := ContactAddressKindLinkedin.Canonical(); got != ContactAddressKindLinkedinId {
		t.Fatalf("legacy linkedin canonicalizes to %q, want linkedin_id", got)
	}
	for _, kind := range []ContactAddressKind{ContactAddressKindLinkedinId, ContactAddressKindLinkedinPublicIdentifier,
		ContactAddressKindEmail, ContactAddressKindSlack} {
		if got := kind.Canonical(); got != kind {
			t.Fatalf("%q canonicalizes to %q, want itself", kind, got)
		}
	}
	if aliases := LegacyContactAddressKinds(ContactAddressKindLinkedinId); len(aliases) != 1 || aliases[0] != ContactAddressKindLinkedin {
		t.Fatalf("linkedin_id aliases = %v, want [linkedin]", aliases)
	}
	if aliases := LegacyContactAddressKinds(ContactAddressKindLinkedinPublicIdentifier); aliases != nil {
		t.Fatalf("linkedin_public_identifier must have no aliases, got %v", aliases)
	}
}

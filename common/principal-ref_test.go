package common

import "testing"

// The entry grammar is a discipline surface: it is safe only while this package
// is the sole builder and parser, and while slugs reject the separator. These
// tests are what pin that.

func TestPrincipalRefString(t *testing.T) {
	tests := []struct {
		name  string
		ref   PrincipalRef
		orgId string
		want  string
	}{
		// person / agent / api share one namespace, so they render identically.
		{"person", PrincipalRef{PrincipalTypePerson, "u-1"}, "acme", "user:u-1"},
		{"agent", PrincipalRef{PrincipalTypeAgent, "u-2"}, "acme", "user:u-2"},
		{"api client", PrincipalRef{PrincipalTypeApi, "u-3"}, "acme", "user:u-3"},

		{"team is org-qualified", PrincipalRef{PrincipalTypeTeam, "sales"}, "acme", "team:acme/sales"},
		{"org names itself", PrincipalRef{PrincipalTypeOrg, "acme"}, "acme", "org:acme"},

		// Fail closed rather than open: an unscoped team renders to nothing, so
		// it matches nothing. The alternative — a bare slug — would match a
		// same-named team in another tenant.
		{"team without an org renders empty", PrincipalRef{PrincipalTypeTeam, "sales"}, "", ""},
		{"empty id renders empty", PrincipalRef{PrincipalTypePerson, ""}, "acme", ""},
		{"unknown type renders empty", PrincipalRef{PrincipalType("wat"), "x"}, "acme", ""},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if got := test.ref.String(test.orgId); got != test.want {
				t.Fatalf("String(%q) = %q, want %q", test.orgId, got, test.want)
			}
		})
	}
}

func TestPrincipalRefEntry(t *testing.T) {
	ref := PrincipalRef{PrincipalTypeTeam, "sales"}
	if got := ref.Entry("acme", "read"); got != "team:acme/sales#read" {
		t.Fatalf("Entry = %q, want team:acme/sales#read", got)
	}
	// A principal that cannot be rendered must not produce a half-formed entry
	// like "#read", which would be a wildcard-shaped string in an array compare.
	if got := ref.Entry("", "read"); got != "" {
		t.Fatalf("Entry with no org = %q, want empty", got)
	}
	if got := ref.Entry("acme", ""); got != "" {
		t.Fatalf("Entry with no permission = %q, want empty", got)
	}
}

func TestParseEntryRoundTrips(t *testing.T) {
	cases := []struct {
		ref        PrincipalRef
		orgId      string
		permission string
		principal  string
	}{
		{PrincipalRef{PrincipalTypePerson, "u-1"}, "acme", "read", "user:u-1"},
		{PrincipalRef{PrincipalTypeTeam, "sales-emea"}, "acme", "write", "team:acme/sales-emea"},
		{PrincipalRef{PrincipalTypeOrg, "acme"}, "acme", "delete", "org:acme"},
	}

	for _, test := range cases {
		entry := test.ref.Entry(test.orgId, test.permission)
		principal, permission, err := ParseEntry(entry)
		if err != nil {
			t.Fatalf("ParseEntry(%q): %v", entry, err)
		}
		if principal != test.principal {
			t.Errorf("principal = %q, want %q", principal, test.principal)
		}
		if permission != test.permission {
			t.Errorf("permission = %q, want %q", permission, test.permission)
		}
	}
}

func TestParseEntryRejectsMalformedInput(t *testing.T) {
	for _, entry := range []string{
		"",
		"user:u-1",  // no separator
		"#read",     // no principal
		"user:u-1#", // no permission
		"#",         // neither
	} {
		if _, _, err := ParseEntry(entry); err == nil {
			t.Errorf("ParseEntry(%q) should have failed", entry)
		}
	}
}

// A `#` inside a team slug would make an entry ambiguous to parse, so it is
// rejected at the type rather than left for each storage path to remember.
func TestPrincipalRefIsValid(t *testing.T) {
	tests := []struct {
		name string
		ref  PrincipalRef
		want bool
	}{
		{"person", PrincipalRef{PrincipalTypePerson, "u-1"}, true},
		{"agent", PrincipalRef{PrincipalTypeAgent, "u-1"}, true},
		{"api", PrincipalRef{PrincipalTypeApi, "u-1"}, true},
		{"team", PrincipalRef{PrincipalTypeTeam, "sales"}, true},
		{"org", PrincipalRef{PrincipalTypeOrg, "acme"}, true},

		{"empty id", PrincipalRef{PrincipalTypePerson, ""}, false},
		{"unknown type", PrincipalRef{PrincipalType("group"), "x"}, false},
		{"team slug containing the separator", PrincipalRef{PrincipalTypeTeam, "sa#les"}, false},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if got := test.ref.IsValid(); got != test.want {
				t.Fatalf("IsValid = %v, want %v", got, test.want)
			}
		})
	}
}

// Why a separator in a slug is rejected, stated precisely.
//
// ParseEntry itself survives one: it splits on the LAST separator, so
// "team:acme/sa#les#read" still yields ("team:acme/sa#les", "read"). The
// validation is not covering for a broken parser.
//
// It exists because the grammar is only unambiguous under last-separator
// parsing, and the obvious wrong implementation — splitting on the separator —
// is one keystroke away and fails SILENTLY: it yields a principal that is a
// prefix of the real one, which in an array-overlap comparison simply matches
// nothing, so access is quietly lost rather than loudly broken. Forbidding the
// character at the type means no future parser can be wrong about it.
func TestSeparatorInASlugIsRejectedToKeepTheGrammarUnambiguous(t *testing.T) {
	ref := PrincipalRef{PrincipalTypeTeam, "sa#les"}
	if ref.IsValid() {
		t.Fatal("a slug containing the separator must be rejected")
	}

	entry := ref.Entry("acme", "read")

	// The canonical parser copes.
	principal, permission, err := ParseEntry(entry)
	if err != nil {
		t.Fatalf("ParseEntry: %v", err)
	}
	if principal != "team:acme/sa#les" || permission != "read" {
		t.Fatalf("ParseEntry(%q) = (%q, %q), want (team:acme/sa#les, read)", entry, principal, permission)
	}

	// The naive one does not — and this is the failure the validation forecloses.
	naive := splitOnFirstSeparator(entry)
	if naive == principal {
		t.Fatal("expected a first-separator split to disagree with the canonical parse")
	}
	if naive != "team:acme/sa" {
		t.Fatalf("first-separator split = %q, want team:acme/sa", naive)
	}
}

func splitOnFirstSeparator(entry string) string {
	for index := range len(entry) {
		if string(entry[index]) == PermissionSeparator {
			return entry[:index]
		}
	}
	return entry
}

// UserType's values are a subset of PrincipalType's, which is what lets a
// principal attribute be seeded from created_by without a conversion table.
func TestPrincipalRefFromUser(t *testing.T) {
	tests := []struct {
		user UserRef
		want PrincipalRef
	}{
		{UserRef{UserTypePerson, "u-1"}, PrincipalRef{PrincipalTypePerson, "u-1"}},
		{UserRef{UserTypeAgent, "u-2"}, PrincipalRef{PrincipalTypeAgent, "u-2"}},
		{UserRef{UserTypeApi, "u-3"}, PrincipalRef{PrincipalTypeApi, "u-3"}},
	}

	for _, test := range tests {
		got := PrincipalRefFromUser(test.user)
		if got != test.want {
			t.Errorf("PrincipalRefFromUser(%+v) = %+v, want %+v", test.user, got, test.want)
		}
		if !got.IsValid() {
			t.Errorf("a converted UserRef must be a valid principal: %+v", got)
		}
	}
}

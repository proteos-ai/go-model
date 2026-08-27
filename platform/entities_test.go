package platform

import "testing"

func TestEntities_CanonicalSet(t *testing.T) {
	entities := Entities()
	if len(entities) != 49 {
		t.Fatalf("expected 49 platform entities, got %d", len(entities))
	}
	seen := make(map[string]bool, len(entities))
	for _, entity := range entities {
		if entity.Slug == "" || entity.Name == "" {
			t.Errorf("entity %+v: Slug and Name must be set", entity)
		}
		if seen[entity.Slug] {
			t.Errorf("duplicate slug %q", entity.Slug)
		}
		seen[entity.Slug] = true
	}
}

func TestSlugs_MatchEntities(t *testing.T) {
	if len(Slugs()) != len(Entities()) {
		t.Fatalf("Slugs() and Entities() must have equal length")
	}
	want := []string{
		"organizations", "users", "roles", "user-role-assignments", "role-entity-permissions",
		"teams", "team-members",
		"entities", "pages", "menu-configurations", "apps", "components", "lists",
		"list-views", "design-references", "modules", "variables", "deployments", "files",
		"hooks", "actions",
		"workflows", "workflow-executions",
		"knowledge-nodes", "knowledge-links", "knowledge-labels", "knowledge-spaces",
		"agents", "prompts", "skills", "tools", "toolsets", "mcp-servers", "agent-sessions",
		"topics", "events",
		"connections", "conversations", "messages", "agent-listeners", "transcriptions",
		"glossary-terms", "mistranscribed-terms", "contacts", "conversation-filters",
		"conversation-types", "contact-groups", "tone-profiles", "connectors",
	}
	got := Slugs()
	if len(got) != len(want) {
		t.Fatalf("expected %d slugs, got %d", len(want), len(got))
	}
	for i, slug := range want {
		if got[i] != slug {
			t.Errorf("slug[%d]: want %q, got %q", i, slug, got[i])
		}
	}
}

func TestIsReserved(t *testing.T) {
	for _, slug := range []string{"users", "files", "list-views", "actions"} {
		if !IsReserved(slug) {
			t.Errorf("%q should be reserved", slug)
		}
	}
	for _, slug := range []string{"customer", "invoice", "Users", "ava-threads", ""} {
		if IsReserved(slug) {
			t.Errorf("%q should NOT be reserved", slug)
		}
	}
}

func TestIsReservedTableName(t *testing.T) {
	if !IsReservedTableName(AccessGrantsTable) {
		t.Fatalf("IsReservedTableName(%q) must be true", AccessGrantsTable)
	}
	// A physical-table claim is not a permission target: the two sets must not
	// overlap, or the table name would surface in the role-permission dropdown.
	if IsReserved(AccessGrantsTable) {
		t.Fatalf("%q must not be a platform ENTITY slug", AccessGrantsTable)
	}
	for _, slug := range Slugs() {
		if IsReservedTableName(slug) {
			t.Fatalf("platform entity slug %q must not also be a reserved table name", slug)
		}
	}
	if IsReservedTableName("deals") {
		t.Fatal("ordinary slugs must not be reserved table names")
	}
}

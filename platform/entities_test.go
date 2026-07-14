package platform

import "testing"

func TestEntities_CanonicalSet(t *testing.T) {
	entities := Entities()
	if len(entities) != 39 {
		t.Fatalf("expected 39 platform entities, got %d", len(entities))
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
		"entities", "pages", "menu-configurations", "apps", "components", "lists",
		"list-views", "design-references", "modules", "variables", "deployments", "files",
		"hooks", "actions",
		"workflows", "workflow-executions",
		"knowledge-nodes", "knowledge-links", "knowledge-labels",
		"agents", "prompts", "skills", "tools", "mcp-servers", "agent-sessions",
		"topics", "events",
		"connections", "conversations", "messages", "agent-listeners", "transcriptions",
		"glossary-terms", "connectors",
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

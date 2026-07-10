// Package platform is the single source of truth for the platform's built-in
// "platform entities" — system resources (users, roles, files, entities, …) that
// are permission targets alongside user-defined entities. It is deliberately
// dependency-neutral (no account/meta coupling) so every service can consume it:
//   - account-service grants the admin role rights over all of them on bootstrap
//   - metadata-service reserves their slugs (users can't create colliding entities)
//   - the frontend mirrors this list to populate the role-permission dropdown
package platform

// PlatformEntity is a built-in system resource that can be a permission target.
// It is referenced by Slug everywhere (permissions key on a free-form entity_slug
// string, so platform and user entities are uniform). Name is the display label
// for the role-permission dropdown.
type PlatformEntity struct {
	Slug string
	Name string
}

// platformEntities is the canonical, complete set of platform entities. Keep this
// list in sync with the TS mirror in packages/sdk-ts (auth/platform-entities.ts).
var platformEntities = []PlatformEntity{
	// Access (account-service)
	{Slug: "organizations", Name: "Organizations"},
	{Slug: "users", Name: "Users"},
	{Slug: "roles", Name: "Roles"},
	{Slug: "user-role-assignments", Name: "User Role Assignments"},
	{Slug: "role-entity-permissions", Name: "Role Entity Permissions"},
	// Schema / content (metadata-service)
	{Slug: "entities", Name: "Entities"},
	{Slug: "pages", Name: "Pages"},
	{Slug: "menu-configurations", Name: "Menu Configurations"},
	{Slug: "apps", Name: "Apps"},
	{Slug: "components", Name: "Components"},
	{Slug: "lists", Name: "Lists"},
	{Slug: "list-views", Name: "List Views"},
	// System (metadata-service)
	{Slug: "modules", Name: "Modules"},
	{Slug: "variables", Name: "Variables"},
	{Slug: "deployments", Name: "Deployments"},
	// Storage (storage-service)
	{Slug: "files", Name: "Files"},
	// Automation (function-service)
	{Slug: "hooks", Name: "Hooks"},
	{Slug: "actions", Name: "Actions"},
	// Knowledge (knowledge-service)
	{Slug: "knowledge-nodes", Name: "Knowledge Nodes"},
	{Slug: "knowledge-links", Name: "Knowledge Links"},
	{Slug: "knowledge-labels", Name: "Knowledge Labels"},
	// Agent suite (agent-service)
	{Slug: "agents", Name: "Agents"},
	{Slug: "prompts", Name: "Prompts"},
	{Slug: "skills", Name: "Skills"},
	{Slug: "tools", Name: "Tools"},
	{Slug: "mcp-servers", Name: "MCP Servers"},
	{Slug: "agent-sessions", Name: "Agent Sessions"},
	// Messaging bus (event-service)
	{Slug: "topics", Name: "Topics"},
	{Slug: "events", Name: "Events"},
	// Conversations (conversation-service).
	// NOTE: `connections` is enforced by BOTH conversation-service and
	// connector-service (2026-07 decision) — a connection is the same concept
	// in both (an installed integration instance holding credentials), and
	// conversation connections migrate onto connector-service later. A role
	// granted connections rights therefore governs both services.
	{Slug: "connections", Name: "Connections"},
	{Slug: "conversations", Name: "Conversations"},
	{Slug: "messages", Name: "Messages"},
	{Slug: "agent-listeners", Name: "Agent Listeners"},
	{Slug: "transcriptions", Name: "Transcriptions"},
	// Connectors (connector-service). `connections` above is shared; this is
	// the manifest catalog.
	{Slug: "connectors", Name: "Connectors"},
}

// reservedSlugs is the membership set behind IsReserved.
var reservedSlugs = func() map[string]struct{} {
	set := make(map[string]struct{}, len(platformEntities))
	for _, entity := range platformEntities {
		set[entity.Slug] = struct{}{}
	}
	return set
}()

// Entities returns the canonical platform entities (a fresh copy), in display
// order, for callers that need the display metadata (e.g. the role-permission
// dropdown).
func Entities() []PlatformEntity {
	out := make([]PlatformEntity, len(platformEntities))
	copy(out, platformEntities)
	return out
}

// Slugs returns just the platform entity slugs, in canonical order — used by the
// admin bootstrap to grant the admin role rights over every platform entity.
func Slugs() []string {
	out := make([]string, len(platformEntities))
	for i, entity := range platformEntities {
		out[i] = entity.Slug
	}
	return out
}

// IsReserved reports whether a slug belongs to a platform entity. Slugs are
// canonical kebab-case lowercase, so the match is exact. Used to reject
// user-defined entities that would collide with a platform entity.
func IsReserved(slug string) bool {
	_, ok := reservedSlugs[slug]
	return ok
}

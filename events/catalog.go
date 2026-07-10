package eventmodel

import "strings"

// displayNames maps known platform topic names to friendly labels. The set of
// topics is ALWAYS discovered live from Redis (SCAN) — this map only supplies
// nicer labels on top. record.<entity>.events is dynamic (one stream per
// entity), so it is matched structurally below rather than enumerated here.
//
// Catalog sources are fragmented across the codebase by design:
//   - organization.events  → accountmodel (Global; excluded from per-org views)
//   - record.<entity>.events → datamodel (PerOrg, dynamic per entity slug)
//   - hooks.events / actions.events → function-service-private
//
// We do not import those (it would invert dependency direction); we re-declare
// the bare names here purely for display.
var displayNames = map[string]string{
	"organization.events": "Organization Events",
	"hooks.events":        "Hook Events",
	"actions.events":      "Action Events",
}

// DisplayNameFor returns a friendly label for a topic name, falling back to the
// raw name when the topic is unknown. The ".dlq" suffix is ignored (a
// dead-letter stream shares its source topic's label; the UI distinguishes them
// via Topic.Kind). record.<entity>.events becomes "<Entity> Record Events".
func DisplayNameFor(name string) string {
	base := strings.TrimSuffix(name, ".dlq")
	if label, ok := displayNames[base]; ok {
		return label
	}
	if strings.HasPrefix(base, "record.") && strings.HasSuffix(base, ".events") {
		entity := strings.TrimSuffix(strings.TrimPrefix(base, "record."), ".events")
		if entity != "" {
			return titleizeSlug(entity) + " Record Events"
		}
	}
	return name
}

// titleizeSlug turns a kebab/snake/dot entity slug into space-separated,
// capitalized words ("purchase-order" → "Purchase Order").
func titleizeSlug(slug string) string {
	replacer := strings.NewReplacer("-", " ", "_", " ", ".", " ")
	words := strings.Fields(replacer.Replace(slug))
	for i, word := range words {
		if word == "" {
			continue
		}
		words[i] = strings.ToUpper(word[:1]) + word[1:]
	}
	return strings.Join(words, " ")
}

package eventmodel

import "strings"

// displayNames maps known platform topic names to friendly labels. Live topic
// DISCOVERY still happens against Redis — this map only supplies nicer labels
// on top, for both the live view and the catalog. record.<entity>.events is
// dynamic (one stream per entity), so it is matched structurally below rather
// than enumerated here.
//
// The topics themselves are declared in go.proteos.ai/events/topics, which this
// package deliberately does NOT import: eventmodel is a leaf the SDK-facing
// shapes live in, and importing the transport's topic catalogue would invert
// the dependency direction. The bare names are re-declared here purely for
// display; event-service performs the actual registry merge.
var displayNames = map[string]string{
	"organization.events":         "Organization Events",
	"hooks.events":                "Hook Events",
	"actions.events":              "Action Events",
	"message.events":              "Message Events",
	"reaction.events":             "Reaction Events",
	"conversation.events":         "Conversation Events",
	"conversation.domain.events":  "Conversation Milestones",
	"contact.events":              "Contact Events",
	"connection.events":           "Connection Events",
	"connector.events":            "Connector Events",
	"connector.inbound.events":    "Connector Inbound Events",
	"file-version-content.events": "File Content Events",
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

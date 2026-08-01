package eventmodel

import "time"

// CatalogTopic is one entry of the event catalog: a topic a caller may
// subscribe to, together with the event types it carries. It merges the static
// platform registry (which knows a topic's event types before anyone has
// published to it) with what is live in Redis right now (which knows about
// topics nobody declared).
//
// Distinct from Topic, which is purely a live projection of one Redis stream.
// Topic answers "what is on the bus"; CatalogTopic answers "what can I listen
// to, and for which events".
type CatalogTopic struct {
	// Name is the logical topic name, e.g. "message.events".
	Name string `json:"name"`
	// DisplayName is a friendly label, falling back to Name.
	DisplayName string `json:"display_name"`
	// Description explains what the topic carries. Empty for live-discovered
	// topics that are not in the static registry.
	Description string `json:"description,omitempty"`
	// EventTypes are the event types carried on the topic. Canonical for
	// registry topics; sampled from recent events for live-only ones, and so
	// possibly incomplete (or empty when the sample found nothing).
	EventTypes []string `json:"event_types"`
	// IsLive reports whether a stream currently exists for the caller's org. A
	// registry topic nobody has published to yet is still catalogued, with
	// IsLive false — it remains a valid subscription target.
	IsLive bool `json:"is_live"`
	// EventCount is the live stream's XLEN (retention window only); 0 when the
	// topic is not live.
	EventCount int64 `json:"event_count"`
	// LastEventAt is the dispatch time of the newest entry, when the topic is
	// live and non-empty.
	LastEventAt *time.Time `json:"last_event_at,omitempty"`
}

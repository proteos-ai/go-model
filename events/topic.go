// Package eventmodel holds the read-only projections of the Redis-Streams
// message bus that event-service exposes: Topic and Event. Nothing here is
// persisted — topics and messages are discovered by reading Redis directly, so
// every field is derived at read time. This package is deliberately neutral (no
// event-bus-adapter, no service coupling) so the SDK-facing JSON shapes live in
// one place.
package eventmodel

import "time"

// TopicKind classifies a discovered stream.
type TopicKind string

const (
	// TopicKindEvent is a normal event stream (e.g. record.contact.events).
	TopicKindEvent TopicKind = "event"
	// TopicKindDeadLetter is a dead-letter stream (the ".dlq" suffix) holding
	// poison-pill messages a consumer group exhausted its retries on.
	TopicKindDeadLetter TopicKind = "dead_letter"
)

// Topic is a read-only projection of a single Redis Streams key scoped to one
// org. Topics are not owned by any service — they are discovered by SCAN-ing
// the bus for the caller's org, so all fields are derived from Redis at read
// time and reflect only the live retention window (~72h / 1M entries).
type Topic struct {
	// Name is the logical topic: the stream key with the "events" prefix and
	// the {org} hashtag stripped — e.g. "record.contact.events" or, for a
	// dead-letter stream, "record.contact.events.dlq".
	Name string `json:"name"`
	// DisplayName is a friendly label from the static catalog, or Name when the
	// topic is not a known platform topic.
	DisplayName string `json:"display_name"`
	// Kind is "event" or "dead_letter", classified by the ".dlq" suffix.
	Kind TopicKind `json:"kind"`
	// SourceTopic is the base topic a dead-letter stream belongs to (Name with
	// the ".dlq" suffix removed); empty for event topics. Drives DLQ redrive.
	SourceTopic string `json:"source_topic,omitempty"`
	// EventCount is the current XLEN of the stream (entries in the retention
	// window, not all-time).
	EventCount int64 `json:"event_count"`
	// LastEventAt is the dispatch time of the newest entry, when the stream is
	// non-empty.
	LastEventAt *time.Time `json:"last_event_at,omitempty"`
	// ConsumerGroups are the consumer groups registered on the stream, for an
	// at-a-glance view of who reads it and how far behind they are.
	ConsumerGroups []ConsumerGroup `json:"consumer_groups"`
}

// ConsumerGroup is a read-only view of a Redis Streams consumer group on a
// topic (from XINFO GROUPS).
type ConsumerGroup struct {
	Name      string `json:"name"`
	Consumers int64  `json:"consumers"`
	Pending   int64  `json:"pending"`
	Lag       int64  `json:"lag"`
}

// RedriveResult reports the outcome of moving dead-letter entries back to their
// source topic for re-processing.
type RedriveResult struct {
	// SourceTopic is the event topic the entries were redriven to.
	SourceTopic string `json:"source_topic"`
	// Redriven is the number of messages moved off the dead-letter stream.
	Redriven int `json:"redriven"`
}

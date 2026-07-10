package eventmodel

import (
	"encoding/json"
	"time"
)

// Event is a read-only projection of a single Redis Streams entry. It
// reconstructs the events.Event envelope every proteos producer writes
// (id/org/topic/type/key/payload/dispatched_at plus h.<name> headers), plus the
// Redis entry id used for cursoring and live tail.
type Event struct {
	// StreamId is the Redis Streams entry id ("<ms>-<seq>"). It is the cursor
	// for paging (?before=) and the resume point for live tail — prefer it over
	// DispatchedAt, which is producer-supplied and not strictly monotonic.
	StreamId string `json:"stream_id"`
	// Id is the application message id (UUID) carried in the envelope.
	Id string `json:"id"`
	// OrgId is the org the event is about.
	OrgId string `json:"org_id"`
	// Topic is the logical topic the entry was read from.
	Topic string `json:"topic"`
	// Type names the event within the topic (e.g. "contact.created").
	Type string `json:"type"`
	// Key is the partition key (e.g. the record id); may be empty.
	Key string `json:"key,omitempty"`
	// Payload is the opaque event body. Producers marshal a struct, so it is
	// JSON in practice and surfaced verbatim for the UI to pretty-print. When a
	// stream somehow holds non-JSON bytes the adapter wraps them as a JSON
	// string so this stays a valid JSON value.
	Payload json.RawMessage `json:"payload"`
	// Headers are the envelope headers (h.<name> fields), e.g. entity_slug.
	Headers map[string]string `json:"headers,omitempty"`
	// DispatchedAt is the producer-stamped dispatch time (from dispatched_at).
	DispatchedAt time.Time `json:"dispatched_at"`
}

// Package eventapi holds the request/response DTOs for event-service's
// HTTP API — the wire contracts the SDK and web client bind to.
package eventapi

import (
	"encoding/json"

	eventmodel "go.proteos.ai/model/events"
)

// GetEventsQuery is the query for reading messages off a topic.
type GetEventsQuery struct {
	// Limit caps how many of the most-recent messages to return (newest first).
	// Zero falls back to the service default.
	Limit int `json:"limit" form:"limit"`
	// Before is a Redis stream-id cursor; only entries strictly older than it
	// are returned, for paging back through history.
	Before string `json:"before" form:"before"`
	// Follow switches the response to a chunked NDJSON live tail (one JSON
	// message per line) instead of a single JSON array.
	Follow bool `json:"follow" form:"follow"`
}

// PublishEventRequest publishes a test message onto a topic. The message is
// written with the standard platform envelope and IS delivered to real
// consumers — a privileged admin/testing action, not a no-op.
type PublishEventRequest struct {
	// Type names the event within the topic (e.g. "contact.created"). Required.
	Type string `json:"type" validate:"required"`
	// Key is an optional partition key (e.g. a record id).
	Key string `json:"key"`
	// Payload is the event body — any JSON value. Surfaced to consumers verbatim.
	Payload json.RawMessage `json:"payload"`
}

// GetTopicsResponse wraps the topic list.
type GetTopicsResponse struct {
	Data []eventmodel.Topic `json:"data"`
}

// GetEventsResponse wraps a (non-follow) message list.
type GetEventsResponse struct {
	Data []eventmodel.Event `json:"data"`
}

package agentmodel

import "encoding/json"

// Content blocks are the self-tagged pieces that survive where a payload is
// genuinely heterogeneous: message bodies (text | file) and tool results
// (text | image | json). Single-content events stay FLAT (see session-payload.go);
// only these slices carry an inner `type`. Shapes mirror example-message-model.json.

// TextBlock is a plain text block. Agent messages are text-only.
type TextBlock struct {
	Type string `json:"type"` // always "text"
	Text string `json:"text"`
}

// FileBlock is a file value object: a reference to a stored File (storage-service)
// plus light metadata. It is embedded inline in a user message ContentBlock
// (type=file) and nested under `image` in a tool ResultBlock.
type FileBlock struct {
	FileId      string  `json:"file_id,omitempty"`
	FileName    *string `json:"file_name,omitempty"`
	MimeType    *string `json:"mime_type,omitempty"`
	SizeInBytes *int64  `json:"size_in_bytes,omitempty"`
}

// ContentBlock is one block of a user/agent message body. For type=text only
// Text is set; for type=file the embedded FileBlock fields are promoted inline
// (file_id, file_name, …) — matching the wire shape in example-message-model.json.
type ContentBlock struct {
	Type      string `json:"type"` // text | file
	Text      string `json:"text,omitempty"`
	FileBlock        // inline: file_id, file_name, mime_type, size_in_bytes
}

// ResultBlock is one block of a tool result. Unlike message file blocks, an
// image nests its file under `image`; json carries an arbitrary raw object.
type ResultBlock struct {
	Type  string          `json:"type"` // text | image | json
	Text  string          `json:"text,omitempty"`
	Image *FileBlock      `json:"image,omitempty"`
	Json  json.RawMessage `json:"json,omitempty"`
}

// TurnRole labels a Turn as originating from the user or the agent.
type TurnRole string

const (
	TurnRoleUser  TurnRole = "user"
	TurnRoleAgent TurnRole = "agent"
)

// Turn is the groupByTurn() projection — NEVER stored, recomputed from the event
// log. Events hold the SessionEvent envelopes verbatim (no reshaping), grouped by
// turn_id and ordered by seq.
type Turn struct {
	Id     string         `json:"id"`
	Role   TurnRole       `json:"role"`
	Events []SessionEvent `json:"events"`
}

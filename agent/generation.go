package agentmodel

import "strings"

// GenerationRequest is a single, stateless model call: send content blocks, get
// one message back. It is the synchronous, no-persistence counterpart to a Session
// — no event log, no tools, no multi-turn state. The same ContentBlock vocabulary
// the session log uses (text | file) carries the request body, so callers compose
// a generation from the exact pieces they already build for user messages.
type GenerationRequest struct {
	// Model selects the model id + sampling knobs. An empty ModelId resolves to the
	// service default at call time.
	Model ModelConfig `json:"model"`
	// System is an optional system prompt steering the response.
	System string `json:"system,omitempty"`
	// Content is the message body the model answers — text and file blocks.
	Content []ContentBlock `json:"content"`
}

// ModelStopReason is why the model stopped generating. An open string so a new
// provider reason ("end_turn" | "max_tokens" | "stop_sequence" | …) never breaks
// decoding — the same forward-compatible stance as StopReason on a session. The
// caller's key distinction is end_turn (complete) vs max_tokens (truncated).
type ModelStopReason string

// GenerationResult is the assistant's reply to a GenerationRequest: the returned
// text blocks, why generation stopped, the model that produced it, and token
// usage. There is no role — a generation always returns the assistant message.
type GenerationResult struct {
	Content    []TextBlock     `json:"content"`
	StopReason ModelStopReason `json:"stop_reason,omitempty"`
	ModelId    string          `json:"model_id"`
	Usage      ModelUsage      `json:"usage"`
}

// Text concatenates the text of all blocks — the common case where the caller
// wants the reply as a single string rather than walking the block slice.
func (result GenerationResult) Text() string {
	parts := make([]string, 0, len(result.Content))
	for _, block := range result.Content {
		parts = append(parts, block.Text)
	}
	return strings.Join(parts, "")
}

package storageapi

// ParseFileRequest is the optional body of POST /files/:id/parse. An empty body
// parses the current version with the auto-selected backend; callers can pin a
// parser ("model" | "document_ai" | "agent") and force a re-parse of cached content.
type ParseFileRequest struct {
	Parser string `json:"parser,omitempty"`
	Force  bool   `json:"force,omitempty"`
}

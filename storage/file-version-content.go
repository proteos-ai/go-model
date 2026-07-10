package storagemodel

import "time"

// ParseStatus tracks a version's content extraction. A row is created pending,
// flips to ready with the extracted content, or to failed with an error message.
type ParseStatus string

const (
	ParseStatusPending    ParseStatus = "pending"
	ParseStatusProcessing ParseStatus = "processing"
	ParseStatusReady      ParseStatus = "ready"
	ParseStatusFailed     ParseStatus = "failed"
)

// FileVersionContent is the parsed/extracted content of a single file version — the
// cache produced by POST /files/:id/parse and read via GET /files/:id/content. It is
// distinct from the file's raw bytes (the /download endpoint): "content" here always
// means the extracted text/markdown, never the original blob.
type FileVersionContent struct {
	Id        string      `json:"id"`
	VersionId string      `json:"version_id"`
	FileId    string      `json:"file_id"`
	OrgId     *string     `json:"org_id,omitempty"`
	Format    string      `json:"format"` // markdown | text
	Content   string      `json:"content"`
	Parser    string      `json:"parser,omitempty"` // model | document_ai | agent
	Model     string      `json:"model,omitempty"`  // provenance: model id / processor id
	Status    ParseStatus `json:"status"`
	Error     *string     `json:"error,omitempty"`
	ParsedAt  *time.Time  `json:"parsed_at,omitempty"`
	CreatedAt time.Time   `json:"created_at"`
	UpdatedAt time.Time   `json:"updated_at"`
}

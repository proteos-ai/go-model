package storagemodel

import "time"

// FileVersionContentEvent is the message payload published on
// FileVersionContentEvents when a file_version_content row is created or updated.
// It mirrors FileVersionContent WITHOUT the (potentially large) content body —
// consumers that need the body fetch it from storage's GET /files/:id/content.
type FileVersionContentEvent struct {
	Id        string      `json:"id"`
	VersionId string      `json:"version_id"`
	FileId    string      `json:"file_id"`
	OrgId     *string     `json:"org_id,omitempty"`
	Format    string      `json:"format"`
	Parser    string      `json:"parser,omitempty"`
	Model     string      `json:"model,omitempty"`
	Status    ParseStatus `json:"status"`
	Error     *string     `json:"error,omitempty"`
	ParsedAt  *time.Time  `json:"parsed_at,omitempty"`
	CreatedAt time.Time   `json:"created_at"`
	UpdatedAt time.Time   `json:"updated_at"`
}

// NewFileVersionContentEvent projects a FileVersionContent onto the event payload,
// dropping the content body.
func NewFileVersionContentEvent(content FileVersionContent) FileVersionContentEvent {
	return FileVersionContentEvent{
		Id:        content.Id,
		VersionId: content.VersionId,
		FileId:    content.FileId,
		OrgId:     content.OrgId,
		Format:    content.Format,
		Parser:    content.Parser,
		Model:     content.Model,
		Status:    content.Status,
		Error:     content.Error,
		ParsedAt:  content.ParsedAt,
		CreatedAt: content.CreatedAt,
		UpdatedAt: content.UpdatedAt,
	}
}

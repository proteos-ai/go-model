package conversationmodel

import (
	"time"

	"go.proteos.ai/model/common"
)

// Attachment is one file that arrived with a message — stored as its own
// entity (table attachment). File is the platform FileRef ({id, name} —
// bytes live in storage-service; size and content type are resolved from
// storage by id on read, per the FileRef convention). Message.Attachments is
// the read-time projection of these rows (same pattern as reactions); the
// attachment is NEVER inlined into message content.
//
// ContentId carries the MIME Content-ID for inline parts (cid: references in
// the html body), so a renderer can later resolve embedded images to their
// stored files. IsInline distinguishes those from "real" attachments a user
// would download.
type Attachment struct {
	Id           string         `json:"id"`
	OrgId        string         `json:"org_id"`
	MessageId    string         `json:"message_id"`
	ConnectionId string         `json:"connection_id"`
	File         common.FileRef `json:"file"`
	ContentId    string         `json:"content_id,omitempty"`
	IsInline     bool           `json:"is_inline"`
	CreatedAt    time.Time      `json:"created_at"`
	CreatedBy    common.UserRef `json:"created_by"`
}

// NormalizedAttachment is the connector→domain wire shape for one attachment
// riding on a NormalizedInboundMessage: decoded bytes plus the MIME identity
// (the mime type travels here because the storage upload needs it; it is not
// persisted on the Attachment row — FileRef convention). The domain uploads
// the bytes to storage-service and persists an Attachment row; connectors
// never talk to storage themselves.
type NormalizedAttachment struct {
	FileName  string `json:"file_name"`
	MimeType  string `json:"mime_type"`
	ContentId string `json:"content_id,omitempty"`
	IsInline  bool   `json:"is_inline,omitempty"`
	Content   []byte `json:"content,omitempty"`
}

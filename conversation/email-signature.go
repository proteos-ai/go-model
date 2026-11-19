package conversationmodel

import (
	"time"

	"go.proteos.ai/model/common"
)

// EmailSignature is the block an email send adds below what the person wrote
// and above the quote. It is never part of a stored message body: the send
// adds it on the wire, and the message keeps a MessageSignature snapshot.
//
// A signature belongs to one person (OwnedBy) or, with OwnedBy nil, to the
// organisation, and can be tied to one mailbox (ConnectionId). Among the
// defaults a send picks the most specific: the sending mailbox's, then the
// sender's own, then the organisation's.
type EmailSignature struct {
	Id    string `json:"id"`
	OrgId string `json:"org_id"`
	Name  string `json:"name" sortable:""`
	// Html is the rich part, sanitized like every stored email html block.
	Html string `json:"html"`
	// Text is the plain-text part. Derived from Html unless IsTextCustom.
	Text         string `json:"text"`
	IsTextCustom bool   `json:"is_text_custom"`
	// OwnedBy is the person the signature belongs to; nil for the organisation.
	OwnedBy *common.UserRef `json:"owned_by,omitempty"`
	// ConnectionId ties the signature to one mailbox; empty for every mailbox.
	ConnectionId string `json:"connection_id,omitempty"`
	// IsDefault: a send that names no signature may pick this one.
	IsDefault bool                    `json:"is_default" sortable:""`
	AppliesTo EmailSignatureAppliesTo `json:"applies_to"`
	CreatedAt time.Time               `json:"created_at" sortable:""`
	CreatedBy common.UserRef          `json:"created_by"`
	UpdatedAt time.Time               `json:"updated_at" sortable:""`
	UpdatedBy common.UserRef          `json:"updated_by"`
}

// MessageSignature is a signature as one email carries it. On a sent message
// it is the snapshot of what went out, so a later edit of the signature never
// rewrites history. On a draft read it is what a send would add now.
type MessageSignature struct {
	Id   string `json:"id"`
	Name string `json:"name"`
	Html string `json:"html"`
	Text string `json:"text"`
}

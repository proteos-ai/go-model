package conversationapi

import (
	"go.proteos.ai/model/common"
	conversationmodel "go.proteos.ai/model/conversation"
)

// SendMessageRequest sends an outbound message. Three addressing modes —
// exactly one, validated by logic.ValidateSendAddressing, not struct tags:
//   - reply: ConversationId set → send into that existing conversation's thread.
//   - thread: ReplyToMessageId set → reply targeted at ONE message. On channels
//     with a second thread level (Slack) this starts/continues the thread rooted
//     at that message (a child conversation carrying parent_conversation_id +
//     started_by_message_id); on email it is an anchored reply into the same
//     conversation (In-Reply-To/References point at the target instead of
//     thread-latest); channels with no message-targeted reply reject with
//     threads_not_supported.
//   - originate: neither set → ConnectionId + at least one To recipient
//     required; a new conversation is created (or an existing one reused when
//     the connector can resolve the destination key up front) and the connector
//     starts a fresh thread. Each recipient is a polymorphic target (kind
//     room|participant + the connector-side external id — a Slack channel or
//     user id, an email address). To may carry MULTIPLE participants (a group
//     thread; Slack opens an MPIM) or exactly one room. Cc/Bcc are email-only.
//     Subject is used by email. Routing always resolves conversation →
//     connection → connector server-side.
type SendMessageRequest struct {
	ConversationId   string                           `json:"conversation_id"`
	ReplyToMessageId string                           `json:"reply_to_message_id"`
	ConnectionId     string                           `json:"connection_id"`
	To               []conversationmodel.Recipient    `json:"to"`
	Cc               []conversationmodel.Recipient    `json:"cc,omitempty"`
	Bcc              []conversationmodel.Recipient    `json:"bcc,omitempty"`
	Subject          string                           `json:"subject"`
	Content          []conversationmodel.ContentBlock `json:"content"`
	// Attachments are files to send WITH the message. The client pre-uploads
	// each file to storage-service and references it here; the server downloads
	// the bytes under the caller's token and hands them to the channel
	// connector. Name is the display/file name (falls back to the id when
	// empty). Content may be empty when attachments are present — "at least
	// text OR attachments" is enforced by logic.ValidateSendContent.
	Attachments []common.FileRef `json:"attachments,omitempty"`
}

// UpdateDraftRequest edits a status=draft message — full replacement of the
// editable surface (PUT semantics), the addressing MODE stays immutable (a
// reply draft cannot become an originate draft). To/Cc/Bcc replace the draft's
// recipients PRESENCE-based: any group present in the JSON (including an
// explicit empty `[]` — a clear) replaces the stored set from exactly the
// provided groups; all three absent keeps the stored recipients. Email
// contact-address drafts only — room drafts address their venue, not people.
// Subject applies only while the draft's conversation is itself status=draft
// (mirroring Send, where subject is originate-only); Attachments replace the
// draft's attachment rows (files pre-uploaded to storage like
// SendMessageRequest).
type UpdateDraftRequest struct {
	To          []conversationmodel.Recipient    `json:"to,omitempty"`
	Cc          []conversationmodel.Recipient    `json:"cc,omitempty"`
	Bcc         []conversationmodel.Recipient    `json:"bcc,omitempty"`
	Subject     string                           `json:"subject"`
	Content     []conversationmodel.ContentBlock `json:"content"`
	Attachments []common.FileRef                 `json:"attachments,omitempty"`
}

// GetManyMessagesQuery pages a conversation's messages; list order is
// occurred_at (when it was said), not created_at.
type GetManyMessagesQuery struct {
	Direction *string `json:"direction" form:"direction" db:"direction"`
	Status    *string `json:"status" form:"status" db:"status"`
	common.Pagination
	common.Sorting
}

type GetManyMessagesResponse struct {
	Meta common.ResponseMeta         `json:"meta"`
	Data []conversationmodel.Message `json:"data"`
}

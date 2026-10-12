package conversationmodel

import (
	"time"

	"go.proteos.ai/model/common"
)

// ConversationNote is what a PERSON recorded about a conversation — during a
// call, after a meeting, while reading a thread: "mentioned a budget freeze
// until Q2", "wants the proposal by Friday". An append log: one row per
// note, stamped with its author, never merged — several people may note the
// same conversation. Notes are the human's half of conversation
// intelligence: the summary reads them to learn what the human found
// important, and summary-basis tag evaluations see them next to the summary.
// A note is as visible as the conversation it is on; only its author edits
// or deletes it.
type ConversationNote struct {
	Id             string `json:"id"`
	OrgId          string `json:"org_id"`
	ConversationId string `json:"conversation_id"`
	// Content is the note text (markdown allowed).
	Content   string         `json:"content"`
	CreatedAt time.Time      `json:"created_at" sortable:""`
	CreatedBy common.UserRef `json:"created_by"`
	UpdatedAt time.Time      `json:"updated_at" sortable:""`
	UpdatedBy common.UserRef `json:"updated_by"`
}

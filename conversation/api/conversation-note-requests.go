package conversationapi

import (
	"go.proteos.ai/model/common"
	conversationmodel "go.proteos.ai/model/conversation"
)

// CreateConversationNoteRequest adds one note to a conversation (POST
// /conversations/:id/notes). The author is the caller; the conversation
// comes from the path.
type CreateConversationNoteRequest struct {
	Content string `json:"content" validate:"required,max=20000"`
}

// UpdateConversationNoteRequest rewrites a note's text (PATCH
// /conversation-notes/:id) — the author's door only (403
// conversation_note_not_author otherwise).
type UpdateConversationNoteRequest struct {
	Content string `json:"content" validate:"required,max=20000"`
}

// GetManyConversationNotesQuery pages a conversation's notes (oldest first
// unless sorted).
type GetManyConversationNotesQuery struct {
	common.Pagination
	common.Sorting
}

type GetManyConversationNotesResponse struct {
	Meta common.ResponseMeta                  `json:"meta"`
	Data []conversationmodel.ConversationNote `json:"data"`
}

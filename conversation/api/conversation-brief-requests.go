package conversationapi

import (
	"go.proteos.ai/model/common"
	conversationmodel "go.proteos.ai/model/conversation"
)

// CreateConversationBriefRequest prepares guidance for a conversation with
// one or more contacts. Any number of briefs may be open for the same people
// at once. TypeKey must name an existing ConversationType when set; Channel
// must be a known channel when set.
type CreateConversationBriefRequest struct {
	Subject    string                    `json:"subject"`
	ContactIds []string                  `json:"contact_ids" validate:"required,min=1"`
	Channel    conversationmodel.Channel `json:"channel"`
	TypeKey    string                    `json:"type_key"`
	Content    string                    `json:"content" validate:"required"`
}

// UpdateConversationBriefRequest revises an open brief (PATCH semantics —
// absent fields keep their value; ContactIds replaces the set wholesale when
// present). Only a status=prepared brief is editable (409
// conversation_brief_not_editable otherwise).
type UpdateConversationBriefRequest struct {
	Subject    *string                    `json:"subject,omitempty"`
	ContactIds *[]string                  `json:"contact_ids,omitempty"`
	Channel    *conversationmodel.Channel `json:"channel,omitempty"`
	Content    *string                    `json:"content,omitempty"`
	TypeKey    *string                    `json:"type_key,omitempty"`
}

// AttachConversationBriefRequest binds an open brief to a conversation
// explicitly (the after-the-fact door; the softphone attaches through the
// call-start callback instead).
type AttachConversationBriefRequest struct {
	ConversationId string `json:"conversation_id" validate:"required"`
}

// GetManyConversationBriefsQuery filters the org's briefs. ContactId = briefs
// naming that person (containment); ConversationId = the brief a
// conversation was held against.
type GetManyConversationBriefsQuery struct {
	ContactId      *string `json:"contact_id" form:"contact_id"`
	ConversationId *string `json:"conversation_id" form:"conversation_id" db:"conversation_id"`
	Channel        *string `json:"channel" form:"channel" db:"channel"`
	TypeKey        *string `json:"type_key" form:"type_key" db:"type_key"`
	Status         *string `json:"status" form:"status" db:"status"`
	common.Pagination
	common.Sorting
}

type GetManyConversationBriefsResponse struct {
	Meta common.ResponseMeta                   `json:"meta"`
	Data []conversationmodel.ConversationBrief `json:"data"`
}

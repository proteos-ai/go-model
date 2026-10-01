package conversationapi

import (
	"time"

	"go.proteos.ai/model/common"
	conversationmodel "go.proteos.ai/model/conversation"
)

// CreateConversationTagRequest is the manual door (POST
// /conversations/:id/tags): a person tags a conversation, a message or a span.
// Anchors follow the definition's scope (conversation: none; message:
// start == end message, no offsets; span: code-point offsets into the
// message's first text block, end exclusive). The service stamps source=user,
// confidence 1, and fills Quote from the anchored text.
type CreateConversationTagRequest struct {
	DefinitionKey  string `json:"definition_key" validate:"required,max=255"`
	Value          string `json:"value" validate:"required,max=255"`
	Why            string `json:"why" validate:"max=1024"`
	StartMessageId string `json:"start_message_id" validate:"max=36"`
	StartOffset    *int   `json:"start_offset,omitempty"`
	EndMessageId   string `json:"end_message_id" validate:"max=36"`
	EndOffset      *int   `json:"end_offset,omitempty"`
}

// UpdateConversationTagRequest flips a tag between active and dismissed (the
// feedback loop) and lets a person annotate why.
type UpdateConversationTagRequest struct {
	Status *conversationmodel.ConversationTagStatus `json:"status,omitempty"`
	Why    *string                                  `json:"why,omitempty" validate:"omitempty,max=1024"`
}

// EvaluateConversationTagsRequest starts a manual evaluation (POST
// /conversations/:id/tags/evaluate, 202). Empty DefinitionKeys = every
// definition the currently matching sets reference (and every model tag on
// the conversation is replaced); a list narrows both the questions and the
// replacement to those keys.
type EvaluateConversationTagsRequest struct {
	DefinitionKeys []string `json:"definition_keys"`
}

// GetManyConversationTagsQuery serves both GET /conversations/:id/tags (the
// controller forces ConversationId from the path) and the org-wide GET
// /conversation-tags.
type GetManyConversationTagsQuery struct {
	ConversationId *string `json:"conversation_id" form:"conversation_id" db:"conversation_id"`
	DefinitionKey  *string `json:"definition_key" form:"definition_key" db:"definition_key"`
	Value          *string `json:"value" form:"value" db:"value"`
	Channel        *string `json:"channel" form:"channel" db:"channel"`
	Source         *string `json:"source" form:"source" db:"source"`
	Status         *string `json:"status" form:"status" db:"status"`
	// MessageId matches tags anchored on the message at either end (no db tag:
	// OR across two columns, repository-applied).
	MessageId *string `json:"message_id" form:"message_id"`
	// ContactId filters to conversations whose roster contains the contact
	// (EXISTS on conversation.participants, repository-applied).
	ContactId *string `json:"contact_id" form:"contact_id"`
	// Since / Until bound observed_at (inclusive / exclusive).
	Since *time.Time `json:"since" form:"since" time_format:"2006-01-02T15:04:05Z07:00"`
	Until *time.Time `json:"until" form:"until" time_format:"2006-01-02T15:04:05Z07:00"`
	common.Pagination
	common.Sorting
}

type GetManyConversationTagsResponse struct {
	Meta common.ResponseMeta                 `json:"meta"`
	Data []conversationmodel.ConversationTag `json:"data"`
}

// CountConversationTagsQuery groups active tags for dashboards. GroupBy is a
// comma-separated whitelist of definition_key, value, channel (order kept).
type CountConversationTagsQuery struct {
	GroupBy       string     `json:"group_by" form:"group_by"`
	DefinitionKey *string    `json:"definition_key" form:"definition_key"`
	Channel       *string    `json:"channel" form:"channel"`
	Source        *string    `json:"source" form:"source"`
	Status        *string    `json:"status" form:"status"`
	Since         *time.Time `json:"since" form:"since" time_format:"2006-01-02T15:04:05Z07:00"`
	Until         *time.Time `json:"until" form:"until" time_format:"2006-01-02T15:04:05Z07:00"`
}

type CountConversationTagsResponse struct {
	Data []conversationmodel.ConversationTagCount `json:"data"`
}

type GetManyConversationTagEvaluationsQuery struct {
	Trigger *string `json:"trigger" form:"trigger" db:"trigger"`
	Status  *string `json:"status" form:"status" db:"status"`
	common.Pagination
	common.Sorting
}

type GetManyConversationTagEvaluationsResponse struct {
	Meta common.ResponseMeta                           `json:"meta"`
	Data []conversationmodel.ConversationTagEvaluation `json:"data"`
}

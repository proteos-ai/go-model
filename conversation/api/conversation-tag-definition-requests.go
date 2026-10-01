package conversationapi

import (
	"go.proteos.ai/model/common"
	conversationmodel "go.proteos.ai/model/conversation"
)

// CreateConversationTagDefinitionRequest creates one tag definition. Key is
// the immutable identity (kebab/snake, like every platform key); Question is
// the decide-shaped question the evaluator answers; Scope and Evaluator say
// where answers anchor and which model path answers. IsEnabled nil ⇒ true
// (manifests may omit it). ModuleSlug is stamped by `pro module deploy` —
// interactive callers leave it empty.
type CreateConversationTagDefinitionRequest struct {
	Key           string                                     `json:"key" validate:"required,max=255"`
	Name          string                                     `json:"name" validate:"max=255"`
	Group         string                                     `json:"group" validate:"max=255"`
	Question      conversationmodel.ConversationTagQuestion  `json:"question"`
	Scope         conversationmodel.ConversationTagScope     `json:"scope" validate:"required"`
	Evaluator     conversationmodel.ConversationTagEvaluator `json:"evaluator"`
	Directions    []conversationmodel.MessageDirection       `json:"directions"`
	MinConfidence float64                                    `json:"min_confidence" validate:"min=0,max=1"`
	Color         string                                     `json:"color" validate:"max=64"`
	IsEnabled     *bool                                      `json:"is_enabled,omitempty"`
	ModuleSlug    string                                     `json:"module_slug" validate:"max=255"`
}

// UpdateConversationTagDefinitionRequest is a partial update — nil leaves the
// stored value untouched. Question, when present, replaces the stored question
// wholesale but must keep its Type (400 conversation_tag_question_type_immutable
// otherwise). Directions replaces wholesale. Key is immutable.
type UpdateConversationTagDefinitionRequest struct {
	Name          *string                                     `json:"name,omitempty" validate:"omitempty,max=255"`
	Group         *string                                     `json:"group,omitempty" validate:"omitempty,max=255"`
	Question      *conversationmodel.ConversationTagQuestion  `json:"question,omitempty"`
	Scope         *conversationmodel.ConversationTagScope     `json:"scope,omitempty"`
	Evaluator     *conversationmodel.ConversationTagEvaluator `json:"evaluator,omitempty"`
	Directions    *[]conversationmodel.MessageDirection       `json:"directions,omitempty"`
	MinConfidence *float64                                    `json:"min_confidence,omitempty" validate:"omitempty,min=0,max=1"`
	Color         *string                                     `json:"color,omitempty" validate:"omitempty,max=64"`
	IsEnabled     *bool                                       `json:"is_enabled,omitempty"`
}

type GetManyConversationTagDefinitionsQuery struct {
	// ModuleSlug filters to definitions deployed by one module — the `pro
	// module` remote-state discovery door.
	ModuleSlug *string `json:"module_slug" form:"module_slug"`
	// Group / Scope / IsEnabled are exact filters (repository-applied; `group`
	// is a reserved word so the generic mapper is bypassed).
	Group     *string `json:"group" form:"group"`
	Scope     *string `json:"scope" form:"scope"`
	IsEnabled *bool   `json:"is_enabled" form:"is_enabled"`
	// Search filters by a case-insensitive substring of key or name.
	Search *string `json:"search" form:"search"`
	common.Pagination
	common.Sorting
}

type GetManyConversationTagDefinitionsResponse struct {
	Meta common.ResponseMeta                           `json:"meta"`
	Data []conversationmodel.ConversationTagDefinition `json:"data"`
}

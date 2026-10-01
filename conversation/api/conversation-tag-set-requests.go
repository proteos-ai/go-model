package conversationapi

import (
	"go.proteos.ai/model/common"
	conversationmodel "go.proteos.ai/model/conversation"
)

// CreateConversationTagSetRequest creates one tag set: which definitions to
// evaluate (DefinitionKeys, ≥ 1, each must exist), for which conversations
// (at least one of Channels / ConnectionIds / TypeKeys non-empty; Directions
// additionally limits which messages arm/are judged — empty = both), when
// (EvaluateOn ⊆ [realtime, completion], ≥ 1) and the realtime timing (nil ⇒
// defaults: quiet 90s, max wait 600s). IsEnabled nil ⇒ true. ModuleSlug is
// stamped by `pro module deploy`.
type CreateConversationTagSetRequest struct {
	Key                string                                     `json:"key" validate:"required,max=255"`
	Name               string                                     `json:"name" validate:"max=255"`
	DefinitionKeys     []string                                   `json:"definition_keys" validate:"required,min=1"`
	Channels           []conversationmodel.Channel                `json:"channels"`
	ConnectionIds      []string                                   `json:"connection_ids"`
	TypeKeys           []string                                   `json:"type_keys"`
	Directions         []conversationmodel.MessageDirection       `json:"directions"`
	EvaluateOn         []conversationmodel.ConversationTagTrigger `json:"evaluate_on" validate:"required,min=1"`
	QuietWindowSeconds *int                                       `json:"quiet_window_seconds,omitempty"`
	MaxWaitSeconds     *int                                       `json:"max_wait_seconds,omitempty"`
	IsEnabled          *bool                                      `json:"is_enabled,omitempty"`
	ModuleSlug         string                                     `json:"module_slug" validate:"max=255"`
}

// UpdateConversationTagSetRequest is a partial update — nil leaves the stored
// value untouched; slices replace wholesale. Key is immutable.
type UpdateConversationTagSetRequest struct {
	Name               *string                                     `json:"name,omitempty" validate:"omitempty,max=255"`
	DefinitionKeys     *[]string                                   `json:"definition_keys,omitempty"`
	Channels           *[]conversationmodel.Channel                `json:"channels,omitempty"`
	ConnectionIds      *[]string                                   `json:"connection_ids,omitempty"`
	TypeKeys           *[]string                                   `json:"type_keys,omitempty"`
	Directions         *[]conversationmodel.MessageDirection       `json:"directions,omitempty"`
	EvaluateOn         *[]conversationmodel.ConversationTagTrigger `json:"evaluate_on,omitempty"`
	QuietWindowSeconds *int                                        `json:"quiet_window_seconds,omitempty"`
	MaxWaitSeconds     *int                                        `json:"max_wait_seconds,omitempty"`
	IsEnabled          *bool                                       `json:"is_enabled,omitempty"`
}

type GetManyConversationTagSetsQuery struct {
	// ModuleSlug filters to sets deployed by one module.
	ModuleSlug *string `json:"module_slug" form:"module_slug"`
	// DefinitionKey filters to sets that reference the definition
	// (definition_keys @> '["key"]', repository-applied).
	DefinitionKey *string `json:"definition_key" form:"definition_key"`
	IsEnabled     *bool   `json:"is_enabled" form:"is_enabled" db:"is_enabled"`
	// Search filters by a case-insensitive substring of key or name.
	Search *string `json:"search" form:"search"`
	common.Pagination
	common.Sorting
}

type GetManyConversationTagSetsResponse struct {
	Meta common.ResponseMeta                    `json:"meta"`
	Data []conversationmodel.ConversationTagSet `json:"data"`
}

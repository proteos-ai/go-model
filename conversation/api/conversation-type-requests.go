package conversationapi

import (
	"go.proteos.ai/model/common"
	conversationmodel "go.proteos.ai/model/conversation"
)

// CreateConversationTypeRequest creates one conversation type. Key is the
// immutable identity (lowercase kebab, like every platform key); Definition is
// what the classifier reads to decide membership. Config optionally links a
// summary prompt. ModuleSlug is stamped by `pro module deploy` — interactive
// callers leave it empty.
type CreateConversationTypeRequest struct {
	Key        string                                   `json:"key" validate:"required,max=255"`
	Name       string                                   `json:"name" validate:"max=255"`
	Definition string                                   `json:"definition" validate:"required,max=4096"`
	Config     conversationmodel.ConversationTypeConfig `json:"config"`
	ModuleSlug string                                   `json:"module_slug" validate:"max=255"`
}

// UpdateConversationTypeRequest is a partial update — nil leaves the stored
// value untouched. Config, when present, replaces the stored config wholesale
// (send an empty object to clear the summary prompt link). Key is immutable.
type UpdateConversationTypeRequest struct {
	Name       *string                                   `json:"name,omitempty" validate:"omitempty,max=255"`
	Definition *string                                   `json:"definition,omitempty" validate:"omitempty,max=4096"`
	Config     *conversationmodel.ConversationTypeConfig `json:"config,omitempty"`
}

type GetManyConversationTypesQuery struct {
	// ModuleSlug filters to types deployed by one module — the `pro module`
	// remote-state discovery door.
	ModuleSlug *string `json:"module_slug" form:"module_slug"`
	// Search filters by a case-insensitive substring of key or name.
	Search *string `json:"search" form:"search"`
	common.Pagination
	common.Sorting
}

type GetManyConversationTypesResponse struct {
	Meta common.ResponseMeta                  `json:"meta"`
	Data []conversationmodel.ConversationType `json:"data"`
}

package conversationapi

import (
	"time"

	"go.proteos.ai/model/common"
	conversationmodel "go.proteos.ai/model/conversation"
)

// CreateConversationObjectiveDefinitionRequest creates one objective
// definition. Key is the immutable identity (kebab/snake, like every
// platform key); TypeKey restricts it to one conversation type (empty =
// every conversation); PhaseKey pins it to one of that type's phases
// (requires TypeKey); Measure says how achievement is told (nil = guidance
// only). IsEnabled nil ⇒ true (manifests may omit it). ModuleSlug is stamped
// by `pro module deploy` — interactive callers leave it empty.
type CreateConversationObjectiveDefinitionRequest struct {
	Key         string                                          `json:"key" validate:"required,max=255"`
	Name        string                                          `json:"name" validate:"required,max=255"`
	Description string                                          `json:"description" validate:"max=4096"`
	TypeKey     string                                          `json:"type_key" validate:"max=255"`
	PhaseKey    string                                          `json:"phase_key" validate:"max=255"`
	Measure     *conversationmodel.ConversationObjectiveMeasure `json:"measure,omitempty"`
	IsEnabled   *bool                                           `json:"is_enabled,omitempty"`
	ModuleSlug  string                                          `json:"module_slug" validate:"max=255"`
}

// Normalized returns the request with the stored defaults filled in:
// is_enabled nil ⇒ true. Shared by the service and `pro module` (see
// CreateConversationTagDefinitionRequest.Normalized).
func (request CreateConversationObjectiveDefinitionRequest) Normalized() CreateConversationObjectiveDefinitionRequest {
	if request.IsEnabled == nil {
		isEnabled := true
		request.IsEnabled = &isEnabled
	}
	return request
}

// UpdateConversationObjectiveDefinitionRequest is a partial update — nil
// leaves the stored value untouched. Measure is tri-state (common.Optional):
// absent keeps the stored measure, JSON null clears it (guidance only), a
// value replaces it wholesale. `omitzero` keeps an unset Optional OUT of the
// encoded body (an SDK patch that touches only `name` must not clear the
// measure); an explicit null still encodes as null. Key is immutable.
type UpdateConversationObjectiveDefinitionRequest struct {
	Name        *string                                                         `json:"name,omitempty" validate:"omitempty,min=1,max=255"`
	Description *string                                                         `json:"description,omitempty" validate:"omitempty,max=4096"`
	TypeKey     *string                                                         `json:"type_key,omitempty" validate:"omitempty,max=255"`
	PhaseKey    *string                                                         `json:"phase_key,omitempty" validate:"omitempty,max=255"`
	Measure     common.Optional[conversationmodel.ConversationObjectiveMeasure] `json:"measure,omitzero"`
	IsEnabled   *bool                                                           `json:"is_enabled,omitempty"`
}

type GetManyConversationObjectiveDefinitionsQuery struct {
	// ModuleSlug filters to definitions deployed by one module.
	ModuleSlug *string `json:"module_slug" form:"module_slug"`
	TypeKey    *string `json:"type_key" form:"type_key" db:"type_key"`
	PhaseKey   *string `json:"phase_key" form:"phase_key" db:"phase_key"`
	IsEnabled  *bool   `json:"is_enabled" form:"is_enabled" db:"is_enabled"`
	// Search filters by a case-insensitive substring of key or name.
	Search *string `json:"search" form:"search"`
	common.Pagination
	common.Sorting
}

type GetManyConversationObjectiveDefinitionsResponse struct {
	Meta common.ResponseMeta                                 `json:"meta"`
	Data []conversationmodel.ConversationObjectiveDefinition `json:"data"`
}

// GetManyConversationObjectivesQuery is the org-wide ledger read (GET
// /conversation-objectives): every filter is optional; Since/Until bound
// recorded_at.
type GetManyConversationObjectivesQuery struct {
	ConversationId      *string    `json:"conversation_id" form:"conversation_id" db:"conversation_id"`
	DefinitionKey       *string    `json:"definition_key" form:"definition_key" db:"definition_key"`
	TypeKey             *string    `json:"type_key" form:"type_key" db:"type_key"`
	PhaseKey            *string    `json:"phase_key" form:"phase_key" db:"phase_key"`
	Channel             *string    `json:"channel" form:"channel" db:"channel"`
	Status              *string    `json:"status" form:"status" db:"status"`
	Source              *string    `json:"source" form:"source" db:"source"`
	IsProvisional       *bool      `json:"is_provisional" form:"is_provisional" db:"is_provisional"`
	IsConversationEnded *bool      `json:"is_conversation_ended" form:"is_conversation_ended" db:"is_conversation_ended"`
	Since               *time.Time `json:"since" form:"since" time_format:"2006-01-02T15:04:05Z07:00"`
	Until               *time.Time `json:"until" form:"until" time_format:"2006-01-02T15:04:05Z07:00"`
	common.Pagination
	common.Sorting
}

type GetManyConversationObjectivesResponse struct {
	Meta common.ResponseMeta                       `json:"meta"`
	Data []conversationmodel.ConversationObjective `json:"data"`
}

// CountConversationObjectivesQuery groups the ledger for dashboards. GroupBy
// is a comma-separated whitelist of definition_key, status, source,
// is_provisional, type_key, phase_key, channel (order kept; default
// definition_key,status).
type CountConversationObjectivesQuery struct {
	GroupBy             string     `json:"group_by" form:"group_by"`
	DefinitionKey       *string    `json:"definition_key" form:"definition_key"`
	TypeKey             *string    `json:"type_key" form:"type_key"`
	PhaseKey            *string    `json:"phase_key" form:"phase_key"`
	Channel             *string    `json:"channel" form:"channel"`
	Status              *string    `json:"status" form:"status"`
	Source              *string    `json:"source" form:"source"`
	IsProvisional       *bool      `json:"is_provisional" form:"is_provisional"`
	IsConversationEnded *bool      `json:"is_conversation_ended" form:"is_conversation_ended"`
	Since               *time.Time `json:"since" form:"since" time_format:"2006-01-02T15:04:05Z07:00"`
	Until               *time.Time `json:"until" form:"until" time_format:"2006-01-02T15:04:05Z07:00"`
}

type CountConversationObjectivesResponse struct {
	Data []conversationmodel.ConversationObjectiveCount `json:"data"`
}

// GetConversationObjectiveProgressResponse is GET
// /conversations/:id/objective-progress.
type GetConversationObjectiveProgressResponse struct {
	Data conversationmodel.ConversationObjectiveProgress `json:"data"`
}

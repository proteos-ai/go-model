package conversationapi

import (
	"go.proteos.ai/model/common"
	conversationmodel "go.proteos.ai/model/conversation"
)

// CreateConversationFilterRequest creates an ingest-time filter rule. Empty
// ConnectionId = global (org-wide); set = scoped to that connection (and
// evaluated BEFORE the global rules — a matching connection rule is final).
// FilterConfig is the raw per-type parameters; the service validates,
// canonicalizes and types it against FilterType. IsEnabled is a tri-state
// pointer — omitted defaults to TRUE. IsSelfDomainAcknowledged opts past the
// own-domain guard: a block rule targeting one of the org's own domains (a
// connected mailbox's domain or a platform user's domain) is rejected with
// self_domain_block unless this is set — blocking your own domain drops every
// thread your team is on; internal_conversations is almost always what's
// wanted instead.
type CreateConversationFilterRequest struct {
	ConnectionId string                                     `json:"connection_id"`
	FilterType   conversationmodel.ConversationFilterType   `json:"filter_type" validate:"required"`
	Action       conversationmodel.ConversationFilterAction `json:"action" validate:"required"`
	FilterConfig map[string]any                             `json:"filter_config"`
	IsEnabled    *bool                                      `json:"is_enabled,omitempty"`

	IsSelfDomainAcknowledged *bool `json:"is_self_domain_acknowledged,omitempty"`
}

// UpdateConversationFilterRequest mutates a filter in place. FilterType and
// FilterConfig must be sent together when either changes (they are mutually
// constraining). No connection_id — rules are never re-scoped; delete and
// recreate instead (mirrors agent-listener).
type UpdateConversationFilterRequest struct {
	FilterType   *conversationmodel.ConversationFilterType   `json:"filter_type,omitempty"`
	Action       *conversationmodel.ConversationFilterAction `json:"action,omitempty"`
	FilterConfig *map[string]any                             `json:"filter_config,omitempty"`
	IsEnabled    *bool                                       `json:"is_enabled,omitempty"`

	IsSelfDomainAcknowledged *bool `json:"is_self_domain_acknowledged,omitempty"`
}

// GetManyConversationFiltersQuery filters the org's rules. ConnectionId is
// tri-state: nil = all rules, pointer to "" = global-only rules, pointer to an
// id = that connection's rules.
type GetManyConversationFiltersQuery struct {
	ConnectionId *string `json:"connection_id" form:"connection_id" db:"connection_id"`
	FilterType   *string `json:"filter_type" form:"filter_type" db:"filter_type"`
	Action       *string `json:"action" form:"action" db:"action"`
	IsEnabled    *bool   `json:"is_enabled" form:"is_enabled" db:"is_enabled"`
	common.Pagination
	common.Sorting
}

type GetManyConversationFiltersResponse struct {
	Meta common.ResponseMeta                    `json:"meta"`
	Data []conversationmodel.ConversationFilter `json:"data"`
}

// GetManyConversationFilterEventsQuery pages one filter's drop-audit trail,
// newest first.
type GetManyConversationFilterEventsQuery struct {
	common.Pagination
	common.Sorting
}

type GetManyConversationFilterEventsResponse struct {
	Meta common.ResponseMeta                         `json:"meta"`
	Data []conversationmodel.ConversationFilterEvent `json:"data"`
}

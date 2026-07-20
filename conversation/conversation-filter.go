package conversationmodel

import (
	"time"

	"go.proteos.ai/model/common"
)

// ConversationFilter is an ingest-time rule that decides whether an inbound
// message enters the platform at all: a matching enabled block rule drops the
// message BEFORE dedupe, contact resolution, and persistence — nothing is
// stored except a content-free ConversationFilterEvent audit row. ConnectionId
// scopes the rule to one connection; empty = global (org-wide).
//
// Evaluation is scope-first, then specificity (no priority field — the order
// is total): connection-scoped rules are evaluated first and any match — allow
// or block — is final; global rules apply only when no connection rule
// matched. Within a scope the specificity order is
// address > domain > role_based > automated > internal_conversations > all,
// and within one class allow beats block. So a connection-scoped all-allow
// bypasses every global filter for that connection, and a connection-scoped
// all-block pauses its ingest entirely (add connection address-allows on top
// for a de-facto allowlist).
type ConversationFilter struct {
	Id           string                   `json:"id"`
	OrgId        string                   `json:"org_id"`
	ConnectionId string                   `json:"connection_id"`
	FilterType   ConversationFilterType   `json:"filter_type" sortable:""`
	Action       ConversationFilterAction `json:"action" sortable:""`
	// FilterConfig is the typed, per-type configuration (a tagged union keyed by
	// FilterType — see conversation-filter-config.go). Serializes to the bare
	// variant; FilterType discriminates.
	FilterConfig ConversationFilterConfig `json:"filter_config,omitempty"`
	IsEnabled    bool                     `json:"is_enabled" sortable:""`
	// RecentEventCount is a read-only projection: how many messages this rule
	// dropped in the last 30 days (logic.RecentEventWindow), rolled up from the
	// conversation_filter_event audit table on CRUD reads only — never persisted
	// on the row, never computed on the hot ingest path.
	RecentEventCount int            `json:"recent_event_count,omitempty"`
	CreatedAt        time.Time      `json:"created_at" sortable:""`
	CreatedBy        common.UserRef `json:"created_by"`
	UpdatedAt        time.Time      `json:"updated_at" sortable:""`
	UpdatedBy        common.UserRef `json:"updated_by"`
}

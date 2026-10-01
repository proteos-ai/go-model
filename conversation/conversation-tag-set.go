package conversationmodel

import (
	"time"

	"go.proteos.ai/model/common"
)

// ConversationTagSet describes WHERE and WHEN a group of tag definitions is
// evaluated. A conversation matches a set when every NON-EMPTY selector
// contains it: Channels, ConnectionIds, TypeKeys (an empty conversation
// type_key never matches a type-scoped set, so type-scoped sets only apply
// once a conversation has been classified). At least one of those three must
// be non-empty — tagging is opt-in, never org-wide by accident. Directions is
// a message-level filter on top (which messages arm / are judged), not a
// conversation selector. Matched sets
// union their definition keys; timing takes the minimum quiet window and max
// wait over the matched sets that evaluate in realtime. Keyed (org_id, key)
// and module-deployable (conversation-tag-sets/<key>.json).
type ConversationTagSet struct {
	OrgId string `json:"org_id"`
	Key   string `json:"key" sortable:""`
	Name  string `json:"name" sortable:""`
	// DefinitionKeys lists the definitions this set evaluates (≥ 1; each must
	// exist in the org at write time — a definition deleted later is skipped).
	DefinitionKeys []string  `json:"definition_keys"`
	Channels       []Channel `json:"channels"`
	ConnectionIds  []string  `json:"connection_ids"`
	TypeKeys       []string  `json:"type_keys"`
	// Directions narrows which NEW messages arm a realtime evaluation (and
	// which messages a completion evaluation judges for message/span scope):
	// inbound only, outbound only, or both when empty. A definition's own
	// Directions narrows further.
	Directions []MessageDirection `json:"directions"`
	// EvaluateOn is a non-empty subset of [realtime, completion]: realtime =
	// the debounced window after new messages arrive; completion = once the
	// conversation completed (final transcript + summary for spoken media, End
	// for the rest). Manual evaluations ignore it.
	EvaluateOn []ConversationTagTrigger `json:"evaluate_on"`
	// QuietWindowSeconds: a realtime evaluation fires this long after the LAST
	// new message (default 90, min 5). MaxWaitSeconds bounds the wait since the
	// FIRST unjudged message (default 600, ≥ quiet window) so a busy stream is
	// still tagged periodically.
	QuietWindowSeconds int  `json:"quiet_window_seconds"`
	MaxWaitSeconds     int  `json:"max_wait_seconds"`
	IsEnabled          bool `json:"is_enabled"`
	// ModuleSlug attributes the set to the module that deployed it; empty =
	// not module-owned.
	ModuleSlug string         `json:"module_slug,omitempty"`
	CreatedAt  time.Time      `json:"created_at" sortable:""`
	CreatedBy  common.UserRef `json:"created_by"`
	UpdatedAt  time.Time      `json:"updated_at" sortable:""`
	UpdatedBy  common.UserRef `json:"updated_by"`
}

// Defaults for a tag set's realtime timing when a request leaves them unset.
const (
	ConversationTagSetDefaultQuietWindowSeconds = 90
	ConversationTagSetMinQuietWindowSeconds     = 5
	ConversationTagSetDefaultMaxWaitSeconds     = 600
)

// EvaluatesOn reports whether the set subscribes to the automatic trigger.
func (set ConversationTagSet) EvaluatesOn(trigger ConversationTagTrigger) bool {
	for _, entry := range set.EvaluateOn {
		if entry == trigger {
			return true
		}
	}
	return false
}

// ArmsOn reports whether a message of the direction may arm a realtime
// evaluation under this set (empty Directions = any).
func (set ConversationTagSet) ArmsOn(direction MessageDirection) bool {
	if len(set.Directions) == 0 {
		return true
	}
	for _, entry := range set.Directions {
		if entry == direction {
			return true
		}
	}
	return false
}

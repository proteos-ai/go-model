package conversationmodel

import (
	"slices"
	"time"

	"go.proteos.ai/model/common"
)

// SendingRule is one org-configured constraint on OUTBOUND sending, evaluated
// synchronously on every Send/SendDraft before a message row is minted: a
// denied send returns a typed error carrying the earliest time it would be
// allowed again. Three rule types (RuleType + the typed RuleConfig union):
//
//   - window        — WHEN sending is allowed: weekday time ranges in the
//     recipient's local time (Contact.Timezone, else the rule's fallback).
//   - limit         — HOW MUCH one sender connection may send per rolling
//     period (deliverability / provider-safety caps, lemlist "sending
//     limits"; presets per account type expand into these rows).
//   - frequency_cap — HOW OFTEN one contact may be contacted per rolling
//     period (Braze/AJO "frequency capping").
//
// A rule is defined once and LINKED to any number of connections and/or
// channels: it applies to a send when the connection is in ConnectionIds, OR
// the channel is in Channels, OR both are empty (org-wide). Per rule type the
// most specific non-empty tier wins — rules naming the connection replace
// rules naming the channel, which replace org-wide rules; inside a tier
// windows union and limits/caps must all pass.
//
// Replies (a send into an existing conversation) are exempt when
// IsReplyExempt is set — answering someone who wrote to us is not outreach.
// Windows and frequency caps default to exempt; limits do not (the provider
// counts every send).
type SendingRule struct {
	Id    string `json:"id"`
	OrgId string `json:"org_id"`
	// Name is a human label ("EU business hours", "LinkedIn – Sales Navigator");
	// presets stamp their own name. Empty allowed.
	Name string `json:"name"`
	// ConnectionIds / Channels are the link targets (see type doc). Both empty =
	// org-wide. Serialize as arrays, never null.
	ConnectionIds []string        `json:"connection_ids"`
	Channels      []Channel       `json:"channels"`
	RuleType      SendingRuleType `json:"rule_type" sortable:""`
	// RuleConfig is the typed per-type configuration (a tagged union keyed by
	// RuleType — see sending-rule-config.go). Serializes to the bare variant.
	RuleConfig    SendingRuleConfig `json:"rule_config,omitempty"`
	IsEnabled     bool              `json:"is_enabled" sortable:""`
	IsReplyExempt bool              `json:"is_reply_exempt"`
	CreatedAt     time.Time         `json:"created_at" sortable:""`
	CreatedBy     common.UserRef    `json:"created_by"`
	UpdatedAt     time.Time         `json:"updated_at" sortable:""`
	UpdatedBy     common.UserRef    `json:"updated_by"`
}

// Scope tier of a rule relative to one send: 0 = names the connection,
// 1 = names the channel, 2 = org-wide, -1 = does not apply.
func (rule SendingRule) TierFor(connectionId string, channel Channel) int {
	if connectionId != "" && slices.Contains(rule.ConnectionIds, connectionId) {
		return 0
	}
	if channel != "" && slices.Contains(rule.Channels, channel) {
		return 1
	}
	if rule.IsOrgWide() {
		return 2
	}
	return -1
}

// IsOrgWide reports whether the rule links to nothing (applies everywhere).
func (rule SendingRule) IsOrgWide() bool {
	return len(rule.ConnectionIds) == 0 && len(rule.Channels) == 0
}

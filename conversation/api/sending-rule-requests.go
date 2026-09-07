package conversationapi

import (
	"go.proteos.ai/model/common"
	conversationmodel "go.proteos.ai/model/conversation"
)

// CreateSendingRuleRequest creates one outbound sending constraint. ConnectionIds
// / Channels link the rule (both empty = org-wide). RuleConfig is the raw
// per-type parameters; the service validates and types it against RuleType.
// IsEnabled defaults to TRUE when omitted; IsReplyExempt defaults per type
// (window + frequency_cap true, limit false) when omitted.
type CreateSendingRuleRequest struct {
	Name          string                            `json:"name"`
	ConnectionIds []string                          `json:"connection_ids"`
	Channels      []conversationmodel.Channel       `json:"channels"`
	RuleType      conversationmodel.SendingRuleType `json:"rule_type" validate:"required"`
	RuleConfig    map[string]any                    `json:"rule_config"`
	IsEnabled     *bool                             `json:"is_enabled,omitempty"`
	IsReplyExempt *bool                             `json:"is_reply_exempt,omitempty"`
}

// UpdateSendingRuleRequest mutates a rule in place. RuleType and RuleConfig
// must be sent together when either changes. ConnectionIds / Channels replace
// the stored links wholesale when present (re-linking a shared rule IS the
// point of the arrays; an explicit empty array clears the tier).
type UpdateSendingRuleRequest struct {
	Name          *string                            `json:"name,omitempty"`
	ConnectionIds *[]string                          `json:"connection_ids,omitempty"`
	Channels      *[]conversationmodel.Channel       `json:"channels,omitempty"`
	RuleType      *conversationmodel.SendingRuleType `json:"rule_type,omitempty"`
	RuleConfig    *map[string]any                    `json:"rule_config,omitempty"`
	IsEnabled     *bool                              `json:"is_enabled,omitempty"`
	IsReplyExempt *bool                              `json:"is_reply_exempt,omitempty"`
}

// GetManySendingRulesQuery filters the org's rules. ConnectionId = rules LINKED
// to that connection (not the rules that would apply to it — tiering is a
// send-time concern); Channel likewise.
type GetManySendingRulesQuery struct {
	ConnectionId *string `json:"connection_id" form:"connection_id"`
	Channel      *string `json:"channel" form:"channel"`
	RuleType     *string `json:"rule_type" form:"rule_type" db:"rule_type"`
	IsEnabled    *bool   `json:"is_enabled" form:"is_enabled" db:"is_enabled"`
	common.Pagination
	common.Sorting
}

type GetManySendingRulesResponse struct {
	Meta common.ResponseMeta             `json:"meta"`
	Data []conversationmodel.SendingRule `json:"data"`
}

// GetManySendingLimitPresetsQuery narrows the static preset catalog to one
// connector (omit for all).
type GetManySendingLimitPresetsQuery struct {
	ConnectorKey *string `json:"connector_key" form:"connector_key"`
}

type GetManySendingLimitPresetsResponse struct {
	Data []conversationmodel.SendingLimitPreset `json:"data"`
}

// ApplySendingLimitPresetRequest expands a preset into limit rules linked to
// the given connections, replacing whatever limit rules those connections were
// linked to before (a connection stripped from a shared rule leaves the rule
// on its other targets; a rule left with no targets is deleted).
type ApplySendingLimitPresetRequest struct {
	ConnectionIds []string `json:"connection_ids" validate:"required,min=1"`
	PresetKey     string   `json:"preset_key" validate:"required"`
}

type ApplySendingLimitPresetResponse struct {
	Data []conversationmodel.SendingRule `json:"data"`
}

// SendEligibilityRequest is the dry-run twin of SendMessageRequest: the same
// three addressing modes (reply: ConversationId; thread: ReplyToMessageId;
// originate: ConnectionId + To), no content. Nothing is minted — an originate
// check resolves the connection and recipients only. ActionType widens the
// check to a channel action (invitation, profile_visit, inmail) — then only
// originate mode applies (ConnectionId + To = the target); empty = message.
type SendEligibilityRequest struct {
	ConversationId   string                              `json:"conversation_id"`
	ReplyToMessageId string                              `json:"reply_to_message_id"`
	ConnectionId     string                              `json:"connection_id"`
	To               []conversationmodel.Recipient       `json:"to"`
	Cc               []conversationmodel.Recipient       `json:"cc,omitempty"`
	Bcc              []conversationmodel.Recipient       `json:"bcc,omitempty"`
	ActionType       conversationmodel.ChannelActionType `json:"action_type,omitempty"`
}

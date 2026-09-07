package conversationmodel

import (
	"encoding/json"
	"fmt"
)

// SendingRuleConfig is the typed, per-type configuration of a SendingRule — a
// tagged union discriminated by the sibling RuleType (mirrors
// ConversationFilterConfig). Every type carries parameters.
type SendingRuleConfig interface {
	isSendingRuleConfig()
	RuleType() SendingRuleType
}

// WindowDay is one open range on one weekday, in the recipient's local
// time. From/Until are "HH:MM" floating wall-clock times, From < Until, same
// day (no midnight crossing in v1). A weekday may carry several ranges — they
// union. A weekday with no range is closed all day.
type WindowDay struct {
	Day   Weekday `json:"day"`
	From  string  `json:"from"`
	Until string  `json:"until"`
}

// WindowRuleConfig — WHEN sending is allowed. Evaluated per To recipient in
// Contact.Timezone; a contact without a timezone uses FallbackTimezone
// (required, IANA). Open iff every To recipient is inside a range.
type WindowRuleConfig struct {
	Days             []WindowDay `json:"days"`
	FallbackTimezone string      `json:"fallback_timezone"`
}

func (WindowRuleConfig) isSendingRuleConfig()      {}
func (WindowRuleConfig) RuleType() SendingRuleType { return SendingRuleTypeWindow }

// LimitRuleConfig — HOW MUCH one connection may send: at most MaxCount sends
// of Action within the rolling Period, plus an optional minimum spacing
// between consecutive sends (MinGapSeconds, 0 = none). Counting reads the
// connection's outbound messages with status pending|sent (pending counts —
// a racing send reserves its slot; failed never counts).
type LimitRuleConfig struct {
	Action        ChannelActionType `json:"action"`
	MaxCount      int               `json:"max_count"`
	Period        SendingPeriod     `json:"period"`
	MinGapSeconds int               `json:"min_gap_seconds,omitempty"`
}

func (LimitRuleConfig) isSendingRuleConfig()      {}
func (LimitRuleConfig) RuleType() SendingRuleType { return SendingRuleTypeLimit }

// FrequencyCapRuleConfig — HOW OFTEN one contact may be contacted: at most
// MaxCount outbound messages addressed to the contact within the rolling
// Period. The rule's Channels link scopes the count to those channels; an
// org-wide cap pools every channel.
type FrequencyCapRuleConfig struct {
	MaxCount int           `json:"max_count"`
	Period   SendingPeriod `json:"period"`
}

func (FrequencyCapRuleConfig) isSendingRuleConfig()      {}
func (FrequencyCapRuleConfig) RuleType() SendingRuleType { return SendingRuleTypeFrequencyCap }

// MarshalSendingRuleConfig encodes a variant to its stored (JSONB) bare shape.
// The discriminator lives in the sibling rule_type column. nil → '{}'.
func MarshalSendingRuleConfig(config SendingRuleConfig) (json.RawMessage, error) {
	if config == nil {
		return json.RawMessage("{}"), nil
	}
	return json.Marshal(config)
}

// DecodeSendingRuleConfig rebuilds the typed variant from the stored bare JSON
// plus the discriminating rule type.
func DecodeSendingRuleConfig(ruleType SendingRuleType, raw []byte) (SendingRuleConfig, error) {
	switch ruleType {
	case SendingRuleTypeWindow:
		config := WindowRuleConfig{}
		if err := unmarshalConfig(raw, &config); err != nil {
			return nil, err
		}
		return config, nil
	case SendingRuleTypeLimit:
		config := LimitRuleConfig{}
		if err := unmarshalConfig(raw, &config); err != nil {
			return nil, err
		}
		return config, nil
	case SendingRuleTypeFrequencyCap:
		config := FrequencyCapRuleConfig{}
		if err := unmarshalConfig(raw, &config); err != nil {
			return nil, err
		}
		return config, nil
	default:
		return nil, fmt.Errorf("unknown sending rule type %q", ruleType)
	}
}

package conversationmodel

import (
	"encoding/json"
	"fmt"
)

// AgentListenerTriggerConfig is the typed, per-trigger configuration of an
// AgentListener — a tagged union discriminated by the sibling TriggerType
// (mirrors ConnectionCredentials, whose Kind is the discriminator). It replaces
// the former untyped map[string]any so each trigger's parameters are explicit
// and validatable. always/mention carry no parameters; channel/keyword do.
type AgentListenerTriggerConfig interface {
	isAgentListenerTriggerConfig()
	TriggerType() AgentListenerTriggerType
}

// AlwaysConfig — fire on every inbound message on the target. No parameters.
type AlwaysConfig struct{}

func (AlwaysConfig) isAgentListenerTriggerConfig()         {}
func (AlwaysConfig) TriggerType() AgentListenerTriggerType { return TriggerTypeAlways }

// MentionConfig — fire when the connection's bot user is @-mentioned. No
// parameters (the bot user comes from the connection).
type MentionConfig struct{}

func (MentionConfig) isAgentListenerTriggerConfig()         {}
func (MentionConfig) TriggerType() AgentListenerTriggerType { return TriggerTypeMention }

// ChannelConfig — fire only for messages in a specific integration-side channel.
type ChannelConfig struct {
	ExternalChannelId string `json:"external_channel_id"`
}

func (ChannelConfig) isAgentListenerTriggerConfig()         {}
func (ChannelConfig) TriggerType() AgentListenerTriggerType { return TriggerTypeChannel }

// KeywordConfig — fire when any configured phrase appears in the message text.
type KeywordConfig struct {
	Phrases []string `json:"phrases"`
}

func (KeywordConfig) isAgentListenerTriggerConfig()         {}
func (KeywordConfig) TriggerType() AgentListenerTriggerType { return TriggerTypeKeyword }

// MarshalTriggerConfig encodes a variant to its stored (JSONB) bare shape. The
// discriminator lives in the sibling trigger_type column, so — unlike
// ConnectionCredentials — no {kind,data} envelope is needed. nil → '{}'.
func MarshalTriggerConfig(config AgentListenerTriggerConfig) (json.RawMessage, error) {
	if config == nil {
		return json.RawMessage("{}"), nil
	}
	return json.Marshal(config)
}

// DecodeTriggerConfig rebuilds the typed variant from the stored bare JSON plus
// the discriminating trigger type. Empty/absent config decodes to the empty
// variant for that type (always/mention have none).
func DecodeTriggerConfig(triggerType AgentListenerTriggerType, raw []byte) (AgentListenerTriggerConfig, error) {
	switch triggerType {
	case TriggerTypeAlways:
		return AlwaysConfig{}, nil
	case TriggerTypeMention:
		return MentionConfig{}, nil
	case TriggerTypeChannel:
		config := ChannelConfig{}
		if err := unmarshalConfig(raw, &config); err != nil {
			return nil, err
		}
		return config, nil
	case TriggerTypeKeyword:
		config := KeywordConfig{}
		if err := unmarshalConfig(raw, &config); err != nil {
			return nil, err
		}
		return config, nil
	default:
		return nil, fmt.Errorf("unknown trigger type %q", triggerType)
	}
}

func unmarshalConfig(raw []byte, target any) error {
	trimmed := string(raw)
	if len(raw) == 0 || trimmed == "null" || trimmed == "{}" {
		return nil
	}
	return json.Unmarshal(raw, target)
}

package conversationmodel

import (
	"encoding/json"
	"fmt"
)

// AgentListenerAcknowledgementConfig is the typed, per-type configuration of a
// listener's immediate acknowledgement — a tagged union discriminated by the
// sibling AcknowledgementType (mirrors AgentListenerTriggerConfig). The
// acknowledgement is executed by the DISPATCHER (never the agent) right after
// the agent turn is queued, so people see the agent is on it while it thinks.
type AgentListenerAcknowledgementConfig interface {
	isAgentListenerAcknowledgementConfig()
	AcknowledgementType() AgentListenerAcknowledgementType
}

// ReactionAcknowledgementConfig — react to the triggering message with an
// emoji. Emoji is the connector-native token (Slack shortcode "eyes", no
// colons; unicode glyph on Unipile channels).
type ReactionAcknowledgementConfig struct {
	Emoji string `json:"emoji"`
}

func (ReactionAcknowledgementConfig) isAgentListenerAcknowledgementConfig() {}
func (ReactionAcknowledgementConfig) AcknowledgementType() AgentListenerAcknowledgementType {
	return AcknowledgementTypeReaction
}

// MessageAcknowledgementConfig — post a short configurable text into the
// conversation (threaded under the triggering message where the channel
// supports targeted replies).
type MessageAcknowledgementConfig struct {
	Text string `json:"text"`
}

func (MessageAcknowledgementConfig) isAgentListenerAcknowledgementConfig() {}
func (MessageAcknowledgementConfig) AcknowledgementType() AgentListenerAcknowledgementType {
	return AcknowledgementTypeMessage
}

// MarshalAcknowledgementConfig encodes a variant to its stored (JSONB) bare
// shape. The discriminator lives in the sibling acknowledgement_type column,
// so no {kind,data} envelope is needed. nil → '{}'.
func MarshalAcknowledgementConfig(config AgentListenerAcknowledgementConfig) (json.RawMessage, error) {
	if config == nil {
		return json.RawMessage("{}"), nil
	}
	return json.Marshal(config)
}

// DecodeAcknowledgementConfig rebuilds the typed variant from the stored bare
// JSON plus the discriminating type. AcknowledgementTypeNone ("") is a legal
// value meaning "no acknowledgement" and decodes to nil (mirrors
// wake-phrase-empty-is-off); unknown types error.
func DecodeAcknowledgementConfig(acknowledgementType AgentListenerAcknowledgementType, raw []byte) (AgentListenerAcknowledgementConfig, error) {
	switch acknowledgementType {
	case AcknowledgementTypeNone:
		return nil, nil
	case AcknowledgementTypeReaction:
		config := ReactionAcknowledgementConfig{}
		if err := unmarshalConfig(raw, &config); err != nil {
			return nil, err
		}
		return config, nil
	case AcknowledgementTypeMessage:
		config := MessageAcknowledgementConfig{}
		if err := unmarshalConfig(raw, &config); err != nil {
			return nil, err
		}
		return config, nil
	default:
		return nil, fmt.Errorf("unknown acknowledgement type %q", acknowledgementType)
	}
}

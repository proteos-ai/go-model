package conversationmodel

import "testing"

func TestAcknowledgementConfigRoundTrip(t *testing.T) {
	cases := []struct {
		name                string
		acknowledgementType AgentListenerAcknowledgementType
		config              AgentListenerAcknowledgementConfig
	}{
		{"reaction", AcknowledgementTypeReaction, ReactionAcknowledgementConfig{Emoji: "eyes"}},
		{"message", AcknowledgementTypeMessage, MessageAcknowledgementConfig{Text: "On it"}},
	}
	for _, testCase := range cases {
		t.Run(testCase.name, func(t *testing.T) {
			raw, err := MarshalAcknowledgementConfig(testCase.config)
			if err != nil {
				t.Fatalf("marshal: %v", err)
			}
			decoded, err := DecodeAcknowledgementConfig(testCase.acknowledgementType, raw)
			if err != nil {
				t.Fatalf("decode: %v", err)
			}
			if decoded != testCase.config {
				t.Fatalf("round trip mismatch: %#v != %#v", decoded, testCase.config)
			}
		})
	}
}

func TestMarshalAcknowledgementConfigNil(t *testing.T) {
	raw, err := MarshalAcknowledgementConfig(nil)
	if err != nil {
		t.Fatalf("marshal nil: %v", err)
	}
	if string(raw) != "{}" {
		t.Fatalf("nil must marshal to {}, got %s", raw)
	}
}

func TestDecodeAcknowledgementConfigNoneIsNil(t *testing.T) {
	config, err := DecodeAcknowledgementConfig(AcknowledgementTypeNone, []byte(`{"emoji":"eyes"}`))
	if err != nil {
		t.Fatalf("none must never error: %v", err)
	}
	if config != nil {
		t.Fatalf("none must decode to nil, got %#v", config)
	}
}

func TestDecodeAcknowledgementConfigUnknownTypeErrors(t *testing.T) {
	if _, err := DecodeAcknowledgementConfig("bogus", []byte(`{}`)); err == nil {
		t.Fatalf("unknown type must error")
	}
}

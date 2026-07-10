package conversationmodel

import "testing"

func TestTriggerConfigRoundTrip(t *testing.T) {
	cases := []struct {
		name string
		typ  AgentListenerTriggerType
		cfg  AgentListenerTriggerConfig
	}{
		{"always", TriggerTypeAlways, AlwaysConfig{}},
		{"mention", TriggerTypeMention, MentionConfig{}},
		{"channel", TriggerTypeChannel, ChannelConfig{ExternalChannelId: "C123"}},
		{"keyword", TriggerTypeKeyword, KeywordConfig{Phrases: []string{"help", "urgent"}}},
	}
	for _, testCase := range cases {
		t.Run(testCase.name, func(t *testing.T) {
			raw, err := MarshalTriggerConfig(testCase.cfg)
			if err != nil {
				t.Fatalf("marshal: %v", err)
			}
			decoded, err := DecodeTriggerConfig(testCase.typ, raw)
			if err != nil {
				t.Fatalf("decode: %v", err)
			}
			if decoded.TriggerType() != testCase.typ {
				t.Fatalf("type = %q, want %q", decoded.TriggerType(), testCase.typ)
			}
			switch want := testCase.cfg.(type) {
			case ChannelConfig:
				if got, ok := decoded.(ChannelConfig); !ok || got != want {
					t.Fatalf("channel round trip lost data: %#v", decoded)
				}
			case KeywordConfig:
				got, ok := decoded.(KeywordConfig)
				if !ok || len(got.Phrases) != len(want.Phrases) {
					t.Fatalf("keyword round trip lost data: %#v", decoded)
				}
			}
		})
	}
}

func TestTriggerConfigNilAndUnknown(t *testing.T) {
	// nil (always/mention with no params) marshals to '{}'.
	raw, err := MarshalTriggerConfig(nil)
	if err != nil || string(raw) != "{}" {
		t.Fatalf("nil config must marshal to '{}': %q %v", raw, err)
	}
	if _, err := DecodeTriggerConfig("carrier_pigeon", []byte("{}")); err == nil {
		t.Fatalf("an unknown trigger type must error, not silently drop")
	}
}

package agentmodel

import (
	"encoding/json"
	"testing"
)

func TestUserMessagePayloadContextRoundTrip(t *testing.T) {
	encoded, err := json.Marshal(UserMessagePayload{Content: []ContentBlock{{Type: "text", Text: "hi"}}})
	if err != nil {
		t.Fatal(err)
	}
	if string(encoded) != `{"content":[{"type":"text","text":"hi"}]}` {
		t.Fatalf("empty context must be omitted, got %s", encoded)
	}
	var decoded UserMessagePayload
	if err := json.Unmarshal([]byte(`{"content":[{"type":"text","text":"hi"}],"context":"ctx"}`), &decoded); err != nil {
		t.Fatal(err)
	}
	if decoded.Context != "ctx" || len(decoded.Content) != 1 {
		t.Fatalf("unexpected decode %+v", decoded)
	}
}

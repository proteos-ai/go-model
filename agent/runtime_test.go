package agentmodel

import (
	"encoding/json"
	"testing"
)

func TestRuntimeUnmarshalJSON(t *testing.T) {
	testCases := []struct {
		name    string
		input   string
		want    Runtime
		wantErr bool
	}{
		{name: "durable", input: `"durable"`, want: RuntimeDurable},
		{name: "embedded", input: `"embedded"`, want: RuntimeEmbedded},
		{name: "sandbox", input: `"sandbox"`, want: RuntimeSandbox},
		{name: "empty means not given", input: `""`, want: ""},
		{name: "null means not given", input: `null`, want: ""},
		{name: "a provider name is not a runtime", input: `"anthropic_managed"`, wantErr: true},
		{name: "unknown value", input: `"in_memory"`, wantErr: true},
		{name: "wrong type", input: `1`, wantErr: true},
	}
	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			// Start from a set value: "not given" must clear it, not leave it behind.
			runtime := RuntimeEmbedded
			err := json.Unmarshal([]byte(testCase.input), &runtime)
			if testCase.wantErr {
				if err == nil {
					t.Fatalf("Unmarshal(%s) = %q, want an error", testCase.input, runtime)
				}
				return
			}
			if err != nil {
				t.Fatalf("Unmarshal(%s): %v", testCase.input, err)
			}
			if runtime != testCase.want {
				t.Errorf("Unmarshal(%s) = %q, want %q", testCase.input, runtime, testCase.want)
			}
		})
	}
}

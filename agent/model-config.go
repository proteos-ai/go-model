package agentmodel

// ModelConfig is the inline LLM configuration an Agent runs with. A provider /
// credentials registry is deferred; v1 carries the model id plus optional
// sampling knobs.
type ModelConfig struct {
	ModelId     string   `json:"model_id"`
	Temperature *float64 `json:"temperature,omitempty"`
	MaxTokens   *int     `json:"max_tokens,omitempty"`
	Thinking    *bool    `json:"thinking,omitempty"`
}

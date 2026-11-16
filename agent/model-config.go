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

// The model providers an agent's model can come from. A ModelInfo names one, and
// so does stored reasoning (ReasoningPayload.Provider).
const (
	ModelProviderAnthropic = "anthropic"
	ModelProviderOpenAi    = "openai"
	ModelProviderGoogle    = "google"
)

// ModelInfo describes one model an agent can be configured with.
type ModelInfo struct {
	Id   string `json:"id"`
	Name string `json:"name"`
	// Provider is the model provider that serves the model (ModelProvider*).
	Provider string `json:"provider"`
	// Runtimes are the runtimes an agent with this model can run on.
	Runtimes []Runtime `json:"runtimes"`
	// ContextWindow is how many tokens one request to the model can hold, input
	// and output together. Zero when it is not known.
	ContextWindow int64 `json:"context_window,omitempty"`
}

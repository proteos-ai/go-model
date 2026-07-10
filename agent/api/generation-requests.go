package agentapi

import agentmodel "go.proteos.ai/model/agent"

// GenerateRequest is the body of POST /agents/v1/models/generate — a single,
// stateless model call. It mirrors agentmodel.GenerationRequest; the controller
// maps it onto the domain value. Model is optional (empty model_id resolves to the
// service default); Content carries the message body as the shared ContentBlock
// vocabulary (text | file). The response is a bare agentmodel.GenerationResult.
type GenerateRequest struct {
	Model   agentmodel.ModelConfig    `json:"model"`
	System  string                    `json:"system,omitempty"`
	Content []agentmodel.ContentBlock `json:"content" validate:"required,min=1"`
}

// ToGenerationRequest maps the wire DTO onto the domain value.
func (request GenerateRequest) ToGenerationRequest() agentmodel.GenerationRequest {
	return agentmodel.GenerationRequest{
		Model:   request.Model,
		System:  request.System,
		Content: request.Content,
	}
}

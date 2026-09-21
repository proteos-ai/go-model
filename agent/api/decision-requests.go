package agentapi

import (
	"encoding/json"

	agentmodel "go.proteos.ai/model/agent"
)

// DecideRequest is the body of POST /agents/v1/models/decide — a single, stateless
// decision call: one text-only state and a map of typed questions (boolean | choice
// | score) answered in parallel with calibrated probabilities. Model is optional
// (empty model_id resolves to the service default). The response is a bare
// agentmodel.DecisionResult.
//
// Questions stay raw on the wire DTO (precedent: CreateToolRequest.Binding): the
// domain decodes them into the typed agentmodel.DecisionQuestion union, so a
// malformed question surfaces as 400 invalid_decision_request naming the key,
// instead of a bind error falling through the reflective error handler as a 500.
type DecideRequest struct {
	Model     agentmodel.ModelConfig     `json:"model"`
	State     json.RawMessage            `json:"state"`
	Questions map[string]json.RawMessage `json:"questions"`
}

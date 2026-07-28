package workflowmodel_test

import (
	"encoding/json"
	"testing"

	workflowmodel "go.proteos.ai/model/workflow"
)

func TestDecodeAgentActionParams(t *testing.T) {
	t.Run("flat message params", func(t *testing.T) {
		params, err := workflowmodel.DecodeAgentActionParams(json.RawMessage(`{
			"agent_key": "triage-agent",
			"kickoff_type": "message",
			"message_source": "manual",
			"message_text": "go"
		}`))
		if err != nil {
			t.Fatal(err)
		}
		if params.AgentKey != "triage-agent" || params.KickoffType != workflowmodel.KickoffTypeMessage ||
			params.MessageText != "go" {
			t.Fatalf("params = %+v", params)
		}
	})

	t.Run("flat outcome params with output schema", func(t *testing.T) {
		params, err := workflowmodel.DecodeAgentActionParams(json.RawMessage(`{
			"agent_key": "triage-agent",
			"kickoff_type": "outcome",
			"description": "the goal",
			"rubric_type": "text",
			"rubric_content": "the rubric",
			"max_iterations": 3,
			"output_schema": [{"name": "sentiment", "type": "string", "is_required": true}],
			"is_output_required": true
		}`))
		if err != nil {
			t.Fatal(err)
		}
		if params.KickoffType != workflowmodel.KickoffTypeOutcome || params.RubricContent != "the rubric" {
			t.Fatalf("params = %+v", params)
		}
		if params.MaxIterations == nil || *params.MaxIterations != 3 {
			t.Fatalf("max_iterations = %v", params.MaxIterations)
		}
		if len(params.OutputSchema) != 1 || params.OutputSchema[0].Name != "sentiment" || !params.IsOutputRequired {
			t.Fatalf("output schema = %+v required = %v", params.OutputSchema, params.IsOutputRequired)
		}
	})

	t.Run("old nested kickoff decodes to zero-value flat fields", func(t *testing.T) {
		// No back-compat: the nested `kickoff` key is an unknown field, so the
		// decode succeeds but kickoff_type stays empty — save-time validation
		// rejects the node until it is re-edited.
		params, err := workflowmodel.DecodeAgentActionParams(json.RawMessage(`{
			"agent_key": "triage-agent",
			"kickoff": {"type": "message", "message": {"content": [{"type": "text", "text": "hi"}]}}
		}`))
		if err != nil {
			t.Fatal(err)
		}
		if params.AgentKey != "triage-agent" || params.KickoffType != "" {
			t.Fatalf("params = %+v", params)
		}
	})

	t.Run("malformed JSON is an error", func(t *testing.T) {
		if _, err := workflowmodel.DecodeAgentActionParams(json.RawMessage(`{`)); err == nil {
			t.Fatal("expected error")
		}
	})
}

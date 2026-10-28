package common

import (
	"encoding/json"
	"testing"
)

func outcomeIn(value string) *FilterGroup {
	return &FilterGroup{
		LogicalOperator: LogicalOperatorAnd,
		Elements:        []FilterElement{{Field: "conversation.outcome_key", Operator: ComparisonOperatorIn, Value: value}},
	}
}

func TestValidatePageActions(t *testing.T) {
	cases := []struct {
		name     string
		actions  []PageAction
		wantCode string
	}{
		{"empty ok", nil, ""},
		{"legacy action (no kind) ok", []PageAction{{Label: "Sign", Action: "sign"}}, ""},
		{"workflow with inputs ok", []PageAction{{Label: "Sync", Kind: PageActionKindWorkflow, Workflow: "sync", Inputs: map[string]string{"id": "{{ record.id }}"}}}, ""},
		{"agent ok", []PageAction{{Label: "Ask", Kind: PageActionKindAgent, Agent: "deal-creator"}}, ""},
		{"agent with prompt, context and visible_when ok", []PageAction{{Label: "Ask", Kind: PageActionKindAgent, Agent: "a", Prompt: "Hi", Context: "{{ conversation | json }}", VisibleWhen: outcomeIn("booked|callback")}}, ""},
		{"agent with context only ok", []PageAction{{Label: "Ask", Kind: PageActionKindAgent, Agent: "a", Context: "c"}}, ""},
		{"intent create ok", []PageAction{{Label: "Create deal", Kind: PageActionKindIntent, Intent: PageActionIntentCreateRecord, Entity: "deal", Prompt: "Create it"}}, ""},
		{"intent update without record ok", []PageAction{{Label: "Update deal", Kind: PageActionKindIntent, Intent: PageActionIntentUpdateRecord, Entity: "deal"}}, ""},
		{"intent update with record ok", []PageAction{{Label: "Update", Kind: PageActionKindIntent, Intent: PageActionIntentUpdateRecord, Entity: "deal", Record: "{{ record.id }}"}}, ""},
		{"link ok", []PageAction{{Label: "Open site", Kind: PageActionKindLink, Url: "https://example.com/{{ record.slug | url_encode }}"}}, ""},
		{"link with visible_when ok", []PageAction{{Label: "Open", Kind: PageActionKindLink, Url: "https://example.com", VisibleWhen: outcomeIn("x")}}, ""},
		{"visible_when on workflow ok", []PageAction{{Label: "Sync", Kind: PageActionKindWorkflow, Workflow: "w", VisibleWhen: outcomeIn("x")}}, ""},
		{"label missing", []PageAction{{Action: "sign"}}, "page_action_label_required"},
		{"action target missing", []PageAction{{Label: "x"}}, "page_action_target_required"},
		{"workflow target missing", []PageAction{{Label: "x", Kind: PageActionKindWorkflow}}, "page_action_target_required"},
		{"agent target missing", []PageAction{{Label: "x", Kind: PageActionKindAgent}}, "page_action_target_required"},
		{"intent missing", []PageAction{{Label: "x", Kind: PageActionKindIntent, Entity: "deal"}}, "page_action_target_required"},
		{"intent entity missing", []PageAction{{Label: "x", Kind: PageActionKindIntent, Intent: PageActionIntentCreateRecord}}, "page_action_target_required"},
		{"link url missing", []PageAction{{Label: "x", Kind: PageActionKindLink}}, "page_action_target_required"},
		{"link url blank", []PageAction{{Label: "x", Kind: PageActionKindLink, Url: "  "}}, "page_action_target_required"},
		{"unknown intent", []PageAction{{Label: "x", Kind: PageActionKindIntent, Intent: "ask_agent", Entity: "deal"}}, "unsupported_page_action_intent"},
		{"action sets workflow", []PageAction{{Label: "x", Action: "a", Workflow: "w"}}, "page_action_target_conflict"},
		{"agent sets action", []PageAction{{Label: "x", Kind: PageActionKindAgent, Agent: "g", Action: "a"}}, "page_action_target_conflict"},
		{"action sets agent", []PageAction{{Label: "x", Action: "a", Agent: "g"}}, "page_action_target_conflict"},
		{"action sets entity", []PageAction{{Label: "x", Action: "a", Entity: "deal"}}, "page_action_target_conflict"},
		{"create intent sets record", []PageAction{{Label: "x", Kind: PageActionKindIntent, Intent: PageActionIntentCreateRecord, Entity: "deal", Record: "r1"}}, "page_action_target_conflict"},
		{"action sets url", []PageAction{{Label: "x", Action: "a", Url: "https://example.com"}}, "page_action_target_conflict"},
		{"link sets action", []PageAction{{Label: "x", Kind: PageActionKindLink, Url: "https://example.com", Action: "a"}}, "page_action_target_conflict"},
		{"link with params", []PageAction{{Label: "x", Kind: PageActionKindLink, Url: "https://example.com", Params: map[string]string{"k": "v"}}}, "page_action_prefill_mismatch"},
		{"link with prompt", []PageAction{{Label: "x", Kind: PageActionKindLink, Url: "https://example.com", Prompt: "p"}}, "page_action_prompt_mismatch"},
		{"action with inputs", []PageAction{{Label: "x", Action: "a", Inputs: map[string]string{"k": "v"}}}, "page_action_prefill_mismatch"},
		{"workflow with params", []PageAction{{Label: "x", Kind: PageActionKindWorkflow, Workflow: "w", Params: map[string]string{"k": "v"}}}, "page_action_prefill_mismatch"},
		{"agent with params", []PageAction{{Label: "x", Kind: PageActionKindAgent, Agent: "g", Params: map[string]string{"k": "v"}}}, "page_action_prefill_mismatch"},
		{"prompt on action kind", []PageAction{{Label: "x", Action: "a", Prompt: "p"}}, "page_action_prompt_mismatch"},
		{"context on workflow kind", []PageAction{{Label: "x", Kind: PageActionKindWorkflow, Workflow: "w", Context: "c"}}, "page_action_prompt_mismatch"},
		{"visible_when unknown operator", []PageAction{{Label: "x", Action: "a", VisibleWhen: &FilterGroup{LogicalOperator: LogicalOperatorAnd, Elements: []FilterElement{{Field: "f", Operator: "like", Value: "v"}}}}}, "page_action_visible_when_invalid"},
		{"visible_when empty group", []PageAction{{Label: "x", Action: "a", VisibleWhen: &FilterGroup{LogicalOperator: LogicalOperatorAnd}}}, "page_action_visible_when_invalid"},
		{"unknown kind", []PageAction{{Label: "x", Kind: "button", Action: "a"}}, "unsupported_page_action_kind"},
	}
	for _, testCase := range cases {
		t.Run(testCase.name, func(t *testing.T) {
			err := ValidatePageActions(testCase.actions)
			if testCase.wantCode == "" {
				if err != nil {
					t.Fatalf("expected success, got %v", err)
				}
				return
			}
			validationError, ok := AsPageActionValidationError(err)
			if !ok {
				t.Fatalf("expected *PageActionValidationError, got %T (%v)", err, err)
			}
			if validationError.Code != testCase.wantCode {
				t.Fatalf("expected code %q, got %q (%s)", testCase.wantCode, validationError.Code, validationError.Message)
			}
		})
	}
}

func TestValidatePageActionsReportsTheFailingIndex(t *testing.T) {
	err := ValidatePageActions([]PageAction{{Label: "ok", Action: "a"}, {Label: "bad"}})
	validationError, ok := AsPageActionValidationError(err)
	if !ok || validationError.Index != 1 {
		t.Fatalf("expected index 1, got %v", err)
	}
	if got := validationError.Error(); got != "actions[1]: an action button requires `action` (the action slug)" {
		t.Fatalf("unexpected message %q", got)
	}
}

func TestPageActionKindNormalized(t *testing.T) {
	if PageActionKind("").Normalized() != PageActionKindAction {
		t.Fatal("empty kind should normalize to action")
	}
	if PageActionKindIntent.Normalized() != PageActionKindIntent {
		t.Fatal("intent should stay intent")
	}
}

// A button written before the new kinds existed must serialize exactly as it
// did — the CLI compares manifests by content.
func TestPageActionJSONOmitsUnsetFields(t *testing.T) {
	encoded, err := json.Marshal(PageAction{Label: "Sign", Icon: "Pen", Action: "sign"})
	if err != nil {
		t.Fatal(err)
	}
	if got, want := string(encoded), `{"label":"Sign","icon":"Pen","action":"sign"}`; got != want {
		t.Fatalf("got %s, want %s", got, want)
	}
}

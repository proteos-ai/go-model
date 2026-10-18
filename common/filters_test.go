package common

import "testing"

func TestValidateFilterGroup(t *testing.T) {
	cases := []struct {
		name     string
		group    *FilterGroup
		wantPath string
		wantCode string
	}{
		{"nil", nil, "rule", "missing"},
		{"bad logical operator", &FilterGroup{LogicalOperator: "xor", Elements: []FilterElement{{Field: "f", Operator: ComparisonOperatorEmpty}}}, "rule.logical_operator", "invalid_value"},
		{"empty group", &FilterGroup{LogicalOperator: LogicalOperatorAnd}, "rule", "missing"},
		{"blank field", &FilterGroup{LogicalOperator: LogicalOperatorAnd, Elements: []FilterElement{{Field: " ", Operator: ComparisonOperatorEmpty}}}, "rule.elements[0].field", "missing"},
		{"unknown operator", &FilterGroup{LogicalOperator: LogicalOperatorAnd, Elements: []FilterElement{{Field: "f", Operator: "like", Value: "v"}}}, "rule.elements[0].operator", "invalid_value"},
		{"missing value", &FilterGroup{LogicalOperator: LogicalOperatorAnd, Elements: []FilterElement{{Field: "f", Operator: ComparisonOperatorEquals}}}, "rule.elements[0].value", "missing"},
		{"nested", &FilterGroup{LogicalOperator: LogicalOperatorOr, Groups: []FilterGroup{{LogicalOperator: LogicalOperatorAnd, Elements: []FilterElement{{Field: "f", Operator: ComparisonOperatorEmpty}, {Field: "g", Operator: ComparisonOperatorIn}}}}}, "rule.groups[0].elements[1].value", "missing"},
	}
	for _, testCase := range cases {
		t.Run(testCase.name, func(t *testing.T) {
			problems := ValidateFilterGroup("rule", testCase.group)
			if len(problems) == 0 {
				t.Fatal("expected a problem")
			}
			if problems[0].Path != testCase.wantPath || problems[0].Code != testCase.wantCode {
				t.Fatalf("got %s/%s, want %s/%s", problems[0].Path, problems[0].Code, testCase.wantPath, testCase.wantCode)
			}
		})
	}
}

func TestValidateFilterGroupAcceptsPresenceRules(t *testing.T) {
	group := &FilterGroup{LogicalOperator: LogicalOperatorAnd, Elements: []FilterElement{
		{Field: "conversation.outcome_key", Operator: ComparisonOperatorIn, Value: "interested|callback"},
		{Field: "records_by_entity.deal", Operator: ComparisonOperatorEmpty},
	}}
	if problems := ValidateFilterGroup("visible_when", group); len(problems) != 0 {
		t.Fatalf("expected no problems, got %v", problems)
	}
}

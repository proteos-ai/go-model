package metamodel

import (
	"testing"

	conversationmodel "go.proteos.ai/model/conversation"
)

func TestContactBindingDefaults(t *testing.T) {
	if got := (ContactBinding{}).EffectiveDuplicatePolicy(); got != conversationmodel.RecordDuplicatePolicyFlag {
		t.Fatalf("unset policy must default to flag, got %q", got)
	}
	if got := (ContactBinding{}).WithDefaults().DuplicatePolicy; got != conversationmodel.RecordDuplicatePolicyFlag {
		t.Fatalf("WithDefaults must fill flag, got %q", got)
	}
	reject := ContactBinding{DuplicatePolicy: conversationmodel.RecordDuplicatePolicyReject}
	if got := reject.WithDefaults().DuplicatePolicy; got != conversationmodel.RecordDuplicatePolicyReject {
		t.Fatalf("WithDefaults must keep an explicit policy, got %q", got)
	}
}

func TestValidateContactBinding(t *testing.T) {
	for _, policy := range []conversationmodel.RecordDuplicatePolicy{"", "reject", "flag"} {
		if err := ValidateContactBinding(ContactBinding{DuplicatePolicy: policy}); err != nil {
			t.Fatalf("policy %q must be accepted, got %v", policy, err)
		}
	}
	if err := ValidateContactBinding(ContactBinding{DuplicatePolicy: "merge"}); err == nil {
		t.Fatalf("unknown policy must be rejected")
	}
}

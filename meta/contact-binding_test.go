package metamodel

import (
	"bytes"
	"encoding/json"
	"testing"

	conversationmodel "go.proteos.ai/model/conversation"
)

func TestEffectiveDuplicatePolicy(t *testing.T) {
	// nil is the common case: an entity that declares no binding at all.
	if got := EffectiveDuplicatePolicy(nil); got != conversationmodel.RecordDuplicatePolicyFlag {
		t.Fatalf("no binding must default to flag, got %q", got)
	}
	if got := EffectiveDuplicatePolicy(&ContactBinding{}); got != conversationmodel.RecordDuplicatePolicyFlag {
		t.Fatalf("unset policy must default to flag, got %q", got)
	}
	reject := ContactBinding{DuplicatePolicy: conversationmodel.RecordDuplicatePolicyReject}
	if got := EffectiveDuplicatePolicy(&reject); got != conversationmodel.RecordDuplicatePolicyReject {
		t.Fatalf("an explicit policy must be kept, got %q", got)
	}
}

func TestContactBindingIsZero(t *testing.T) {
	if !(ContactBinding{}).IsZero() {
		t.Fatalf("an unset binding must report zero")
	}
	if (ContactBinding{DuplicatePolicy: conversationmodel.RecordDuplicatePolicyFlag}).IsZero() {
		t.Fatalf("an explicit flag policy must not report zero")
	}
}

// An entity that declares no binding must serialize WITHOUT the key — that
// absence is what keeps a manifest omitting contact_binding free of drift.
func TestEntityOmitsAbsentContactBinding(t *testing.T) {
	encoded, err := json.Marshal(Entity{Slug: "invoice"})
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	if bytes.Contains(encoded, []byte("contact_binding")) {
		t.Fatalf("absent contact_binding must not be serialized, got %s", encoded)
	}

	declared := Entity{Slug: "lead", ContactBinding: &ContactBinding{
		DuplicatePolicy: conversationmodel.RecordDuplicatePolicyReject,
	}}
	encoded, err = json.Marshal(declared)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	if !bytes.Contains(encoded, []byte(`"contact_binding":{"duplicate_policy":"reject"}`)) {
		t.Fatalf("declared contact_binding must be serialized, got %s", encoded)
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

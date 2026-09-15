package datamodel

import (
	"encoding/json"
	"testing"

	conversationmodel "go.proteos.ai/model/conversation"
)

func TestContactBindingEvidenceRoundTrip(t *testing.T) {
	evidence := ContactBindingEvidence{ContactId: "c-1", AddressKeys: []conversationmodel.ContactAddressKey{
		{Kind: conversationmodel.ContactAddressKindEmail, Value: "ada@example.com"}}}
	inProcess, ok := ParseContactBindingEvidence(evidence.ToEvidence())
	if !ok || inProcess.ContactId != "c-1" || len(inProcess.AddressKeys) != 1 {
		t.Fatalf("in-process bag must parse, got %+v %v", inProcess, ok)
	}
	raw, _ := json.Marshal(evidence.ToEvidence())
	var decoded map[string]any
	if err := json.Unmarshal(raw, &decoded); err != nil {
		t.Fatal(err)
	}
	afterJSON, ok := ParseContactBindingEvidence(decoded)
	if !ok || afterJSON.AddressKeys[0].Value != "ada@example.com" {
		t.Fatalf("JSON round trip must parse, got %+v %v", afterJSON, ok)
	}
	if _, ok := ParseContactBindingEvidence(map[string]any{"address_keys": []any{}}); ok {
		t.Fatalf("a bag without contact_id is not contact-binding evidence")
	}
}

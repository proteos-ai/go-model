package datamodel

import (
	"encoding/json"
	"strings"
	"testing"
)

func TestRecordDraftEnums(t *testing.T) {
	for _, intent := range []RecordDraftIntent{RecordDraftIntentCreate, RecordDraftIntentUpdate, RecordDraftIntentDelete, RecordDraftIntentUpsert} {
		if !intent.IsValid() {
			t.Fatalf("%q must be a valid intent", intent)
		}
	}
	if RecordDraftIntent("patch").IsValid() {
		t.Fatal("unknown intent must be invalid")
	}
	for _, status := range []RecordDraftStatus{RecordDraftStatusProposed, RecordDraftStatusApproved, RecordDraftStatusApprovedWithChanges, RecordDraftStatusRejected} {
		if !status.IsValid() {
			t.Fatalf("%q must be a valid status", status)
		}
	}
	if RecordDraftStatusProposed.IsTerminal() {
		t.Fatal("proposed is the open state")
	}
	for _, status := range []RecordDraftStatus{RecordDraftStatusApproved, RecordDraftStatusApprovedWithChanges, RecordDraftStatusRejected} {
		if !status.IsTerminal() {
			t.Fatalf("%q must be terminal", status)
		}
	}
	if RecordDraftStatus("done").IsTerminal() || RecordDraftStatus("done").IsValid() {
		t.Fatal("unknown status is neither valid nor terminal")
	}
}

// A change's three values are never omitted: a null is "cleared", and a
// missing key would read the same.
func TestRecordDraftChangeKeepsNullValues(t *testing.T) {
	raw, err := json.Marshal(RecordDraftChange{TransactionId: "1", Kind: RecordDraftChangeAttributeChanged, Attribute: "stage", Proposed: "won", Applied: nil})
	if err != nil {
		t.Fatal(err)
	}
	body := string(raw)
	for _, key := range []string{`"current":null`, `"proposed":"won"`, `"applied":null`} {
		if !strings.Contains(body, key) {
			t.Fatalf("change JSON must carry %s, got %s", key, body)
		}
	}
	dropped, _ := json.Marshal(RecordDraftChange{TransactionId: "2", Kind: RecordDraftChangeItemDropped})
	if strings.Contains(string(dropped), `"attribute"`) {
		t.Fatalf("an item_dropped change has no attribute, got %s", dropped)
	}
}

func TestRecordDraftItemRoundTrip(t *testing.T) {
	item := RecordDraftItem{TransactionId: "1", Id: "r-1", Data: map[string]any{"amount": "12.50", "is_active": true}}
	raw, _ := json.Marshal(item)
	var decoded RecordDraftItem
	if err := json.Unmarshal(raw, &decoded); err != nil {
		t.Fatal(err)
	}
	if decoded.Id != "r-1" || decoded.Data["amount"] != "12.50" || decoded.Data["is_active"] != true {
		t.Fatalf("round trip lost fields: %+v", decoded)
	}
	create, _ := json.Marshal(RecordDraftItem{TransactionId: "1", Data: map[string]any{"name": "x"}})
	if strings.Contains(string(create), `"id"`) {
		t.Fatalf("a create item carries no id key, got %s", create)
	}
}

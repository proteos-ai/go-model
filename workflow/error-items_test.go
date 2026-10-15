package workflowmodel_test

import (
	"testing"

	workflowmodel "go.proteos.ai/model/workflow"
)

func TestErrorItems(t *testing.T) {
	inputs := []workflowmodel.Item{
		{Json: map[string]any{"id": "a", "error": "stale"}, Binary: map[string]workflowmodel.BinaryRef{"file": {FileId: "f-1"}}},
		{Json: map[string]any{"id": "b"}},
	}
	failure := workflowmodel.ExecutionError{
		Code:        "upstream_unavailable",
		Message:     "conversations: POST /channel-actions status 502",
		Details:     map[string]any{"status": 502},
		IsUserError: false,
	}

	items := workflowmodel.ErrorItems("input_2", inputs, failure)
	if len(items) != 2 {
		t.Fatalf("expected 2 error items, got %d", len(items))
	}
	for index, item := range items {
		if item.PairedItem == nil || item.PairedItem.Item != index || item.PairedItem.Input != "input_2" {
			t.Errorf("item %d: paired_item = %+v, want item %d", index, item.PairedItem, index)
		}
		errorValue, ok := item.Json["error"].(map[string]any)
		if !ok {
			t.Fatalf("item %d: json.error = %#v, want object", index, item.Json["error"])
		}
		if errorValue["code"] != failure.Code || errorValue["message"] != failure.Message ||
			errorValue["is_user_error"] != false {
			t.Errorf("item %d: json.error = %#v", index, errorValue)
		}
		details, _ := errorValue["details"].(map[string]any)
		if details["status"] != 502 {
			t.Errorf("item %d: json.error.details = %#v", index, errorValue["details"])
		}
	}
	if items[0].Json["id"] != "a" || items[1].Json["id"] != "b" {
		t.Errorf("input json not preserved: %#v / %#v", items[0].Json, items[1].Json)
	}
	if items[0].Binary["file"].FileId != "f-1" {
		t.Errorf("binary refs not carried: %#v", items[0].Binary)
	}
	// The input is never mutated.
	if inputs[0].Json["error"] != "stale" {
		t.Errorf("input mutated: %#v", inputs[0].Json)
	}

	if got := workflowmodel.ErrorItems("main", nil, failure); len(got) != 0 {
		t.Errorf("nil input should yield no items, got %d", len(got))
	}
}

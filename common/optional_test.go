package common

import (
	"encoding/json"
	"testing"
	"time"
)

func TestOptionalTimeTriState(t *testing.T) {
	type patch struct {
		ValidUntil Optional[time.Time] `json:"valid_until"`
	}

	t.Run("absent key leaves the field not-present", func(t *testing.T) {
		var got patch
		if err := json.Unmarshal([]byte(`{}`), &got); err != nil {
			t.Fatalf("unmarshal: %v", err)
		}
		if got.ValidUntil.Present {
			t.Fatalf("expected Present=false for an absent key")
		}
		if got.ValidUntil.Value != nil {
			t.Fatalf("expected Value=nil for an absent key")
		}
	})

	t.Run("explicit null marks present with nil value (clear)", func(t *testing.T) {
		var got patch
		if err := json.Unmarshal([]byte(`{"valid_until": null}`), &got); err != nil {
			t.Fatalf("unmarshal: %v", err)
		}
		if !got.ValidUntil.Present {
			t.Fatalf("expected Present=true for an explicit null")
		}
		if got.ValidUntil.Value != nil {
			t.Fatalf("expected Value=nil for an explicit null")
		}
	})

	t.Run("value marks present with the decoded time (set)", func(t *testing.T) {
		var got patch
		if err := json.Unmarshal([]byte(`{"valid_until": "2026-06-22T00:00:00Z"}`), &got); err != nil {
			t.Fatalf("unmarshal: %v", err)
		}
		if !got.ValidUntil.Present {
			t.Fatalf("expected Present=true for a value")
		}
		if got.ValidUntil.Value == nil {
			t.Fatalf("expected a non-nil Value for a value")
		}
		want := time.Date(2026, 6, 22, 0, 0, 0, 0, time.UTC)
		if !got.ValidUntil.Value.Equal(want) {
			t.Fatalf("expected %v, got %v", want, *got.ValidUntil.Value)
		}
	})

	t.Run("invalid datetime surfaces an error", func(t *testing.T) {
		var got patch
		if err := json.Unmarshal([]byte(`{"valid_until": "not-a-date"}`), &got); err == nil {
			t.Fatalf("expected an error for an invalid datetime")
		}
	})
}

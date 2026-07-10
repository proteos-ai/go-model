package common

import (
	"bytes"
	"encoding/json"
)

// Optional is a presence-aware nullable wrapper for tri-state PATCH fields, where
// "absent", "explicit null", and "a value" must be told apart. A plain pointer
// cannot: both an omitted key and a JSON `null` unmarshal to nil.
//
//	absent        → {Present: false, Value: nil}  → leave the field unchanged
//	"key": null   → {Present: true,  Value: nil}  → clear the column to NULL
//	"key": <val>  → {Present: true,  Value: &val} → set the column to <val>
//
// encoding/json only calls UnmarshalJSON for keys actually present in the body,
// so the zero value ({false, nil}) correctly means "absent". Use it on update
// request structs for columns that can be both set and cleared; mark the field
// `bun:"-"` so the ORM ignores it and the repository applies it explicitly.
type Optional[T any] struct {
	Present bool
	Value   *T
}

// UnmarshalJSON records that the key was present and decodes its value. A literal
// JSON `null` leaves Value nil (a clear); any other value is decoded into *Value.
func (optional *Optional[T]) UnmarshalJSON(data []byte) error {
	optional.Present = true
	if bytes.Equal(bytes.TrimSpace(data), []byte("null")) {
		optional.Value = nil
		return nil
	}
	var value T
	if err := json.Unmarshal(data, &value); err != nil {
		return err
	}
	optional.Value = &value
	return nil
}

// MarshalJSON renders nil as `null` and otherwise the wrapped value. Optional is
// a request-side type; this exists mainly so round-tripping in tests is honest.
func (optional Optional[T]) MarshalJSON() ([]byte, error) {
	return json.Marshal(optional.Value)
}

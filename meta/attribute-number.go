package metamodel

// NumberAttributeMeta holds number-specific validation options
// Used for both "number" (float) and "integer" types
type NumberAttributeMeta struct {
	Minimum          *float64 `json:"minimum,omitempty"`
	Maximum          *float64 `json:"maximum,omitempty"`
	ExclusiveMinimum *float64 `json:"exclusive_minimum,omitempty"`
	ExclusiveMaximum *float64 `json:"exclusive_maximum,omitempty"`
	MultipleOf       *float64 `json:"multiple_of,omitempty"`
}

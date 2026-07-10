package metamodel

// EnumValue represents a single allowed value in an enum
type EnumValue struct {
	Value       string `json:"value"`                 // The actual value stored
	Label       string `json:"label,omitempty"`       // Display label (defaults to Value if empty)
	Description string `json:"description,omitempty"` // Optional description/help text
}

// EnumAttributeMeta holds enum-specific options with inline values
type EnumAttributeMeta struct {
	Values []EnumValue `json:"values"`
}

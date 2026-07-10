package metamodel

// JSONSchema represents a JSON Schema document following draft-07 specification.
// See: https://json-schema.org/draft-07/json-schema-validation
type JSONSchema struct {
	// Meta
	Schema      string `json:"$schema,omitempty"`
	ID          string `json:"$id,omitempty"`
	Title       string `json:"title,omitempty"`
	Description string `json:"description,omitempty"`

	// Type
	Type string `json:"type,omitempty"`

	// Object validation
	Properties           map[string]*JSONSchema `json:"properties,omitempty"`
	Required             []string               `json:"required,omitempty"`
	AdditionalProperties *bool                  `json:"additionalProperties,omitempty"`

	// Array validation
	Items       *JSONSchema `json:"items,omitempty"`
	MinItems    *int        `json:"minItems,omitempty"`
	MaxItems    *int        `json:"maxItems,omitempty"`
	UniqueItems *bool       `json:"uniqueItems,omitempty"`

	// String validation
	MinLength *int   `json:"minLength,omitempty"`
	MaxLength *int   `json:"maxLength,omitempty"`
	Pattern   string `json:"pattern,omitempty"`
	Format    string `json:"format,omitempty"`

	// Number validation
	Minimum          *float64 `json:"minimum,omitempty"`
	Maximum          *float64 `json:"maximum,omitempty"`
	ExclusiveMinimum *float64 `json:"exclusiveMinimum,omitempty"`
	ExclusiveMaximum *float64 `json:"exclusiveMaximum,omitempty"`
	MultipleOf       *float64 `json:"multipleOf,omitempty"`

	// Enum
	Enum []any `json:"enum,omitempty"`

	// Default value
	Default any `json:"default,omitempty"`

	// Read-only
	ReadOnly *bool `json:"readOnly,omitempty"`
}

// Draft07SchemaURI is the URI for JSON Schema draft-07
const Draft07SchemaURI = "http://json-schema.org/draft-07/schema#"

package metamodel

// AttributeType represents JSON Schema compatible types plus extensions
type AttributeType string

const (
	AttributeTypeString  AttributeType = "string"
	AttributeTypeNumber  AttributeType = "number"
	AttributeTypeInteger AttributeType = "integer"
	AttributeTypeBoolean AttributeType = "boolean"

	AttributeTypeDatetime AttributeType = "datetime" // Date time values
	AttributeTypeEnum     AttributeType = "enum"     // Inline enum definition
	AttributeTypeArray    AttributeType = "array"    // Array of items
	AttributeTypeObject   AttributeType = "object"   // Nested object

	AttributeTypeRelation      AttributeType = "relation"       // Foreign-key reference to another entity
	AttributeTypeUser          AttributeType = "user"           // Reference to a platform user (account-service)
	AttributeTypeCurrency      AttributeType = "currency"       // Monetary amount in a currency: {amount, currency_code}
	AttributeTypeKnowledgeText AttributeType = "knowledge-text" // Long text stored as a knowledge node; record holds {id}
	AttributeTypeFile          AttributeType = "file"           // Reference to a stored file (storage-service); record holds {id, name}
)

// Attribute represents a typed property of an entity. Mirrors the JSON
// Schema-ish shape the platform uses; `Meta` carries the type-specific
// validation/config (see the AttributeMeta variants).
type Attribute struct {
	Name        string        `json:"name"`
	Type        AttributeType `json:"type"`
	Label       string        `json:"label,omitempty"`
	Description string        `json:"description,omitempty"`

	// Validation flags
	IsRequired bool `json:"is_required"`         // Is this field required?
	IsNullable bool `json:"is_nullable"`         // Can be null?
	IsUnique   bool `json:"is_unique,omitempty"` // Must be unique across records

	// IsReadOnly marks server-managed fields (id, created_at, updated_at). These
	// are surfaced in schemas but rejected on create/update payloads.
	IsReadOnly bool `json:"is_read_only,omitempty"` // Field cannot be modified after creation

	// IsPlatformManaged marks fields whose lifecycle the platform controls — the
	// canonical platform attributes (id, created_at, updated_at, created_by,
	// updated_by). They are auto-added to every entity and cannot be removed or
	// redefined by users without elevated permission. This is the cross-cutting
	// "platform-managed" flag (the same concept will later mark platform-managed
	// entities, permissions, etc.).
	IsPlatformManaged bool `json:"is_platform_managed,omitempty"`

	// === Type-specific metadata ===
	// The concrete type is determined by the Type field
	Meta any `json:"meta,omitempty"`

	// === Defaults ===
	// DefaultValue carries the attribute's default; serialised as `default_value`
	// on the wire (the meta map preserves the original).
	DefaultValue any `json:"default_value,omitempty"`

	// Options is the legacy/raw attribute config map (kept for compatibility
	// with callers that read the unstructured shape).
	Options map[string]any `json:"options,omitempty"`
}

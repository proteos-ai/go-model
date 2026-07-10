package metamodel

import (
	"encoding/json"
)

// EntityToJSONSchema converts an Entity to a JSON Schema draft-07 document.
// The resulting schema validates data instances that conform to the entity's attribute definitions.
func EntityToJSONSchema(entity Entity) *JSONSchema {
	schema := &JSONSchema{
		Schema:      Draft07SchemaURI,
		ID:          entity.Slug,
		Title:       entity.Name,
		Description: entity.Description,
		Type:        "object",
		Properties:  make(map[string]*JSONSchema),
		Required:    []string{},
	}

	for _, attr := range entity.Attributes {
		propSchema := AttributeToJSONSchema(attr)
		schema.Properties[attr.Name] = propSchema

		if attr.IsRequired {
			schema.Required = append(schema.Required, attr.Name)
		}
	}

	// Only include required array if there are required fields
	if len(schema.Required) == 0 {
		schema.Required = nil
	}

	return schema
}

// AttributesToJSONSchema converts a bare slice of attributes (no enclosing Entity) into a
// JSON Schema object document — { type: "object", properties: <each attr>, required: <IsRequired names> }.
// Used for entity-less attribute sets such as an action's params/returns or a prompt's inputs.
func AttributesToJSONSchema(attrs []Attribute) *JSONSchema {
	schema := &JSONSchema{
		Type:       "object",
		Properties: make(map[string]*JSONSchema),
		Required:   []string{},
	}

	for _, attr := range attrs {
		schema.Properties[attr.Name] = AttributeToJSONSchema(attr)

		if attr.IsRequired {
			schema.Required = append(schema.Required, attr.Name)
		}
	}

	if len(schema.Required) == 0 {
		schema.Required = nil
	}

	return schema
}

// AttributeToJSONSchema converts a single attribute to a JSON Schema property definition.
func AttributeToJSONSchema(attr Attribute) *JSONSchema {
	schema := &JSONSchema{
		Title:       attr.Label,
		Description: attr.Description,
		Default:     attr.DefaultValue,
	}

	if attr.IsReadOnly {
		readOnly := true
		schema.ReadOnly = &readOnly
	}

	switch attr.Type {
	case AttributeTypeString:
		schema.Type = "string"
		applyStringMeta(schema, attr.Meta)

	case AttributeTypeInteger:
		schema.Type = "integer"
		applyNumberMeta(schema, attr.Meta)

	case AttributeTypeNumber:
		schema.Type = "number"
		applyNumberMeta(schema, attr.Meta)

	case AttributeTypeBoolean:
		schema.Type = "boolean"

	case AttributeTypeDatetime:
		schema.Type = "string"
		applyDatetimeMeta(schema, attr.Meta)

	case AttributeTypeEnum:
		schema.Type = "string"
		applyEnumMeta(schema, attr.Meta)

	case AttributeTypeArray:
		schema.Type = "array"
		applyArrayMeta(schema, attr.Meta)

	case AttributeTypeObject:
		schema.Type = "object"
		applyObjectMeta(schema, attr.Meta)

	case AttributeTypeRelation:
		// A relation is stored as the scalar value of the related entity's
		// referenced attribute (default `id`, which today is a string UUID).
		// The host attribute's wire type is therefore "string" by default;
		// future work can resolve the related attribute and inherit its type.
		schema.Type = "string"

	case AttributeTypeUser:
		// A user attribute stores the composite { type, id } (common.UserRef) —
		// the account-service identity plus its kind. The JSON Schema is a
		// backwards-compat export only — authoritative record-value enforcement
		// lives in data-service's recordvalidation.
		schema.Type = "object"
		applyUserMeta(schema)

	case AttributeTypeCurrency:
		// A currency value is the composite { amount, currency_code }. The
		// JSON Schema is a backwards-compat export only — authoritative
		// record-value enforcement lives in data-service's recordvalidation.
		schema.Type = "object"
		applyCurrencyMeta(schema)

	case AttributeTypeKnowledgeText:
		// A knowledge-text value is the composite { id, content }
		// (common.KnowledgeNodeRef) — the stored shape is { id } and content
		// is filled on single-record reads. The JSON Schema is a
		// backwards-compat export only — authoritative record-value
		// enforcement lives in data-service's recordvalidation.
		schema.Type = "object"
		applyKnowledgeTextMeta(schema)

	case AttributeTypeFile:
		// A file value is the composite { id, name } (common.FileRef) — id is
		// the storage-service file id, name the denormalised filename. The JSON
		// Schema is a backwards-compat export only — authoritative record-value
		// enforcement lives in data-service's recordvalidation.
		schema.Type = "object"
		applyFileMeta(schema)
	}

	// Handle nullable by wrapping type
	if attr.IsNullable {
		schema = wrapNullable(schema)
	}

	return schema
}

// Regex patterns for formats not natively supported in JSON Schema draft-07
const (
	// UUIDPattern matches UUID v1-5 format (8-4-4-4-12 hex digits)
	UUIDPattern = "^[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}$"

	// DurationPattern matches ISO 8601 duration format (e.g., P1Y2M3DT4H5M6S)
	// At least one component must be present. Uses alternation to avoid lookahead.
	// Matches: P followed by (date parts with optional time | just time parts)
	DurationPattern = "^P(([0-9]+Y([0-9]+M)?([0-9]+W)?([0-9]+D)?(T([0-9]+H)?([0-9]+M)?([0-9]+(\\.[0-9]+)?S)?)?)|([0-9]+M([0-9]+W)?([0-9]+D)?(T([0-9]+H)?([0-9]+M)?([0-9]+(\\.[0-9]+)?S)?)?)|([0-9]+W([0-9]+D)?(T([0-9]+H)?([0-9]+M)?([0-9]+(\\.[0-9]+)?S)?)?)|([0-9]+D(T([0-9]+H)?([0-9]+M)?([0-9]+(\\.[0-9]+)?S)?)?)|T(([0-9]+H([0-9]+M)?([0-9]+(\\.[0-9]+)?S)?)|([0-9]+M([0-9]+(\\.[0-9]+)?S)?)|([0-9]+(\\.[0-9]+)?S)))$"
)

// formatsNotInDraft07 lists formats that need regex fallback patterns
var formatsNotInDraft07 = map[StringFormat]string{
	StringFormatUUID: UUIDPattern,
}

// applyStringMeta applies string-specific metadata to the schema
func applyStringMeta(schema *JSONSchema, metaValue any) {
	if metaValue == nil {
		return
	}

	stringMeta := ParseMetaAs[StringAttributeMeta](metaValue)
	if stringMeta == nil {
		return
	}

	schema.MinLength = stringMeta.MinLength
	schema.MaxLength = stringMeta.MaxLength
	schema.Pattern = stringMeta.Pattern

	if stringMeta.Format != "" {
		// Check if this format needs a regex fallback for draft-07 compatibility
		if pattern, needsFallback := formatsNotInDraft07[stringMeta.Format]; needsFallback {
			// Use regex pattern instead of unsupported format for draft-07 validation
			// Only set pattern if user hasn't already specified one
			if schema.Pattern == "" {
				schema.Pattern = pattern
			}
			// Note: We don't set the format field for unsupported formats because
			// strict draft-07 validators will reject unknown formats.
			// The pattern provides equivalent validation.
		} else {
			schema.Format = string(stringMeta.Format)
		}
	}
}

// applyNumberMeta applies number-specific metadata to the schema
func applyNumberMeta(schema *JSONSchema, metaValue any) {
	if metaValue == nil {
		return
	}

	numberMeta := ParseMetaAs[NumberAttributeMeta](metaValue)
	if numberMeta == nil {
		return
	}

	schema.Minimum = numberMeta.Minimum
	schema.Maximum = numberMeta.Maximum
	schema.ExclusiveMinimum = numberMeta.ExclusiveMinimum
	schema.ExclusiveMaximum = numberMeta.ExclusiveMaximum
	schema.MultipleOf = numberMeta.MultipleOf
}

// applyDatetimeMeta applies datetime-specific metadata to the schema
func applyDatetimeMeta(schema *JSONSchema, metaValue any) {
	if metaValue == nil {
		// Default to date-time format
		schema.Format = "date-time"
		return
	}

	datetimeMeta := ParseMetaAs[DatetimeAttributeMeta](metaValue)
	if datetimeMeta == nil {
		schema.Format = "date-time"
		return
	}

	// Map our datetime formats to JSON Schema formats
	switch datetimeMeta.Format {
	case DatetimeFormatDateTime:
		schema.Format = "date-time"
	case DatetimeFormatDate:
		schema.Format = "date"
	case DatetimeFormatTime:
		schema.Format = "time"
	case DatetimeFormatDuration:
		// Duration is not in draft-07, use regex pattern fallback instead
		// We don't set format because strict draft-07 validators reject unknown formats
		schema.Pattern = DurationPattern
	default:
		schema.Format = "date-time"
	}

	// Note: JSON Schema draft-07 doesn't have native min/max for date-time strings
	// These would need to be validated at application level or using custom formats
}

// applyEnumMeta applies enum-specific metadata to the schema
func applyEnumMeta(schema *JSONSchema, metaValue any) {
	if metaValue == nil {
		return
	}

	enumMeta := ParseMetaAs[EnumAttributeMeta](metaValue)
	if enumMeta == nil {
		return
	}

	// Extract just the values for the enum constraint
	enumValues := make([]any, len(enumMeta.Values))
	for i, v := range enumMeta.Values {
		enumValues[i] = v.Value
	}
	schema.Enum = enumValues
}

// applyArrayMeta applies array-specific metadata to the schema
func applyArrayMeta(schema *JSONSchema, metaValue any) {
	if metaValue == nil {
		return
	}

	arrayMeta := ParseMetaAs[ArrayAttributeMeta](metaValue)
	if arrayMeta == nil {
		return
	}

	schema.MinItems = arrayMeta.MinItems
	schema.MaxItems = arrayMeta.MaxItems

	if arrayMeta.ItemsMustBeUnique {
		uniqueItems := true
		schema.UniqueItems = &uniqueItems
	}

	// Recursively convert items schema
	if arrayMeta.Items != nil {
		schema.Items = AttributeToJSONSchema(*arrayMeta.Items)
	}
}

// applyObjectMeta applies object-specific metadata to the schema
func applyObjectMeta(schema *JSONSchema, metaValue any) {
	if metaValue == nil {
		return
	}

	objectMeta := ParseMetaAs[ObjectAttributeMeta](metaValue)
	if objectMeta == nil {
		return
	}

	if len(objectMeta.Attributes) > 0 {
		schema.Properties = make(map[string]*JSONSchema)
		schema.Required = []string{}

		for _, nestedAttr := range objectMeta.Attributes {
			propSchema := AttributeToJSONSchema(nestedAttr)
			schema.Properties[nestedAttr.Name] = propSchema

			if nestedAttr.IsRequired {
				schema.Required = append(schema.Required, nestedAttr.Name)
			}
		}

		if len(schema.Required) == 0 {
			schema.Required = nil
		}
	}
}

// applyCurrencyMeta builds the object sub-schema for a currency value —
// { amount: decimal-string, currency_code: 3-letter code }, both required.
// This is a backwards-compat export only: ISO-code validity and the
// allowed-currencies allow-list are enforced in data-service's
// recordvalidation, so no ISO enum is embedded here.
func applyCurrencyMeta(schema *JSONSchema) {
	schema.Properties = map[string]*JSONSchema{
		"amount":        {Type: "string", Pattern: `^-?\d+(\.\d+)?$`},
		"currency_code": {Type: "string", Pattern: "^[A-Z]{3}$"},
	}
	schema.Required = []string{"amount", "currency_code"}
}

// applyUserMeta shapes the JSON Schema for a `user` value — the composite
// { type, id } (common.UserRef): id is the platform user id, type is the user
// kind (person | agent | api).
func applyUserMeta(schema *JSONSchema) {
	schema.Properties = map[string]*JSONSchema{
		"type": {Type: "string", Enum: []any{"person", "agent", "api"}},
		"id":   {Type: "string"},
	}
	schema.Required = []string{"type", "id"}
}

// applyKnowledgeTextMeta shapes the JSON Schema for a `knowledge-text` value —
// the composite { id, content } (common.KnowledgeNodeRef): id is the knowledge
// node id (always present once stored), content is the transient text body
// (present on single-record reads and on writes).
func applyKnowledgeTextMeta(schema *JSONSchema) {
	schema.Properties = map[string]*JSONSchema{
		"id":      {Type: "string"},
		"content": {Type: "string"},
	}
	schema.Required = []string{"id"}
}

// applyFileMeta shapes the JSON Schema for a `file` value — the composite
// { id, name } (common.FileRef): id is the storage-service file id and name is
// the denormalised filename. Both are required.
func applyFileMeta(schema *JSONSchema) {
	schema.Properties = map[string]*JSONSchema{
		"id":   {Type: "string"},
		"name": {Type: "string"},
	}
	schema.Required = []string{"id", "name"}
}

// wrapNullable wraps the schema to allow null values.
// In JSON Schema draft-07, this is done using anyOf with null type.
// However, for simplicity and wider compatibility, we just note that
// the schema allows null - validators should handle this appropriately.
// A more complete implementation would use: {"anyOf": [<schema>, {"type": "null"}]}
func wrapNullable(schema *JSONSchema) *JSONSchema {
	// For draft-07, nullable is often handled by validators reading this annotation
	// or by using anyOf. We keep it simple here.
	return schema
}

// ParseMetaAs attempts to parse the meta field as the specified type.
// The meta field can be either the struct directly or a map from JSON unmarshaling.
func ParseMetaAs[T any](metaValue any) *T {
	if metaValue == nil {
		return nil
	}

	// Try direct type assertion first
	if typed, ok := metaValue.(T); ok {
		return &typed
	}

	// Try pointer type assertion
	if typed, ok := metaValue.(*T); ok {
		return typed
	}

	// If it's a map (from JSON unmarshaling), re-marshal and unmarshal
	if _, ok := metaValue.(map[string]any); ok {
		jsonBytes, err := json.Marshal(metaValue)
		if err != nil {
			return nil
		}

		var result T
		if err := json.Unmarshal(jsonBytes, &result); err != nil {
			return nil
		}
		return &result
	}

	return nil
}

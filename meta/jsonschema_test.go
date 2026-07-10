package metamodel

import (
	"encoding/json"
	"strings"
	"testing"

	"github.com/santhosh-tekuri/jsonschema"
)

// compileSchema compiles the generated JSONSchema using the jsonschema library
// to ensure it's a valid JSON Schema draft-07 document
func compileSchema(t *testing.T, schema *JSONSchema) *jsonschema.Schema {
	t.Helper()

	jsonBytes, err := json.Marshal(schema)
	if err != nil {
		t.Fatalf("failed to marshal schema: %v", err)
	}

	compiler := jsonschema.NewCompiler()
	if err := compiler.AddResource("schema.json", strings.NewReader(string(jsonBytes))); err != nil {
		t.Fatalf("failed to add schema resource: %v", err)
	}

	compiled, err := compiler.Compile("schema.json")
	if err != nil {
		t.Fatalf("failed to compile schema: %v\nSchema JSON: %s", err, string(jsonBytes))
	}

	return compiled
}

// validateData validates JSON data against a compiled schema
func validateData(t *testing.T, compiled *jsonschema.Schema, data string) error {
	t.Helper()
	return compiled.Validate(strings.NewReader(data))
}

func TestEntityToJSONSchema_BasicEntity(t *testing.T) {
	entity := Entity{
		Slug:        "user",
		Name:        "User",
		Description: "A user entity",
		Attributes: []Attribute{
			{
				Name:       "name",
				Type:       AttributeTypeString,
				Label:      "Name",
				IsRequired: true,
			},
			{
				Name:  "age",
				Type:  AttributeTypeInteger,
				Label: "Age",
			},
		},
	}

	schema := EntityToJSONSchema(entity)

	if schema.Schema != Draft07SchemaURI {
		t.Errorf("expected $schema to be %s, got %s", Draft07SchemaURI, schema.Schema)
	}
	if schema.ID != "user" {
		t.Errorf("expected $id to be 'user', got %s", schema.ID)
	}
	if schema.Title != "User" {
		t.Errorf("expected title to be 'User', got %s", schema.Title)
	}
	if schema.Type != "object" {
		t.Errorf("expected type to be 'object', got %s", schema.Type)
	}
	if len(schema.Properties) != 2 {
		t.Errorf("expected 2 properties, got %d", len(schema.Properties))
	}
	if len(schema.Required) != 1 || schema.Required[0] != "name" {
		t.Errorf("expected required to contain 'name', got %v", schema.Required)
	}

	// Validate schema compiles and can validate data
	compiled := compileSchema(t, schema)

	// Valid data
	if err := validateData(t, compiled, `{"name": "John", "age": 30}`); err != nil {
		t.Errorf("expected valid data to pass: %v", err)
	}

	// Missing required field
	if err := validateData(t, compiled, `{"age": 30}`); err == nil {
		t.Error("expected validation to fail for missing required field 'name'")
	}

	// Wrong type for age
	if err := validateData(t, compiled, `{"name": "John", "age": "thirty"}`); err == nil {
		t.Error("expected validation to fail for wrong type")
	}
}

func TestEntityToJSONSchema_StringAttribute(t *testing.T) {
	minLen := 1
	maxLen := 100
	entity := Entity{
		Slug: "test",
		Attributes: []Attribute{
			{
				Name: "username",
				Type: AttributeTypeString,
				Meta: StringAttributeMeta{
					MinLength: &minLen,
					MaxLength: &maxLen,
					Pattern:   "^[a-z]+$",
				},
			},
		},
	}

	schema := EntityToJSONSchema(entity)
	prop := schema.Properties["username"]

	if prop.Type != "string" {
		t.Errorf("expected type 'string', got %s", prop.Type)
	}
	if prop.MinLength == nil || *prop.MinLength != 1 {
		t.Errorf("expected minLength 1, got %v", prop.MinLength)
	}
	if prop.MaxLength == nil || *prop.MaxLength != 100 {
		t.Errorf("expected maxLength 100, got %v", prop.MaxLength)
	}
	if prop.Pattern != "^[a-z]+$" {
		t.Errorf("expected pattern '^[a-z]+$', got %s", prop.Pattern)
	}

	// Validate schema compiles and validates data
	compiled := compileSchema(t, schema)

	// Valid data matching pattern
	if err := validateData(t, compiled, `{"username": "abc"}`); err != nil {
		t.Errorf("expected valid data to pass: %v", err)
	}

	// Too short
	if err := validateData(t, compiled, `{"username": ""}`); err == nil {
		t.Error("expected validation to fail for empty string (minLength)")
	}

	// Pattern mismatch (uppercase)
	if err := validateData(t, compiled, `{"username": "ABC"}`); err == nil {
		t.Error("expected validation to fail for pattern mismatch")
	}
}

func TestEntityToJSONSchema_StringAttributeWithFormat(t *testing.T) {
	entity := Entity{
		Slug: "test",
		Attributes: []Attribute{
			{
				Name: "email",
				Type: AttributeTypeString,
				Meta: StringAttributeMeta{
					Format: StringFormatEmail,
				},
			},
		},
	}

	schema := EntityToJSONSchema(entity)
	prop := schema.Properties["email"]

	if prop.Format != "email" {
		t.Errorf("expected format 'email', got %s", prop.Format)
	}

	// Validate schema compiles and validates data
	compiled := compileSchema(t, schema)

	// Valid email
	if err := validateData(t, compiled, `{"email": "test@example.com"}`); err != nil {
		t.Errorf("expected valid email to pass: %v", err)
	}

	// Invalid email
	if err := validateData(t, compiled, `{"email": "not-an-email"}`); err == nil {
		t.Error("expected validation to fail for invalid email format")
	}
}

func TestEntityToJSONSchema_NumberAttribute(t *testing.T) {
	min := 0.0
	max := 100.0
	multipleOf := 0.5
	entity := Entity{
		Slug: "test",
		Attributes: []Attribute{
			{
				Name: "score",
				Type: AttributeTypeNumber,
				Meta: NumberAttributeMeta{
					Minimum:    &min,
					Maximum:    &max,
					MultipleOf: &multipleOf,
				},
			},
		},
	}

	schema := EntityToJSONSchema(entity)
	prop := schema.Properties["score"]

	if prop.Type != "number" {
		t.Errorf("expected type 'number', got %s", prop.Type)
	}
	if prop.Minimum == nil || *prop.Minimum != 0.0 {
		t.Errorf("expected minimum 0, got %v", prop.Minimum)
	}
	if prop.Maximum == nil || *prop.Maximum != 100.0 {
		t.Errorf("expected maximum 100, got %v", prop.Maximum)
	}
	if prop.MultipleOf == nil || *prop.MultipleOf != 0.5 {
		t.Errorf("expected multipleOf 0.5, got %v", prop.MultipleOf)
	}

	// Validate schema compiles and validates data
	compiled := compileSchema(t, schema)

	// Valid data
	if err := validateData(t, compiled, `{"score": 50.5}`); err != nil {
		t.Errorf("expected valid data to pass: %v", err)
	}

	// Below minimum
	if err := validateData(t, compiled, `{"score": -1}`); err == nil {
		t.Error("expected validation to fail for value below minimum")
	}

	// Above maximum
	if err := validateData(t, compiled, `{"score": 101}`); err == nil {
		t.Error("expected validation to fail for value above maximum")
	}

	// Not multiple of 0.5
	if err := validateData(t, compiled, `{"score": 50.3}`); err == nil {
		t.Error("expected validation to fail for value not multiple of 0.5")
	}
}

func TestEntityToJSONSchema_IntegerAttribute(t *testing.T) {
	entity := Entity{
		Slug: "test",
		Attributes: []Attribute{
			{
				Name: "count",
				Type: AttributeTypeInteger,
			},
		},
	}

	schema := EntityToJSONSchema(entity)
	prop := schema.Properties["count"]

	if prop.Type != "integer" {
		t.Errorf("expected type 'integer', got %s", prop.Type)
	}

	// Validate schema compiles and validates data
	compiled := compileSchema(t, schema)

	// Valid integer
	if err := validateData(t, compiled, `{"count": 42}`); err != nil {
		t.Errorf("expected valid data to pass: %v", err)
	}

	// Float should fail for integer type
	if err := validateData(t, compiled, `{"count": 42.5}`); err == nil {
		t.Error("expected validation to fail for float value")
	}
}

func TestEntityToJSONSchema_BooleanAttribute(t *testing.T) {
	entity := Entity{
		Slug: "test",
		Attributes: []Attribute{
			{
				Name:         "active",
				Type:         AttributeTypeBoolean,
				DefaultValue: true,
			},
		},
	}

	schema := EntityToJSONSchema(entity)
	prop := schema.Properties["active"]

	if prop.Type != "boolean" {
		t.Errorf("expected type 'boolean', got %s", prop.Type)
	}
	if prop.Default != true {
		t.Errorf("expected default true, got %v", prop.Default)
	}

	// Validate schema compiles and validates data
	compiled := compileSchema(t, schema)

	// Valid boolean
	if err := validateData(t, compiled, `{"active": true}`); err != nil {
		t.Errorf("expected valid data to pass: %v", err)
	}
	if err := validateData(t, compiled, `{"active": false}`); err != nil {
		t.Errorf("expected valid data to pass: %v", err)
	}

	// String should fail
	if err := validateData(t, compiled, `{"active": "true"}`); err == nil {
		t.Error("expected validation to fail for string value")
	}
}

func TestEntityToJSONSchema_UserAttribute(t *testing.T) {
	entity := Entity{
		Slug: "task",
		Attributes: []Attribute{
			{
				Name: "assignee_id",
				Type: AttributeTypeUser,
			},
		},
	}

	schema := EntityToJSONSchema(entity)
	prop := schema.Properties["assignee_id"]

	// A user attribute stores the composite { type, id } (common.UserRef).
	if prop.Type != "object" {
		t.Errorf("expected type 'object', got %s", prop.Type)
	}
	if prop.Properties["id"] == nil || prop.Properties["type"] == nil {
		t.Errorf("expected 'id' and 'type' properties, got %v", prop.Properties)
	}

	compiled := compileSchema(t, schema)

	// A { type, id } object is valid.
	if err := validateData(t, compiled, `{"assignee_id": {"type": "person", "id": "8f1b3c2a-0d4e-4f6a-9b2c-1e2d3f4a5b6c"}}`); err != nil {
		t.Errorf("expected valid user ref to pass: %v", err)
	}

	// A bare string is no longer a valid user reference.
	if err := validateData(t, compiled, `{"assignee_id": "8f1b3c2a-0d4e-4f6a-9b2c-1e2d3f4a5b6c"}`); err == nil {
		t.Error("expected validation to fail for bare string value")
	}
}

func TestEntityToJSONSchema_KnowledgeTextAttribute(t *testing.T) {
	entity := Entity{
		Slug: "meeting",
		Attributes: []Attribute{
			{
				Name: "notes",
				Type: AttributeTypeKnowledgeText,
			},
		},
	}

	schema := EntityToJSONSchema(entity)
	prop := schema.Properties["notes"]

	// A knowledge-text attribute stores the composite { id, content }
	// (common.KnowledgeNodeRef) — id is the stored ref, content is the
	// transient body filled on single-record reads.
	if prop.Type != "object" {
		t.Errorf("expected type 'object', got %s", prop.Type)
	}
	if prop.Properties["id"] == nil || prop.Properties["content"] == nil {
		t.Errorf("expected 'id' and 'content' properties, got %v", prop.Properties)
	}

	compiled := compileSchema(t, schema)

	// The stored ref { id } is valid, with or without content.
	if err := validateData(t, compiled, `{"notes": {"id": "n-1"}}`); err != nil {
		t.Errorf("expected bare ref to pass: %v", err)
	}
	if err := validateData(t, compiled, `{"notes": {"id": "n-1", "content": "# hello"}}`); err != nil {
		t.Errorf("expected enriched ref to pass: %v", err)
	}

	// Content without the id misses the required stored shape.
	if err := validateData(t, compiled, `{"notes": {"content": "# hello"}}`); err == nil {
		t.Error("expected validation to fail without id")
	}
}

func TestEntityToJSONSchema_DatetimeAttribute(t *testing.T) {
	// Test standard JSON Schema draft-07 formats
	tests := []struct {
		name           string
		format         DatetimeFormat
		expectedFormat string
	}{
		{"date-time", DatetimeFormatDateTime, "date-time"},
		{"date", DatetimeFormatDate, "date"},
		{"time", DatetimeFormatTime, "time"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			entity := Entity{
				Slug: "test",
				Attributes: []Attribute{
					{
						Name: "timestamp",
						Type: AttributeTypeDatetime,
						Meta: DatetimeAttributeMeta{
							Format: tt.format,
						},
					},
				},
			}

			schema := EntityToJSONSchema(entity)
			prop := schema.Properties["timestamp"]

			if prop.Type != "string" {
				t.Errorf("expected type 'string', got %s", prop.Type)
			}
			if prop.Format != tt.expectedFormat {
				t.Errorf("expected format '%s', got '%s'", tt.expectedFormat, prop.Format)
			}

			// Validate schema compiles
			compileSchema(t, schema)
		})
	}
}

func TestEntityToJSONSchema_DatetimeDurationAttribute(t *testing.T) {
	// Duration format is not in draft-07, so we add a regex pattern fallback
	entity := Entity{
		Slug: "test",
		Attributes: []Attribute{
			{
				Name: "duration",
				Type: AttributeTypeDatetime,
				Meta: DatetimeAttributeMeta{
					Format: DatetimeFormatDuration,
				},
			},
		},
	}

	schema := EntityToJSONSchema(entity)
	prop := schema.Properties["duration"]

	if prop.Type != "string" {
		t.Errorf("expected type 'string', got %s", prop.Type)
	}
	// No format field - draft-07 doesn't support "duration", so we use pattern only
	if prop.Format != "" {
		t.Errorf("expected no format for duration (draft-07 doesn't support it), got '%s'", prop.Format)
	}
	if prop.Pattern != DurationPattern {
		t.Errorf("expected duration pattern, got %s", prop.Pattern)
	}

	// Now we can compile and validate thanks to the regex pattern
	compiled := compileSchema(t, schema)

	// Valid ISO 8601 durations should pass
	if err := validateData(t, compiled, `{"duration": "P1Y"}`); err != nil {
		t.Errorf("expected P1Y to pass: %v", err)
	}
	if err := validateData(t, compiled, `{"duration": "P1M"}`); err != nil {
		t.Errorf("expected P1M to pass: %v", err)
	}
	if err := validateData(t, compiled, `{"duration": "P1D"}`); err != nil {
		t.Errorf("expected P1D to pass: %v", err)
	}
	if err := validateData(t, compiled, `{"duration": "PT1H"}`); err != nil {
		t.Errorf("expected PT1H to pass: %v", err)
	}
	if err := validateData(t, compiled, `{"duration": "PT1M"}`); err != nil {
		t.Errorf("expected PT1M to pass: %v", err)
	}
	if err := validateData(t, compiled, `{"duration": "PT1S"}`); err != nil {
		t.Errorf("expected PT1S to pass: %v", err)
	}
	if err := validateData(t, compiled, `{"duration": "P1Y2M3DT4H5M6S"}`); err != nil {
		t.Errorf("expected full duration to pass: %v", err)
	}
	if err := validateData(t, compiled, `{"duration": "PT1.5S"}`); err != nil {
		t.Errorf("expected fractional seconds to pass: %v", err)
	}
	if err := validateData(t, compiled, `{"duration": "P1W"}`); err != nil {
		t.Errorf("expected P1W (weeks) to pass: %v", err)
	}

	// Invalid durations should fail
	if err := validateData(t, compiled, `{"duration": "1Y"}`); err == nil {
		t.Error("expected duration without P prefix to fail")
	}
	if err := validateData(t, compiled, `{"duration": "not-a-duration"}`); err == nil {
		t.Error("expected invalid duration to fail")
	}
	if err := validateData(t, compiled, `{"duration": "P"}`); err == nil {
		t.Error("expected empty duration to fail")
	}
}

func TestEntityToJSONSchema_EnumAttribute(t *testing.T) {
	entity := Entity{
		Slug: "test",
		Attributes: []Attribute{
			{
				Name: "status",
				Type: AttributeTypeEnum,
				Meta: EnumAttributeMeta{
					Values: []EnumValue{
						{Value: "active", Label: "Active"},
						{Value: "inactive", Label: "Inactive"},
						{Value: "pending", Label: "Pending"},
					},
				},
			},
		},
	}

	schema := EntityToJSONSchema(entity)
	prop := schema.Properties["status"]

	if prop.Type != "string" {
		t.Errorf("expected type 'string', got %s", prop.Type)
	}
	if len(prop.Enum) != 3 {
		t.Errorf("expected 3 enum values, got %d", len(prop.Enum))
	}

	expectedValues := []string{"active", "inactive", "pending"}
	for i, expected := range expectedValues {
		if prop.Enum[i] != expected {
			t.Errorf("expected enum[%d] to be '%s', got '%v'", i, expected, prop.Enum[i])
		}
	}

	// Validate schema compiles and validates data
	compiled := compileSchema(t, schema)

	// Valid enum value
	if err := validateData(t, compiled, `{"status": "active"}`); err != nil {
		t.Errorf("expected valid data to pass: %v", err)
	}

	// Invalid enum value
	if err := validateData(t, compiled, `{"status": "unknown"}`); err == nil {
		t.Error("expected validation to fail for invalid enum value")
	}
}

func TestEntityToJSONSchema_ArrayAttribute(t *testing.T) {
	minItems := 1
	maxItems := 10
	entity := Entity{
		Slug: "test",
		Attributes: []Attribute{
			{
				Name: "tags",
				Type: AttributeTypeArray,
				Meta: ArrayAttributeMeta{
					MinItems:          &minItems,
					MaxItems:          &maxItems,
					ItemsMustBeUnique: true,
					Items: &Attribute{
						Name: "tag",
						Type: AttributeTypeString,
					},
				},
			},
		},
	}

	schema := EntityToJSONSchema(entity)
	prop := schema.Properties["tags"]

	if prop.Type != "array" {
		t.Errorf("expected type 'array', got %s", prop.Type)
	}
	if prop.MinItems == nil || *prop.MinItems != 1 {
		t.Errorf("expected minItems 1, got %v", prop.MinItems)
	}
	if prop.MaxItems == nil || *prop.MaxItems != 10 {
		t.Errorf("expected maxItems 10, got %v", prop.MaxItems)
	}
	if prop.UniqueItems == nil || !*prop.UniqueItems {
		t.Errorf("expected uniqueItems true, got %v", prop.UniqueItems)
	}
	if prop.Items == nil {
		t.Fatal("expected items schema, got nil")
	}
	if prop.Items.Type != "string" {
		t.Errorf("expected items type 'string', got %s", prop.Items.Type)
	}

	// Validate schema compiles and validates data
	compiled := compileSchema(t, schema)

	// Valid array
	if err := validateData(t, compiled, `{"tags": ["a", "b", "c"]}`); err != nil {
		t.Errorf("expected valid data to pass: %v", err)
	}

	// Empty array (below minItems)
	if err := validateData(t, compiled, `{"tags": []}`); err == nil {
		t.Error("expected validation to fail for empty array (minItems)")
	}

	// Duplicate items
	if err := validateData(t, compiled, `{"tags": ["a", "a"]}`); err == nil {
		t.Error("expected validation to fail for duplicate items")
	}

	// Wrong item type
	if err := validateData(t, compiled, `{"tags": [1, 2, 3]}`); err == nil {
		t.Error("expected validation to fail for wrong item type")
	}
}

func TestEntityToJSONSchema_ObjectAttribute(t *testing.T) {
	entity := Entity{
		Slug: "test",
		Attributes: []Attribute{
			{
				Name: "address",
				Type: AttributeTypeObject,
				Meta: ObjectAttributeMeta{
					Attributes: []Attribute{
						{
							Name:       "street",
							Type:       AttributeTypeString,
							IsRequired: true,
						},
						{
							Name: "city",
							Type: AttributeTypeString,
						},
						{
							Name: "zipCode",
							Type: AttributeTypeString,
						},
					},
				},
			},
		},
	}

	schema := EntityToJSONSchema(entity)
	prop := schema.Properties["address"]

	if prop.Type != "object" {
		t.Errorf("expected type 'object', got %s", prop.Type)
	}
	if len(prop.Properties) != 3 {
		t.Errorf("expected 3 properties, got %d", len(prop.Properties))
	}
	if len(prop.Required) != 1 || prop.Required[0] != "street" {
		t.Errorf("expected required to contain 'street', got %v", prop.Required)
	}

	// Validate schema compiles and validates data
	compiled := compileSchema(t, schema)

	// Valid nested object
	if err := validateData(t, compiled, `{"address": {"street": "123 Main St", "city": "NYC"}}`); err != nil {
		t.Errorf("expected valid data to pass: %v", err)
	}

	// Missing required nested field
	if err := validateData(t, compiled, `{"address": {"city": "NYC"}}`); err == nil {
		t.Error("expected validation to fail for missing required nested field 'street'")
	}
}

func TestEntityToJSONSchema_ReadOnlyAttribute(t *testing.T) {
	entity := Entity{
		Slug: "test",
		Attributes: []Attribute{
			{
				Name:       "id",
				Type:       AttributeTypeString,
				IsReadOnly: true,
			},
		},
	}

	schema := EntityToJSONSchema(entity)
	prop := schema.Properties["id"]

	if prop.ReadOnly == nil || !*prop.ReadOnly {
		t.Errorf("expected readOnly true, got %v", prop.ReadOnly)
	}

	// Validate schema compiles
	compileSchema(t, schema)
}

func TestEntityToJSONSchema_MetaFromJSON(t *testing.T) {
	// Test that meta works when unmarshaled from JSON (as map[string]any)
	entityJSON := `{
		"slug": "test",
		"name": "Test",
		"attributes": [
			{
				"name": "username",
				"type": "string",
				"meta": {
					"min_length": 5,
					"max_length": 100,
					"pattern": "^[a-z]+$"
				}
			}
		]
	}`

	var entity Entity
	if err := json.Unmarshal([]byte(entityJSON), &entity); err != nil {
		t.Fatalf("failed to unmarshal entity: %v", err)
	}

	schema := EntityToJSONSchema(entity)
	prop := schema.Properties["username"]

	if prop.Type != "string" {
		t.Errorf("expected type 'string', got %s", prop.Type)
	}
	if prop.MinLength == nil || *prop.MinLength != 5 {
		t.Errorf("expected minLength 5, got %v", prop.MinLength)
	}
	if prop.MaxLength == nil || *prop.MaxLength != 100 {
		t.Errorf("expected maxLength 100, got %v", prop.MaxLength)
	}
	if prop.Pattern != "^[a-z]+$" {
		t.Errorf("expected pattern '^[a-z]+$', got %s", prop.Pattern)
	}

	// Validate schema compiles and validates data
	compiled := compileSchema(t, schema)

	// Valid data
	if err := validateData(t, compiled, `{"username": "hello"}`); err != nil {
		t.Errorf("expected valid data to pass: %v", err)
	}

	// Too short
	if err := validateData(t, compiled, `{"username": "abc"}`); err == nil {
		t.Error("expected validation to fail for string below minLength")
	}
}

func TestEntityToJSONSchema_NestedArrayOfObjects(t *testing.T) {
	// Test deeply nested structure: array of objects
	entity := Entity{
		Slug: "order",
		Name: "Order",
		Attributes: []Attribute{
			{
				Name: "items",
				Type: AttributeTypeArray,
				Meta: ArrayAttributeMeta{
					Items: &Attribute{
						Name: "item",
						Type: AttributeTypeObject,
						Meta: ObjectAttributeMeta{
							Attributes: []Attribute{
								{
									Name:       "productId",
									Type:       AttributeTypeString,
									IsRequired: true,
								},
								{
									Name: "quantity",
									Type: AttributeTypeInteger,
								},
							},
						},
					},
				},
			},
		},
	}

	schema := EntityToJSONSchema(entity)
	items := schema.Properties["items"]

	if items.Type != "array" {
		t.Fatalf("expected items type 'array', got %s", items.Type)
	}
	if items.Items == nil {
		t.Fatal("expected items.Items, got nil")
	}
	if items.Items.Type != "object" {
		t.Errorf("expected items.Items type 'object', got %s", items.Items.Type)
	}
	if len(items.Items.Properties) != 2 {
		t.Errorf("expected 2 nested properties, got %d", len(items.Items.Properties))
	}
	if len(items.Items.Required) != 1 || items.Items.Required[0] != "productId" {
		t.Errorf("expected nested required to contain 'productId', got %v", items.Items.Required)
	}

	// Validate schema compiles and validates data
	compiled := compileSchema(t, schema)

	// Valid nested array of objects
	if err := validateData(t, compiled, `{"items": [{"productId": "abc", "quantity": 2}]}`); err != nil {
		t.Errorf("expected valid data to pass: %v", err)
	}

	// Missing required field in nested object
	if err := validateData(t, compiled, `{"items": [{"quantity": 2}]}`); err == nil {
		t.Error("expected validation to fail for missing required nested field 'productId'")
	}

	// Wrong type in nested object
	if err := validateData(t, compiled, `{"items": [{"productId": "abc", "quantity": "two"}]}`); err == nil {
		t.Error("expected validation to fail for wrong type in nested object")
	}
}

func TestEntityToJSONSchema_EmptyEntity(t *testing.T) {
	entity := Entity{
		Slug: "empty",
		Name: "Empty",
	}

	schema := EntityToJSONSchema(entity)

	if schema.Type != "object" {
		t.Errorf("expected type 'object', got %s", schema.Type)
	}
	if schema.Properties == nil {
		t.Error("expected properties map, got nil")
	}
	if len(schema.Properties) != 0 {
		t.Errorf("expected 0 properties, got %d", len(schema.Properties))
	}
	if schema.Required != nil {
		t.Errorf("expected required to be nil, got %v", schema.Required)
	}

	// Validate schema compiles and validates data
	compiled := compileSchema(t, schema)

	// Any object should be valid
	if err := validateData(t, compiled, `{}`); err != nil {
		t.Errorf("expected empty object to be valid: %v", err)
	}
	if err := validateData(t, compiled, `{"anyField": "anyValue"}`); err != nil {
		t.Errorf("expected object with extra fields to be valid: %v", err)
	}
}

func TestEntityToJSONSchema_JSONOutput(t *testing.T) {
	// Verify the schema produces valid JSON output
	entity := Entity{
		Slug:        "product",
		Name:        "Product",
		Description: "A product entity",
		Attributes: []Attribute{
			{
				Name:       "name",
				Type:       AttributeTypeString,
				Label:      "Product Name",
				IsRequired: true,
			},
			{
				Name:  "price",
				Type:  AttributeTypeNumber,
				Label: "Price",
			},
		},
	}

	schema := EntityToJSONSchema(entity)

	jsonBytes, err := json.MarshalIndent(schema, "", "  ")
	if err != nil {
		t.Fatalf("failed to marshal schema to JSON: %v", err)
	}

	// Verify it can be unmarshaled back
	var parsed JSONSchema
	if err := json.Unmarshal(jsonBytes, &parsed); err != nil {
		t.Fatalf("failed to unmarshal schema from JSON: %v", err)
	}

	if parsed.Schema != Draft07SchemaURI {
		t.Errorf("expected $schema to be %s after round-trip", Draft07SchemaURI)
	}

	// Validate schema compiles
	compileSchema(t, schema)
}

func TestEntityToJSONSchema_ComplexEntity(t *testing.T) {
	// Test a complex real-world entity with multiple attribute types
	minLen := 1
	maxLen := 255
	minItems := 1
	minPrice := 0.0

	entity := Entity{
		Slug:        "product",
		Name:        "Product",
		Description: "A product in the catalog",
		Attributes: []Attribute{
			{
				Name:       "name",
				Type:       AttributeTypeString,
				Label:      "Product Name",
				IsRequired: true,
				Meta: StringAttributeMeta{
					MinLength: &minLen,
					MaxLength: &maxLen,
				},
			},
			{
				Name:  "price",
				Type:  AttributeTypeNumber,
				Label: "Price",
				Meta: NumberAttributeMeta{
					Minimum: &minPrice,
				},
			},
			{
				Name:       "status",
				Type:       AttributeTypeEnum,
				Label:      "Status",
				IsRequired: true,
				Meta: EnumAttributeMeta{
					Values: []EnumValue{
						{Value: "draft", Label: "Draft"},
						{Value: "active", Label: "Active"},
						{Value: "discontinued", Label: "Discontinued"},
					},
				},
			},
			{
				Name:  "tags",
				Type:  AttributeTypeArray,
				Label: "Tags",
				Meta: ArrayAttributeMeta{
					MinItems:          &minItems,
					ItemsMustBeUnique: true,
					Items: &Attribute{
						Name: "tag",
						Type: AttributeTypeString,
					},
				},
			},
			{
				Name:       "createdAt",
				Type:       AttributeTypeDatetime,
				Label:      "Created At",
				IsReadOnly: true,
				Meta: DatetimeAttributeMeta{
					Format: DatetimeFormatDateTime,
				},
			},
			{
				Name:  "inStock",
				Type:  AttributeTypeBoolean,
				Label: "In Stock",
			},
			{
				Name: "metadata",
				Type: AttributeTypeObject,
				Meta: ObjectAttributeMeta{
					Attributes: []Attribute{
						{
							Name: "sku",
							Type: AttributeTypeString,
						},
						{
							Name: "weight",
							Type: AttributeTypeNumber,
						},
					},
				},
			},
		},
	}

	schema := EntityToJSONSchema(entity)

	// Validate schema compiles
	compiled := compileSchema(t, schema)

	// Valid complete product
	validProduct := `{
		"name": "Widget",
		"price": 29.99,
		"status": "active",
		"tags": ["electronics", "gadget"],
		"createdAt": "2024-01-15T10:30:00Z",
		"inStock": true,
		"metadata": {
			"sku": "WGT-001",
			"weight": 0.5
		}
	}`
	if err := validateData(t, compiled, validProduct); err != nil {
		t.Errorf("expected valid product to pass: %v", err)
	}

	// Minimal valid product (only required fields)
	minimalProduct := `{"name": "Widget", "status": "draft"}`
	if err := validateData(t, compiled, minimalProduct); err != nil {
		t.Errorf("expected minimal valid product to pass: %v", err)
	}

	// Invalid: missing required field
	if err := validateData(t, compiled, `{"name": "Widget"}`); err == nil {
		t.Error("expected validation to fail for missing required 'status'")
	}

	// Invalid: empty name
	if err := validateData(t, compiled, `{"name": "", "status": "active"}`); err == nil {
		t.Error("expected validation to fail for empty name (minLength)")
	}

	// Invalid: negative price
	if err := validateData(t, compiled, `{"name": "Widget", "status": "active", "price": -1}`); err == nil {
		t.Error("expected validation to fail for negative price")
	}

	// Invalid: wrong enum value
	if err := validateData(t, compiled, `{"name": "Widget", "status": "invalid"}`); err == nil {
		t.Error("expected validation to fail for invalid status enum")
	}
}

func TestEntityToJSONSchema_ExclusiveMinMax(t *testing.T) {
	exclMin := 0.0
	exclMax := 100.0
	entity := Entity{
		Slug: "test",
		Attributes: []Attribute{
			{
				Name: "score",
				Type: AttributeTypeNumber,
				Meta: NumberAttributeMeta{
					ExclusiveMinimum: &exclMin,
					ExclusiveMaximum: &exclMax,
				},
			},
		},
	}

	schema := EntityToJSONSchema(entity)
	prop := schema.Properties["score"]

	if prop.ExclusiveMinimum == nil || *prop.ExclusiveMinimum != 0.0 {
		t.Errorf("expected exclusiveMinimum 0, got %v", prop.ExclusiveMinimum)
	}
	if prop.ExclusiveMaximum == nil || *prop.ExclusiveMaximum != 100.0 {
		t.Errorf("expected exclusiveMaximum 100, got %v", prop.ExclusiveMaximum)
	}

	compiled := compileSchema(t, schema)

	// Value at boundary should fail (exclusive)
	if err := validateData(t, compiled, `{"score": 0}`); err == nil {
		t.Error("expected validation to fail for value at exclusive minimum")
	}
	if err := validateData(t, compiled, `{"score": 100}`); err == nil {
		t.Error("expected validation to fail for value at exclusive maximum")
	}

	// Value just inside boundary should pass
	if err := validateData(t, compiled, `{"score": 0.001}`); err != nil {
		t.Errorf("expected value just above exclusive minimum to pass: %v", err)
	}
	if err := validateData(t, compiled, `{"score": 99.999}`); err != nil {
		t.Errorf("expected value just below exclusive maximum to pass: %v", err)
	}
}

func TestEntityToJSONSchema_NilMeta(t *testing.T) {
	// Test that nil meta doesn't cause issues for any type
	entity := Entity{
		Slug: "test",
		Attributes: []Attribute{
			{Name: "str", Type: AttributeTypeString, Meta: nil},
			{Name: "num", Type: AttributeTypeNumber, Meta: nil},
			{Name: "int", Type: AttributeTypeInteger, Meta: nil},
			{Name: "bool", Type: AttributeTypeBoolean, Meta: nil},
			{Name: "arr", Type: AttributeTypeArray, Meta: nil},
			{Name: "obj", Type: AttributeTypeObject, Meta: nil},
			{Name: "enum", Type: AttributeTypeEnum, Meta: nil},
			{Name: "dt", Type: AttributeTypeDatetime, Meta: nil},
		},
	}

	schema := EntityToJSONSchema(entity)

	// Should compile without errors
	compiled := compileSchema(t, schema)

	// Valid data for all nil-meta fields
	validData := `{
		"str": "hello",
		"num": 3.14,
		"int": 42,
		"bool": true,
		"arr": [1, 2, 3],
		"obj": {"key": "value"},
		"enum": "anything",
		"dt": "2024-01-15T10:30:00Z"
	}`
	if err := validateData(t, compiled, validData); err != nil {
		t.Errorf("expected valid data to pass: %v", err)
	}
}

func TestEntityToJSONSchema_AllStringFormats(t *testing.T) {
	// Formats supported by draft-07 jsonschema validator
	tests := []struct {
		format       StringFormat
		validValue   string
		invalidValue string
	}{
		{StringFormatEmail, "test@example.com", "not-an-email"},
		{StringFormatURI, "https://example.com/path", "not a uri"},
		{StringFormatIPv4, "192.168.1.1", "999.999.999.999"},
		{StringFormatIPv6, "2001:0db8:85a3:0000:0000:8a2e:0370:7334", "not-ipv6"},
		{StringFormatHostname, "example.com", "invalid hostname!"},
	}

	for _, tt := range tests {
		t.Run(string(tt.format), func(t *testing.T) {
			entity := Entity{
				Slug: "test",
				Attributes: []Attribute{
					{
						Name: "field",
						Type: AttributeTypeString,
						Meta: StringAttributeMeta{
							Format: tt.format,
						},
					},
				},
			}

			schema := EntityToJSONSchema(entity)
			if schema.Properties["field"].Format != string(tt.format) {
				t.Errorf("expected format %s, got %s", tt.format, schema.Properties["field"].Format)
			}

			compiled := compileSchema(t, schema)

			// Valid value should pass
			if err := validateData(t, compiled, `{"field": "`+tt.validValue+`"}`); err != nil {
				t.Errorf("expected valid %s to pass: %v", tt.format, err)
			}

			// Invalid value should fail
			if err := validateData(t, compiled, `{"field": "`+tt.invalidValue+`"}`); err == nil {
				t.Errorf("expected invalid %s to fail", tt.format)
			}
		})
	}
}

func TestEntityToJSONSchema_UUIDFormat(t *testing.T) {
	// UUID format is not in draft-07, so we use regex pattern instead
	entity := Entity{
		Slug: "test",
		Attributes: []Attribute{
			{
				Name: "id",
				Type: AttributeTypeString,
				Meta: StringAttributeMeta{
					Format: StringFormatUUID,
				},
			},
		},
	}

	schema := EntityToJSONSchema(entity)
	prop := schema.Properties["id"]

	// Should have pattern for validation but no format (draft-07 rejects unknown formats)
	if prop.Format != "" {
		t.Errorf("expected no format for UUID (draft-07 doesn't support it), got %s", prop.Format)
	}
	if prop.Pattern != UUIDPattern {
		t.Errorf("expected UUID pattern, got %s", prop.Pattern)
	}

	// Now we can compile and validate thanks to the regex pattern
	compiled := compileSchema(t, schema)

	// Valid UUIDs should pass
	if err := validateData(t, compiled, `{"id": "123e4567-e89b-12d3-a456-426614174000"}`); err != nil {
		t.Errorf("expected valid UUID to pass: %v", err)
	}
	if err := validateData(t, compiled, `{"id": "550e8400-e29b-41d4-a716-446655440000"}`); err != nil {
		t.Errorf("expected valid UUID to pass: %v", err)
	}

	// Invalid UUIDs should fail
	if err := validateData(t, compiled, `{"id": "not-a-uuid"}`); err == nil {
		t.Error("expected invalid UUID to fail")
	}
	if err := validateData(t, compiled, `{"id": "123e4567-e89b-12d3-a456"}`); err == nil {
		t.Error("expected incomplete UUID to fail")
	}
	if err := validateData(t, compiled, `{"id": "123e4567-e89b-12d3-a456-42661417400g"}`); err == nil {
		t.Error("expected UUID with invalid character to fail")
	}
}

func TestEntityToJSONSchema_ArrayWithoutItemsSchema(t *testing.T) {
	entity := Entity{
		Slug: "test",
		Attributes: []Attribute{
			{
				Name: "mixedArray",
				Type: AttributeTypeArray,
				// No items schema - any items allowed
			},
		},
	}

	schema := EntityToJSONSchema(entity)
	prop := schema.Properties["mixedArray"]

	if prop.Type != "array" {
		t.Errorf("expected type 'array', got %s", prop.Type)
	}
	if prop.Items != nil {
		t.Errorf("expected items to be nil, got %v", prop.Items)
	}

	compiled := compileSchema(t, schema)

	// Mixed types should be valid
	if err := validateData(t, compiled, `{"mixedArray": [1, "two", true, null]}`); err != nil {
		t.Errorf("expected mixed array to pass: %v", err)
	}
}

func TestEntityToJSONSchema_DeeplyNestedObject(t *testing.T) {
	// 3 levels deep: company -> department -> team -> members
	entity := Entity{
		Slug: "company",
		Attributes: []Attribute{
			{
				Name: "department",
				Type: AttributeTypeObject,
				Meta: ObjectAttributeMeta{
					Attributes: []Attribute{
						{
							Name:       "name",
							Type:       AttributeTypeString,
							IsRequired: true,
						},
						{
							Name: "team",
							Type: AttributeTypeObject,
							Meta: ObjectAttributeMeta{
								Attributes: []Attribute{
									{
										Name:       "teamName",
										Type:       AttributeTypeString,
										IsRequired: true,
									},
									{
										Name: "memberCount",
										Type: AttributeTypeInteger,
									},
								},
							},
						},
					},
				},
			},
		},
	}

	schema := EntityToJSONSchema(entity)
	compiled := compileSchema(t, schema)

	// Valid deeply nested data
	validData := `{
		"department": {
			"name": "Engineering",
			"team": {
				"teamName": "Backend",
				"memberCount": 5
			}
		}
	}`
	if err := validateData(t, compiled, validData); err != nil {
		t.Errorf("expected valid deeply nested data to pass: %v", err)
	}

	// Missing required at level 2
	if err := validateData(t, compiled, `{"department": {"name": "Engineering", "team": {"memberCount": 5}}}`); err == nil {
		t.Error("expected validation to fail for missing 'teamName'")
	}
}

func TestEntityToJSONSchema_DefaultValues(t *testing.T) {
	entity := Entity{
		Slug: "test",
		Attributes: []Attribute{
			{Name: "strField", Type: AttributeTypeString, DefaultValue: "default text"},
			{Name: "numField", Type: AttributeTypeNumber, DefaultValue: 42.5},
			{Name: "intField", Type: AttributeTypeInteger, DefaultValue: 10},
			{Name: "boolField", Type: AttributeTypeBoolean, DefaultValue: false},
			{Name: "arrField", Type: AttributeTypeArray, DefaultValue: []string{"a", "b"}},
		},
	}

	schema := EntityToJSONSchema(entity)

	if schema.Properties["strField"].Default != "default text" {
		t.Errorf("expected string default 'default text', got %v", schema.Properties["strField"].Default)
	}
	if schema.Properties["numField"].Default != 42.5 {
		t.Errorf("expected number default 42.5, got %v", schema.Properties["numField"].Default)
	}
	if schema.Properties["intField"].Default != 10 {
		t.Errorf("expected integer default 10, got %v", schema.Properties["intField"].Default)
	}
	if schema.Properties["boolField"].Default != false {
		t.Errorf("expected boolean default false, got %v", schema.Properties["boolField"].Default)
	}

	// Schema should still compile
	compileSchema(t, schema)
}

func TestEntityToJSONSchema_MultipleRequiredFields(t *testing.T) {
	entity := Entity{
		Slug: "test",
		Attributes: []Attribute{
			{Name: "firstName", Type: AttributeTypeString, IsRequired: true},
			{Name: "lastName", Type: AttributeTypeString, IsRequired: true},
			{Name: "email", Type: AttributeTypeString, IsRequired: true},
			{Name: "phone", Type: AttributeTypeString, IsRequired: false},
		},
	}

	schema := EntityToJSONSchema(entity)

	if len(schema.Required) != 3 {
		t.Errorf("expected 3 required fields, got %d", len(schema.Required))
	}

	compiled := compileSchema(t, schema)

	// All required fields present
	if err := validateData(t, compiled, `{"firstName": "John", "lastName": "Doe", "email": "john@example.com"}`); err != nil {
		t.Errorf("expected valid data to pass: %v", err)
	}

	// Missing one required field
	if err := validateData(t, compiled, `{"firstName": "John", "lastName": "Doe"}`); err == nil {
		t.Error("expected validation to fail for missing 'email'")
	}

	// Missing multiple required fields
	if err := validateData(t, compiled, `{"firstName": "John"}`); err == nil {
		t.Error("expected validation to fail for missing required fields")
	}
}

func TestEntityToJSONSchema_ArrayOfArrays(t *testing.T) {
	entity := Entity{
		Slug: "matrix",
		Attributes: []Attribute{
			{
				Name: "grid",
				Type: AttributeTypeArray,
				Meta: ArrayAttributeMeta{
					Items: &Attribute{
						Name: "row",
						Type: AttributeTypeArray,
						Meta: ArrayAttributeMeta{
							Items: &Attribute{
								Name: "cell",
								Type: AttributeTypeInteger,
							},
						},
					},
				},
			},
		},
	}

	schema := EntityToJSONSchema(entity)
	compiled := compileSchema(t, schema)

	// Valid 2D array
	if err := validateData(t, compiled, `{"grid": [[1, 2, 3], [4, 5, 6]]}`); err != nil {
		t.Errorf("expected valid 2D array to pass: %v", err)
	}

	// Invalid: wrong inner type
	if err := validateData(t, compiled, `{"grid": [[1, 2, "three"]]}`); err == nil {
		t.Error("expected validation to fail for wrong inner type")
	}

	// Invalid: not nested array
	if err := validateData(t, compiled, `{"grid": [1, 2, 3]}`); err == nil {
		t.Error("expected validation to fail for non-nested array")
	}
}

func TestEntityToJSONSchema_ObjectWithAllFieldsRequired(t *testing.T) {
	entity := Entity{
		Slug: "strictObject",
		Attributes: []Attribute{
			{
				Name: "config",
				Type: AttributeTypeObject,
				Meta: ObjectAttributeMeta{
					Attributes: []Attribute{
						{Name: "host", Type: AttributeTypeString, IsRequired: true},
						{Name: "port", Type: AttributeTypeInteger, IsRequired: true},
						{Name: "secure", Type: AttributeTypeBoolean, IsRequired: true},
					},
				},
			},
		},
	}

	schema := EntityToJSONSchema(entity)
	nested := schema.Properties["config"]

	if len(nested.Required) != 3 {
		t.Errorf("expected 3 required nested fields, got %d", len(nested.Required))
	}

	compiled := compileSchema(t, schema)

	// All fields present
	if err := validateData(t, compiled, `{"config": {"host": "localhost", "port": 8080, "secure": true}}`); err != nil {
		t.Errorf("expected valid data to pass: %v", err)
	}

	// Missing any field should fail
	if err := validateData(t, compiled, `{"config": {"host": "localhost", "port": 8080}}`); err == nil {
		t.Error("expected validation to fail for missing 'secure'")
	}
}

func TestEntityToJSONSchema_MaxItemsExceeded(t *testing.T) {
	maxItems := 3
	entity := Entity{
		Slug: "test",
		Attributes: []Attribute{
			{
				Name: "tags",
				Type: AttributeTypeArray,
				Meta: ArrayAttributeMeta{
					MaxItems: &maxItems,
					Items: &Attribute{
						Name: "tag",
						Type: AttributeTypeString,
					},
				},
			},
		},
	}

	schema := EntityToJSONSchema(entity)
	compiled := compileSchema(t, schema)

	// Exactly at limit
	if err := validateData(t, compiled, `{"tags": ["a", "b", "c"]}`); err != nil {
		t.Errorf("expected array at max limit to pass: %v", err)
	}

	// Over limit
	if err := validateData(t, compiled, `{"tags": ["a", "b", "c", "d"]}`); err == nil {
		t.Error("expected validation to fail for array exceeding maxItems")
	}
}

func TestEntityToJSONSchema_StringMaxLengthExceeded(t *testing.T) {
	maxLen := 5
	entity := Entity{
		Slug: "test",
		Attributes: []Attribute{
			{
				Name: "code",
				Type: AttributeTypeString,
				Meta: StringAttributeMeta{
					MaxLength: &maxLen,
				},
			},
		},
	}

	schema := EntityToJSONSchema(entity)
	compiled := compileSchema(t, schema)

	// Exactly at limit
	if err := validateData(t, compiled, `{"code": "abcde"}`); err != nil {
		t.Errorf("expected string at max length to pass: %v", err)
	}

	// Over limit
	if err := validateData(t, compiled, `{"code": "abcdef"}`); err == nil {
		t.Error("expected validation to fail for string exceeding maxLength")
	}
}

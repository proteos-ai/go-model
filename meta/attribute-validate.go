package metamodel

import (
	"fmt"
	"regexp"
)

// attributeNamePattern is the snake_case identifier rule every attribute name
// must satisfy (matches the platform's wire-format convention).
var attributeNamePattern = regexp.MustCompile(`^[a-z][a-z0-9_]*$`)

// validAttributeTypes is the closed set of attribute types a bare schema
// definition may use.
var validAttributeTypes = map[AttributeType]bool{
	AttributeTypeString:        true,
	AttributeTypeNumber:        true,
	AttributeTypeInteger:       true,
	AttributeTypeBoolean:       true,
	AttributeTypeDatetime:      true,
	AttributeTypeEnum:          true,
	AttributeTypeArray:         true,
	AttributeTypeObject:        true,
	AttributeTypeRelation:      true,
	AttributeTypeUser:          true,
	AttributeTypeCurrency:      true,
	AttributeTypeKnowledgeText: true,
	AttributeTypeFile:          true,
}

// ValidateAttributeDefinitions checks that a bare attribute list is a
// well-formed schema definition: non-empty snake_case names, unique at each
// nesting level, known types, and recursively valid object/array meta. Pure —
// entity-level rules (platform attributes, relations against org entities)
// stay in their owning services; this validates only the shape of the
// definition itself, for entity-less attribute sets such as a workflow node's
// declared output.
func ValidateAttributeDefinitions(attrs []Attribute) error {
	return validateAttributeDefinitions(attrs, "")
}

func validateAttributeDefinitions(attrs []Attribute, path string) error {
	seen := make(map[string]bool, len(attrs))
	for _, attr := range attrs {
		location := attr.Name
		if path != "" {
			location = path + "." + attr.Name
		}
		if attr.Name == "" {
			return fmt.Errorf("attribute at %q has an empty name", path)
		}
		if !attributeNamePattern.MatchString(attr.Name) {
			return fmt.Errorf("attribute name %q must be snake_case (^[a-z][a-z0-9_]*$)", location)
		}
		if seen[attr.Name] {
			return fmt.Errorf("duplicate attribute name %q", location)
		}
		seen[attr.Name] = true
		if !validAttributeTypes[attr.Type] {
			return fmt.Errorf("attribute %q has unknown type %q", location, attr.Type)
		}
		if err := validateAttributeMeta(attr, location); err != nil {
			return err
		}
	}
	return nil
}

// validateAttributeMeta recurses into the type-specific meta of object and
// array attributes so nested definitions obey the same rules.
func validateAttributeMeta(attr Attribute, location string) error {
	switch attr.Type {
	case AttributeTypeObject:
		if attr.Meta == nil {
			return nil
		}
		objectMeta := ParseMetaAs[ObjectAttributeMeta](attr.Meta)
		if objectMeta == nil {
			return fmt.Errorf("attribute %q has invalid object meta", location)
		}
		return validateAttributeDefinitions(objectMeta.Attributes, location)
	case AttributeTypeArray:
		if attr.Meta == nil {
			return nil
		}
		arrayMeta := ParseMetaAs[ArrayAttributeMeta](attr.Meta)
		if arrayMeta == nil {
			return fmt.Errorf("attribute %q has invalid array meta", location)
		}
		if arrayMeta.Items != nil {
			return validateAttributeDefinitions([]Attribute{*arrayMeta.Items}, location)
		}
	}
	return nil
}

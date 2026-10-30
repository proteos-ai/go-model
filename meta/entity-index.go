package metamodel

import (
	"fmt"
	"regexp"
	"strings"
)

// UniqueConstraint declares that no two records of an entity may share the
// same combination of values for the listed attributes. It is the entity's
// business rule — the one index concern a module declares; performance indexes
// are the platform's job and never appear in a manifest.
//
// A one-attribute constraint is spelled on the attribute itself as
// `is_unique: true`; this struct is for the composite case (2+ attributes).
// EffectiveUniqueConstraints folds both spellings into one list. Postgres
// enforces each as a unique index on the record table; a write that violates
// one fails with 409 record_unique_violation naming the constraint.
type UniqueConstraint struct {
	// Name identifies the constraint within the entity (snake_case, like an
	// attribute name). It is what a 409 reports back.
	Name string `json:"name"`
	// Attributes are top-level attribute names or object paths
	// (`address.city`). Platform columns are not allowed: `id` is the primary
	// key already and the timestamps are never unique.
	Attributes []string `json:"attributes"`
}

// UniqueConstraintNameOf is the name EffectiveUniqueConstraints gives to the
// constraint implied by `is_unique: true` on one attribute. Exported so a 409
// and the index status can be traced back to the attribute flag.
func UniqueConstraintNameOf(attributeName string) string {
	return "unique_" + attributeName
}

// MaxUniqueConstraintAttributes caps a constraint's width. A unique key over
// more columns than this is almost always a modeling mistake, and every column
// costs on every write.
const MaxUniqueConstraintAttributes = 8

var uniqueConstraintNamePattern = regexp.MustCompile(`^[a-z][a-z0-9_]*$`)

// IsIndexableAttributeType reports whether attributes of this type can carry
// an index (unique or otherwise). The projection the query builder uses for
// each type is a single, deterministic SQL expression (text extraction, a cast,
// or the IMMUTABLE datetime wrapper), which is what makes it indexable:
//
//   - array: served by the GIN index on the whole payload; no scalar to index
//   - object: a container, not a value — index its leaves by path instead
//   - knowledge-text: large text, never compared for equality
//
// Datetime attributes with a `format` of date, time or duration are refused
// separately (IsIndexableAttribute): ::date / ::time are not IMMUTABLE, and a
// duration ("PT1H") is not a timestamp at all.
func IsIndexableAttributeType(attributeType AttributeType) bool {
	switch attributeType {
	case AttributeTypeString, AttributeTypeEnum, AttributeTypeRelation, AttributeTypeContactAddress,
		AttributeTypeNumber, AttributeTypeInteger, AttributeTypeBoolean, AttributeTypeCurrency,
		AttributeTypeUser, AttributeTypePrincipal, AttributeTypeFile, AttributeTypeDatetime:
		return true
	}
	return false
}

// IsIndexableAttribute is IsIndexableAttributeType plus the per-attribute
// exceptions: a datetime attribute whose format is `date` or `time` is not
// indexable (its SQL projection casts to ::date / ::time, which Postgres does
// not accept in an index expression), nor is a `duration` (an ISO-8601
// duration string, which the timestamp projection cannot parse).
func IsIndexableAttribute(attribute Attribute) bool {
	if !IsIndexableAttributeType(attribute.Type) {
		return false
	}
	if attribute.Type == AttributeTypeDatetime {
		if meta := ParseMetaAs[DatetimeAttributeMeta](attribute.Meta); meta != nil {
			switch meta.Format {
			case DatetimeFormatDate, DatetimeFormatTime, DatetimeFormatDuration:
				return false
			}
		}
	}
	return true
}

// FindAttributeByPath resolves a dot path (`address.city`) against an attribute
// list: every segment but the last must be an `object` attribute, and the leaf
// must be declared. Returns nil when any segment does not resolve.
func FindAttributeByPath(attributes []Attribute, path []string) *Attribute {
	if len(path) == 0 {
		return nil
	}
	current := attributes
	for depth, segment := range path {
		var found *Attribute
		for index := range current {
			if current[index].Name == segment {
				found = &current[index]
				break
			}
		}
		if found == nil {
			return nil
		}
		if depth == len(path)-1 {
			return found
		}
		if found.Type != AttributeTypeObject {
			return nil
		}
		meta := ParseMetaAs[ObjectAttributeMeta](found.Meta)
		if meta == nil {
			return nil
		}
		current = meta.Attributes
	}
	return nil
}

// EffectiveUniqueConstraints is the entity's full unique-constraint list: one
// constraint per top-level attribute flagged is_unique (named by
// UniqueConstraintNameOf), followed by the declared composite constraints, in
// order. Platform attributes are skipped — `id` is the primary key.
//
// This is THE list data-service builds unique indexes from; metadata-service
// validates the same list, so the two can never disagree about what "unique"
// means for an entity.
func EffectiveUniqueConstraints(attributes []Attribute, declared []UniqueConstraint) []UniqueConstraint {
	effective := make([]UniqueConstraint, 0, len(declared))
	for _, attribute := range attributes {
		if attribute.IsUnique && !attribute.IsPlatformManaged {
			effective = append(effective, UniqueConstraint{
				Name:       UniqueConstraintNameOf(attribute.Name),
				Attributes: []string{attribute.Name},
			})
		}
	}
	return append(effective, declared...)
}

// ValidateUniqueConstraints checks the effective constraint list (see
// EffectiveUniqueConstraints) against the entity's attributes:
//
//   - a declared composite constraint has 2..MaxUniqueConstraintAttributes
//     attributes (a one-attribute constraint is `is_unique` on the attribute —
//     declaring it twice is refused as redundant);
//   - names are snake_case and unique within the entity;
//   - every attribute path resolves to a declared, indexable attribute;
//   - no two constraints cover the same attribute set.
//
// Pure: callers run it over the MERGED attribute list so a constraint may span
// attributes an entity extension contributes.
func ValidateUniqueConstraints(attributes []Attribute, declared []UniqueConstraint) error {
	for _, constraint := range declared {
		if len(constraint.Attributes) < 2 {
			if len(constraint.Attributes) == 1 {
				return fmt.Errorf("unique constraint %q: a one-attribute constraint is declared as is_unique on attribute %q, not in unique_constraints",
					constraint.Name, constraint.Attributes[0])
			}
			return fmt.Errorf("unique constraint %q: at least two attributes are required", constraint.Name)
		}
	}

	effective := EffectiveUniqueConstraints(attributes, declared)
	names := make(map[string]struct{}, len(effective))
	keys := make(map[string]string, len(effective))
	for _, constraint := range effective {
		if !uniqueConstraintNamePattern.MatchString(constraint.Name) {
			return fmt.Errorf("unique constraint %q: name must be snake_case (lowercase letters, digits, underscores)", constraint.Name)
		}
		if _, seen := names[constraint.Name]; seen {
			return fmt.Errorf("unique constraint %q is declared more than once", constraint.Name)
		}
		names[constraint.Name] = struct{}{}

		if len(constraint.Attributes) > MaxUniqueConstraintAttributes {
			return fmt.Errorf("unique constraint %q: at most %d attributes are allowed, got %d",
				constraint.Name, MaxUniqueConstraintAttributes, len(constraint.Attributes))
		}
		seenAttributes := make(map[string]struct{}, len(constraint.Attributes))
		for _, path := range constraint.Attributes {
			if _, dup := seenAttributes[path]; dup {
				return fmt.Errorf("unique constraint %q: attribute %q is listed twice", constraint.Name, path)
			}
			seenAttributes[path] = struct{}{}
			if err := validateIndexableAttributePath(attributes, path); err != nil {
				return fmt.Errorf("unique constraint %q: %w", constraint.Name, err)
			}
		}

		key := strings.Join(constraint.Attributes, "\x00")
		if other, dup := keys[key]; dup {
			return fmt.Errorf("unique constraint %q covers the same attributes as %q", constraint.Name, other)
		}
		keys[key] = constraint.Name
	}
	return nil
}

// validateIndexableAttributePath is the one rule for "can this path be indexed":
// it must resolve to a declared attribute that is not platform-managed, and
// that attribute's type must be indexable.
func validateIndexableAttributePath(attributes []Attribute, path string) error {
	if path == "" {
		return fmt.Errorf("attribute path is empty")
	}
	segments := strings.Split(path, ".")
	for _, segment := range segments {
		if !attributeNamePattern.MatchString(segment) {
			return fmt.Errorf("attribute path %q is not a valid attribute path", path)
		}
	}
	if IsPlatformAttributeName(segments[0]) {
		return fmt.Errorf("attribute %q is a platform attribute and cannot be part of a unique constraint", path)
	}
	leaf := FindAttributeByPath(attributes, segments)
	if leaf == nil {
		return fmt.Errorf("attribute %q is not declared on the entity", path)
	}
	if !IsIndexableAttribute(*leaf) {
		return fmt.Errorf("attribute %q (%s) cannot be indexed", path, leaf.Type)
	}
	return nil
}

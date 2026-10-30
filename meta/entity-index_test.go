package metamodel_test

import (
	"strings"
	"testing"

	metamodel "go.proteos.ai/model/meta"
)

func uniqueTestAttributes() []metamodel.Attribute {
	return metamodel.EnsurePlatformAttributes([]metamodel.Attribute{
		{Name: "external_id", Type: metamodel.AttributeTypeString, IsUnique: true},
		{Name: "org_ref", Type: metamodel.AttributeTypeRelation},
		{Name: "number", Type: metamodel.AttributeTypeInteger},
		{Name: "tags", Type: metamodel.AttributeTypeArray},
		{Name: "notes", Type: metamodel.AttributeTypeKnowledgeText},
		{Name: "birthday", Type: metamodel.AttributeTypeDatetime,
			Meta: metamodel.DatetimeAttributeMeta{Format: metamodel.DatetimeFormatDate}},
		{Name: "address", Type: metamodel.AttributeTypeObject, Meta: metamodel.ObjectAttributeMeta{
			Attributes: []metamodel.Attribute{{Name: "city", Type: metamodel.AttributeTypeString}},
		}},
	})
}

func TestEffectiveUniqueConstraintsFoldsIsUniqueAndSkipsPlatformAttributes(t *testing.T) {
	declared := []metamodel.UniqueConstraint{{Name: "by_org_number", Attributes: []string{"org_ref", "number"}}}
	effective := metamodel.EffectiveUniqueConstraints(uniqueTestAttributes(), declared)

	if len(effective) != 2 {
		t.Fatalf("expected 2 effective constraints (is_unique + declared), got %d: %+v", len(effective), effective)
	}
	// `id` is flagged is_unique as a platform attribute and must NOT appear.
	if effective[0].Name != "unique_external_id" || effective[0].Attributes[0] != "external_id" {
		t.Fatalf("is_unique was not folded first: %+v", effective[0])
	}
	if effective[1].Name != "by_org_number" {
		t.Fatalf("declared constraint lost: %+v", effective[1])
	}
}

func TestValidateUniqueConstraints(t *testing.T) {
	attributes := uniqueTestAttributes()
	tests := []struct {
		name     string
		declared []metamodel.UniqueConstraint
		wantErr  string
	}{
		{name: "none declared, is_unique alone is valid"},
		{
			name:     "a composite constraint over indexable attributes",
			declared: []metamodel.UniqueConstraint{{Name: "by_org_number", Attributes: []string{"org_ref", "number"}}},
		},
		{
			name:     "an object path resolves to its leaf",
			declared: []metamodel.UniqueConstraint{{Name: "by_city_number", Attributes: []string{"address.city", "number"}}},
		},
		{
			name:     "one attribute belongs on is_unique",
			declared: []metamodel.UniqueConstraint{{Name: "by_number", Attributes: []string{"number"}}},
			wantErr:  "declared as is_unique",
		},
		{
			name:     "an empty attribute list",
			declared: []metamodel.UniqueConstraint{{Name: "empty"}},
			wantErr:  "at least two attributes",
		},
		{
			name:     "a name colliding with the folded is_unique name",
			declared: []metamodel.UniqueConstraint{{Name: "unique_external_id", Attributes: []string{"org_ref", "number"}}},
			wantErr:  "declared more than once",
		},
		{
			name:     "a non snake_case name",
			declared: []metamodel.UniqueConstraint{{Name: "ByOrg", Attributes: []string{"org_ref", "number"}}},
			wantErr:  "snake_case",
		},
		{
			name:     "an undeclared attribute",
			declared: []metamodel.UniqueConstraint{{Name: "x", Attributes: []string{"org_ref", "nope"}}},
			wantErr:  "not declared",
		},
		{
			name:     "a platform attribute",
			declared: []metamodel.UniqueConstraint{{Name: "x", Attributes: []string{"id", "number"}}},
			wantErr:  "platform attribute",
		},
		{
			name:     "an array attribute",
			declared: []metamodel.UniqueConstraint{{Name: "x", Attributes: []string{"tags", "number"}}},
			wantErr:  "cannot be indexed",
		},
		{
			name:     "a knowledge-text attribute",
			declared: []metamodel.UniqueConstraint{{Name: "x", Attributes: []string{"notes", "number"}}},
			wantErr:  "cannot be indexed",
		},
		{
			name:     "a date-formatted datetime",
			declared: []metamodel.UniqueConstraint{{Name: "x", Attributes: []string{"birthday", "number"}}},
			wantErr:  "cannot be indexed",
		},
		{
			name:     "the object container itself",
			declared: []metamodel.UniqueConstraint{{Name: "x", Attributes: []string{"address", "number"}}},
			wantErr:  "cannot be indexed",
		},
		{
			name:     "an attribute listed twice in one constraint",
			declared: []metamodel.UniqueConstraint{{Name: "x", Attributes: []string{"number", "number"}}},
			wantErr:  "listed twice",
		},
		{
			name: "two constraints over the same attributes",
			declared: []metamodel.UniqueConstraint{
				{Name: "a", Attributes: []string{"org_ref", "number"}},
				{Name: "b", Attributes: []string{"org_ref", "number"}},
			},
			wantErr: "covers the same attributes",
		},
		{
			name: "too wide",
			declared: []metamodel.UniqueConstraint{{Name: "wide", Attributes: []string{
				"org_ref", "number", "external_id", "address.city", "org_ref", "number", "external_id", "address.city", "org_ref",
			}}},
			wantErr: "at most 8",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			err := metamodel.ValidateUniqueConstraints(attributes, test.declared)
			if test.wantErr == "" {
				if err != nil {
					t.Fatalf("unexpected error: %v", err)
				}
				return
			}
			if err == nil || !strings.Contains(err.Error(), test.wantErr) {
				t.Fatalf("error = %v, want it to contain %q", err, test.wantErr)
			}
		})
	}
}

func TestIsIndexableAttribute(t *testing.T) {
	indexable := []metamodel.AttributeType{
		metamodel.AttributeTypeString, metamodel.AttributeTypeEnum, metamodel.AttributeTypeRelation,
		metamodel.AttributeTypeContactAddress, metamodel.AttributeTypeNumber, metamodel.AttributeTypeInteger,
		metamodel.AttributeTypeBoolean, metamodel.AttributeTypeCurrency, metamodel.AttributeTypeUser,
		metamodel.AttributeTypePrincipal, metamodel.AttributeTypeFile, metamodel.AttributeTypeDatetime,
	}
	for _, attributeType := range indexable {
		if !metamodel.IsIndexableAttribute(metamodel.Attribute{Type: attributeType}) {
			t.Errorf("%s should be indexable", attributeType)
		}
	}
	for _, attributeType := range []metamodel.AttributeType{
		metamodel.AttributeTypeArray, metamodel.AttributeTypeObject, metamodel.AttributeTypeKnowledgeText,
	} {
		if metamodel.IsIndexableAttribute(metamodel.Attribute{Type: attributeType}) {
			t.Errorf("%s should not be indexable", attributeType)
		}
	}
	for _, format := range []metamodel.DatetimeFormat{metamodel.DatetimeFormatDate, metamodel.DatetimeFormatTime, metamodel.DatetimeFormatDuration} {
		attribute := metamodel.Attribute{Type: metamodel.AttributeTypeDatetime, Meta: metamodel.DatetimeAttributeMeta{Format: format}}
		if metamodel.IsIndexableAttribute(attribute) {
			t.Errorf("datetime format %q should not be indexable", format)
		}
	}
}

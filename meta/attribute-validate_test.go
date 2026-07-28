package metamodel_test

import (
	"strings"
	"testing"

	metamodel "go.proteos.ai/model/meta"
)

func TestValidateAttributeDefinitions(t *testing.T) {
	tests := []struct {
		name    string
		attrs   []metamodel.Attribute
		wantErr string
	}{
		{
			name: "valid flat attributes",
			attrs: []metamodel.Attribute{
				{Name: "sentiment", Type: metamodel.AttributeTypeEnum},
				{Name: "score_2", Type: metamodel.AttributeTypeNumber},
			},
		},
		{
			name: "valid nested object and array",
			attrs: []metamodel.Attribute{
				{Name: "summary", Type: metamodel.AttributeTypeObject, Meta: metamodel.ObjectAttributeMeta{
					Attributes: []metamodel.Attribute{{Name: "headline", Type: metamodel.AttributeTypeString}},
				}},
				{Name: "tags", Type: metamodel.AttributeTypeArray, Meta: metamodel.ArrayAttributeMeta{
					Items: &metamodel.Attribute{Name: "tag", Type: metamodel.AttributeTypeString},
				}},
			},
		},
		{
			name:    "empty name",
			attrs:   []metamodel.Attribute{{Name: "", Type: metamodel.AttributeTypeString}},
			wantErr: "empty name",
		},
		{
			name:    "non-snake-case name",
			attrs:   []metamodel.Attribute{{Name: "Sentiment", Type: metamodel.AttributeTypeString}},
			wantErr: "snake_case",
		},
		{
			name:    "kebab-case name",
			attrs:   []metamodel.Attribute{{Name: "my-field", Type: metamodel.AttributeTypeString}},
			wantErr: "snake_case",
		},
		{
			name: "duplicate names",
			attrs: []metamodel.Attribute{
				{Name: "score", Type: metamodel.AttributeTypeNumber},
				{Name: "score", Type: metamodel.AttributeTypeString},
			},
			wantErr: "duplicate",
		},
		{
			name:    "unknown type",
			attrs:   []metamodel.Attribute{{Name: "score", Type: "decimal"}},
			wantErr: "unknown type",
		},
		{
			name: "invalid nested attribute",
			attrs: []metamodel.Attribute{
				{Name: "summary", Type: metamodel.AttributeTypeObject, Meta: metamodel.ObjectAttributeMeta{
					Attributes: []metamodel.Attribute{{Name: "Bad Name", Type: metamodel.AttributeTypeString}},
				}},
			},
			wantErr: "snake_case",
		},
		{
			name: "invalid array item type",
			attrs: []metamodel.Attribute{
				{Name: "tags", Type: metamodel.AttributeTypeArray, Meta: metamodel.ArrayAttributeMeta{
					Items: &metamodel.Attribute{Name: "tag", Type: "nope"},
				}},
			},
			wantErr: "unknown type",
		},
		{
			name: "object meta as a map (post-JSON) still recurses",
			attrs: []metamodel.Attribute{
				{Name: "summary", Type: metamodel.AttributeTypeObject, Meta: map[string]any{
					"attributes": []any{map[string]any{"name": "UPPER", "type": "string"}},
				}},
			},
			wantErr: "snake_case",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			err := metamodel.ValidateAttributeDefinitions(test.attrs)
			if test.wantErr == "" {
				if err != nil {
					t.Fatalf("unexpected error: %v", err)
				}
				return
			}
			if err == nil || !strings.Contains(err.Error(), test.wantErr) {
				t.Fatalf("error = %v, want containing %q", err, test.wantErr)
			}
		})
	}
}

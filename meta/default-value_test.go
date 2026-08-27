package metamodel

import (
	"encoding/json"
	"testing"
)

func TestIsCurrentUserDefaultRecognisesEveryShape(t *testing.T) {
	// The decoded shape: what a stored schema or a request body yields.
	var decoded Attribute
	raw := `{"name":"owner","type":"principal","default_value":{"type":"current_user"}}`
	if err := json.Unmarshal([]byte(raw), &decoded); err != nil {
		t.Fatal(err)
	}
	if !IsCurrentUserDefault(decoded.DefaultValue) {
		t.Fatalf("decoded sentinel not recognised: %#v", decoded.DefaultValue)
	}
	if !IsCurrentUserDefault(CurrentUserDefault()) {
		t.Fatal("the builder's own shape must be recognised")
	}

	for name, value := range map[string]any{
		"nil":            nil,
		"bare id":        "current_user",
		"a real ref":     map[string]any{"type": "person", "id": "alice"},
		"extra keys":     map[string]any{"type": "current_user", "id": "x"},
		"wrong type key": map[string]any{"kind": "current_user"},
	} {
		if IsCurrentUserDefault(value) {
			t.Errorf("%s must not be read as the sentinel", name)
		}
	}
}

func TestAcceptsCurrentUserDefault(t *testing.T) {
	if !AcceptsCurrentUserDefault(AttributeTypeUser) || !AcceptsCurrentUserDefault(AttributeTypePrincipal) {
		t.Fatal("user and principal attributes carry an identity and accept the sentinel")
	}
	for _, other := range []AttributeType{AttributeTypeString, AttributeTypeRelation, AttributeTypeObject} {
		if AcceptsCurrentUserDefault(other) {
			t.Errorf("%s must not accept the sentinel", other)
		}
	}
}

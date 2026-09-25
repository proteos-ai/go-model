package common

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// The list operators are the ones whose wire value is pipe-joined. Every
// parser splits on this predicate, so it has to name exactly the operators
// whose value is a list and nothing else.
func TestComparisonOperator_IsMultiValue(t *testing.T) {
	multi := map[ComparisonOperator]bool{
		ComparisonOperatorIn:          true,
		ComparisonOperatorNotIn:       true,
		ComparisonOperatorContainsAll: true,
	}
	for _, operator := range ComparisonOperators {
		assert.Equal(t, multi[operator], operator.IsMultiValue(), string(operator))
	}
}

func TestComparisonOperator_UnmarshalJSON(t *testing.T) {
	for _, operator := range ComparisonOperators {
		var decoded ComparisonOperator
		require.NoError(t, json.Unmarshal([]byte(`"`+string(operator)+`"`), &decoded), string(operator))
		assert.Equal(t, operator, decoded)
	}

	// Negation is a prefix on this platform (not_in, not_empty, not_contains);
	// the suffix spelling is not an alias.
	var decoded ComparisonOperator
	require.Error(t, json.Unmarshal([]byte(`"contains_not"`), &decoded))
}

package common

import (
	"encoding/json"
	"fmt"

	"golang.org/x/exp/slices"
)

type ComparisonOperator string

const (
	ComparisonOperatorEquals             ComparisonOperator = "eq"
	ComparisonOperatorGreaterThan        ComparisonOperator = "gt"
	ComparisonOperatorLessThan           ComparisonOperator = "lt"
	ComparisonOperatorGreaterThanOrEqual ComparisonOperator = "gte"
	ComparisonOperatorLessThanOrEqual    ComparisonOperator = "lte"
	ComparisonOperatorNotEquals          ComparisonOperator = "ne"
	ComparisonOperatorIn                 ComparisonOperator = "in"
	ComparisonOperatorNotIn              ComparisonOperator = "not_in"
	ComparisonOperatorContains           ComparisonOperator = "contains"
	ComparisonOperatorNotContains        ComparisonOperator = "not_contains"
	ComparisonOperatorContainsAll        ComparisonOperator = "contains_all"
	ComparisonOperatorStartsWith         ComparisonOperator = "starts_with"
	ComparisonOperatorEndsWith           ComparisonOperator = "ends_with"
	ComparisonOperatorEmpty              ComparisonOperator = "empty"
	ComparisonOperatorNotEmpty           ComparisonOperator = "not_empty"
)

var ComparisonOperators = []ComparisonOperator{
	ComparisonOperatorEquals,
	ComparisonOperatorGreaterThan,
	ComparisonOperatorLessThan,
	ComparisonOperatorGreaterThanOrEqual,
	ComparisonOperatorLessThanOrEqual,
	ComparisonOperatorNotEquals,
	ComparisonOperatorIn,
	ComparisonOperatorNotIn,
	ComparisonOperatorContains,
	ComparisonOperatorNotContains,
	ComparisonOperatorContainsAll,
	ComparisonOperatorStartsWith,
	ComparisonOperatorEndsWith,
	ComparisonOperatorEmpty,
	ComparisonOperatorNotEmpty,
}

// IsMultiValue reports whether the operator compares against a LIST of values.
// On the wire such a value travels as one pipe-joined string ("a|b|c"); every
// parser splits on this predicate rather than on a hand-kept operator list, so
// a new list operator cannot be split in one place and missed in another.
func (comparisonOperator ComparisonOperator) IsMultiValue() bool {
	switch comparisonOperator {
	case ComparisonOperatorIn, ComparisonOperatorNotIn, ComparisonOperatorContainsAll:
		return true
	default:
		return false
	}
}

func (ComparisonOperator) Enum() []interface{} {
	enums := []interface{}{}
	for _, element := range ComparisonOperators {
		enums = append(enums, element)
	}
	return enums
}

func (comparisonOperator *ComparisonOperator) UnmarshalJSON(byteArray []byte) error {
	str := string(byteArray)
	if str == `null` {
		*comparisonOperator = ""
		return nil
	}
	type ComparisonOperatorType ComparisonOperator
	var stringValue *ComparisonOperatorType = (*ComparisonOperatorType)(comparisonOperator)
	err := json.Unmarshal(byteArray, &stringValue)
	if err != nil {
		return err
	}

	if slices.Contains(ComparisonOperators, *comparisonOperator) {
		return nil
	}

	return fmt.Errorf("invalid comparisonOperator: %s", *stringValue)
}

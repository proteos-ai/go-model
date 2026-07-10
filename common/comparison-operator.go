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
	ComparisonOperatorStartsWith,
	ComparisonOperatorEndsWith,
	ComparisonOperatorEmpty,
	ComparisonOperatorNotEmpty,
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

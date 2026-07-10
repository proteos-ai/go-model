package common

import (
	"encoding/json"
	"fmt"

	"golang.org/x/exp/slices"
)

type LogicalOperator string

const (
	LogicalOperatorAnd LogicalOperator = "and"
	LogicalOperatorOr  LogicalOperator = "or"
)

var LogicalOperators = []LogicalOperator{
	LogicalOperatorAnd,
	LogicalOperatorOr,
}

func (LogicalOperator) Enum() []interface{} {
	enums := []interface{}{}
	for _, element := range LogicalOperators {
		enums = append(enums, element)
	}
	return enums
}

func (logicalOperator *LogicalOperator) UnmarshalJSON(byteArray []byte) error {
	str := string(byteArray)
	if str == `null` {
		*logicalOperator = ""
		return nil
	}
	type LogicalOperatorType LogicalOperator
	var stringValue *LogicalOperatorType = (*LogicalOperatorType)(logicalOperator)
	err := json.Unmarshal(byteArray, &stringValue)
	if err != nil {
		return err
	}

	if slices.Contains(LogicalOperators, *logicalOperator) {
		return nil
	}

	return fmt.Errorf("invalid logicalOperator: %s", *stringValue)
}

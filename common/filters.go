package common

import (
	"fmt"
	"strings"

	"golang.org/x/exp/slices"
)

type FilterGroup struct {
	LogicalOperator LogicalOperator `json:"logical_operator"`
	Elements        []FilterElement `json:"elements,omitempty"`
	Groups          []FilterGroup   `json:"groups,omitempty"`
}

type FilterElement struct {
	Field    string             `json:"field"`
	Value    string             `json:"value"`
	Operator ComparisonOperator `json:"operator"`
}

// IsValueless reports whether the operator takes no value (empty / not_empty).
// Every other ComparisonOperator expects one.
func (comparisonOperator ComparisonOperator) IsValueless() bool {
	return comparisonOperator == ComparisonOperatorEmpty || comparisonOperator == ComparisonOperatorNotEmpty
}

// FilterGroupValidationError is one problem in a FilterGroup tree. Path points
// at the offending node (`visible_when.elements[1].operator`); Code is one of
// `missing` / `invalid_value`.
type FilterGroupValidationError struct {
	Path    string
	Code    string
	Message string
}

// ValidateFilterGroup checks the shape of a predicate tree (the model list
// filters, layout `visible_when` and action `visible_when` share):
//   - logical_operator is and | or
//   - a group holds at least one element or nested group
//   - every element names a field and a known ComparisonOperator
//   - every operator but empty / not_empty carries a value (set operators
//     are pipe-joined strings — non-empty is enough)
//
// Field names are not checked: what a predicate may reference depends on the
// surface evaluating it. Pure; returns every problem, in tree order.
func ValidateFilterGroup(path string, group *FilterGroup) []FilterGroupValidationError {
	var problems []FilterGroupValidationError
	add := func(at, code, message string) {
		problems = append(problems, FilterGroupValidationError{Path: at, Code: code, Message: message})
	}
	var walk func(at string, current *FilterGroup)
	walk = func(at string, current *FilterGroup) {
		if current == nil {
			add(at, "missing", "filter group is required")
			return
		}
		if current.LogicalOperator != LogicalOperatorAnd && current.LogicalOperator != LogicalOperatorOr {
			add(at+".logical_operator", "invalid_value",
				fmt.Sprintf("logical_operator %q is not one of and|or", current.LogicalOperator))
		}
		if len(current.Elements) == 0 && len(current.Groups) == 0 {
			add(at, "missing", "filter group must contain at least one element or nested group")
		}
		for index, element := range current.Elements {
			elementPath := fmt.Sprintf("%s.elements[%d]", at, index)
			if strings.TrimSpace(element.Field) == "" {
				add(elementPath+".field", "missing", "field is required")
			}
			if !slices.Contains(ComparisonOperators, element.Operator) {
				add(elementPath+".operator", "invalid_value",
					fmt.Sprintf("operator %q is not a known ComparisonOperator", element.Operator))
			}
			if !element.Operator.IsValueless() && element.Value == "" {
				add(elementPath+".value", "missing",
					fmt.Sprintf("value is required for operator %q", element.Operator))
			}
		}
		for index := range current.Groups {
			walk(fmt.Sprintf("%s.groups[%d]", at, index), &current.Groups[index])
		}
	}
	walk(path, group)
	return problems
}

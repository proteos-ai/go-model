package common

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

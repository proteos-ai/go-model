package common

type ClientQueryElement struct {
	Attribute          string
	ComparisonOperator ComparisonOperator
	Value              any
}

type ClientQuery struct {
	LogicalOperator LogicalOperator
	Elements        []ClientQueryElement
}

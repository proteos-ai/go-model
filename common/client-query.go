package common

type ClientQueryElement struct {
	Attribute          string
	ComparisonOperator ComparisonOperator
	Value              any
}

// ClientQuery is the parsed record-list filter. It is a TREE: Elements are the
// leaves joined at this level by LogicalOperator, and Groups are nested
// sub-trees joined into the same level by the same operator — so
// `(a AND b) OR c` is a group with LogicalOperator=OR, one element (c) and one
// sub-group (a AND b).
//
// The flat URL-param form (`?stage=won&_logicalOperator=OR`) only ever produces
// a root with Elements and no Groups; the `_filter` JSON param produces the
// full tree. Both land in this one shape, so everything downstream — SQL
// building, sort extraction — has a single structure to walk.
type ClientQuery struct {
	LogicalOperator LogicalOperator
	Elements        []ClientQueryElement
	Groups          []ClientQuery
}

// IsEmpty reports whether the query carries no conditions at any depth. A
// group holding only empty sub-groups is empty too, which is what lets the
// caller skip the WHERE clause entirely rather than emit "()".
func (query *ClientQuery) IsEmpty() bool {
	if query == nil {
		return true
	}
	if len(query.Elements) > 0 {
		return false
	}
	for i := range query.Groups {
		if !query.Groups[i].IsEmpty() {
			return false
		}
	}
	return true
}

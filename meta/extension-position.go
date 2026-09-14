package metamodel

// ExtensionPosition says where an extension's block goes relative to its
// anchor — shared by page extensions (elements / tabs) and list extensions
// (columns). `first`/`last` mean "inside this container" (for a list: at the
// head / tail of the column list); `before`/`after` mean "next to this node"
// (a sibling of the same kind: element next to element, tab next to tab,
// column next to column).
type ExtensionPosition string

const (
	ExtensionPositionBefore ExtensionPosition = "before"
	ExtensionPositionAfter  ExtensionPosition = "after"
	ExtensionPositionFirst  ExtensionPosition = "first"
	ExtensionPositionLast   ExtensionPosition = "last"
)

// SupportedExtensionPositions is the set of positions this version accepts.
var SupportedExtensionPositions = []ExtensionPosition{
	ExtensionPositionBefore,
	ExtensionPositionAfter,
	ExtensionPositionFirst,
	ExtensionPositionLast,
}

// IsInside reports whether the position targets the anchor's own content /
// the list's ends (first/last) rather than a sibling slot (before/after).
func (position ExtensionPosition) IsInside() bool {
	return position == ExtensionPositionFirst || position == ExtensionPositionLast
}

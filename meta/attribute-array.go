package metamodel

// ArrayAttributeMeta holds array-specific validation options
type ArrayAttributeMeta struct {
	Items             *Attribute `json:"items,omitempty"`
	MinItems          *int       `json:"min_items,omitempty"`
	MaxItems          *int       `json:"max_items,omitempty"`
	ItemsMustBeUnique bool       `json:"items_must_be_unique,omitempty"`
}

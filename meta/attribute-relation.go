package metamodel

import (
	"encoding/json"
	"fmt"

	"golang.org/x/exp/slices"
)

// OnDeleteAction is the policy applied to the host row when the related row is
// deleted. The runtime semantic is enforced at the data-service level; this
// metadata declares the intent.
type OnDeleteAction string

const (
	OnDeleteCascade  OnDeleteAction = "cascade"
	OnDeleteRestrict OnDeleteAction = "restrict"
	OnDeleteSetNull  OnDeleteAction = "set-null"
)

var OnDeleteActions = []OnDeleteAction{OnDeleteCascade, OnDeleteRestrict, OnDeleteSetNull}

func (a *OnDeleteAction) UnmarshalJSON(b []byte) error {
	if string(b) == `null` {
		*a = ""
		return nil
	}
	var s string
	if err := json.Unmarshal(b, &s); err != nil {
		return err
	}
	v := OnDeleteAction(s)
	if !slices.Contains(OnDeleteActions, v) {
		return fmt.Errorf("invalid onDelete: %q", s)
	}
	*a = v
	return nil
}

// RelationAttributeMeta holds the metadata for an attribute of type `relation`.
// A relation attribute is, physically, a foreign-key column on the host entity
// pointing at `RelatedAttribute` on `RelatedEntitySlug`. The predicate reads
// from the host outward — e.g. on `Order.customerId` with predicate
// "is placed by", the sentence is "Order is placed by Customer".
type RelationAttributeMeta struct {
	RelatedEntitySlug string         `json:"related_entity_slug"`
	RelatedAttribute  string         `json:"related_attribute"`
	Predicate         string         `json:"predicate"`
	Description       string         `json:"description,omitempty"`
	OnDelete          OnDeleteAction `json:"on_delete"`
}

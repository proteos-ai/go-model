package metamodel

import (
	"go.proteos.ai/model/common"
	"time"
)

// SelectionMode controls a list's row-selection affordance:
//
//   - on_demand (default): a Select toggle in the toolbar reveals the
//     checkboxes. Reading is the common visit, so nothing is shown until the
//     user says they want to act.
//   - always:              checkboxes are showing from the start — for lists
//     whose whole point is picking rows and running an action.
//   - off:                 rows can never be checked, even when the list
//     carries actions. The escape hatch for a list that must stay read-only.
type SelectionMode string

const (
	SelectionModeOnDemand SelectionMode = "on_demand"
	SelectionModeAlways   SelectionMode = "always"
	SelectionModeOff      SelectionMode = "off"
)

// SupportedSelectionModes is the accepted set; empty normalizes to on_demand.
var SupportedSelectionModes = []SelectionMode{
	SelectionModeOnDemand,
	SelectionModeAlways,
	SelectionModeOff,
}

// Normalized returns the mode with the empty value defaulting to on_demand,
// so lists persisted before the field existed keep working.
func (mode SelectionMode) Normalized() SelectionMode {
	if mode == "" {
		return SelectionModeOnDemand
	}
	return mode
}

func (SelectionMode) Enum() []interface{} {
	out := make([]interface{}, len(SupportedSelectionModes))
	for i, mode := range SupportedSelectionModes {
		out[i] = mode
	}
	return out
}

type List struct {
	Slug       string   `json:"slug"`
	OrgId      string   `json:"org_id" sortable:""`
	ModuleSlug string   `json:"module_slug" sortable:""`
	EntitySlug string   `json:"entity_slug"`
	Name       string   `json:"name"`
	Columns    []Column `json:"columns"`
	// Actions are the list's toolbar buttons, the same shape a page carries.
	// They act on the rows SELECTED in the list, so an `action` button names
	// an `entity_batch` action (invoked once with every selected record id)
	// and a `workflow` button starts one run for the selection. Prefill
	// templates resolve against the list scope {selection, entity, user}.
	Actions []PageAction `json:"actions"`
	// SelectionMode decides whether the list's rows can be checked, and
	// whether the checkboxes are showing from the start. Empty normalizes to
	// SelectionModeOnDemand.
	SelectionMode SelectionMode `json:"selection_mode,omitempty"`
	// DefaultPageSlug optionally names the record page to render when a
	// record is opened from this list. Empty = the org's default page for
	// the entity.
	DefaultPageSlug string               `json:"default_page_slug,omitempty"`
	Sorting         []SortConfig         `json:"sorting"`
	Filters         []common.FilterGroup `json:"filters"`
	CreatedAt       time.Time            `json:"created_at"`
	UpdatedAt       time.Time            `json:"updated_at"`
	CreatedBy       common.UserRef       `json:"created_by"`
	UpdatedBy       common.UserRef       `json:"updated_by"`
}

type Column struct {
	// Attribute names what the column shows — an attribute name or a dot
	// path, told apart by the type of the first segment:
	//
	//   name                 an attribute on the list's own entity
	//   address.city         a leaf inside one of its object attributes
	//   company_id.name      a field of the RELATED record, reached
	//                        through a relation attribute (the first
	//                        segment is the FK)
	//
	// A bare object attribute is not a valid column — it carries no value
	// of its own, only leaves. Sorting follows the same grammar minus the
	// relation hop: an attribute or an object path can be ordered by, a
	// related field cannot (it would need a join). See
	// dbutils.SQLSorting.ToOrderByClause.
	Attribute string `json:"attribute"`
	Label     string `json:"label"`
	Width     int    `json:"width"`
}

type SortConfig struct {
	// Attribute is an attribute name or an object path ("address.city").
	// Relation paths are not sortable — ordering by a related field would
	// need a join the records query does not do.
	Attribute string               `json:"attribute"`
	Direction common.SortDirection `json:"direction"`
}

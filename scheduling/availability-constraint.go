package schedulingmodel

import (
	"time"

	"go.proteos.ai/model/common"
)

// AvailabilityConstraint is one rule about a user's time on top of their
// weekly hours: a recurring window (RecurrenceRule = "DTSTART;TZID=…\nRRULE:…",
// DurationMinutes long, or IsAllDay) that either opens time (available — also
// outside the weekly hours) or blocks it (unavailable). CalendarEventTypeKeys
// scopes the rule to some types: an available window "for intro calls" opens
// it for them AND keeps every other type out ("Tuesday afternoons are for
// intro calls only"); an unavailable window "for demos" blocks demos only.
// Empty = every type.
type AvailabilityConstraint struct {
	Id                    string           `json:"id"`
	OrgId                 string           `json:"org_id"`
	UserId                string           `json:"user_id"`
	Kind                  AvailabilityKind `json:"kind" sortable:""`
	RecurrenceRule        string           `json:"recurrence_rule"`
	DurationMinutes       int              `json:"duration_minutes"`
	IsAllDay              bool             `json:"is_all_day"`
	CalendarEventTypeKeys []string         `json:"calendar_event_type_keys"`
	Label                 string           `json:"label,omitempty" sortable:""`
	CreatedAt             time.Time        `json:"created_at" sortable:""`
	CreatedBy             common.UserRef   `json:"created_by"`
	UpdatedAt             time.Time        `json:"updated_at" sortable:""`
	UpdatedBy             common.UserRef   `json:"updated_by"`
}

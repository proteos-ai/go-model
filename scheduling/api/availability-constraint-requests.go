package schedulingapi

import (
	"go.proteos.ai/model/common"
	schedulingmodel "go.proteos.ai/model/scheduling"
)

// CreateAvailabilityConstraintRequest adds one rule about a user's time.
// RecurrenceRule is the audibene form "DTSTART;TZID=<zone>:<local
// time>\nRRULE:<rule>" — the DTSTART minute must sit on the quarter-hour
// grid. DurationMinutes is the window's length after each occurrence (0 only
// with IsAllDay, which opens / blocks the whole occurrence day in the rule's
// zone). CalendarEventTypeKeys narrows an `available` window to some types
// (empty = every type); the keys must exist in the org.
type CreateAvailabilityConstraintRequest struct {
	Kind                  schedulingmodel.AvailabilityKind `json:"kind" validate:"required,oneof=available unavailable"`
	RecurrenceRule        string                           `json:"recurrence_rule" validate:"required,max=2048"`
	DurationMinutes       int                              `json:"duration_minutes" validate:"min=0,max=10080"`
	IsAllDay              bool                             `json:"is_all_day"`
	CalendarEventTypeKeys []string                         `json:"calendar_event_type_keys" validate:"max=32,dive,required,max=128"`
	Label                 string                           `json:"label" validate:"max=255"`
}

// UpdateAvailabilityConstraintRequest is a partial update; nil leaves a field
// as is. CalendarEventTypeKeys non-nil replaces the list (empty = every type).
type UpdateAvailabilityConstraintRequest struct {
	Kind                  *schedulingmodel.AvailabilityKind `json:"kind,omitempty" validate:"omitempty,oneof=available unavailable"`
	RecurrenceRule        *string                           `json:"recurrence_rule,omitempty" validate:"omitempty,max=2048"`
	DurationMinutes       *int                              `json:"duration_minutes,omitempty" validate:"omitempty,min=0,max=10080"`
	IsAllDay              *bool                             `json:"is_all_day,omitempty"`
	CalendarEventTypeKeys *[]string                         `json:"calendar_event_type_keys,omitempty" validate:"omitempty,max=32,dive,required,max=128"`
	Label                 *string                           `json:"label,omitempty" validate:"omitempty,max=255"`
}

type GetManyAvailabilityConstraintsQuery struct {
	Kind *string `json:"kind" form:"kind" validate:"omitempty,oneof=available unavailable"`
	common.Pagination
	common.Sorting
}

type GetManyAvailabilityConstraintsResponse struct {
	Meta common.ResponseMeta                      `json:"meta"`
	Data []schedulingmodel.AvailabilityConstraint `json:"data"`
}

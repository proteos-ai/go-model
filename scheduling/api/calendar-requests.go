package schedulingapi

import (
	"go.proteos.ai/model/common"
	schedulingmodel "go.proteos.ai/model/scheduling"
)

// UpdateCalendarRequest changes the platform roles of a calendar. Every field
// is a tri-state pointer. IsSynced flips the mirror on/off (backfill + watch,
// or stop watch); IsTarget=true demotes the user's current target calendar in
// the same transaction — there is exactly one target per user.
type UpdateCalendarRequest struct {
	IsSynced             *bool                               `json:"is_synced,omitempty"`
	IsAvailabilitySource *bool                               `json:"is_availability_source,omitempty"`
	IsTarget             *bool                               `json:"is_target,omitempty"`
	Visibility           *schedulingmodel.CalendarVisibility `json:"visibility,omitempty" validate:"omitempty,oneof=private busy_only full"`
	Color                *string                             `json:"color,omitempty" validate:"omitempty,max=32"`
}

type GetManyCalendarsQuery struct {
	CalendarConnectionId *string `json:"calendar_connection_id" form:"calendar_connection_id"`
	IsSynced             *bool   `json:"is_synced" form:"is_synced"`
	IsAvailabilitySource *bool   `json:"is_availability_source" form:"is_availability_source"`
	IsTarget             *bool   `json:"is_target" form:"is_target"`
	common.Pagination
	common.Sorting
}

type GetManyCalendarsResponse struct {
	Meta common.ResponseMeta        `json:"meta"`
	Data []schedulingmodel.Calendar `json:"data"`
}

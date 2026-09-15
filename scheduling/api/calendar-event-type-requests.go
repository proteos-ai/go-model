package schedulingapi

import (
	"go.proteos.ai/model/common"
	schedulingmodel "go.proteos.ai/model/scheduling"
)

// PutCalendarEventTypeRequest creates or replaces one calendar event type
// (PUT /calendar-event-types/:key — the key comes from the path; every field
// is the new value). The three policy flags default to TRUE when absent, the
// way the columns do — a minimal PUT yields an enabled, cancellable,
// reschedulable type; send false to switch one off. SlotIntervalMinutes is the candidate grid (5|10|15|20|
// 30|60); the buffers pad the busy check around each candidate; the notice /
// window / daily cap bound WHEN a slot may be taken (0 = no cap).
type PutCalendarEventTypeRequest struct {
	Name                 string                        `json:"name" validate:"required,max=255"`
	Description          string                        `json:"description" validate:"max=4096"`
	Color                string                        `json:"color" validate:"max=32"`
	DurationMinutes      int                           `json:"duration_minutes" validate:"required,min=5,max=1440"`
	SlotIntervalMinutes  int                           `json:"slot_interval_minutes" validate:"omitempty,oneof=5 10 15 20 30 60"`
	BufferBeforeMinutes  int                           `json:"buffer_before_minutes" validate:"min=0,max=1440"`
	BufferAfterMinutes   int                           `json:"buffer_after_minutes" validate:"min=0,max=1440"`
	MinimumNoticeMinutes int                           `json:"minimum_notice_minutes" validate:"min=0"`
	SchedulingWindowDays int                           `json:"scheduling_window_days" validate:"min=0,max=3650"`
	MaxPerDay            int                           `json:"max_per_day" validate:"min=0"`
	Location             schedulingmodel.EventLocation `json:"location"`
	ConversationTypeKey  string                        `json:"conversation_type_key" validate:"max=255"`
	TitleTemplate        string                        `json:"title_template" validate:"max=1024"`
	IsCancelAllowed      *bool                         `json:"is_cancel_allowed,omitempty"`
	IsRescheduleAllowed  *bool                         `json:"is_reschedule_allowed,omitempty"`
	CancelNoticeMinutes  int                           `json:"cancel_notice_minutes" validate:"min=0"`
	IsEnabled            *bool                         `json:"is_enabled,omitempty"`
}

// UpdateCalendarEventTypeRequest is a partial update; nil leaves a field as
// is. Location, when present, replaces the stored location wholesale.
type UpdateCalendarEventTypeRequest struct {
	Name                 *string                        `json:"name,omitempty" validate:"omitempty,max=255"`
	Description          *string                        `json:"description,omitempty" validate:"omitempty,max=4096"`
	Color                *string                        `json:"color,omitempty" validate:"omitempty,max=32"`
	DurationMinutes      *int                           `json:"duration_minutes,omitempty" validate:"omitempty,min=5,max=1440"`
	SlotIntervalMinutes  *int                           `json:"slot_interval_minutes,omitempty" validate:"omitempty,oneof=5 10 15 20 30 60"`
	BufferBeforeMinutes  *int                           `json:"buffer_before_minutes,omitempty" validate:"omitempty,min=0,max=1440"`
	BufferAfterMinutes   *int                           `json:"buffer_after_minutes,omitempty" validate:"omitempty,min=0,max=1440"`
	MinimumNoticeMinutes *int                           `json:"minimum_notice_minutes,omitempty" validate:"omitempty,min=0"`
	SchedulingWindowDays *int                           `json:"scheduling_window_days,omitempty" validate:"omitempty,min=0,max=3650"`
	MaxPerDay            *int                           `json:"max_per_day,omitempty" validate:"omitempty,min=0"`
	Location             *schedulingmodel.EventLocation `json:"location,omitempty"`
	ConversationTypeKey  *string                        `json:"conversation_type_key,omitempty" validate:"omitempty,max=255"`
	TitleTemplate        *string                        `json:"title_template,omitempty" validate:"omitempty,max=1024"`
	IsCancelAllowed      *bool                          `json:"is_cancel_allowed,omitempty"`
	IsRescheduleAllowed  *bool                          `json:"is_reschedule_allowed,omitempty"`
	CancelNoticeMinutes  *int                           `json:"cancel_notice_minutes,omitempty" validate:"omitempty,min=0"`
	IsEnabled            *bool                          `json:"is_enabled,omitempty"`
}

type GetManyCalendarEventTypesQuery struct {
	IsEnabled *bool `json:"is_enabled" form:"is_enabled"`
	// Search filters by a case-insensitive substring of key or name.
	Search *string `json:"search" form:"search"`
	common.Pagination
	common.Sorting
}

type GetManyCalendarEventTypesResponse struct {
	Meta common.ResponseMeta                 `json:"meta"`
	Data []schedulingmodel.CalendarEventType `json:"data"`
}

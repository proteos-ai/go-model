package schedulingapi

import (
	"time"

	schedulingmodel "go.proteos.ai/model/scheduling"
)

// GetManyCalendarEventsQuery reads a window of events. UserIds / CalendarIds
// pick whose calendars (default: the caller's); colleagues' events come back
// projected by their calendar's visibility. From/To are required and capped by
// the service (≤ 93 days).
type GetManyCalendarEventsQuery struct {
	UserIds           []string  `json:"user_ids" form:"user_ids" validate:"max=64,dive,required,max=64"`
	CalendarIds       []string  `json:"calendar_ids" form:"calendar_ids" validate:"max=64,dive,required,max=64"`
	From              time.Time `json:"from" form:"from" validate:"required"`
	To                time.Time `json:"to" form:"to" validate:"required,gtfield=From"`
	Timezone          string    `json:"timezone" form:"timezone" validate:"max=64"`
	TypeKey           string    `json:"type_key" form:"type_key"`
	RecordEntitySlug  string    `json:"record_entity_slug" form:"record_entity_slug"`
	RecordId          string    `json:"record_id" form:"record_id"`
	ContactId         string    `json:"contact_id" form:"contact_id"`
	IsDeletedIncluded bool      `json:"is_deleted_included" form:"is_deleted_included"`
}

type GetManyCalendarEventsResponse struct {
	Data []schedulingmodel.CalendarEvent `json:"data"`
}

// CalendarAttendeeRequest is one person to invite on a platform write.
type CalendarAttendeeRequest struct {
	Email      string `json:"email" validate:"required,email,max=320"`
	Name       string `json:"name" validate:"max=200"`
	IsOptional bool   `json:"is_optional"`
}

// HostsRequest schedules the event for platform users instead of one explicit
// calendar: the service resolves each host's target calendar, assigns per
// Assignment, and writes one provider event owned by the organizer.
type HostsRequest struct {
	UserIds         []string                   `json:"user_ids" validate:"required,min=1,max=64,dive,required,max=64"`
	Assignment      schedulingmodel.Assignment `json:"assignment" validate:"required,oneof=single round_robin collective"`
	OrganizerUserId string                     `json:"organizer_user_id,omitempty" validate:"omitempty,max=64"`
}

// CreateCalendarEventRequest writes an event through to the provider. CalendarId
// names the caller's own calendar; omitted, the caller's target calendar is
// used (409 no_target_calendar when none is set). Hosts is the alternative
// for scheduling across platform users. IsExclusive
// makes the row claim the owner's time (409 slot_unavailable on overlap).
type CreateCalendarEventRequest struct {
	CalendarId              string                            `json:"calendar_id,omitempty" validate:"max=64"`
	Hosts                   *HostsRequest                     `json:"hosts,omitempty"`
	Title                   string                            `json:"title" validate:"required,max=1024"`
	Description             string                            `json:"description" validate:"max=32768"`
	Location                string                            `json:"location" validate:"max=1024"`
	StartAt                 time.Time                         `json:"start_at" validate:"required"`
	EndAt                   time.Time                         `json:"end_at" validate:"required,gtfield=StartAt"`
	IsAllDay                bool                              `json:"is_all_day"`
	Timezone                string                            `json:"timezone" validate:"max=64"`
	Attendees               []CalendarAttendeeRequest         `json:"attendees" validate:"dive"`
	IsConferencingRequested bool                              `json:"is_conferencing_requested"`
	Transparency            schedulingmodel.EventTransparency `json:"transparency" validate:"omitempty,oneof=busy free"`
	Visibility              schedulingmodel.EventVisibility   `json:"visibility" validate:"omitempty,oneof=default private public"`
	Notify                  schedulingmodel.NotifyMode        `json:"notify,omitempty" validate:"omitempty,oneof=all external_only none"`
	IsExclusive             bool                              `json:"is_exclusive"`
	TypeKey                 string                            `json:"type_key" validate:"max=128"`
	LinkKey                 string                            `json:"link_key" validate:"max=128"`
	Record                  *schedulingmodel.RecordRef        `json:"record,omitempty"`
	// ManageTokenHash is NEVER bound from the wire (json:"-"): the booking
	// orchestrator sets it so a publicly booked event carries its manage
	// token hash from the first insert.
	ManageTokenHash string `json:"-"`
}

// UpdateCalendarEventRequest is a partial update; nil leaves a field as is.
// Attendees non-nil replaces the whole list (an empty list clears it).
type UpdateCalendarEventRequest struct {
	Title                   *string                            `json:"title,omitempty" validate:"omitempty,max=1024"`
	Description             *string                            `json:"description,omitempty" validate:"omitempty,max=32768"`
	Location                *string                            `json:"location,omitempty" validate:"omitempty,max=1024"`
	StartAt                 *time.Time                         `json:"start_at,omitempty"`
	EndAt                   *time.Time                         `json:"end_at,omitempty"`
	IsAllDay                *bool                              `json:"is_all_day,omitempty"`
	Timezone                *string                            `json:"timezone,omitempty" validate:"omitempty,max=64"`
	Attendees               *[]CalendarAttendeeRequest         `json:"attendees,omitempty" validate:"omitempty,dive"`
	IsConferencingRequested *bool                              `json:"is_conferencing_requested,omitempty"`
	Transparency            *schedulingmodel.EventTransparency `json:"transparency,omitempty" validate:"omitempty,oneof=busy free"`
	Visibility              *schedulingmodel.EventVisibility   `json:"visibility,omitempty" validate:"omitempty,oneof=default private public"`
	Notify                  schedulingmodel.NotifyMode         `json:"notify,omitempty" validate:"omitempty,oneof=all external_only none"`
	TypeKey                 *string                            `json:"type_key,omitempty" validate:"omitempty,max=128"`
	Record                  *schedulingmodel.RecordRef         `json:"record,omitempty"`
	// IsRecordCleared removes the record back-link (a nil Record cannot say so).
	IsRecordCleared bool `json:"is_record_cleared,omitempty"`
}

// RespondCalendarEventRequest records the caller's RSVP on an event they were
// invited to.
type RespondCalendarEventRequest struct {
	Status  schedulingmodel.ParticipationStatus `json:"status" validate:"required,oneof=accepted declined tentative"`
	Comment string                              `json:"comment" validate:"max=1024"`
}

// UpdateCalendarSyncRequest is the body of POST /me/calendar-connections/:id/sync
// (currently empty; reserved for a forced re-baseline flag).
type UpdateCalendarSyncRequest struct {
	IsRebaseline bool `json:"is_rebaseline"`
}

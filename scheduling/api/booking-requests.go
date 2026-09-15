package schedulingapi

import (
	"time"

	schedulingmodel "go.proteos.ai/model/scheduling"
)

// BookingContact is the person a link booking is for. Id names an existing
// contact (the page opened with contact_id: no name / email asked); otherwise
// Email is required and Name / Phone are observations that mint or enrich
// the contact. On the authenticated link route an empty contact means the
// caller's own platform contact.
type BookingContact struct {
	Id    string `json:"id,omitempty" validate:"omitempty,max=64"`
	Name  string `json:"name,omitempty" validate:"max=200"`
	Email string `json:"email,omitempty" validate:"omitempty,email,max=320"`
	Phone string `json:"phone,omitempty" validate:"max=64"`
}

// BookSchedulingLinkRequest takes one slot through a link
// (POST /scheduling-links/:key/calendar-events and the public mirror).
// StartAt must be a slot the link currently offers; the type's duration
// sets the end. HostUserId narrows to one of the link's hosts (needed with
// is_host_selectable; a plain round_robin link elects the host itself).
type BookSchedulingLinkRequest struct {
	// TypeKey picks the type on a multi-type link (optional when the link
	// offers one).
	TypeKey    string         `json:"type_key,omitempty" validate:"omitempty,max=128"`
	StartAt    time.Time      `json:"start_at" validate:"required"`
	Timezone   string         `json:"timezone" validate:"required,max=64"`
	Contact    BookingContact `json:"contact"`
	Notes      string         `json:"notes" validate:"max=4096"`
	HostUserId string         `json:"host_user_id,omitempty" validate:"omitempty,max=64"`
}

// BookSchedulingLinkResponse is the authenticated booking answer: the
// stored event plus the one-time manage token (never readable again — only
// its hash is kept) and, when the deployment knows its web origin, the
// manage page URL to hand to the contact.
type BookSchedulingLinkResponse struct {
	schedulingmodel.CalendarEvent
	ManageToken string `json:"manage_token"`
	ManageUrl   string `json:"manage_url,omitempty"`
}

// PublicCalendarEvent is a booked event as its contact sees it on the manage
// page: the when, the how (full location — it is their meeting), the host
// and the policy that governs cancel / reschedule. No provider facts, no
// attendee roster.
type PublicCalendarEvent struct {
	Id       string                        `json:"id"`
	OrgId    string                        `json:"org_id"`
	Title    string                        `json:"title"`
	StartAt  time.Time                     `json:"start_at"`
	EndAt    time.Time                     `json:"end_at"`
	Timezone string                        `json:"timezone,omitempty"`
	Status   schedulingmodel.EventStatus   `json:"status"`
	Location schedulingmodel.EventLocation `json:"location"`
	JoinUrl  string                        `json:"join_url,omitempty"`
	Link     PublicSchedulingLink          `json:"link"`
	Host     PublicHost                    `json:"host"`
	Contact  BookingContact                `json:"contact"`
	// IsCancelled mirrors a tombstoned / cancelled row (the manage page shows
	// the outcome instead of the actions).
	IsCancelled bool       `json:"is_cancelled"`
	DeletedAt   *time.Time `json:"deleted_at,omitempty"`
}

type PublicBookSchedulingLinkResponse struct {
	PublicCalendarEvent
	ManageToken string `json:"manage_token"`
	ManageUrl   string `json:"manage_url,omitempty"`
}

// GetPublicSlotsQuery is the public slots window (query string).
type GetPublicSlotsQuery struct {
	TypeKey    string    `json:"type_key" form:"type_key" validate:"omitempty,max=128"`
	From       time.Time `json:"from" form:"from" validate:"required"`
	To         time.Time `json:"to" form:"to" validate:"required,gtfield=From"`
	Timezone   string    `json:"timezone" form:"timezone" validate:"required,max=64"`
	HostUserId string    `json:"host_user_id" form:"host_user_id" validate:"omitempty,max=64"`
}

// GetPublicSchedulingLinkQuery carries the optional contact the page was
// opened for.
type GetPublicSchedulingLinkQuery struct {
	ContactId string `json:"contact_id" form:"contact_id" validate:"omitempty,max=64"`
}

// CancelPublicCalendarEventRequest / ReschedulePublicCalendarEventRequest
// authenticate by the manage token the booking answered with.
type CancelPublicCalendarEventRequest struct {
	Token  string `json:"token" validate:"required,max=128"`
	Reason string `json:"reason" validate:"max=1024"`
}

type ReschedulePublicCalendarEventRequest struct {
	Token      string    `json:"token" validate:"required,max=128"`
	StartAt    time.Time `json:"start_at" validate:"required"`
	Timezone   string    `json:"timezone" validate:"required,max=64"`
	HostUserId string    `json:"host_user_id,omitempty" validate:"omitempty,max=64"`
}

// GetPublicCalendarEventQuery is the manage-page read (?token=).
type GetPublicCalendarEventQuery struct {
	Token string `json:"token" form:"token" validate:"required,max=128"`
}

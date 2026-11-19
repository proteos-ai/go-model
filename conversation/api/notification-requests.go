package conversationapi

import "time"

// NotificationTemplateKey names one platform notification: a mail the platform
// mailbox sends a platform user about something that happened in their org.
// Each key has exactly one renderer in conversation-service and one data shape
// below; an unknown key is rejected (notification_template_unknown).
type NotificationTemplateKey string

const (
	// NotificationTemplateSchedulingLinkBooked tells a host that someone booked
	// time with them through a scheduling link. Data: SchedulingLinkBookedData.
	NotificationTemplateSchedulingLinkBooked NotificationTemplateKey = "scheduling_link_booked"
	// NotificationTemplateSchedulingLinkRescheduled tells a host that the
	// person who booked them moved the meeting to another time. Data:
	// SchedulingLinkRescheduledData.
	NotificationTemplateSchedulingLinkRescheduled NotificationTemplateKey = "scheduling_link_rescheduled"
	// NotificationTemplateSchedulingLinkCancelled tells a host that a booking
	// is off their calendar: the person cancelled it, or moved it to a time
	// another host takes. Data: SchedulingLinkCancelledData.
	NotificationTemplateSchedulingLinkCancelled NotificationTemplateKey = "scheduling_link_cancelled"
)

// SendNotificationRequest is the body of the internal (X-Internal-Token)
// POST /conversations/v1/internal/notifications. The org rides the body: the
// caller is a first-party service with no user token.
//
// DedupeKey makes the send at-most-once per (org, key): a producer retrying
// the same event gets is_duplicate back instead of a second mail. Pick a key
// that names the event, not the attempt ("scheduling_link_booked:<event id>:
// <host id>").
type SendNotificationRequest struct {
	OrgId          string                  `json:"org_id" validate:"required"`
	PlatformUserId string                  `json:"platform_user_id" validate:"required"`
	TemplateKey    NotificationTemplateKey `json:"template_key" validate:"required"`
	Data           map[string]any          `json:"data"`
	DedupeKey      string                  `json:"dedupe_key" validate:"required"`
}

// SendNotificationResponse reports what the send produced. A duplicate carries
// no ids: the first delivery of the key owns them.
type SendNotificationResponse struct {
	IsDuplicate    bool   `json:"is_duplicate"`
	ConversationId string `json:"conversation_id,omitempty"`
	MessageId      string `json:"message_id,omitempty"`
}

// SchedulingLinkBookedData is the data of NotificationTemplateSchedulingLinkBooked.
// Times are instants; Timezone is the IANA zone they are shown in (the
// host's), UTC when empty or unknown.
type SchedulingLinkBookedData struct {
	Title        string    `json:"title"`
	StartAt      time.Time `json:"start_at"`
	EndAt        time.Time `json:"end_at"`
	Timezone     string    `json:"timezone"`
	ContactName  string    `json:"contact_name"`
	ContactEmail string    `json:"contact_email"`
	Location     string    `json:"location"`
	Notes        string    `json:"notes"`
	// EventUrl opens the booked event in the web app.
	EventUrl string `json:"event_url"`
}

// SchedulingLinkRescheduledData is the data of
// NotificationTemplateSchedulingLinkRescheduled: the booking as it stands
// now, plus the time it had before the move.
type SchedulingLinkRescheduledData struct {
	SchedulingLinkBookedData
	PreviousStartAt time.Time `json:"previous_start_at"`
	PreviousEndAt   time.Time `json:"previous_end_at"`
}

// SchedulingLinkCancelledData is the data of
// NotificationTemplateSchedulingLinkCancelled. StartAt / EndAt are when the
// meeting would have been, shown in Timezone (the host's).
type SchedulingLinkCancelledData struct {
	Title        string    `json:"title"`
	StartAt      time.Time `json:"start_at"`
	EndAt        time.Time `json:"end_at"`
	Timezone     string    `json:"timezone"`
	ContactName  string    `json:"contact_name"`
	ContactEmail string    `json:"contact_email"`
	// Reason is what the person wrote when cancelling, if anything.
	Reason string `json:"reason"`
	// IsMovedToAnotherHost: the booking was not dropped but moved to a time
	// another host takes, so it left THIS host's calendar.
	IsMovedToAnotherHost bool `json:"is_moved_to_another_host"`
	// CalendarUrl opens the host's calendar on the day the meeting was on.
	CalendarUrl string `json:"calendar_url"`
}

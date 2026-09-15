package schedulingmodel

import (
	"time"

	"go.proteos.ai/model/common"
	conversationmodel "go.proteos.ai/model/conversation"
)

// CalendarAttendee is one person or resource on an event. It embeds the
// platform's person snapshot (ContactRef: external_id/name/email/contact_id/
// platform_user, hydrated from conversation-service's contact directory) —
// the same shape a conversation participant and a message recipient carry —
// plus the calendar facts. Hosts of a platform-scheduled event are attendees
// too, marked IsHost; the calendar owner's own row is IsSelf.
type CalendarAttendee struct {
	conversationmodel.ContactRef
	ResponseStatus ParticipationStatus `json:"response_status"`
	IsOptional     bool                `json:"is_optional"`
	IsOrganizer    bool                `json:"is_organizer"`
	IsSelf         bool                `json:"is_self"`
	// IsHost marks a platform user whose time this event claims (set at
	// creation from the hosts list; always false on provider-mirrored events).
	IsHost bool `json:"is_host"`
	// IsResource marks rooms and equipment — never resolved to a contact.
	IsResource bool `json:"is_resource"`
}

// RecordRef is the one seam between the platform's calendar and a module's
// business object: the record this event belongs to (an appointment, an
// interview, a meeting). Named like contact_record_link's fields.
type RecordRef struct {
	EntitySlug string `json:"entity_slug"`
	RecordId   string `json:"record_id"`
}

// EventLocation is how a calendar event type meets.
type EventLocation struct {
	Kind LocationKind `json:"kind"`
	// IsIntegratedConferencing asks the provider to mint a Meet / Teams link on
	// each event of this type.
	IsIntegratedConferencing bool   `json:"is_integrated_conferencing"`
	JoinUrl                  string `json:"join_url,omitempty"`
	Phone                    string `json:"phone,omitempty"`
	Address                  string `json:"address,omitempty"`
	Instructions             string `json:"instructions,omitempty"`
}

// CalendarEvent is ONE entry in ONE provider calendar, mirrored two-way. It is
// the scheduling primitive: a platform-created row with IsExclusive claims its
// owner's time (enforced by a DB exclusion constraint); mirrored provider rows
// never do. All copies of one meeting across calendars share IcalUid.
type CalendarEvent struct {
	Id         string `json:"id"`
	OrgId      string `json:"org_id"`
	CalendarId string `json:"calendar_id" sortable:""`
	// Owner is the calendar's owner (denormalized for per-owner window reads and the exclusion constraint, which keys on the scalar owner_id column).
	Owner common.UserRef `json:"owner"`

	// ProviderEventId is the id in THIS calendar; empty while a platform write
	// is pending. IcalUid is the RFC 5545 UID shared by every copy of the
	// meeting. SeriesId is the recurring master's id on expanded instances;
	// Recurrence holds the RRULE / EXDATE lines on a series master only.
	ProviderEventId     string   `json:"provider_event_id,omitempty"`
	IcalUid             string   `json:"ical_uid,omitempty"`
	SeriesId            string   `json:"series_id,omitempty"`
	IsRecurringInstance bool     `json:"is_recurring_instance"`
	Recurrence          []string `json:"recurrence,omitempty"`

	Title       string    `json:"title" sortable:""`
	Description string    `json:"description,omitempty"`
	StartAt     time.Time `json:"start_at" sortable:""`
	EndAt       time.Time `json:"end_at" sortable:""`
	IsAllDay    bool      `json:"is_all_day"`
	// Timezone is the provider's zone for the event (IANA on Google, the raw
	// Windows name on Graph) — display and all-day boundaries only.
	Timezone             string               `json:"timezone,omitempty"`
	Location             string               `json:"location,omitempty"`
	JoinUrl              string               `json:"join_url,omitempty"`
	ConferencingProvider ConferencingProvider `json:"conferencing_provider,omitempty"`

	Organizer        conversationmodel.ContactRef `json:"organizer"`
	Attendees        []CalendarAttendee           `json:"attendees"`
	OwnerResponseStatus ParticipationStatus          `json:"owner_response_status"`

	Transparency      EventTransparency `json:"transparency"`
	Status            EventStatus       `json:"status" sortable:""`
	Visibility        EventVisibility   `json:"visibility"`
	Etag              string            `json:"etag,omitempty"`
	ProviderUpdatedAt *time.Time        `json:"provider_updated_at,omitempty"`

	Source      EventSource `json:"source"`
	IsExclusive bool        `json:"is_exclusive"`
	// TypeKey / LinkKey are set when the platform scheduled the event through a
	// calendar_event_type / scheduling_link; empty on mirrored provider events
	// until a module or agent classifies them.
	TypeKey string     `json:"type_key,omitempty"`
	LinkKey string     `json:"link_key,omitempty"`
	Record  *RecordRef `json:"record,omitempty"`
	// ManageTokenHash backs the public manage link of a publicly booked event.
	// Never on the wire.
	ManageTokenHash string `json:"-"`
	// IsProjected marks a read where the viewer only sees busy blocks of a
	// colleague's calendar (title "Busy", everything else stripped).
	IsProjected bool `json:"is_projected,omitempty"`

	WriteStatus EventWriteStatus `json:"write_status"`
	WriteError  string           `json:"write_error,omitempty"`
	// DeletedAt is the tombstone: the provider (or a platform delete) removed
	// the event. Rows are kept so consumers can purge derived state.
	DeletedAt *time.Time     `json:"deleted_at,omitempty"`
	CreatedAt time.Time      `json:"created_at" sortable:""`
	CreatedBy common.UserRef `json:"created_by"`
	UpdatedAt time.Time      `json:"updated_at" sortable:""`
	UpdatedBy common.UserRef `json:"updated_by"`
}

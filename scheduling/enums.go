// Package schedulingmodel is the shared model of scheduling-service: calendar
// connections, calendars, the two-way calendar_event mirror, calendar event
// types, scheduling links, scheduling profiles and availability constraints.
// The platform knows time, people, calendars and the act of taking a slot; it
// never knows what a slot is for — business objects live in modules and link
// back through CalendarEvent.Record.
package schedulingmodel

// CalendarProvider is the calendar backend behind a connection.
type CalendarProvider string

const (
	CalendarProviderGoogle    CalendarProvider = "google"
	CalendarProviderMicrosoft CalendarProvider = "microsoft"
)

// connectorKeyToProvider is the ONLY place connector keys are spelled: a
// connector-service grant of connector `google-calendar` drives the google
// provider, `microsoft-calendar` the microsoft one.
var connectorKeyToProvider = map[string]CalendarProvider{
	"google-calendar":    CalendarProviderGoogle,
	"microsoft-calendar": CalendarProviderMicrosoft,
}

// CalendarProviderForConnectorKey maps a connector-service connector key onto
// the calendar provider it grants access to; false for non-calendar connectors.
func CalendarProviderForConnectorKey(connectorKey string) (CalendarProvider, bool) {
	provider, isKnown := connectorKeyToProvider[connectorKey]
	return provider, isKnown
}

// CalendarConnectionStatus is the health of our binding to a grant.
type CalendarConnectionStatus string

const (
	CalendarConnectionStatusActive  CalendarConnectionStatus = "active"
	CalendarConnectionStatusError   CalendarConnectionStatus = "error"
	CalendarConnectionStatusRevoked CalendarConnectionStatus = "revoked"
)

// CalendarVisibility is what COLLEAGUES see of a calendar's events: nothing,
// only busy blocks, or everything.
type CalendarVisibility string

const (
	CalendarVisibilityPrivate  CalendarVisibility = "private"
	CalendarVisibilityBusyOnly CalendarVisibility = "busy_only"
	CalendarVisibilityFull     CalendarVisibility = "full"
)

// CalendarSyncStatus is the mirror state of one calendar.
type CalendarSyncStatus string

const (
	CalendarSyncStatusIdle    CalendarSyncStatus = "idle"
	CalendarSyncStatusSyncing CalendarSyncStatus = "syncing"
	CalendarSyncStatusError   CalendarSyncStatus = "error"
)

// ParticipationStatus is an attendee's RSVP.
type ParticipationStatus string

const (
	ParticipationStatusNeedsAction ParticipationStatus = "needs_action"
	ParticipationStatusAccepted    ParticipationStatus = "accepted"
	ParticipationStatusDeclined    ParticipationStatus = "declined"
	ParticipationStatusTentative   ParticipationStatus = "tentative"
)

// LocationKind is how a calendar event type meets.
type LocationKind string

const (
	LocationKindVideo    LocationKind = "video"
	LocationKindPhone    LocationKind = "phone"
	LocationKindInPerson LocationKind = "in_person"
	LocationKindCustom   LocationKind = "custom"
)

// EventTransparency says whether an event blocks its owner's time.
type EventTransparency string

const (
	EventTransparencyBusy EventTransparency = "busy"
	EventTransparencyFree EventTransparency = "free"
)

// EventStatus is the provider lifecycle state of an event.
type EventStatus string

const (
	EventStatusConfirmed EventStatus = "confirmed"
	EventStatusTentative EventStatus = "tentative"
	EventStatusCancelled EventStatus = "cancelled"
)

// EventVisibility is the provider-side visibility of one event (distinct from
// the platform's per-calendar CalendarVisibility).
type EventVisibility string

const (
	EventVisibilityDefault EventVisibility = "default"
	EventVisibilityPrivate EventVisibility = "private"
	EventVisibilityPublic  EventVisibility = "public"
)

// EventSource says who created the row first: the provider (mirrored) or the
// platform (written through).
type EventSource string

const (
	EventSourceProvider EventSource = "provider"
	EventSourcePlatform EventSource = "platform"
)

// EventWriteStatus is the delivery state of a platform write to the provider.
// pending_mirror marks a host's copy of a meeting the platform wrote to the
// ORGANIZER's calendar: the provider fans the invite out and the host's own
// sync later merges the copy onto this row by ical_uid.
type EventWriteStatus string

const (
	EventWriteStatusSynced        EventWriteStatus = "synced"
	EventWriteStatusPending       EventWriteStatus = "pending"
	EventWriteStatusPendingMirror EventWriteStatus = "pending_mirror"
	EventWriteStatusFailed        EventWriteStatus = "failed"
)

// ConferencingProvider names integrated conferencing minted by the provider.
type ConferencingProvider string

const (
	ConferencingProviderGoogleMeet ConferencingProvider = "google_meet"
	ConferencingProviderTeams      ConferencingProvider = "teams"
)

// Assignment says how many of a link's hosts end up in the meeting.
type Assignment string

const (
	AssignmentSingle     Assignment = "single"
	AssignmentRoundRobin Assignment = "round_robin"
	AssignmentCollective Assignment = "collective"
)

// AvailabilityKind is what a constraint does to the window it covers.
type AvailabilityKind string

const (
	AvailabilityKindAvailable   AvailabilityKind = "available"
	AvailabilityKindUnavailable AvailabilityKind = "unavailable"
)

// NotifyMode controls whom the provider emails about a platform write.
type NotifyMode string

const (
	NotifyModeAll          NotifyMode = "all"
	NotifyModeExternalOnly NotifyMode = "external_only"
	NotifyModeNone         NotifyMode = "none"
)

// RecurrenceScope selects what an edit of a recurring instance targets.
type RecurrenceScope string

const (
	RecurrenceScopeInstance RecurrenceScope = "instance"
	RecurrenceScopeSeries   RecurrenceScope = "series"
)

// SlotUnavailabilityReason explains why a host is not offered at a time
// (verbose slot responses).
type SlotUnavailabilityReason string

const (
	SlotUnavailabilityNoCalendar       SlotUnavailabilityReason = "no_calendar"
	SlotUnavailabilityOutsideHours     SlotUnavailabilityReason = "outside_hours"
	SlotUnavailabilityBlocked          SlotUnavailabilityReason = "blocked"
	SlotUnavailabilityOtherEventType   SlotUnavailabilityReason = "other_event_type"
	SlotUnavailabilityCalendarBusy     SlotUnavailabilityReason = "calendar_busy"
	SlotUnavailabilityExclusiveEvent   SlotUnavailabilityReason = "exclusive_event"
	SlotUnavailabilityMinimumNotice    SlotUnavailabilityReason = "minimum_notice"
	SlotUnavailabilitySchedulingWindow SlotUnavailabilityReason = "scheduling_window"
	SlotUnavailabilityDailyLimit       SlotUnavailabilityReason = "daily_limit"
)

// RescheduleHost says which host a rescheduled link booking lands on: same
// keeps the booked host (slots of that host only); any re-runs the link's
// assignment (a host the booker picked is allowed).
type RescheduleHost string

const (
	RescheduleHostSame RescheduleHost = "same"
	RescheduleHostAny  RescheduleHost = "any"
)

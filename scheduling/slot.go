package schedulingmodel

import "time"

// Slot is one bookable start for the given hosts (all of them for collective
// links, the union's contributors otherwise).
type Slot struct {
	StartAt     time.Time `json:"start_at"`
	EndAt       time.Time `json:"end_at"`
	HostUserIds []string  `json:"host_user_ids"`
}

// SlotUnavailability explains why a host is not offered at a time. It is
// returned only on verbose slot reads; CalendarEventId is set when a mirrored
// event is the cause.
type SlotUnavailability struct {
	StartAt         time.Time                `json:"start_at"`
	EndAt           time.Time                `json:"end_at"`
	HostUserId      string                   `json:"host_user_id"`
	Reason          SlotUnavailabilityReason `json:"reason"`
	CalendarEventId string                   `json:"calendar_event_id,omitempty"`
}

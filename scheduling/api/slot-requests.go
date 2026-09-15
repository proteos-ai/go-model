package schedulingapi

import (
	"time"

	schedulingmodel "go.proteos.ai/model/scheduling"
)

// SlotParameters are the slot rules of one computation, spelled explicitly
// when no calendar_event_type governs the request (the same fields a type
// carries). Zero window / notice / cap = no cap.
type SlotParameters struct {
	DurationMinutes      int `json:"duration_minutes" validate:"required,min=5,max=1440"`
	SlotIntervalMinutes  int `json:"slot_interval_minutes" validate:"omitempty,oneof=5 10 15 20 30 60"`
	BufferBeforeMinutes  int `json:"buffer_before_minutes" validate:"min=0,max=1440"`
	BufferAfterMinutes   int `json:"buffer_after_minutes" validate:"min=0,max=1440"`
	MinimumNoticeMinutes int `json:"minimum_notice_minutes" validate:"min=0"`
	SchedulingWindowDays int `json:"scheduling_window_days" validate:"min=0,max=3650"`
	MaxPerDay            int `json:"max_per_day" validate:"min=0"`
}

// GetSlotsRequest asks for the bookable starts of some hosts in [From, To)
// (at most 31 days). Exactly one of TypeKey (the type's rules), LinkKey (the
// link's type + hosts + assignment; phase 6) or Parameters governs the rules.
// HostUserIds are the platform users to offer; Assignment decides how their
// free time combines: single / round_robin = the union (a slot is offered when
// ANY host is free), collective = the intersection (ALL hosts). Timezone is
// the requester's IANA zone: slots are grouped by its local day and the
// candidate grid is aligned to its midnight. IsVerbose adds the typed
// unavailability windows per host.
type GetSlotsRequest struct {
	TypeKey     string                     `json:"type_key,omitempty" validate:"max=128"`
	LinkKey     string                     `json:"link_key,omitempty" validate:"max=128"`
	Parameters  *SlotParameters            `json:"parameters,omitempty"`
	HostUserIds []string                   `json:"host_user_ids" validate:"max=64,dive,required,max=64"`
	Assignment  schedulingmodel.Assignment `json:"assignment,omitempty" validate:"omitempty,oneof=single round_robin collective"`
	From        time.Time                  `json:"from" validate:"required"`
	To          time.Time                  `json:"to" validate:"required,gtfield=From"`
	Timezone    string                     `json:"timezone" validate:"required,max=64"`
	IsVerbose   bool                       `json:"is_verbose"`
}

// GetSlotsResponse groups the offered slots by the requester's local day
// ("YYYY-MM-DD"); Unavailability is filled only for verbose requests.
type GetSlotsResponse struct {
	Slots          map[string][]schedulingmodel.Slot    `json:"slots"`
	Unavailability []schedulingmodel.SlotUnavailability `json:"unavailability,omitempty"`
}

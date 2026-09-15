package schedulingmodel

import (
	"time"

	"go.proteos.ai/model/common"
)

// CalendarEventType is the *what* of a scheduled event — the reusable preset
// of slot rules a link, an internal booking, or a module books against.
// Config, so it is keyed by (org_id, key). It links availability constraints
// (which types a constraint applies to), a conversation type (how a
// conversation born from this event is handled), and modules' own types
// (interview_type, appointment_type records point at a key).
type CalendarEventType struct {
	OrgId       string `json:"org_id"`
	Key         string `json:"key" sortable:""`
	Name        string `json:"name" sortable:""`
	Description string `json:"description,omitempty"`
	Color       string `json:"color,omitempty"`

	DurationMinutes      int `json:"duration_minutes"`
	SlotIntervalMinutes  int `json:"slot_interval_minutes"`
	BufferBeforeMinutes  int `json:"buffer_before_minutes"`
	BufferAfterMinutes   int `json:"buffer_after_minutes"`
	MinimumNoticeMinutes int `json:"minimum_notice_minutes"`
	// SchedulingWindowDays caps how far ahead a slot may be taken; 0 = no cap.
	SchedulingWindowDays int `json:"scheduling_window_days"`
	// MaxPerDay caps exclusive events of this type per host per day; 0 = no cap.
	MaxPerDay int `json:"max_per_day"`

	Location EventLocation `json:"location"`
	// ConversationTypeKey is a soft reference to conversation-service's
	// conversation_type: Ava / summaries pick it up from the event.
	ConversationTypeKey string `json:"conversation_type_key,omitempty"`
	// TitleTemplate names a booked event (Liquid): variables type.name,
	// contact.name, contact.email, host.display_name, link.name. Empty =
	// the platform default "{{ type.name }} — {{ contact.name }}".
	TitleTemplate string `json:"title_template,omitempty"`

	IsCancelAllowed     bool `json:"is_cancel_allowed"`
	IsRescheduleAllowed bool `json:"is_reschedule_allowed"`
	CancelNoticeMinutes int  `json:"cancel_notice_minutes"`
	IsEnabled           bool `json:"is_enabled" sortable:""`

	CreatedAt time.Time      `json:"created_at" sortable:""`
	CreatedBy common.UserRef `json:"created_by"`
	UpdatedAt time.Time      `json:"updated_at" sortable:""`
	UpdatedBy common.UserRef `json:"updated_by"`
}

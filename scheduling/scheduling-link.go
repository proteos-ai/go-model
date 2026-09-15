package schedulingmodel

import (
	"time"

	"go.proteos.ai/model/common"
)

// SchedulingLink is the entry point for taking a slot: which hosts, how they
// are assigned, and which calendar event type governs the rules. Public links
// resolve on the unauthenticated scheduling page (/s/$orgId/$key); private
// ones need a platform token (customer portals, internal tools). Config, keyed
// by (org_id, key).
type SchedulingLink struct {
	OrgId string `json:"org_id"`
	Key   string `json:"key" sortable:""`
	Name  string `json:"name" sortable:""`
	// CalendarEventTypeKeys are the calendar event types offered: one = pre-defined, several
	// = the booker picks (type_key on slots / booking then required).
	CalendarEventTypeKeys    []string   `json:"calendar_event_type_keys"`
	HostUserIds []string   `json:"host_user_ids"`
	Assignment  Assignment `json:"assignment"`
	IsPublic    bool       `json:"is_public" sortable:""`
	IsEnabled   bool       `json:"is_enabled" sortable:""`
	// IsHostSelectable lets the booker pick one of the hosts on the page
	// (slots narrowed to that host); off = the assignment decides.
	IsHostSelectable bool `json:"is_host_selectable"`
	// RescheduleHost: same keeps the booked host on a reschedule, any re-runs
	// the assignment.
	RescheduleHost RescheduleHost `json:"reschedule_host"`
	// Url is COMPUTED on reads: the scheduling page address for this link
	// (empty when the deployment has no WEB_PUBLIC_URL).
	Url       string         `json:"url,omitempty"`
	CreatedAt time.Time      `json:"created_at" sortable:""`
	CreatedBy common.UserRef `json:"created_by"`
	UpdatedAt time.Time      `json:"updated_at" sortable:""`
	UpdatedBy common.UserRef `json:"updated_by"`
}

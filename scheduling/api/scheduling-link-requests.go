package schedulingapi

import (
	"go.proteos.ai/model/common"
	schedulingmodel "go.proteos.ai/model/scheduling"
)

// CreateSchedulingLinkRequest creates a link (POST /scheduling-links). The
// key is the URL segment (/s/:orgId/:key), lowercase kebab / snake like a
// type key. CalendarEventTypeKeys: one type = pre-defined, several = the booker picks.
// `single` needs exactly one host. IsEnabled defaults to true,
// RescheduleHost to same.
type CreateSchedulingLinkRequest struct {
	Key              string                         `json:"key" validate:"required,max=128"`
	Name             string                         `json:"name" validate:"required,max=255"`
	CalendarEventTypeKeys         []string                       `json:"calendar_event_type_keys" validate:"required,min=1,max=16,dive,required,max=128"`
	HostUserIds      []string                       `json:"host_user_ids" validate:"required,min=1,max=64,dive,required,max=64"`
	Assignment       schedulingmodel.Assignment     `json:"assignment" validate:"required,oneof=single round_robin collective"`
	IsPublic         bool                           `json:"is_public"`
	IsEnabled        *bool                          `json:"is_enabled,omitempty"`
	IsHostSelectable bool                           `json:"is_host_selectable"`
	RescheduleHost   schedulingmodel.RescheduleHost `json:"reschedule_host,omitempty" validate:"omitempty,oneof=same any"`
}

// UpdateSchedulingLinkRequest is a partial update; nil leaves a field as is.
// HostUserIds non-nil replaces the whole list.
type UpdateSchedulingLinkRequest struct {
	Name             *string                         `json:"name,omitempty" validate:"omitempty,max=255"`
	CalendarEventTypeKeys         *[]string                       `json:"calendar_event_type_keys,omitempty" validate:"omitempty,min=1,max=16,dive,required,max=128"`
	HostUserIds      *[]string                       `json:"host_user_ids,omitempty" validate:"omitempty,min=1,max=64,dive,required,max=64"`
	Assignment       *schedulingmodel.Assignment     `json:"assignment,omitempty" validate:"omitempty,oneof=single round_robin collective"`
	IsPublic         *bool                           `json:"is_public,omitempty"`
	IsEnabled        *bool                           `json:"is_enabled,omitempty"`
	IsHostSelectable *bool                           `json:"is_host_selectable,omitempty"`
	RescheduleHost   *schedulingmodel.RescheduleHost `json:"reschedule_host,omitempty" validate:"omitempty,oneof=same any"`
}

type GetManySchedulingLinksQuery struct {
	TypeKey   *string `json:"type_key" form:"type_key"`
	IsPublic  *bool   `json:"is_public" form:"is_public"`
	IsEnabled *bool   `json:"is_enabled" form:"is_enabled"`
	// Search filters by a case-insensitive substring of key or name.
	Search *string `json:"search" form:"search"`
	common.Pagination
	common.Sorting
}

type GetManySchedulingLinksResponse struct {
	Meta common.ResponseMeta              `json:"meta"`
	Data []schedulingmodel.SchedulingLink `json:"data"`
}

// PublicSchedulingLink is what the unauthenticated scheduling page sees of a
// link: no user ids beyond what host selection needs, the type's booking
// facts, never the provider-side location details (those come with the
// booked event).
type PublicSchedulingLink struct {
	OrgId string `json:"org_id"`
	Key   string `json:"key"`
	Name  string `json:"name"`
	// Types offered; one = pre-defined, several = the page shows a picker.
	Types            []PublicCalendarEventType      `json:"types"`
	Hosts            []PublicHost                   `json:"hosts"`
	Assignment       schedulingmodel.Assignment     `json:"assignment"`
	IsHostSelectable bool                           `json:"is_host_selectable"`
	RescheduleHost   schedulingmodel.RescheduleHost `json:"reschedule_host"`
	// IsContactBound: the page was opened with a contact_id the org knows —
	// the booker enters no name / email.
	IsContactBound bool `json:"is_contact_bound"`
}

type PublicCalendarEventType struct {
	Key                 string              `json:"key"`
	Name                string              `json:"name"`
	Description         string              `json:"description,omitempty"`
	DurationMinutes     int                 `json:"duration_minutes"`
	Location            PublicEventLocation `json:"location"`
	IsCancelAllowed     bool                `json:"is_cancel_allowed"`
	IsRescheduleAllowed bool                `json:"is_reschedule_allowed"`
	CancelNoticeMinutes int                 `json:"cancel_notice_minutes"`
}

// PublicEventLocation is the location as shown BEFORE booking: how the
// meeting happens, not where exactly (join links / phone numbers travel with
// the booked event).
type PublicEventLocation struct {
	Kind         schedulingmodel.LocationKind `json:"kind"`
	Address      string                       `json:"address,omitempty"`
	Instructions string                       `json:"instructions,omitempty"`
}

type PublicHost struct {
	UserId      string `json:"user_id"`
	DisplayName string `json:"display_name"`
	Headline    string `json:"headline,omitempty"`
}

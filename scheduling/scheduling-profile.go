package schedulingmodel

import (
	"time"

	"go.proteos.ai/model/common"
	conversationmodel "go.proteos.ai/model/conversation"
)

// SchedulingProfile is a user's scheduling settings: the zone their hours are
// expressed in, the weekly hours themselves, what a scheduling page shows
// about them, and whether they can be booked at all. One per (org, user).
//
// WeeklyHours is the base availability (the same {day, from, until} windows
// conversation-service's sending rules use), read in Timezone. Empty = no
// weekly hours = bookable any time (within the engine's 05:00–22:00 sane
// hours); availability constraints then add openings and blocks on top.
type SchedulingProfile struct {
	Id          string                        `json:"id"`
	OrgId       string                        `json:"org_id"`
	UserId      string                        `json:"user_id"`
	Timezone    string                        `json:"timezone"`
	WeeklyHours []conversationmodel.WindowDay `json:"weekly_hours"`
	DisplayName string                        `json:"display_name,omitempty"`
	Headline    string                        `json:"headline,omitempty"`
	IsBookable  bool                          `json:"is_bookable"`
	CreatedAt   time.Time                     `json:"created_at" sortable:""`
	CreatedBy   common.UserRef                `json:"created_by"`
	UpdatedAt   time.Time                     `json:"updated_at" sortable:""`
	UpdatedBy   common.UserRef                `json:"updated_by"`
}

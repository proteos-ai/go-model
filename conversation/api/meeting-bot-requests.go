package conversationapi

import (
	"time"

	"go.proteos.ai/model/common"
	conversationmodel "go.proteos.ai/model/conversation"
)

// GetManyMeetingBotsQuery lists the org's Ava bots the caller may see. The
// start window filters by the meeting's start; CalendarEventIcalUid names one
// calendar meeting (pair with the window to narrow to an instance).
type GetManyMeetingBotsQuery struct {
	StartAfter           *time.Time `json:"start_after" form:"start_after"`
	StartBefore          *time.Time `json:"start_before" form:"start_before"`
	CalendarEventIcalUid *string    `json:"calendar_event_ical_uid" form:"calendar_event_ical_uid" db:"calendar_event_ical_uid"`
	Status               *string    `json:"status" form:"status" db:"status"`
	common.Pagination
	common.Sorting
}

type GetManyMeetingBotsResponse struct {
	Data []conversationmodel.MeetingBot `json:"data"`
	Meta common.ResponseMeta            `json:"meta"`
}

package conversationapi

import (
	"time"

	"go.proteos.ai/model/common"
	conversationmodel "go.proteos.ai/model/conversation"
)

// CreateConversationBriefRequest prepares guidance for a conversation, with
// zero or more contacts. Any number of briefs may be open for the same people
// at once. TypeKey must name an existing ConversationType when set; Channel
// must be a known channel when set.
type CreateConversationBriefRequest struct {
	Subject    string                    `json:"subject"`
	ContactIds []string                  `json:"contact_ids"`
	Channel    conversationmodel.Channel `json:"channel"`
	TypeKey    string                    `json:"type_key"`
	Content    string                    `json:"content" validate:"required"`
	// CalendarEventId ties the brief to the caller's view of a calendar event
	// (the caller must be able to read it); the meeting-bot start attaches it
	// to that meeting's conversation.
	CalendarEventId string `json:"calendar_event_id"`
	// IsSharedWithParticipants makes the brief every invitee's (requires
	// CalendarEventId).
	IsSharedWithParticipants bool `json:"is_shared_with_participants"`
}

// UpdateConversationBriefRequest revises an open brief (PATCH semantics —
// absent fields keep their value; ContactIds replaces the set wholesale when
// present). Only a status=prepared brief is editable (409
// conversation_brief_not_editable otherwise).
type UpdateConversationBriefRequest struct {
	Subject    *string                    `json:"subject,omitempty"`
	ContactIds *[]string                  `json:"contact_ids,omitempty"`
	Channel    *conversationmodel.Channel `json:"channel,omitempty"`
	Content    *string                    `json:"content,omitempty"`
	TypeKey    *string                    `json:"type_key,omitempty"`
	// CalendarEventId re-ties the brief ("" unties it — and drops the sharing).
	CalendarEventId          *string `json:"calendar_event_id,omitempty"`
	IsSharedWithParticipants *bool   `json:"is_shared_with_participants,omitempty"`
}

// AttachConversationBriefRequest binds an open brief to a conversation
// explicitly (the after-the-fact door; the softphone attaches through the
// call-start callback instead).
type AttachConversationBriefRequest struct {
	ConversationId string `json:"conversation_id" validate:"required"`
}

// GetManyConversationBriefsQuery filters the org's briefs. ContactId = briefs
// naming that person (containment); ConversationId = the brief a
// conversation was held against.
type GetManyConversationBriefsQuery struct {
	ContactId      *string `json:"contact_id" form:"contact_id"`
	ConversationId *string `json:"conversation_id" form:"conversation_id" db:"conversation_id"`
	// CalendarEventId = the briefs for one person's copy of an event.
	CalendarEventId *string `json:"calendar_event_id" form:"calendar_event_id" db:"calendar_event_id"`
	// CalendarEventIcalUid (+ CalendarEventStartAt to pin one instance of a
	// series) = every visible brief of a meeting, across invitees' copies.
	CalendarEventIcalUid *string    `json:"calendar_event_ical_uid" form:"calendar_event_ical_uid" db:"calendar_event_ical_uid"`
	CalendarEventStartAt *time.Time `json:"calendar_event_start_at" form:"calendar_event_start_at"`
	// CalendarEventStartAfter / CalendarEventStartBefore = the briefs of the
	// meetings starting in a window (a calendar grid's visible range):
	// start_at >= after AND start_at < before.
	CalendarEventStartAfter  *time.Time `json:"calendar_event_start_after" form:"calendar_event_start_after"`
	CalendarEventStartBefore *time.Time `json:"calendar_event_start_before" form:"calendar_event_start_before"`
	Channel                  *string    `json:"channel" form:"channel" db:"channel"`
	TypeKey                  *string    `json:"type_key" form:"type_key" db:"type_key"`
	Status                   *string    `json:"status" form:"status" db:"status"`
	common.Pagination
	common.Sorting
}

type GetManyConversationBriefsResponse struct {
	Meta common.ResponseMeta                   `json:"meta"`
	Data []conversationmodel.ConversationBrief `json:"data"`
}

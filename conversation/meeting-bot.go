package conversationmodel

import (
	"time"

	"go.proteos.ai/model/common"
)

// MeetingBot is the platform's memory of one Ava notetaker bot: scheduled on a
// calendar event (or dispatched adhoc), then joining, recording and ended —
// or failed, or removed by a person. One row per provider bot.
//
// A calendar meeting is named by its INSTANCE key (CalendarEventIcalUid +
// StartAt — Google's UID names the series, Microsoft's one occurrence), which
// is shared by every invitee's calendar. With bot deduplication on there is
// one bot per instance per org; with it off each connected invitee's calendar
// schedules its own, so several rows can share a key.
type MeetingBot struct {
	Id           string `json:"id"`
	OrgId        string `json:"org_id"`
	ConnectionId string `json:"connection_id"`
	// ExternalBotId is the provider's bot id; it becomes the meeting
	// conversation's ExternalConversationId once the meeting begins.
	ExternalBotId string `json:"external_bot_id"`
	// ProviderEventId is the provider's calendar-event id the bot was scheduled
	// on ("" = an adhoc bot, dispatched by URL).
	ProviderEventId      string           `json:"provider_event_id,omitempty"`
	CalendarEventIcalUid string           `json:"calendar_event_ical_uid,omitempty"`
	StartAt              *time.Time       `json:"start_at,omitempty"`
	MeetingUrl           string           `json:"meeting_url,omitempty"`
	Title                string           `json:"title,omitempty"`
	Status               MeetingBotStatus `json:"status" sortable:""`
	// InviteeUserIds are the platform users among the meeting's invitees —
	// they see the bot (and the conversation it creates) without owning the
	// connection. Serializes as an array, never null.
	InviteeUserIds []string `json:"invitee_user_ids"`
	// ConversationId is the meeting conversation, "" until the bot starts.
	ConversationId string         `json:"conversation_id,omitempty"`
	CreatedAt      time.Time      `json:"created_at" sortable:""`
	CreatedBy      common.UserRef `json:"created_by"`
	UpdatedAt      time.Time      `json:"updated_at" sortable:""`
	UpdatedBy      common.UserRef `json:"updated_by"`
}

// MeetingBotStatusFromBotEvent maps a provider bot lifecycle code to the
// MeetingBot status it implies; isKnown=false for codes that carry no status
// change (waiting rooms, unknown codes). call_ended and done both read as
// ended — call_ended precedes done, and the status write is monotonic.
func MeetingBotStatusFromBotEvent(code string) (status MeetingBotStatus, isKnown bool) {
	switch code {
	case "joining_call":
		return MeetingBotStatusJoining, true
	case "in_call_recording":
		return MeetingBotStatusRecording, true
	case "call_ended", "done":
		return MeetingBotStatusEnded, true
	case "fatal":
		return MeetingBotStatusFailed, true
	}
	return "", false
}

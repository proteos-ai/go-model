package conversationmodel

import (
	"time"

	"go.proteos.ai/model/common"
)

// ConversationBrief is guidance prepared AHEAD of a conversation — "talk to
// these people about this": the conversation sibling of a draft message. An
// agent (or a person) writes what the conversation should cover, a human
// opens the brief (a phone brief opens the softphone, which shows it before
// and during the call), and the brief binds to the conversation it was used
// on so the conversation can later be judged against it.
//
// Binding is always EXPLICIT — the human's button press carries the brief id
// into the call (the softphone hands it to the provider, the call-start
// callback attaches it), an attach call names the conversation, or the brief
// names its calendar event and the meeting-bot start attaches it to that
// meeting's conversation. Nothing matches briefs to conversations by contact
// or number, and any number of briefs may be open for the same people,
// channel, type and meeting at once: a brief is a prepared task, not a rule.
// A conversation may carry several briefs (every invitee of a meeting can
// bring their own); each brief's conversation_id is the link. Records that
// want a brief (an activity) store its id; every UI surface renders from the
// id.
//
// A brief on a calendar event is for ONE person's copy of the meeting by
// default — the event owner reads it, the other invitees do not; one brief
// for everybody is the explicit IsSharedWithParticipants.
//
// v1 is text only. Structured topics (and the adherence rating computed from
// the transcript) are the planned extension — they will land as additional
// fields, never as a replacement of Content.
type ConversationBrief struct {
	Id    string `json:"id"`
	OrgId string `json:"org_id"`
	// Subject is the brief's human label ("Renewal call with Dana") — the
	// headline wherever the brief is listed; empty allowed.
	Subject string `json:"subject" sortable:""`
	// ContactIds are the people the conversation is with (deduplicated; may
	// be empty). Serializes as an array, never null.
	ContactIds []string `json:"contact_ids"`
	// Channel says how the conversation should happen (phone opens the
	// softphone); "" = unspecified — the brief is read-only guidance.
	Channel Channel `json:"channel,omitempty" sortable:""`
	// TypeKey is the ConversationType the conversation is meant to be; "" =
	// untyped. Stamped onto the conversation when the brief attaches.
	TypeKey string `json:"type_key,omitempty"`
	// CalendarEventId is the calendar event (one person's mirror row) the
	// brief prepares; "" = not tied to a meeting. The meeting-bot start
	// attaches the brief to the meeting's conversation.
	CalendarEventId string `json:"calendar_event_id,omitempty"`
	// IsSharedWithParticipants makes the brief every invitee's, not only the event
	// owner's. Requires CalendarEventId.
	IsSharedWithParticipants bool `json:"is_shared_with_participants"`
	// CalendarEventIcalUid + CalendarEventStartAt are the meeting instance the
	// event belongs to (the key every invitee's copy shares), resolved from
	// the event on write. Read-only.
	CalendarEventIcalUid string     `json:"calendar_event_ical_uid,omitempty"`
	CalendarEventStartAt *time.Time `json:"calendar_event_start_at,omitempty"`
	// InviteeUserIds are the platform users the brief is for — the event
	// owner, plus every resolved invitee when shared. They read and edit the
	// brief next to its creator, before and after the meeting: a brief on an
	// event is never widened to whoever can see the meeting conversation.
	// Read-only; serializes as an array, never null.
	InviteeUserIds []string `json:"invitee_user_ids"`
	// ConversationId is the conversation the brief was used on; "" while
	// prepared or discarded.
	ConversationId string                  `json:"conversation_id,omitempty"`
	Status         ConversationBriefStatus `json:"status" sortable:""`
	// Content is the guidance text (markdown allowed).
	Content string `json:"content"`
	// UsedAt stamps the moment the brief attached to a conversation.
	UsedAt    *time.Time     `json:"used_at,omitempty"`
	CreatedAt time.Time      `json:"created_at" sortable:""`
	CreatedBy common.UserRef `json:"created_by"`
	UpdatedAt time.Time      `json:"updated_at" sortable:""`
	UpdatedBy common.UserRef `json:"updated_by"`
}

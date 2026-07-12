package conversationapi

import (
	"time"
)

// DispatchMeetingBotRequest sends a meeting bot (Ava) into a meeting through
// an active meeting connection (adhoc-meeting for URL dispatch; the calendar
// members auto-schedule instead and never hit this endpoint). JoinAt schedules
// the bot ahead of time (a scheduled bot joins reliably on time; omit to join
// now). Title seeds the meeting conversation's subject.
type DispatchMeetingBotRequest struct {
	ConnectionId string     `json:"connection_id" validate:"required"`
	MeetingUrl   string     `json:"meeting_url" validate:"required,url"`
	Title        string     `json:"title"`
	JoinAt       *time.Time `json:"join_at,omitempty"`
	// Language is the BCP-47 spoken-language hint for transcription (e.g. "en",
	// "de"). Empty means auto-detect: on the deepgram_streaming provider the
	// bot transcribes with language "multi" (Deepgram's multilingual detection
	// + code-switching). Only the deepgram_streaming provider consumes this;
	// meeting_captions takes the platform's own caption language.
	Language string `json:"language"`
}

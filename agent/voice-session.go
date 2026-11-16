package agentmodel

import (
	"encoding/json"
	"fmt"
	"time"

	"go.proteos.ai/model/common"
	"golang.org/x/exp/slices"
)

// VoiceSession records that a person talked to a session's agent through a voice
// assistant: one row per time voice was entered. The voice assistant listens and
// speaks; what was said is in the session's event log (user.speech, agent.speech),
// and everything the person asked for runs as turns of the session. A session has
// at most one active voice session.
//
// The voice provider's own session id is not here: it lives in the provider
// adapter's store, like a session's provider id.
type VoiceSession struct {
	Id    string `json:"id"`
	OrgId string `json:"org_id"`
	// SessionId is the agent session the person talks to.
	SessionId string             `json:"session_id"`
	Status    VoiceSessionStatus `json:"status"`
	// EndReason is why it ended; nil while it is active.
	EndReason *VoiceEndReason `json:"end_reason"`
	// DurationSeconds is how long the voice provider says the session lasted, as
	// the page reported it or as this server read it when it closed the session
	// itself. A page that ended the session without it (it went away) leaves the
	// time since voice was entered. Nil when this server ended the session and
	// the provider no longer says. The provider's usage stays the billing truth.
	DurationSeconds *int64         `json:"duration_seconds"`
	CreatedBy       common.UserRef `json:"created_by"`
	CreatedAt       time.Time      `json:"created_at"`
	EndedAt         *time.Time     `json:"ended_at"`
}

// VoiceSessionStatus is where a voice session stands.
type VoiceSessionStatus string

const (
	VoiceSessionStatusActive VoiceSessionStatus = "active"
	VoiceSessionStatusEnded  VoiceSessionStatus = "ended"
)

// VoiceEndReason is why a voice session ended.
type VoiceEndReason string

const (
	// VoiceEndReasonUser: the person left voice.
	VoiceEndReasonUser VoiceEndReason = "user"
	// VoiceEndReasonIdle: nobody spoke for the idle limit.
	VoiceEndReasonIdle VoiceEndReason = "idle"
	// VoiceEndReasonMaxDuration: the session reached its longest allowed length.
	VoiceEndReasonMaxDuration VoiceEndReason = "max_duration"
	// VoiceEndReasonSessionEnded: the agent session it belonged to was ended.
	VoiceEndReasonSessionEnded VoiceEndReason = "session_ended"
	// VoiceEndReasonBrowserLeft: the page went away (closed, or its connection
	// to the voice provider was lost).
	VoiceEndReasonBrowserLeft VoiceEndReason = "browser_left"
	// VoiceEndReasonProviderClosed: the voice provider closed the session.
	VoiceEndReasonProviderClosed VoiceEndReason = "provider_closed"
	// VoiceEndReasonSuperseded: voice was entered again while this session was
	// still active, and this server closed it. Never reported by a page.
	VoiceEndReasonSuperseded VoiceEndReason = "superseded"
)

var VoiceEndReasons = []VoiceEndReason{
	VoiceEndReasonUser,
	VoiceEndReasonIdle,
	VoiceEndReasonMaxDuration,
	VoiceEndReasonSessionEnded,
	VoiceEndReasonBrowserLeft,
	VoiceEndReasonProviderClosed,
	VoiceEndReasonSuperseded,
}

func (VoiceEndReason) Enum() []interface{} {
	enums := make([]interface{}, 0, len(VoiceEndReasons))
	for _, element := range VoiceEndReasons {
		enums = append(enums, element)
	}
	return enums
}

func (reason *VoiceEndReason) UnmarshalJSON(byteArray []byte) error {
	type _VoiceEndReason VoiceEndReason
	value := (*_VoiceEndReason)(reason)
	if err := json.Unmarshal(byteArray, value); err != nil {
		return err
	}
	if slices.Contains(VoiceEndReasons, *reason) {
		return nil
	}
	return fmt.Errorf("invalid voice end reason: %s", *reason)
}

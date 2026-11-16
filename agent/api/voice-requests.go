package agentapi

import agentmodel "go.proteos.ai/model/agent"

// StartVoiceRequest enters voice on a session: the page opens a WebRTC
// connection to the voice provider, and this server creates the provider's
// session from the page's offer so the page never holds the provider's key.
type StartVoiceRequest struct {
	// SdpOffer is the WebRTC offer of the page's peer connection.
	SdpOffer string `json:"sdp_offer"`
}

// StartVoiceResponse is the voice session that started, the provider's answer to
// the page's offer, and the limits the page enforces by closing the session.
type StartVoiceResponse struct {
	VoiceSession agentmodel.VoiceSession `json:"voice_session"`
	SdpAnswer    string                  `json:"sdp_answer"`
	// IdleTimeoutSeconds is how long nobody may speak before the page leaves voice.
	IdleTimeoutSeconds int `json:"idle_timeout_seconds"`
	// MaxDurationSeconds is how long a voice session may last at most.
	MaxDurationSeconds int `json:"max_duration_seconds"`
}

// SpeechSpeaker is who said an utterance of a voice session.
type SpeechSpeaker string

const (
	// SpeechSpeakerUser is the person: stored as a user.speech event.
	SpeechSpeakerUser SpeechSpeaker = "user"
	// SpeechSpeakerAgent is the voice assistant: stored as an agent.speech event.
	SpeechSpeakerAgent SpeechSpeaker = "agent"
)

// SpeechSegment is one finished utterance as the page heard it from the voice
// provider: the fields of agentmodel.SpeechPayload, and who spoke.
type SpeechSegment struct {
	Speaker SpeechSpeaker `json:"speaker"`
	Text    string        `json:"text"`
	StartMs int64         `json:"start_ms"`
	EndMs   int64         `json:"end_ms"`
}

// AppendSpeechRequest reports what was said in the session's active voice
// session since the last report, in the order said. A report may repeat what an
// earlier one held (the page sends again what it has no answer for): an
// utterance already in the log is not written twice.
type AppendSpeechRequest struct {
	// VoiceSessionId is the voice session the page reports for: the one it
	// started. A report for one that was replaced since is refused.
	VoiceSessionId string          `json:"voice_session_id"`
	Segments       []SpeechSegment `json:"segments"`
	// StartsTurn says the voice assistant handed a request over: after the
	// segments, a user.voice_request starts a turn of the session as the caller,
	// or is queued behind the turn that runs.
	StartsTurn bool `json:"starts_turn"`
	// DelegationId is the voice provider's delegation the request is of, with
	// StartsTurn only. A report that is sent again (its answer was lost) names
	// the same delegation, and the request is recorded once.
	DelegationId string `json:"delegation_id,omitempty"`
	// ClientTools is the page's client-tool catalog, honoured with StartsTurn
	// only (see AppendEventRequest.ClientTools).
	ClientTools []agentmodel.ClientToolSpec `json:"client_tools,omitempty"`
	// Context is model-only text for the request, with StartsTurn only: what a
	// typed message carries as its context (the application context of a
	// session whose first turn is a spoken one).
	Context string `json:"context,omitempty"`
}

// AppendSpeechResponse is the events the report wrote to the session's log.
type AppendSpeechResponse struct {
	Data []agentmodel.SessionEvent `json:"data"`
}

// EndVoiceRequest leaves voice: why, and how long the voice provider says the
// session lasted (the provider reports it to the page when the session closes).
type EndVoiceRequest struct {
	// VoiceSessionId is the voice session the page leaves: the one it started.
	VoiceSessionId  string                    `json:"voice_session_id"`
	EndReason       agentmodel.VoiceEndReason `json:"end_reason"`
	DurationSeconds *int64                    `json:"duration_seconds"`
	// Segments is what was said since the page's last speech report. A page that
	// is being closed has one request left to send: this one carries both.
	Segments []SpeechSegment `json:"segments,omitempty"`
}

package conversationmodel

import (
	"time"

	"go.proteos.ai/model/common"
)

// Transcription is a standalone batch-transcription artifact: a storage-service
// audio file turned into full text plus diarized speaker turns. Channel-agnostic
// — it becomes conversational only when materialized (POST
// /transcriptions/:id/materialize), which creates an adhoc/meeting Conversation
// with one Message per turn and links it back via ConversationId.
type Transcription struct {
	Id    string `json:"id"`
	OrgId string `json:"org_id"`
	// SourceFileId is the storage-service file the transcription was created from.
	// AudioFileId/TranscriptFileId are optional derived artifacts (a normalized
	// audio copy, the rendered transcript document).
	SourceFileId     string              `json:"source_file_id"`
	AudioFileId      string              `json:"audio_file_id"`
	TranscriptFileId string              `json:"transcript_file_id"`
	Status           TranscriptionStatus `json:"status" sortable:""`
	Language         string              `json:"language"`
	DurationSeconds  float64             `json:"duration_seconds"`
	// Model is the provider model that ran (e.g. nova-3), recorded for
	// reproducibility.
	Model      string `json:"model"`
	IsDiarized bool   `json:"is_diarized"`
	// The flat transcript text is NOT stored inline — it lives in the
	// TranscriptFileId artifact (Deepgram's smart-formatted rendering) and is
	// fully reconstructable from Turns, the canonical structured form.
	Turns        []TranscriptTurn `json:"turns"`
	SpeakerCount int              `json:"speaker_count"`
	// ProviderRequestId is Deepgram's request id, for support/debugging.
	ProviderRequestId string `json:"provider_request_id"`
	// ConversationId is set once the transcription is materialized/linked.
	ConversationId string         `json:"conversation_id"`
	Error          string         `json:"error"`
	CreatedAt      time.Time      `json:"created_at" sortable:""`
	CreatedBy      common.UserRef `json:"created_by"`
	UpdatedAt      time.Time      `json:"updated_at" sortable:""`
	UpdatedBy      common.UserRef `json:"updated_by"`
}

// TranscriptTurn is one diarized utterance: who (speaker index + label), what,
// and when (milliseconds from audio start). Confidence is the provider's
// word-averaged confidence for the utterance.
type TranscriptTurn struct {
	Speaker      int     `json:"speaker"`
	SpeakerLabel string  `json:"speaker_label"`
	Text         string  `json:"text"`
	StartMs      int64   `json:"start_ms"`
	EndMs        int64   `json:"end_ms"`
	Confidence   float64 `json:"confidence"`
}

package conversationapi

import (
	"go.proteos.ai/model/common"
	conversationmodel "go.proteos.ai/model/conversation"
)

// CreateTranscriptionRequest creates a transcription one of two ways, and
// EXACTLY one of FileId / Turns must be set (logic.ValidateTranscriptionCreateMode
// enforces it — struct tags cannot express an either-or).
//
//   - FileId  — transcribe a storage-service audio file through the provider.
//     Asynchronous: the row lands `processing` and the caller polls.
//   - Turns   — an already-written transcript, supplied verbatim. No provider
//     runs; the row lands `completed` synchronously. This is how a transcript
//     (and, via materialize, a conversation) is authored from text with nothing
//     to upload.
//
// IsDiarized defaults to true on the file path (the service applies the default
// when the pointer is nil — exactOptionalPropertyTypes-style tri-state); on the
// turns path it is derived from the turns unless given explicitly.
type CreateTranscriptionRequest struct {
	FileId string `json:"file_id"`
	// Turns is the authored transcript. Timings are optional — text has none,
	// and materialization then sequences turns by index instead.
	Turns      []conversationmodel.TranscriptTurn `json:"turns,omitempty"`
	Language   string                             `json:"language"`
	Model      string                             `json:"model"`
	IsDiarized *bool                              `json:"is_diarized,omitempty"`
}

// MaterializeTranscriptionRequest turns a completed transcription into a
// conversation: Channel defaults to adhoc (meeting allowed — both are
// connector-less); SpeakerNames maps diarized speaker indexes to display names
// (e.g. {"0": "Tonio", "1": "Dirk"}); unmapped speakers keep "Speaker N".
type MaterializeTranscriptionRequest struct {
	Channel      conversationmodel.Channel `json:"channel"`
	Subject      string                    `json:"subject"`
	SpeakerNames map[string]string         `json:"speaker_names"`
}

// UpdateTranscriptionRequest edits a COMPLETED transcription in place. Turns,
// when present, replace the diarized turns wholesale (the canonical structured
// form — the flat-text artifact is regenerated from them). SpeakerLabels maps
// diarized speaker indexes to display labels (e.g. {"0": "Tonio"}) and is
// applied across all turns after any replacement. Editing does NOT retro-update
// messages of an already-materialized conversation.
type UpdateTranscriptionRequest struct {
	Turns         []conversationmodel.TranscriptTurn `json:"turns,omitempty"`
	SpeakerLabels map[string]string                  `json:"speaker_labels,omitempty"`
	Language      *string                            `json:"language,omitempty"`
}

type GetManyTranscriptionsQuery struct {
	Status *string `json:"status" form:"status" db:"status"`
	// ConversationId lists a conversation's transcription artifacts (the
	// Conversation doc's "fetch them with GET /transcriptions?conversation_id=").
	ConversationId *string `json:"conversation_id" form:"conversation_id" db:"conversation_id"`
	// ProviderRequestId narrows to the provider-side transcript identity — the
	// import-idempotency lookup (a webhook redelivery finds its prior import).
	ProviderRequestId *string `json:"provider_request_id" form:"provider_request_id" db:"provider_request_id"`
	common.Pagination
	common.Sorting
}

type GetManyTranscriptionsResponse struct {
	Meta common.ResponseMeta               `json:"meta"`
	Data []conversationmodel.Transcription `json:"data"`
}

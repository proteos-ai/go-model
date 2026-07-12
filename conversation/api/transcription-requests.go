package conversationapi

import (
	"go.proteos.ai/model/common"
	conversationmodel "go.proteos.ai/model/conversation"
)

// CreateTranscriptionRequest transcribes a storage-service file synchronously.
// IsDiarized defaults to true (the service applies the default when the pointer
// is nil — exactOptionalPropertyTypes-style tri-state).
type CreateTranscriptionRequest struct {
	FileId     string `json:"file_id" validate:"required"`
	Language   string `json:"language"`
	Model      string `json:"model"`
	IsDiarized *bool  `json:"is_diarized,omitempty"`
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

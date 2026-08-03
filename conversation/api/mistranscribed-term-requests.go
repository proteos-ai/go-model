package conversationapi

import (
	"go.proteos.ai/model/common"
	conversationmodel "go.proteos.ai/model/conversation"
)

// AcceptMistranscribedTermRequest resolves a proposed finding by applying a
// replacement to the transcript. Replacement overrides the reviewer's
// suggestion when set (a human supplying the actually-correct word); empty
// falls back to the stored SuggestedReplacement. CreateGlossaryTerm, when
// present, promotes the applied replacement into the org glossary and links
// the finding to it — feeding both keyterm boosting and future review passes.
type AcceptMistranscribedTermRequest struct {
	Replacement        string                                `json:"replacement,omitempty" validate:"max=100"`
	CreateGlossaryTerm *PromoteMistranscriptionGlossaryTerm `json:"create_glossary_term,omitempty"`
}

// PromoteMistranscriptionGlossaryTerm carries the glossary fields for the
// promoted term — the term text itself is the applied replacement.
type PromoteMistranscriptionGlossaryTerm struct {
	Definition string `json:"definition,omitempty" validate:"max=4096"`
	Priority   *int   `json:"priority,omitempty" validate:"omitempty,gte=0"`
}

type GetManyMistranscribedTermsQuery struct {
	Status          *string `json:"status" form:"status" db:"status"`
	TranscriptionId *string `json:"transcription_id" form:"transcription_id" db:"transcription_id"`
	ConversationId  *string `json:"conversation_id" form:"conversation_id" db:"conversation_id"`
	common.Pagination
	common.Sorting
}

type GetManyMistranscribedTermsResponse struct {
	Meta common.ResponseMeta                    `json:"meta"`
	Data []conversationmodel.MistranscribedTerm `json:"data"`
}

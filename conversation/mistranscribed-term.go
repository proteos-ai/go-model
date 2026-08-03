package conversationmodel

import (
	"time"

	"go.proteos.ai/model/common"
)

// MistranscribedTerm is one word/phrase the post-transcription review pass
// judged likely misheard by the transcription provider ("flurbo" where the
// speaker said "Fibaro"). High-confidence findings are applied automatically
// (status auto_replaced); the rest are proposed for human review. Resolved
// rows persist as the learning corpus: an accepted term that references a
// glossary term (GlossaryTermId) becomes a known "misheard → glossary term"
// mapping the next review pass replaces deterministically, without the LLM.
type MistranscribedTerm struct {
	Id    string `json:"id"`
	OrgId string `json:"org_id"`
	// TranscriptionId is the transcript the term was found in; ConversationId
	// is the transcription's linked conversation (denormalized for listing
	// findings under a conversation), empty until/unless linked.
	TranscriptionId string `json:"transcription_id"`
	ConversationId  string `json:"conversation_id,omitempty"`
	// Term is the text as transcribed; ContextSnippet is the surrounding
	// utterance text the reviewer judged it in; TurnIndexes are the positions
	// (indexes into Transcription.Turns) where it occurs.
	Term           string `json:"term" sortable:""`
	ContextSnippet string `json:"context_snippet,omitempty"`
	TurnIndexes    []int  `json:"turn_indexes"`
	// SuggestedReplacement is what the reviewer thinks was actually said —
	// empty when the term is flagged as suspicious with no candidate.
	// SuggestionSource says whether it came from the org glossary or the
	// model's own judgement.
	SuggestedReplacement string                             `json:"suggested_replacement,omitempty"`
	SuggestionSource     MistranscriptionSuggestionSource   `json:"suggestion_source,omitempty"`
	// GlossaryTermId references the glossary term the suggestion matched — or,
	// after an accept that promoted the replacement into the glossary, the term
	// this mistranscription became. It is what feeds future passes.
	GlossaryTermId string `json:"glossary_term_id,omitempty"`
	// MisheardConfidence is the reviewer's confidence (0..1) that the term IS
	// mistranscribed; ReplacementConfidence that the suggestion fits. Both must
	// reach the auto-replace threshold for an automatic fix.
	MisheardConfidence    float64 `json:"misheard_confidence"`
	ReplacementConfidence float64 `json:"replacement_confidence"`
	Status                MistranscribedTermStatus `json:"status" sortable:""`
	// AppliedReplacement is the text actually written into the transcript (an
	// auto-replace or an accept — which may override the suggestion).
	AppliedReplacement string          `json:"applied_replacement,omitempty"`
	ResolvedBy         *common.UserRef `json:"resolved_by,omitempty"`
	ResolvedAt         *time.Time      `json:"resolved_at,omitempty"`
	CreatedAt          time.Time       `json:"created_at" sortable:""`
	CreatedBy          common.UserRef  `json:"created_by"`
	UpdatedAt          time.Time       `json:"updated_at" sortable:""`
	UpdatedBy          common.UserRef  `json:"updated_by"`
}

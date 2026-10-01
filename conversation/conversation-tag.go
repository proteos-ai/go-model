package conversationmodel

import (
	"encoding/json"
	"time"

	"go.proteos.ai/model/common"
)

// ConversationTag is one observation a definition made on a conversation.
//
// Anchors encode the definition's scope: conversation scope = every anchor
// field empty; message scope = StartMessageId == EndMessageId with nil
// offsets; span = offsets set (a span may cross messages, then Start/End
// message ids differ). Offsets are UNICODE CODE POINT indexes into the
// message's first text content block, end exclusive — clients working in
// UTF-16 must convert. Quote (and EndQuote for a span that continues into a
// later message) is the verbatim text at the anchor: the reason evaluator
// emits quotes because a model cannot count characters, the service resolves
// them to offsets, and afterwards the quote is the evidence that survives a
// message row change (the final transcript swaps realtime message ids; a term
// correction edits text) — a tag whose anchors no longer resolve still shows
// its quote and can be re-anchored by search.
//
// Value is normalized per question type: "true" for a boolean (false is never
// stored), the option key for a choice, the 0-based level index for a score.
// Status pending marks a provisional realtime observation on a live
// spoken-medium conversation; a completion evaluation replaces it. Dismissed
// tags are feedback: a later evaluation never re-creates a tag matching a
// dismissed one on (definition_key, value, start/end message).
type ConversationTag struct {
	Id             string `json:"id"`
	OrgId          string `json:"org_id"`
	ConversationId string `json:"conversation_id"`
	// Channel is denormalized from the conversation so org-wide lists and
	// counts group by channel without a join.
	Channel       Channel `json:"channel" sortable:""`
	DefinitionKey string  `json:"definition_key" sortable:""`
	Value         string  `json:"value" sortable:""`
	Confidence    float64 `json:"confidence" sortable:""`
	// Probabilities is the raw decide distribution ({"true":p,"false":1-p},
	// option → p, or {"levels":[…],"mean":x}); absent for reason and user tags.
	Probabilities json.RawMessage `json:"probabilities,omitempty"`
	// Why is the reason evaluator's one-line rationale; empty for decide.
	Why    string                `json:"why,omitempty"`
	Source ConversationTagSource `json:"source"`
	Status ConversationTagStatus `json:"status" sortable:""`
	// EvaluationId links a model tag to the evaluation that produced it.
	EvaluationId    string     `json:"evaluation_id,omitempty"`
	StartMessageId  string     `json:"start_message_id,omitempty"`
	StartOffset     *int       `json:"start_offset,omitempty"`
	EndMessageId    string     `json:"end_message_id,omitempty"`
	EndOffset       *int       `json:"end_offset,omitempty"`
	StartOccurredAt *time.Time `json:"start_occurred_at,omitempty"`
	EndOccurredAt   *time.Time `json:"end_occurred_at,omitempty"`
	Quote           string     `json:"quote,omitempty"`
	EndQuote        string     `json:"end_quote,omitempty"`
	// ObservedAt is when the observed thing happened (the end message's
	// occurred_at, else the last message's) — the time-series axis.
	ObservedAt time.Time      `json:"observed_at" sortable:""`
	CreatedAt  time.Time      `json:"created_at" sortable:""`
	CreatedBy  common.UserRef `json:"created_by"`
	UpdatedAt  time.Time      `json:"updated_at" sortable:""`
	UpdatedBy  common.UserRef `json:"updated_by"`
}

// ConversationTagCount is one row of GET /conversation-tags/counts: the
// grouped columns that were requested plus the count.
type ConversationTagCount struct {
	DefinitionKey string  `json:"definition_key,omitempty"`
	Value         string  `json:"value,omitempty"`
	Channel       Channel `json:"channel,omitempty"`
	Count         int     `json:"count"`
}

// ConversationTaggedEvent is the conversation.tagged payload: the conversation,
// the evaluation that landed, and every tag it produced (status pending or
// active — the event header carries the status so a consumer can filter
// provisional realtime observations without decoding).
type ConversationTaggedEvent struct {
	Conversation Conversation              `json:"conversation"`
	Evaluation   ConversationTagEvaluation `json:"evaluation"`
	Tags         []ConversationTag         `json:"tags"`
}

// IsConversationScoped reports whether the tag carries no anchor at all.
func (tag ConversationTag) IsConversationScoped() bool {
	return tag.StartMessageId == "" && tag.EndMessageId == ""
}

// IsSpan reports whether the tag anchors to character offsets.
func (tag ConversationTag) IsSpan() bool {
	return tag.StartOffset != nil || tag.EndOffset != nil
}

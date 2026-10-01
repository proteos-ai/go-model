package conversationmodel

import (
	"time"

	"go.proteos.ai/model/common"
)

// ConversationTagCriteria is populated per question type: boolean → Boolean
// (optional, both descriptions or neither), choice → Choice (option key →
// description, ≥ 2), score → Score (ordered level descriptions, ≥ 2; the tag
// value is the 0-based level index). Descriptions are handed to the evaluator
// verbatim, so they should describe the observable content, not internal
// process.
type ConversationTagCriteria struct {
	Boolean *BooleanTagCriteria `json:"boolean,omitempty"`
	Choice  map[string]string   `json:"choice,omitempty"`
	Score   []string            `json:"score,omitempty"`
}

// BooleanTagCriteria spells out what makes a boolean question hold (true) or
// not (false) when the boundary is subtle.
type BooleanTagCriteria struct {
	True  string `json:"true"`
	False string `json:"false"`
}

// ConversationTagQuestion is the decide-shaped question a definition asks
// over a conversation — the same three types as agentmodel.DecisionQuestion,
// with criteria in an explicit struct so manifests and SDKs need no tagged
// union decoder. Type is immutable after create: tag values are typed by it.
type ConversationTagQuestion struct {
	Type         ConversationTagQuestionType `json:"type"`
	Instructions string                      `json:"instructions"`
	Criteria     ConversationTagCriteria     `json:"criteria"`
}

// ConversationTagScreening makes a reason definition ask the cheap decide
// model first whether it applies anywhere in the judged window (its own
// boolean question, one answer per window); only a definition that passes
// (probability ≥ MinProbability) reaches the reasoning call, which is skipped
// entirely when nothing passes. Boolean questions only — choice and score
// always have an answer. MinProbability nil = the service default.
type ConversationTagScreening struct {
	IsEnabled      bool     `json:"is_enabled"`
	MinProbability *float64 `json:"min_probability,omitempty"`
}

// ConversationTagEvaluator selects the model path. ModelId is only meaningful
// for reason (empty = the service default); decide always uses the
// service-wide decision model. Screening is only meaningful for reason with
// a boolean question; nil = screening on with the service default threshold.
type ConversationTagEvaluator struct {
	Kind      ConversationTagEvaluatorKind `json:"kind"`
	ModelId   string                       `json:"model_id,omitempty"`
	Screening *ConversationTagScreening    `json:"screening,omitempty"`
}

// ConversationTagDefinition describes WHAT to look for in a conversation: one
// typed question plus where its answers anchor (Scope) and which model path
// answers it (Evaluator). WHERE and WHEN a definition is evaluated is a
// ConversationTagSet's business — a definition that no enabled set references
// is never evaluated. Keyed (org_id, key) and module-deployable
// (conversation-tag-definitions/<key>.json), like conversation types.
type ConversationTagDefinition struct {
	OrgId string `json:"org_id"`
	// Key is the immutable identity within the org — the value every tag
	// carries as definition_key.
	Key  string `json:"key" sortable:""`
	Name string `json:"name" sortable:""`
	// Group is a free grouping for UIs ("content", "sentiment", "signal").
	Group    string                  `json:"group,omitempty" sortable:""`
	Question ConversationTagQuestion `json:"question"`
	// Scope: conversation | message | span. span requires Evaluator.Kind =
	// reason (only a reasoning generation can localise verbatim quotes).
	Scope     ConversationTagScope     `json:"scope"`
	Evaluator ConversationTagEvaluator `json:"evaluator"`
	// Directions narrows message/span-scope judging to messages of these
	// directions (e.g. inbound only for customer sentiment); empty = all.
	// Ignored for conversation scope.
	Directions []MessageDirection `json:"directions"`
	// MinConfidence drops observations below it, for every question type — an
	// absent tag means "not determined", not "false".
	MinConfidence float64 `json:"min_confidence"`
	// Color is an optional UI hint (any CSS color token the client accepts).
	Color     string `json:"color,omitempty"`
	IsEnabled bool   `json:"is_enabled"`
	// ModuleSlug attributes the definition to the module that deployed it;
	// empty = not module-owned.
	ModuleSlug string         `json:"module_slug,omitempty"`
	CreatedAt  time.Time      `json:"created_at" sortable:""`
	CreatedBy  common.UserRef `json:"created_by"`
	UpdatedAt  time.Time      `json:"updated_at" sortable:""`
	UpdatedBy  common.UserRef `json:"updated_by"`
}

package conversationmodel

import (
	"time"

	"go.proteos.ai/model/common"
)

// ConversationTagUsage is the model spend of one evaluation, summed over its
// decide and generate calls.
type ConversationTagUsage struct {
	InputTokens  int64 `json:"input_tokens"`
	OutputTokens int64 `json:"output_tokens"`
	RequestCount int   `json:"request_count"`
}

// ConversationTagScreeningResult is one screening verdict of an evaluation:
// the decide probability that the definition applies in the judged window
// and whether that passed the threshold. Probability is nil (and Error set)
// when the screening could not be asked — it then passes, fail-open.
type ConversationTagScreeningResult struct {
	DefinitionKey string   `json:"definition_key"`
	Probability   *float64 `json:"probability,omitempty"`
	IsPassed      bool     `json:"is_passed"`
	Error         string   `json:"error,omitempty"`
}

// ConversationTagEvaluation records one time the tagger ran on a conversation:
// why (Trigger), what it judged (DefinitionKeys, the message window), with
// which models, what it cost, and how it ended. Model tags reference it
// through evaluation_id (provenance), the UI polls it after a manual
// evaluate (202), and the latest done realtime/completion evaluation's
// WindowEndOccurredAt is the watermark the next realtime window starts after.
//
// Window = the messages JUDGED by this evaluation. realtime: everything after
// the previous watermark (older messages ride along as context only);
// completion/manual: all messages. Lifecycle: manual creates pending on the
// 202 then flips to running; realtime/completion create running directly;
// ends done or failed.
type ConversationTagEvaluation struct {
	Id             string                          `json:"id"`
	OrgId          string                          `json:"org_id"`
	ConversationId string                          `json:"conversation_id"`
	Trigger        ConversationTagTrigger          `json:"trigger" sortable:""`
	Status         ConversationTagEvaluationStatus `json:"status" sortable:""`
	DefinitionKeys []string                        `json:"definition_keys"`
	ModelIds       []string                        `json:"model_ids"`
	Usage          ConversationTagUsage            `json:"usage"`
	// IsTruncated marks that the rendered state exceeded the cap and older
	// context (or, for a huge window, the newest window messages — deferred to
	// the next evaluation) was dropped.
	IsTruncated           bool       `json:"is_truncated"`
	WindowStartOccurredAt *time.Time `json:"window_start_occurred_at,omitempty"`
	WindowEndOccurredAt   *time.Time `json:"window_end_occurred_at,omitempty"`
	WindowMessageIds      []string   `json:"window_message_ids"`
	// Screenings lists every screened reason definition's verdict, so "why was
	// nothing tagged" is answerable and thresholds can be tuned.
	Screenings []ConversationTagScreeningResult `json:"screenings"`
	Error      string                           `json:"error,omitempty"`
	StartedAt  time.Time                        `json:"started_at" sortable:""`
	FinishedAt *time.Time                       `json:"finished_at,omitempty"`
	CreatedAt  time.Time                        `json:"created_at" sortable:""`
	CreatedBy  common.UserRef                   `json:"created_by"`
	UpdatedAt  time.Time                        `json:"updated_at"`
	UpdatedBy  common.UserRef                   `json:"updated_by"`
}

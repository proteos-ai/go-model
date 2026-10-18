package conversationmodel

import (
	"time"

	"go.proteos.ai/model/common"
)

// ConversationTypeConfig is the optional per-type behavior configuration.
// Stored as JSONB so later knobs (model overrides, auto-actions, …) need no
// DDL.
type ConversationTypeConfig struct {
	// SummaryPromptKey references an agent-service prompt by key. When a
	// conversation is classified as this type, the prompt's CURRENT version
	// body replaces the built-in summary system prompt (the instruction —
	// metadata, glossary, transcript — is unchanged). Empty ⇒ built-in prompt.
	SummaryPromptKey string `json:"summary_prompt_key,omitempty"`
}

// ConversationOutcome is one way a conversation of a type can end — "booked",
// "not_interested", "callback". The type's list is the outcome vocabulary a
// person picks from (the softphone after a call, the conversation header)
// and the choice a role=outcome tag definition detects among. Description is
// what the model reads to tell the outcomes apart (Name when empty).
type ConversationOutcome struct {
	Key         string `json:"key"`
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
}

// ConversationPhase is one stage a conversation of a type moves through —
// "intro", "discovery", "close". The type's list is ordered; progress
// through it is DERIVED from the objectives pinned to each phase (every
// measured objective of a phase achieved ⇒ the phase is done), never
// detected on its own. Description is for humans (the checklist, the
// summary prompt).
type ConversationPhase struct {
	Key         string `json:"key"`
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
}

// ConversationOutcomeOption is one outcome a conversation may take, as GET
// /conversations/:id/outcomes lists them: the type's outcomes (Source type)
// plus the channel's system outcomes (Source system — the ones the service
// sets by itself when it can tell, such as a call nobody answered, and a
// person may pick by hand).
type ConversationOutcomeOption struct {
	Key         string                          `json:"key"`
	Name        string                          `json:"name"`
	Description string                          `json:"description,omitempty"`
	Source      ConversationOutcomeOptionSource `json:"source"`
}

// ConversationOutcomeOptionSource says where an outcome option comes from.
type ConversationOutcomeOptionSource string

const (
	ConversationOutcomeOptionSourceType   ConversationOutcomeOptionSource = "type"
	ConversationOutcomeOptionSourceSystem ConversationOutcomeOptionSource = "system"
)

// ConversationType is a per-org conversation taxonomy entry. Before a meeting
// summary is generated, a cheap classifier model reads the transcript head
// plus every type's Definition and picks the matching Key (or none); the
// result lands on Conversation.TypeKey and may swap the summary system prompt
// via Config.SummaryPromptKey. Keyed (org_id, key) and module-deployable
// (conversation-types/<key>.json), like agent-service prompts.
type ConversationType struct {
	OrgId string `json:"org_id"`
	// Key is the immutable identity within the org — the value the classifier
	// answers with and the value stored on Conversation.TypeKey.
	Key  string `json:"key" sortable:""`
	Name string `json:"name,omitempty"`
	// Definition tells the classifier when a conversation IS this type — it is
	// injected verbatim into the classification prompt, so it should describe
	// the conversation's content/participants/purpose, not internal process.
	Definition string                 `json:"definition"`
	Config     ConversationTypeConfig `json:"config"`
	// Outcomes is the ordered list a conversation of this type can end with;
	// empty = the type tracks no outcome. Serializes as an array, never null.
	Outcomes []ConversationOutcome `json:"outcomes"`
	// Phases is the ordered list of stages a conversation of this type moves
	// through; empty = the type tracks no phases. Serializes as an array,
	// never null.
	Phases []ConversationPhase `json:"phases"`
	// Actions are what to do next once a conversation of this type has an
	// outcome — the softphone's post-call "Next steps" (record proposals,
	// agent hand-offs, actions, workflows). Same contract as Page.Actions;
	// each button's visible_when reads the wrap-up scope (conversation,
	// contact, records, records_by_entity, …). Serializes as an array, never
	// null.
	Actions []common.PageAction `json:"actions"`
	// ModuleSlug attributes the type to the module that deployed it; empty =
	// not module-owned.
	ModuleSlug string         `json:"module_slug,omitempty"`
	CreatedAt  time.Time      `json:"created_at" sortable:""`
	CreatedBy  common.UserRef `json:"created_by"`
	UpdatedAt  time.Time      `json:"updated_at" sortable:""`
	UpdatedBy  common.UserRef `json:"updated_by"`
}

// OutcomeByKey finds one of the type's outcomes.
func (conversationType ConversationType) OutcomeByKey(key string) (ConversationOutcome, bool) {
	for _, outcome := range conversationType.Outcomes {
		if outcome.Key == key {
			return outcome, true
		}
	}
	return ConversationOutcome{}, false
}

// PhaseByKey finds one of the type's phases.
func (conversationType ConversationType) PhaseByKey(key string) (ConversationPhase, bool) {
	for _, phase := range conversationType.Phases {
		if phase.Key == key {
			return phase, true
		}
	}
	return ConversationPhase{}, false
}

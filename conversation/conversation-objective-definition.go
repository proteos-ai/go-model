package conversationmodel

import (
	"time"

	"go.proteos.ai/model/common"
)

// ConversationObjectiveDefinition is one thing a conversation should achieve
// — "confirm the budget", "agree a next step", "book the demo". It applies to
// every conversation (TypeKey empty) or to conversations of one type, may be
// pinned to one of that type's phases (the checklist groups by phase, and a
// phase counts as done once every measured objective in it is achieved),
// and may be MEASURED so the platform tells whether it was achieved (nil
// Measure = guidance only, a line on the checklist). Its instance on a
// conversation is a ConversationObjective — the same definition → instance
// split as ConversationTagDefinition → ConversationTag. Keyed (org_id, key)
// and module-deployable (conversation-objective-definitions/<key>.json).
type ConversationObjectiveDefinition struct {
	OrgId       string `json:"org_id"`
	Key         string `json:"key" sortable:""`
	Name        string `json:"name" sortable:""`
	Description string `json:"description,omitempty"`
	// TypeKey restricts the objective to conversations of one type; empty =
	// every conversation. A type deleted later just never matches.
	TypeKey string `json:"type_key,omitempty" sortable:""`
	// PhaseKey pins the objective to one of the type's phases (requires
	// TypeKey; must be one of the type's phase keys at write time); empty =
	// ungrouped.
	PhaseKey string `json:"phase_key,omitempty"`
	// Measure says how the platform tells the objective was achieved; nil =
	// it cannot (guidance only).
	Measure   *ConversationObjectiveMeasure `json:"measure,omitempty"`
	IsEnabled bool                          `json:"is_enabled"`
	// ModuleSlug attributes the definition to the module that deployed it;
	// empty = not module-owned.
	ModuleSlug string         `json:"module_slug,omitempty"`
	CreatedAt  time.Time      `json:"created_at" sortable:""`
	CreatedBy  common.UserRef `json:"created_by"`
	UpdatedAt  time.Time      `json:"updated_at" sortable:""`
	UpdatedBy  common.UserRef `json:"updated_by"`
}

// ConversationObjectiveMeasure is how an objective's achievement is told:
//
//	kind tag      a non-dismissed tag on TagDefinitionKey whose value meets
//	              AchievedWhen (a pending tag ⇒ achieved, provisional). The
//	              tag definition must be conversation-scoped and role-less —
//	              a person ticks the objective by creating a user tag on it,
//	              unticks by deleting that tag.
//	kind outcome  the conversation's outcome_key is one of
//	              AchievedWhen.ChoiceIn (the type's outcomes and the
//	              channel's system outcomes) — ticked through the outcome
//	              picker.
type ConversationObjectiveMeasure struct {
	Kind             ConversationObjectiveMeasureKind  `json:"kind"`
	TagDefinitionKey string                            `json:"tag_definition_key,omitempty"`
	AchievedWhen     ConversationObjectiveAchievedWhen `json:"achieved_when"`
}

// ConversationObjectiveAchievedWhen is the value condition, exactly one of
// its fields set, matched to the question type of the measured tag
// definition (boolean / choice / score) or, for kind outcome, ChoiceIn
// over outcome keys. ScoreMin is the 0-based level index a score tag must
// reach (tag values store the level index).
type ConversationObjectiveAchievedWhen struct {
	Boolean  *bool    `json:"boolean,omitempty"`
	ChoiceIn []string `json:"choice_in,omitempty"`
	ScoreMin *int     `json:"score_min,omitempty"`
}

// IsMeasured reports that the platform can tell whether the objective was
// achieved (a Measure is set).
func (definition ConversationObjectiveDefinition) IsMeasured() bool {
	return definition.Measure != nil
}

// AppliesTo reports whether the definition applies to a conversation of the
// type key (empty TypeKey = every conversation).
func (definition ConversationObjectiveDefinition) AppliesTo(typeKey string) bool {
	return definition.TypeKey == "" || definition.TypeKey == typeKey
}

package conversationmodel

import "time"

// ConversationObjective is one objective definition's standing on one
// conversation: open or achieved, by whom, when, and on which tag. It is
// resolved live from the conversation's tags and outcome for the
// per-conversation read (ConversationObjectiveProgress) AND persisted as the
// ledger row of the same shape — one row per APPLICABLE definition, open
// rows included, rewritten by the service at every change point (tags
// landed, user tag changed, outcome set, conversation ended). Rows outlive a
// deleted definition: they are the record of what was and was not achieved,
// never user-written.
type ConversationObjective struct {
	OrgId          string `json:"org_id"`
	ConversationId string `json:"conversation_id"`
	DefinitionKey  string `json:"definition_key" sortable:""`
	// Channel / TypeKey are the conversation's, PhaseKey the definition's,
	// denormalized at resolution time so the ledger groups without a join.
	Channel  Channel                     `json:"channel" sortable:""`
	TypeKey  string                      `json:"type_key,omitempty" sortable:""`
	PhaseKey string                      `json:"phase_key,omitempty" sortable:""`
	Status   ConversationObjectiveStatus `json:"status" sortable:""`
	// Source says who achieved it: user (a person's tag or outcome pick),
	// model (a detected tag or outcome), system (a channel-level outcome such
	// as not_reached). Empty while open.
	Source ConversationObjectiveSource `json:"source,omitempty" sortable:""`
	// IsProvisional: the backing tag is still pending — a realtime detection
	// on a live spoken conversation the completion pass has not confirmed.
	// Never true for a user or system source.
	IsProvisional bool       `json:"is_provisional"`
	AchievedAt    *time.Time `json:"achieved_at,omitempty" sortable:""`
	// TagId is the tag that achieved a tag-measured objective (the one to
	// delete to untick when Source is user); empty for outcome measures.
	TagId string `json:"tag_id,omitempty"`
	// IsConversationEnded is the conversation's status at resolution time —
	// an open objective on an ended conversation was missed.
	IsConversationEnded bool `json:"is_conversation_ended"`
	// RecordedAt is when this standing was resolved (the ledger write).
	RecordedAt time.Time `json:"recorded_at" sortable:""`
}

// IsAchieved reports the objective's status.
func (objective ConversationObjective) IsAchieved() bool {
	return objective.Status == ConversationObjectiveStatusAchieved
}

// ConversationObjectiveCount is one row of GET /conversation-objectives/counts:
// the grouped columns that were requested plus the count.
type ConversationObjectiveCount struct {
	DefinitionKey string                      `json:"definition_key,omitempty"`
	Status        ConversationObjectiveStatus `json:"status,omitempty"`
	Source        ConversationObjectiveSource `json:"source,omitempty"`
	IsProvisional *bool                       `json:"is_provisional,omitempty"`
	TypeKey       string                      `json:"type_key,omitempty"`
	PhaseKey      string                      `json:"phase_key,omitempty"`
	Channel       Channel                     `json:"channel,omitempty"`
	Count         int                         `json:"count"`
}

// ConversationPhaseProgress is one of the type's phases as the conversation
// stands in it, derived from the objectives pinned to the phase: done =
// every measured objective achieved (≥ 1 measured); current = the first
// phase in type order that is not done and has a measured objective;
// upcoming = after the current one; untracked = the phase has no measured
// objective, so progress cannot be told (its guidance objectives still
// list). CompletedAt is the latest achieved_at among the phase's measured
// objectives once done.
type ConversationPhaseProgress struct {
	ConversationPhase
	Status      ConversationPhaseStatus `json:"status"`
	CompletedAt *time.Time              `json:"completed_at,omitempty"`
}

// ConversationObjectiveProgress is GET /conversations/:id/objective-progress:
// where the conversation stands — the current phase key ("" when the type
// has no phases or every tracked phase is done), every phase with its
// status, every applicable objective with its standing, and the applicable
// definitions (name, description, measure) the objectives reference by
// definition_key, so a client renders the checklist from one read.
type ConversationObjectiveProgress struct {
	PhaseKey    string                            `json:"phase_key"`
	Phases      []ConversationPhaseProgress       `json:"phases"`
	Objectives  []ConversationObjective           `json:"objectives"`
	Definitions []ConversationObjectiveDefinition `json:"definitions"`
}

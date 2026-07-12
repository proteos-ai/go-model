package conversationapi

import (
	"go.proteos.ai/model/common"
	conversationmodel "go.proteos.ai/model/conversation"
)

// CreateGlossaryTermRequest adds one custom-vocabulary term to the org glossary.
// Term is the boost phrase injected into transcription; Definition is optional
// documentation (never sent to the provider). Priority is an optional rank —
// prioritized terms are selected ahead of unprioritized ones, lower number
// first — used to choose which terms win when the org exceeds the provider's
// per-request keyterm cap.
type CreateGlossaryTermRequest struct {
	Term       string `json:"term" validate:"required,max=100"`
	Definition string `json:"definition" validate:"max=4096"`
	Priority   *int   `json:"priority,omitempty" validate:"omitempty,gte=0"`
}

// UpdateGlossaryTermRequest is a partial update — every field is a tri-state
// pointer: nil leaves the stored value untouched. Priority additionally
// distinguishes "leave as-is" (omitted) from "clear the priority" (explicit
// null) via IsPriorityCleared, since a nil *int cannot itself carry that intent.
type UpdateGlossaryTermRequest struct {
	Term       *string `json:"term,omitempty" validate:"omitempty,max=100"`
	Definition *string `json:"definition,omitempty" validate:"omitempty,max=4096"`
	Priority   *int    `json:"priority,omitempty" validate:"omitempty,gte=0"`
	// IsPriorityCleared, when true, removes the term's priority (sets it NULL)
	// regardless of Priority. It exists because an omitted `priority` and an
	// explicit `"priority": null` both arrive as a nil pointer; callers that
	// mean "unprioritize this term" set this flag.
	IsPriorityCleared bool `json:"is_priority_cleared,omitempty"`
}

type GetManyGlossaryTermsQuery struct {
	// Search filters terms by a case-insensitive substring of the term text.
	Search *string `json:"search" form:"search"`
	common.Pagination
	common.Sorting
}

type GetManyGlossaryTermsResponse struct {
	Meta common.ResponseMeta              `json:"meta"`
	Data []conversationmodel.GlossaryTerm `json:"data"`
}

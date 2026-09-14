package conversationapi

import (
	"go.proteos.ai/model/common"
	conversationmodel "go.proteos.ai/model/conversation"
)

// CreateToneProfileSetupRequest opts one platform user into tone-of-voice
// synthesis. The client sends the bare user id; the service resolves it to a
// full UserRef against the account directory.
type CreateToneProfileSetupRequest struct {
	OwnedById string `json:"owned_by_id" validate:"required,max=36"`
}

type GetManyToneProfileSetupsQuery struct {
	common.Pagination
	common.Sorting
}

type GetManyToneProfileSetupsResponse struct {
	Meta common.ResponseMeta                  `json:"meta"`
	Data []conversationmodel.ToneProfileSetup `json:"data"`
}

// UpsertToneProfileRequest writes ONE hand-authored tone profile. Identity is
// the scope tuple (profiled user × channel × contact group × contact), not an
// id — so the same request re-sent updates the row it created. The written row
// is stamped source=manual, which locks that scope against the synthesis sweep;
// it deliberately overwrites a model row sitting at the same scope.
//
// Scope grammar (mirrors the resolve precedence): omit Channel for the
// cross-channel base voice; ContactGroupKey requires a Channel; ContactId
// requires both, because resolution only reaches a contact row through its
// group.
type UpsertToneProfileRequest struct {
	OwnedById       string                    `json:"owned_by_id" validate:"required,max=36"`
	Channel         conversationmodel.Channel `json:"channel"`
	ContactGroupKey string                    `json:"contact_group_key" validate:"max=255"`
	ContactId       string                    `json:"contact_id" validate:"max=36"`
	// Instructions is the COMPLETE markdown instruction set for this tier —
	// served verbatim to a drafting model, never concatenated with another tier.
	Instructions string `json:"instructions" validate:"required"`
	// Differences is the optional human-facing note on how this tier deviates
	// from the one above; never sent to a drafter.
	Differences string `json:"differences"`
}

type GetManyToneProfilesQuery struct {
	// OwnedById filters to one profiled user's rows (the detail view).
	OwnedById *string `json:"owned_by_id" form:"owned_by_id"`
	// Channel filters to one medium's rows ('' has no special meaning here —
	// omit the filter to get every tier including the user aggregate).
	Channel *conversationmodel.Channel `json:"channel" form:"channel"`
	// Scope filters to one tier (user | channel | group | contact).
	Scope *conversationmodel.ToneProfileScope `json:"scope" form:"scope"`
	// Source filters to hand-authored (manual) or synthesized (model) rows.
	Source *conversationmodel.ToneProfileSource `json:"source" form:"source"`
	common.Pagination
	common.Sorting
}

type GetManyToneProfilesResponse struct {
	Meta common.ResponseMeta             `json:"meta"`
	Data []conversationmodel.ToneProfile `json:"data"`
}

// ResolveToneProfileQuery asks for the single most-specific profile for a
// drafting context: contact row (when ContactId is set and the contact is
// grouped) → group row → (user, channel) base → user aggregate. Channel empty
// = "no channel context", which resolves straight to the user aggregate.
type ResolveToneProfileQuery struct {
	OwnedById string                     `json:"owned_by_id" form:"owned_by_id" validate:"required"`
	Channel   *conversationmodel.Channel `json:"channel" form:"channel"`
	ContactId *string                    `json:"contact_id" form:"contact_id"`
}

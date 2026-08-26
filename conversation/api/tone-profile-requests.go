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

type GetManyToneProfilesQuery struct {
	// OwnedById filters to one profiled user's rows (the detail view).
	OwnedById *string `json:"owned_by_id" form:"owned_by_id"`
	// Channel filters to one medium's rows ('' has no special meaning here —
	// omit the filter to get every tier including the user aggregate).
	Channel *conversationmodel.Channel `json:"channel" form:"channel"`
	// Scope filters to one tier (user | channel | group | contact).
	Scope *conversationmodel.ToneProfileScope `json:"scope" form:"scope"`
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

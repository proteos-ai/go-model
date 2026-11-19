package conversationapi

import (
	"go.proteos.ai/model/common"
	conversationmodel "go.proteos.ai/model/conversation"
)

// CreateEmailSignatureRequest creates one signature. Html or Text is required.
// Text left empty is derived from Html. OwnedById empty makes it the
// organisation's; ConnectionId empty makes it apply to every mailbox.
// IsDefault defaults to false and AppliesTo to all.
type CreateEmailSignatureRequest struct {
	Name         string                                    `json:"name" validate:"required,max=120"`
	Html         string                                    `json:"html" validate:"max=100000"`
	Text         string                                    `json:"text,omitempty" validate:"max=20000"`
	OwnedById    string                                    `json:"owned_by_id,omitempty" validate:"max=36"`
	ConnectionId string                                    `json:"connection_id,omitempty" validate:"max=36"`
	IsDefault    *bool                                     `json:"is_default,omitempty"`
	AppliesTo    conversationmodel.EmailSignatureAppliesTo `json:"applies_to,omitempty"`
}

// UpdateEmailSignatureRequest changes a signature in place; absent fields keep
// what is stored. Text sent empty goes back to being derived from the html.
// OwnedById / ConnectionId sent empty hand the signature to the organisation /
// to every mailbox.
type UpdateEmailSignatureRequest struct {
	Name         *string                                    `json:"name,omitempty" validate:"omitempty,min=1,max=120"`
	Html         *string                                    `json:"html,omitempty" validate:"omitempty,max=100000"`
	Text         *string                                    `json:"text,omitempty" validate:"omitempty,max=20000"`
	OwnedById    *string                                    `json:"owned_by_id,omitempty" validate:"omitempty,max=36"`
	ConnectionId *string                                    `json:"connection_id,omitempty" validate:"omitempty,max=36"`
	IsDefault    *bool                                      `json:"is_default,omitempty"`
	AppliesTo    *conversationmodel.EmailSignatureAppliesTo `json:"applies_to,omitempty"`
}

// GetManyEmailSignaturesQuery filters the org's signatures. OwnedById is one
// person's; IsOrganisation true is the organisation's own (no owner), false
// is everyone's personal ones.
type GetManyEmailSignaturesQuery struct {
	OwnedById      *string `json:"owned_by_id" form:"owned_by_id"`
	IsOrganisation *bool   `json:"is_organisation" form:"is_organisation"`
	ConnectionId   *string `json:"connection_id" form:"connection_id" db:"connection_id"`
	IsDefault      *bool   `json:"is_default" form:"is_default" db:"is_default"`
	common.Pagination
	common.Sorting
}

type GetManyEmailSignaturesResponse struct {
	Meta common.ResponseMeta                `json:"meta"`
	Data []conversationmodel.EmailSignature `json:"data"`
}

// ResolveEmailSignatureQuery asks which default signature an email would get:
// the mailbox's (ConnectionId), then the sender's own (OwnedById, the caller
// when empty and a person), then the organisation's. IsReply picks between
// the defaults for new emails and for replies.
type ResolveEmailSignatureQuery struct {
	ConnectionId string `json:"connection_id" form:"connection_id"`
	OwnedById    string `json:"owned_by_id" form:"owned_by_id"`
	IsReply      bool   `json:"is_reply" form:"is_reply"`
}

// ResolveEmailSignatureResponse carries the signature that applies; Data is
// null when none does, which is an ordinary answer, not an error.
type ResolveEmailSignatureResponse struct {
	Data *conversationmodel.MessageSignature `json:"data"`
}

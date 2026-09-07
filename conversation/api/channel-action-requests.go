package conversationapi

import (
	"go.proteos.ai/model/common"
	conversationmodel "go.proteos.ai/model/conversation"
)

// PerformChannelActionRequest performs one act through a connection: send a
// LinkedIn invitation, visit a profile, send an InMail. Target is the person
// acted on — the same polymorphic recipient shape a send's `to` carries (kind
// contact-address + the connector-side external id, a LinkedIn member id).
// Params is the raw per-type input; the service validates and types it
// against ActionType (see conversationmodel.ChannelActionParams).
type PerformChannelActionRequest struct {
	ConnectionId string                              `json:"connection_id" validate:"required"`
	ActionType   conversationmodel.ChannelActionType `json:"action_type" validate:"required"`
	Target       conversationmodel.Recipient         `json:"target"`
	Params       map[string]any                      `json:"params"`
}

// RespondChannelActionRequest answers an INBOUND action (a received
// invitation): accept or decline.
type RespondChannelActionRequest struct {
	Response conversationmodel.ChannelActionResponse `json:"response" validate:"required"`
}

// GetManyChannelActionsQuery filters the org's ledger.
type GetManyChannelActionsQuery struct {
	// Channel narrows to one medium (the hub's channel filter) — the ledger
	// carries it denormalized, so no connection lookup is needed.
	Channel      *string `json:"channel" form:"channel" db:"channel"`
	ConnectionId *string `json:"connection_id" form:"connection_id" db:"connection_id"`
	ActionType   *string `json:"action_type" form:"action_type" db:"action_type"`
	Direction    *string `json:"direction" form:"direction" db:"direction"`
	Status       *string `json:"status" form:"status" db:"status"`
	ContactId    *string `json:"contact_id" form:"contact_id" db:"contact_id"`
	common.Pagination
	common.Sorting
}

type GetManyChannelActionsResponse struct {
	Meta common.ResponseMeta               `json:"meta"`
	Data []conversationmodel.ChannelAction `json:"data"`
}

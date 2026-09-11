package conversationapi

import (
	"time"

	"go.proteos.ai/model/common"
	conversationmodel "go.proteos.ai/model/conversation"
)

// GetManyChannelEventsQuery pages the channel_event ledger. MessageId narrows
// to one message's events (the per-message list also lives under
// /messages/:id/events); ConnectionId + EventType + the From/Until window are
// the connection-health slices.
type GetManyChannelEventsQuery struct {
	ConnectionId   *string    `json:"connection_id" form:"connection_id" db:"connection_id"`
	MessageId      *string    `json:"message_id" form:"message_id" db:"message_id"`
	ConversationId *string    `json:"conversation_id" form:"conversation_id" db:"conversation_id"`
	ContactId      *string    `json:"contact_id" form:"contact_id" db:"contact_id"`
	EventType      *string    `json:"event_type" form:"event_type" db:"event_type"`
	From           *time.Time `json:"from" form:"from"`
	Until          *time.Time `json:"until" form:"until"`
	common.Pagination
	common.Sorting
}

type GetManyChannelEventsResponse struct {
	Meta common.ResponseMeta              `json:"meta"`
	Data []conversationmodel.ChannelEvent `json:"data"`
}

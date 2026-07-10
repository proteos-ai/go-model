package conversationapi

import (
	conversationmodel "go.proteos.ai/model/conversation"
)

// AddReactionRequest places an emoji reaction on a message. Emoji is the
// connector-native token (Slack shortcode like "thumbsup"); on a fixed-set
// channel it must be one of the connection capability's allowed tokens.
type AddReactionRequest struct {
	Emoji string `json:"emoji" validate:"required"`
}

// ListReactionsResponse is the detailed per-emoji aggregation for one message
// (the same shape that rides on Message.reactions).
type ListReactionsResponse struct {
	Data []conversationmodel.Reaction `json:"data"`
}

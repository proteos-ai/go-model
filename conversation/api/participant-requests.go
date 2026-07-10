package conversationapi

import (
	"go.proteos.ai/model/common"
	conversationmodel "go.proteos.ai/model/conversation"
)

// SearchParticipantsQuery filters a connection's participant directory. Q is a
// free-text needle matched (case-insensitive) against display name AND external
// id (handled directly in the repo, not via the generic op-tag filter, since it
// spans two columns).
type SearchParticipantsQuery struct {
	Q *string `json:"q" form:"q"`
	common.Pagination
	common.Sorting
}

type SearchParticipantsResponse struct {
	Meta common.ResponseMeta             `json:"meta"`
	Data []conversationmodel.Participant `json:"data"`
}

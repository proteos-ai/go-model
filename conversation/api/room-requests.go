package conversationapi

import (
	"go.proteos.ai/model/common"
	conversationmodel "go.proteos.ai/model/conversation"
)

// SearchRoomsQuery filters a connection's room directory. Q is a free-text
// needle matched (case-insensitive) against name AND external id.
type SearchRoomsQuery struct {
	Q *string `json:"q" form:"q"`
	common.Pagination
	common.Sorting
}

type SearchRoomsResponse struct {
	Meta common.ResponseMeta      `json:"meta"`
	Data []conversationmodel.Room `json:"data"`
}

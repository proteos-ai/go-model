package conversationapi

import (
	"go.proteos.ai/model/common"
	conversationmodel "go.proteos.ai/model/conversation"
)

// SearchContactAddressesQuery filters the org's contact-address directory for
// the recipient picker. Q is a free-text needle matched (case-insensitive)
// against display name AND canonical value (handled directly in the repo, not
// via the generic op-tag filter, since it spans two columns). Kind/Scope
// narrow to what one connection can reach (the controller derives them from
// the connection's channel + external workspace).
type SearchContactAddressesQuery struct {
	Q     *string                                `json:"q" form:"q"`
	Kinds []conversationmodel.ContactAddressKind `json:"kinds" form:"kinds"`
	Scope *string                                `json:"scope" form:"scope"`
	common.Pagination
	common.Sorting
}

type SearchContactAddressesResponse struct {
	Meta common.ResponseMeta                `json:"meta"`
	Data []conversationmodel.ContactAddress `json:"data"`
}

// AttachContactAddressRequest manually attaches one address to a contact. The
// caller provides the raw identifier; the domain canonicalizes it (lowercased
// email, E.164 phone) and rejects duplicates owned by another contact.
type AttachContactAddressRequest struct {
	Kind  conversationmodel.ContactAddressKind `json:"kind" binding:"required"`
	Scope string                               `json:"scope"`
	Value string                               `json:"value" binding:"required"`
	Name  string                               `json:"name"`
}

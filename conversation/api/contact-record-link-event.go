package conversationapi

import (
	"go.proteos.ai/model/common"
	conversationmodel "go.proteos.ai/model/conversation"
)

// ContactRecordLinkEventPayload is the payload of contact_record_link.events
// (created | updated | deleted): the full link plus what data-service needs
// to maintain record_duplicate rows without a second round-trip.
//
//   - DuplicateRecordIds are the OTHER records of the link's entity that are
//     linked to the same contact — the record-duplicate proof. Empty when the
//     binding found no sibling.
//   - SharedAddressKeys are the canonical keys the record shares with the
//     contact (the evidence shown on the duplicate banner).
//   - Actor is who caused the change (the writing user, or system for
//     merges/sweeps).
//
// Mirrors RecordEventPayload{record, previous_record, actor}: the entity plus
// the context a consumer would otherwise have to re-derive.
type ContactRecordLinkEventPayload struct {
	Link               conversationmodel.ContactRecordLink   `json:"link"`
	DuplicateRecordIds []string                              `json:"duplicate_record_ids"`
	SharedAddressKeys  []conversationmodel.ContactAddressKey `json:"shared_address_keys"`
	Actor              common.UserRef                        `json:"actor"`
}

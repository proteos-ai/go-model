package datamodel

import (
	"time"

	"go.proteos.ai/model/common"
	conversationmodel "go.proteos.ai/model/conversation"
)

// RecordDuplicateSource is HOW a duplicate pair was detected. contact_binding
// is the v1 method (two records of one entity bound to the same conversation
// contact through their contact-address attributes); the axis exists so
// rule-based, AI and manual detection can share the same resource later.
type RecordDuplicateSource string

const (
	RecordDuplicateSourceContactBinding RecordDuplicateSource = "contact_binding"
)

// RecordDuplicateStatus is the review lifecycle of a pair: open (shown on both
// records), dismissed (a human said "not a duplicate"), resolved (the pair
// stopped holding — a record was deleted or rebound elsewhere).
type RecordDuplicateStatus string

const (
	RecordDuplicateStatusOpen      RecordDuplicateStatus = "open"
	RecordDuplicateStatusDismissed RecordDuplicateStatus = "dismissed"
	RecordDuplicateStatusResolved  RecordDuplicateStatus = "resolved"
)

// IsValid reports whether the status is a known value.
func (status RecordDuplicateStatus) IsValid() bool {
	switch status {
	case RecordDuplicateStatusOpen, RecordDuplicateStatusDismissed, RecordDuplicateStatusResolved:
		return true
	}
	return false
}

// RecordDuplicate is one "these two records of ONE entity may be the same
// thing" fact, owned by data-service because it is a record concern.
//
// The pair is asymmetric: PrimaryRecordId is the incumbent — the record to
// keep — and SecondaryRecordId the newcomer that revealed the collision. It
// plays the role suggested_winner_id plays on a contact merge proposal.
//
// Star topology: a group of records that collide for one reason (for
// contact_binding: the same entity bound to the same contact) has exactly ONE
// primary — the oldest record, or the primary the group's existing open rows
// already name — and every other member gets one row against it. Four
// colliding records are three rows, never six; the primary's banner lists all
// secondaries, each secondary's banner names the primary. At most one OPEN row
// per (entity, source, primary, secondary).
//
// Evidence is shaped by Source: for contact_binding it is
// ContactBindingEvidence{contact_id, address_keys}.
type RecordDuplicate struct {
	Id                string                `json:"id" sortable:""`
	OrgId             string                `json:"org_id"`
	EntitySlug        string                `json:"entity_slug" sortable:""`
	PrimaryRecordId   string                `json:"primary_record_id" sortable:""`
	SecondaryRecordId string                `json:"secondary_record_id" sortable:""`
	Source            RecordDuplicateSource `json:"source" sortable:""`
	Evidence          map[string]any        `json:"evidence"`
	Status            RecordDuplicateStatus `json:"status" sortable:""`
	ResolvedBy        *common.UserRef       `json:"resolved_by,omitempty"`
	ResolvedAt        *time.Time            `json:"resolved_at,omitempty"`
	CreatedAt         time.Time             `json:"created_at" sortable:""`
	CreatedBy         common.UserRef        `json:"created_by"`
	UpdatedAt         time.Time             `json:"updated_at" sortable:""`
	UpdatedBy         common.UserRef        `json:"updated_by"`
}

// ContactBindingEvidence is the Evidence shape of a contact_binding duplicate:
// the contact both records are bound to and the canonical address keys the
// secondary shares with it.
type ContactBindingEvidence struct {
	ContactId   string                                `json:"contact_id"`
	AddressKeys []conversationmodel.ContactAddressKey `json:"address_keys"`
}

// ToEvidence renders the evidence as the generic bag stored on the row.
func (evidence ContactBindingEvidence) ToEvidence() map[string]any {
	keys := make([]map[string]any, 0, len(evidence.AddressKeys))
	for _, key := range evidence.AddressKeys {
		keys = append(keys, map[string]any{"kind": string(key.Kind), "value": key.Value})
	}
	return map[string]any{
		"contact_id":   evidence.ContactId,
		"address_keys": keys,
	}
}

// ParseContactBindingEvidence reads the contact_binding shape back out of a
// row's Evidence bag. ok=false when the bag carries no contact id.
func ParseContactBindingEvidence(evidence map[string]any) (ContactBindingEvidence, bool) {
	contactId, _ := evidence["contact_id"].(string)
	if contactId == "" {
		return ContactBindingEvidence{}, false
	}
	parsed := ContactBindingEvidence{ContactId: contactId}
	// The bag is []any after a JSON round trip and []map[string]any when it
	// was built in-process by ToEvidence — read both.
	var rawKeys []map[string]any
	switch keys := evidence["address_keys"].(type) {
	case []map[string]any:
		rawKeys = keys
	case []any:
		for _, rawKey := range keys {
			if key, ok := rawKey.(map[string]any); ok {
				rawKeys = append(rawKeys, key)
			}
		}
	}
	for _, key := range rawKeys {
		kind, _ := key["kind"].(string)
		value, _ := key["value"].(string)
		if kind == "" || value == "" {
			continue
		}
		parsed.AddressKeys = append(parsed.AddressKeys, conversationmodel.ContactAddressKey{
			Kind:  conversationmodel.ContactAddressKind(kind),
			Value: value,
		})
	}
	return parsed, true
}

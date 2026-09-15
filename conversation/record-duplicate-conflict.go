package conversationmodel

// RecordDuplicateConflict is the evidence that a record's contact-address
// values collide with another record of the SAME entity: the contact the
// addresses resolve to, the other records already linked to that contact for
// the entity, and the address keys that were hit. It is the `details` of a
// 409 record_duplicate (duplicate_policy reject) and the `duplicate` field of
// a rejected resolution inside a batch. Under duplicate_policy flag the same
// facts travel on contact_record_link.events as duplicate_record_ids +
// shared_address_keys, and data-service turns them into record_duplicate rows.
type RecordDuplicateConflict struct {
	EntitySlug  string              `json:"entity_slug"`
	ContactId   string              `json:"contact_id"`
	RecordIds   []string            `json:"record_ids"`
	AddressKeys []ContactAddressKey `json:"address_keys"`
}

// Details renders the conflict as the `details` map of an error body.
func (conflict RecordDuplicateConflict) Details() map[string]any {
	keys := make([]map[string]any, 0, len(conflict.AddressKeys))
	for _, key := range conflict.AddressKeys {
		keys = append(keys, map[string]any{"kind": key.Kind, "value": key.Value})
	}
	return map[string]any{
		"entity_slug":  conflict.EntitySlug,
		"contact_id":   conflict.ContactId,
		"record_ids":   conflict.RecordIds,
		"address_keys": keys,
	}
}

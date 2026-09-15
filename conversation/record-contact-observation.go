package conversationmodel

import "go.proteos.ai/model/common"

// RecordContactAddress is one typed `contact-address` attribute value of a
// business record, as data-service canonicalized it on the write
// (CanonicalizeContactAddress). Kind is one of RecordContactAddressKinds.
type RecordContactAddress struct {
	Kind  ContactAddressKind `json:"kind"`
	Value string             `json:"value"`
}

// RecordContactObservation is ONE sighting of a person through a business
// record — the record-side twin of ContactObservation (a sighting through a
// channel). It is the payload of record_contact_observation.created/deleted
// (the org also rides the envelope) and the per-record item of
// POST /conversations/v1/contact-record-links/resolve.
//
// Addresses are in attribute order: the first key that already has an owner
// wins resolution, exactly like the channel-primary key of an ingest sighting.
// Name is the record's materialized title — the name a contact minted from
// this record carries. DuplicatePolicy is the entity's contact_binding
// setting at the time of the write (data-service is the source of truth;
// conversation-service only applies it).
type RecordContactObservation struct {
	OrgId           string                 `json:"org_id"`
	EntitySlug      string                 `json:"entity_slug"`
	RecordId        string                 `json:"record_id"`
	Name            string                 `json:"name,omitempty"`
	Addresses       []RecordContactAddress `json:"addresses"`
	DuplicatePolicy RecordDuplicatePolicy  `json:"duplicate_policy"`
	// Actor is the user whose write produced the observation; the async
	// consumer stamps it as created_by on anything it mints. Absent = system.
	Actor *common.UserRef `json:"actor,omitempty"`
}

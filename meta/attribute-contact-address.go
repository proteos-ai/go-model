package metamodel

import conversationmodel "go.proteos.ai/model/conversation"

// ContactAddressAttributeMeta holds the metadata for an attribute of type
// `contact-address`: ONE canonical digital endpoint of the person the record
// represents. The stored value is a scalar string canonicalized on write
// (lowercased email, E.164 phone, LinkedIn public identifier — see
// conversationmodel.CanonicalizeContactAddress); records carrying such values
// are bound to exactly one conversation contact, and that binding is where
// contact-based deduplication happens (see Entity.ContactBinding).
//
// Kind is REQUIRED — an address without a kind cannot be canonicalized — and
// must be one of conversationmodel.RecordContactAddressKinds (email | phone |
// linkedin). Two attributes may share a kind ("phone" and "mobile" are two
// phone-kind attributes).
type ContactAddressAttributeMeta struct {
	Description string                               `json:"description,omitempty"`
	Kind        conversationmodel.ContactAddressKind `json:"kind"`
}

// ContactAddressMetaOf returns the contact-address meta for an attribute, in
// whatever shape the decoder that produced it left behind (struct, pointer,
// or — the normal case for anything that crossed JSONB/HTTP — map). Uses
// ParseMetaAs for the same reason PrincipalMetaOf does: a bare type switch
// silently misses the map shape every real request carries.
func ContactAddressMetaOf(attribute Attribute) (ContactAddressAttributeMeta, bool) {
	if attribute.Type != AttributeTypeContactAddress {
		return ContactAddressAttributeMeta{}, false
	}
	meta := ParseMetaAs[ContactAddressAttributeMeta](attribute.Meta)
	if meta == nil {
		return ContactAddressAttributeMeta{}, false
	}
	return *meta, true
}

// ContactAddressAttributes returns the top-level contact-address attributes in
// attribute order — the order that decides which address wins resolution.
func ContactAddressAttributes(attributes []Attribute) []Attribute {
	result := make([]Attribute, 0)
	for _, attribute := range attributes {
		if attribute.Type == AttributeTypeContactAddress {
			result = append(result, attribute)
		}
	}
	return result
}

// HasContactAddressAttributes is the short-circuit that keeps the common case
// (an entity without any contact-address attribute) free of binding work.
func HasContactAddressAttributes(attributes []Attribute) bool {
	for _, attribute := range attributes {
		if attribute.Type == AttributeTypeContactAddress {
			return true
		}
	}
	return false
}

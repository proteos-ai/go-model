package metamodel

import (
	"fmt"

	conversationmodel "go.proteos.ai/model/conversation"
)

// ContactBinding is the entity-level setting that governs how records of the
// entity bind to conversation-service contacts through their `contact-address`
// attributes. It is meaningful only for entities that declare at least one
// such attribute, and it is stored ONLY when explicitly declared — an entity
// that never mentions it keeps a null contact_binding rather than a
// materialized default, so a manifest that omits the field stays in sync with
// what the server reports.
//
// DuplicatePolicy decides what happens when a record's addresses resolve to a
// contact that ANOTHER record of the same entity is already linked to:
//
//   - reject: the write fails with 409 record_duplicate (data-service resolves
//     the binding inline, inside the write transaction).
//   - flag (the effective default when unset): the write succeeds; the binding
//     is resolved asynchronously and data-service opens a record_duplicate the
//     record page surfaces until it is dismissed.
type ContactBinding struct {
	DuplicatePolicy conversationmodel.RecordDuplicatePolicy `json:"duplicate_policy"`
}

// IsZero reports whether the binding carries no setting at all. A non-nil
// pointer to a zero ContactBinding is how a writer says "clear this back to
// unset" (see UpdateEntityRequest.ContactBinding); it is stored as null.
func (binding ContactBinding) IsZero() bool {
	return binding.DuplicatePolicy == ""
}

// EffectiveDuplicatePolicy is the policy a write applies: the stored value, or
// flag when the entity declares no binding (the common case — most entities
// have no contact-address attributes at all).
func EffectiveDuplicatePolicy(binding *ContactBinding) conversationmodel.RecordDuplicatePolicy {
	if binding == nil || binding.DuplicatePolicy == "" {
		return conversationmodel.RecordDuplicatePolicyFlag
	}
	return binding.DuplicatePolicy
}

// ValidateContactBinding accepts an unset policy (which means "no binding
// declared") or one of the known values.
func ValidateContactBinding(binding ContactBinding) error {
	if binding.DuplicatePolicy == "" || binding.DuplicatePolicy.IsValid() {
		return nil
	}
	return fmt.Errorf("unknown duplicate_policy %q (expected one of %s, %s)",
		binding.DuplicatePolicy, conversationmodel.RecordDuplicatePolicyReject, conversationmodel.RecordDuplicatePolicyFlag)
}

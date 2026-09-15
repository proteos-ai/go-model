package metamodel

import (
	"fmt"

	conversationmodel "go.proteos.ai/model/conversation"
)

// ContactBinding is the entity-level setting that governs how records of the
// entity bind to conversation-service contacts through their `contact-address`
// attributes. It is meaningful only for entities that declare at least one
// such attribute; it is stored for every entity so the default is explicit
// rather than implied by absence.
//
// DuplicatePolicy decides what happens when a record's addresses resolve to a
// contact that ANOTHER record of the same entity is already linked to:
//
//   - reject: the write fails with 409 record_duplicate (data-service resolves
//     the binding inline, inside the write transaction).
//   - flag (default): the write succeeds; the binding is resolved
//     asynchronously and data-service opens a record_duplicate the record
//     page surfaces until it is dismissed.
type ContactBinding struct {
	DuplicatePolicy conversationmodel.RecordDuplicatePolicy `json:"duplicate_policy"`
}

// WithDefaults fills an unset policy with flag — the stored shape.
func (binding ContactBinding) WithDefaults() ContactBinding {
	if binding.DuplicatePolicy == "" {
		binding.DuplicatePolicy = conversationmodel.RecordDuplicatePolicyFlag
	}
	return binding
}

// EffectiveDuplicatePolicy is the policy a write applies: the stored value,
// or flag when the entity predates the setting.
func (binding ContactBinding) EffectiveDuplicatePolicy() conversationmodel.RecordDuplicatePolicy {
	return binding.WithDefaults().DuplicatePolicy
}

// ValidateContactBinding accepts an unset policy (defaulted to flag by the
// writer) or one of the known values.
func ValidateContactBinding(binding ContactBinding) error {
	if binding.DuplicatePolicy == "" || binding.DuplicatePolicy.IsValid() {
		return nil
	}
	return fmt.Errorf("unknown duplicate_policy %q (expected one of %s, %s)",
		binding.DuplicatePolicy, conversationmodel.RecordDuplicatePolicyReject, conversationmodel.RecordDuplicatePolicyFlag)
}

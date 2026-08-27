package metamodel

import "slices"

// Attribute-level permissions: restricting a FIELD rather than a record.
//
// A salary column readable only by the people team, a margin column writable
// only by finance. Orthogonal to instance-level access — a caller who may see a
// record may still be blocked from one of its fields, and the two checks compose
// without knowing about each other.
//
// Stored as JSONB inside the attribute definition, so adding a restriction needs
// no migration. nil means TODAY'S BEHAVIOUR: unrestricted.
//
// Deliberately an application-layer concept rather than an FGA one. FGA has no
// field concept, and modelling attributes as objects would explode the tuple
// count by the number of attributes times the number of entities.

// AttributeAccessRule is an allow-list of roles and teams.
//
// FAILS CLOSED in two ways worth stating: an empty rule (no roles, no teams)
// matches nobody rather than everybody, and an unknown role or team slug simply
// never matches — a typo removes access, it does not grant it.
type AttributeAccessRule struct {
	// Roles are role slugs, matched against the caller's resolved roles.
	Roles []string `json:"roles,omitempty"`
	// Teams are team slugs, matched against the caller's team closure — so a
	// member of a CHILD team is matched by a rule naming an ancestor, the same
	// roll-up that governs grants.
	Teams []string `json:"teams,omitempty"`
}

// IsSatisfiedBy reports whether a caller holding these roles and team principals
// passes the rule.
//
// teamPrincipals are the caller's rendered team strings ("team:<orgId>/<slug>"),
// which is what the closure carries; orgId is needed to render the rule's bare
// slugs into the same form so a same-named team in another org cannot match.
func (rule *AttributeAccessRule) IsSatisfiedBy(roles []string, teamPrincipals []string, orgId string) bool {
	if rule == nil {
		// No rule means no restriction.
		return true
	}

	for _, required := range rule.Roles {
		if slices.Contains(roles, required) {
			return true
		}
	}
	for _, required := range rule.Teams {
		if slices.Contains(teamPrincipals, "team:"+orgId+"/"+required) {
			return true
		}
	}

	// An empty rule reaches here and denies, which is the intended reading of
	// "restricted to nobody in particular" — it is a lock, not a wildcard.
	return false
}

// AttributeRestrictions gates reading and writing one attribute independently.
//
// A nil arm is unrestricted, so `Write` alone yields a field everyone can see
// and only some can change — the common case for a status or approval field.
type AttributeRestrictions struct {
	Read  *AttributeAccessRule `json:"read,omitempty"`
	Write *AttributeAccessRule `json:"write,omitempty"`
}

// IsEmpty reports whether the restrictions carry no rule at all, which is
// indistinguishable from being absent.
func (restrictions *AttributeRestrictions) IsEmpty() bool {
	return restrictions == nil || (restrictions.Read == nil && restrictions.Write == nil)
}

// CanReadAttribute reports whether a caller may see this attribute's value.
func CanReadAttribute(attribute Attribute, roles []string, teamPrincipals []string, orgId string) bool {
	if attribute.Restrictions == nil {
		return true
	}
	return attribute.Restrictions.Read.IsSatisfiedBy(roles, teamPrincipals, orgId)
}

// CanWriteAttribute reports whether a caller may set this attribute's value.
//
// Read access is NOT implied by write here — the two arms are independent, and a
// write-only field (settable but not readable) is a legitimate, if unusual,
// configuration. The read mask handles visibility separately.
func CanWriteAttribute(attribute Attribute, roles []string, teamPrincipals []string, orgId string) bool {
	if attribute.Restrictions == nil {
		return true
	}
	return attribute.Restrictions.Write.IsSatisfiedBy(roles, teamPrincipals, orgId)
}

// HasRestrictedAttributes reports whether ANY attribute carries a restriction.
//
// The short-circuit that keeps the common case free: an entity with no
// restrictions never resolves the caller's roles, which would otherwise be a
// second remote call on every request.
func HasRestrictedAttributes(attributes []Attribute) bool {
	for _, attribute := range attributes {
		if !attribute.Restrictions.IsEmpty() {
			return true
		}
	}
	return false
}

// RestrictedAttributeNames returns the attributes a caller may NOT read.
// Used by the query guard to reject filters and sorts that would otherwise turn
// a hidden field into an oracle.
func RestrictedAttributeNames(attributes []Attribute, roles []string, teamPrincipals []string, orgId string) []string {
	var hidden []string
	for _, attribute := range attributes {
		if !CanReadAttribute(attribute, roles, teamPrincipals, orgId) {
			hidden = append(hidden, attribute.Name)
		}
	}
	return hidden
}

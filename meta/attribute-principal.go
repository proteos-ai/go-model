package metamodel

import "go.proteos.ai/model/common"

// PrincipalAttributeMeta holds the metadata for an attribute of type
// `principal`.
//
// A principal attribute stores a reference to anything that can HOLD ACCESS —
// a user (person, agent or api client), a team, or the whole organization. The
// value is the composite common.PrincipalRef `{ type, id }`, stored as JSONB
// and filtered on the nested `id` (`data->'field'->>'id'`), exactly like a
// `user` attribute.
//
// Why it is a separate type from `user` rather than a widening of it:
// `user` means "a person did this" — it is authorship, and its value can never
// be a team. `principal` means "this party may be granted access", which is a
// different question with a strictly larger vocabulary. Widening `user` would
// have made every existing user attribute silently accept a team, including the
// created_by / updated_by audit columns where a team makes no sense.
//
// A principal attribute whose meta declares `grants` is a GRANTING attribute:
// setting `deals.account_manager` to a principal confers access directly, with
// the business field acting as the ACL. There is no separate ownership
// attribute — an attribute granting `share` IS what "owner" means.
type PrincipalAttributeMeta struct {
	Description string `json:"description,omitempty"`

	// Grants turns the attribute into a GRANTING attribute: the verbs its value
	// confers on the record — any of read, write, delete, share.
	//
	// `deals.account_manager` with grants ["read","write"] means setting it to a
	// principal gives that principal access with no share call — the business
	// field IS the ACL, which is how access is usually meant to work in a CRM.
	//
	// `share` is the right to hand access on: share the record, set other
	// granting attributes, revoke shares. ["read","write","delete","share"] is
	// full ownership. It is the ONLY route to `share` on a record — it is
	// neither an entity-level permission nor storable as a share — so a
	// write-share recipient can never re-share, and can never set a
	// share-granting attribute to themselves.
	//
	// Empty (the default) means the attribute is an ordinary reference that
	// confers nothing, so an existing principal attribute cannot start granting
	// access by accident.
	//
	// Writing a granting attribute runs the SAME ValidateShareGrant as the share
	// endpoint. Without that guard the attribute is a backdoor and the endpoint's
	// rules are decorative.
	Grants []string `json:"grants,omitempty"`

	// AllowedTypes narrows which principal kinds this attribute accepts. Empty
	// means every kind — except that an attribute conferring `share` never
	// accepts `org`, listed or not (see AllowsType).
	AllowedTypes []common.PrincipalType `json:"allowed_types,omitempty"`
}

// GrantShare is the verb that makes a granting attribute an OWNER attribute.
// Mirrors grants.PermissionShare; spelled here because model/meta must not
// import auth.
const GrantShare = "share"

// AllowsType reports whether this attribute accepts a principal of the given
// kind. An empty AllowedTypes accepts everything.
//
// One rule sits above the list: an attribute that confers `share` refuses
// `org` unconditionally. An org-held `share` would hand management rights —
// re-granting, revoking shares — to every scoped member of the organization,
// which is never what "owner" can mean. Org-wide VISIBILITY is expressed by
// sharing read with the org principal instead. Metadata-service rejects such a
// config at save; this is the runtime backstop for a schema that predates it.
func (meta PrincipalAttributeMeta) AllowsType(principalType common.PrincipalType) bool {
	if principalType == common.PrincipalTypeOrg && meta.ConfersShare() {
		return false
	}
	if len(meta.AllowedTypes) == 0 {
		return true
	}
	for _, allowed := range meta.AllowedTypes {
		if allowed == principalType {
			return true
		}
	}
	return false
}

// IsGranting reports whether setting this attribute confers access.
func (meta PrincipalAttributeMeta) IsGranting() bool { return len(meta.Grants) > 0 }

// ConfersShare reports whether this attribute's principal may hand access on —
// the property that used to be called ownership.
func (meta PrincipalAttributeMeta) ConfersShare() bool {
	for _, verb := range meta.Grants {
		if verb == GrantShare {
			return true
		}
	}
	return false
}

// PrincipalMetaOf returns the principal meta for an attribute, in whatever shape
// the decoder that produced it happened to leave behind.
//
// Attribute.Meta is `any`, so the shape depends entirely on how the entity was
// obtained. In-process construction gives a struct or a pointer to one; anything
// that crossed JSON gives a map[string]any — and CROSSING JSON IS THE NORMAL
// CASE. Entity schemas are persisted as JSONB and data-service fetches them from
// metadata-service over HTTP, so by the time a granting attribute is evaluated
// on a real request its meta is always a map.
//
// A bare type switch therefore matched only the in-memory shape, which is also
// the only shape unit tests build — it passed every test while returning false
// for every production request, silently disabling granting attributes and the
// allowed_types narrowing. ParseMetaAs handles all three shapes; use it, and do
// not reintroduce a type switch here.
func PrincipalMetaOf(attribute Attribute) (PrincipalAttributeMeta, bool) {
	if attribute.Type != AttributeTypePrincipal {
		return PrincipalAttributeMeta{}, false
	}
	meta := ParseMetaAs[PrincipalAttributeMeta](attribute.Meta)
	if meta == nil {
		return PrincipalAttributeMeta{}, false
	}
	return *meta, true
}

// HasGrantingAttributes reports whether any attribute confers access. The
// short-circuit that keeps the common case free.
func HasGrantingAttributes(attributes []Attribute) bool {
	for _, attribute := range attributes {
		if meta, ok := PrincipalMetaOf(attribute); ok && meta.IsGranting() {
			return true
		}
	}
	return false
}

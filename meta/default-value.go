package metamodel

// DefaultValueCurrentUser is the `type` of the one default-value SENTINEL the
// platform recognises: `"default_value": {"type": "current_user"}` on a `user`
// or `principal` attribute means "whoever is writing the record". It is
// resolved by data-service at write time from the caller's identity — never
// stored as-is — so a `created_by`-shaped business field (an owner attribute
// granting `share`, an `assignee`, a `requested_by`) fills itself without a
// hook and without the client knowing who it is.
//
// An object rather than a bare string because a bare string IS a valid user
// value (the client sends bare ids), so `"current_user"` could not be told
// apart from a user whose id happens to be that string. The object shape
// mirrors the ref values these attributes store, minus an id.
const DefaultValueCurrentUser = "current_user"

// CurrentUserDefault renders the sentinel in its wire shape.
func CurrentUserDefault() map[string]any {
	return map[string]any{"type": DefaultValueCurrentUser}
}

// IsCurrentUserDefault reports whether a default_value is the current-user
// sentinel, in any shape it can arrive in: the map a JSON body or JSONB row
// decodes to, or the in-process map CurrentUserDefault builds. Anything else
// — a bare id, a full {type, id} ref, nil — is an ordinary literal default.
func IsCurrentUserDefault(value any) bool {
	object, ok := value.(map[string]any)
	if !ok || len(object) != 1 {
		return false
	}
	kind, _ := object["type"].(string)
	return kind == DefaultValueCurrentUser
}

// AcceptsCurrentUserDefault reports whether an attribute type can carry the
// sentinel: only the two identity-valued types.
func AcceptsCurrentUserDefault(attributeType AttributeType) bool {
	return attributeType == AttributeTypeUser || attributeType == AttributeTypePrincipal
}

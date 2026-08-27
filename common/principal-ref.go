package common

import (
	"fmt"
	"strings"
)

// PrincipalType is the kind of party that can hold access.
//
// It is a strict superset of UserType's vocabulary, so every existing
// {"type":"person","id":"…"} value stays valid when read as a PrincipalRef —
// which is what lets a principal attribute be seeded from `created_by` with a
// straight copy.
type PrincipalType string

const (
	// PrincipalTypePerson is a human user.
	PrincipalTypePerson PrincipalType = "person"
	// PrincipalTypeAgent is an AI agent acting on its own.
	PrincipalTypeAgent PrincipalType = "agent"
	// PrincipalTypeApi is a non-interactive API / service client.
	PrincipalTypeApi PrincipalType = "api"
	// PrincipalTypeTeam is a team; its Id is the team slug.
	PrincipalTypeTeam PrincipalType = "team"
	// PrincipalTypeOrg is the whole organization; its Id is the org id.
	//
	// This is what makes the access rule total. "Visible to everyone in the org"
	// is not a flag or an enum value — it is an org entry sitting in the
	// resource's grants, matched by the same set intersection as every other
	// principal.
	PrincipalTypeOrg PrincipalType = "org"
)

// PrincipalRef references anything that can hold access: a user (person, agent
// or api client), a team, or a whole organization.
//
// Id means different things per type — a user uuid, a team slug, an org id —
// exactly as UserRef's Id already does across its own types. The type
// disambiguates.
type PrincipalRef struct {
	Type PrincipalType `json:"type"`
	Id   string        `json:"id"`
}

// principal string prefixes. person, agent and api ALL render as "user:<id>":
// FGA has a single `user` subject type and the platform-user id namespace is
// shared across the three, so the distinction is authorship metadata, not an
// access-control distinction. Two agents can no more collide with a person's id
// than two people can with each other.
const (
	principalPrefixUser = "user"
	principalPrefixTeam = "team"
	principalPrefixOrg  = "org"
)

// PermissionSeparator splits a principal from the verb it is granted in an
// access_grants entry: "user:<uuid>#read".
//
// `#` is safe because the principal grammar is `type:id`, with `/` for org
// qualification, and team slugs reject `#` at validation (see IsValid).
//
// Note precisely what that validation buys. ParseEntry splits on the LAST
// separator, so it would in fact survive a `#` inside a slug. The rule exists
// because the grammar is unambiguous ONLY under last-separator parsing, and the
// obvious wrong implementation — splitting on the first — fails silently: it
// yields a prefix of the real principal, which in an array-overlap comparison
// matches nothing and quietly drops access instead of erroring. Forbidding the
// character means no future parser can get it wrong.
//
// This package is the ONLY place entries are built or parsed; ad-hoc string
// handling of an entry anywhere else is a review reject.
const PermissionSeparator = "#"

// String renders the principal in its canonical, org-qualified form:
//
//	person | agent | api  ->  "user:<uuid>"
//	team                  ->  "team:<orgId>/<slug>"
//	org                   ->  "org:<orgId>"
//
// Teams are org-qualified for the same reason FGA object ids are: every row is
// already org-scoped, so a bare slug would work in practice, but qualifying it
// means a query that loses its `WHERE org_id` fails closed instead of silently
// matching a same-named team in another tenant.
//
// An empty orgId yields an empty string for the types that need it, so a caller
// that forgot to scope produces an entry that matches nothing rather than one
// that matches everything.
func (ref PrincipalRef) String(orgId string) string {
	switch ref.Type {
	case PrincipalTypePerson, PrincipalTypeAgent, PrincipalTypeApi:
		if ref.Id == "" {
			return ""
		}
		return principalPrefixUser + ":" + ref.Id
	case PrincipalTypeTeam:
		if orgId == "" || ref.Id == "" {
			return ""
		}
		return principalPrefixTeam + ":" + orgId + "/" + ref.Id
	case PrincipalTypeOrg:
		// An org principal names itself; Id is the org id.
		if ref.Id == "" {
			return ""
		}
		return principalPrefixOrg + ":" + ref.Id
	default:
		return ""
	}
}

// Entry builds one perm-qualified access_grants entry, e.g. "user:<uuid>#read".
// Empty when the principal cannot be rendered (see String).
func (ref PrincipalRef) Entry(orgId string, permission string) string {
	principal := ref.String(orgId)
	if principal == "" || permission == "" {
		return ""
	}
	return principal + PermissionSeparator + permission
}

// IsValid reports whether the ref is well formed enough to be stored or granted.
func (ref PrincipalRef) IsValid() bool {
	if ref.Id == "" {
		return false
	}
	switch ref.Type {
	case PrincipalTypePerson, PrincipalTypeAgent, PrincipalTypeApi, PrincipalTypeOrg:
		return true
	case PrincipalTypeTeam:
		// A `#` in a slug would break the entry grammar irrecoverably: parsing
		// "team:acme/sa#les#read" is ambiguous. Rejected at the type so no
		// storage path has to remember.
		return !strings.Contains(ref.Id, PermissionSeparator)
	default:
		return false
	}
}

// PrincipalRefFromUser widens a UserRef into a PrincipalRef. UserType's values
// are a subset of PrincipalType's, so this is a straight copy — the conversion
// that lets a principal attribute be seeded from a user ref.
func PrincipalRefFromUser(ref UserRef) PrincipalRef {
	return PrincipalRef{Type: PrincipalType(ref.Type), Id: ref.Id}
}

// ParseEntry splits an access_grants entry back into its principal string and
// permission. The principal is returned as the rendered string (e.g.
// "team:acme/sales"), not a PrincipalRef, because that is what set intersection
// compares — reconstructing a ref would need the org id back out of the string.
func ParseEntry(entry string) (principal string, permission string, err error) {
	index := strings.LastIndex(entry, PermissionSeparator)
	if index <= 0 || index == len(entry)-1 {
		return "", "", fmt.Errorf("malformed access grant entry %q: expected <principal>%s<permission>", entry, PermissionSeparator)
	}
	return entry[:index], entry[index+1:], nil
}

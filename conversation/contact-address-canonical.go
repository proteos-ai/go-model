package conversationmodel

import (
	"regexp"
	"slices"
	"strings"
)

// The canonical forms of the record-capable address kinds live here, in the
// shared model, because two services must apply the SAME rule: data-service
// canonicalizes `contact-address` attribute values on record writes, and
// conversation-service canonicalizes every sighting before it probes the
// contact_address table. A value canonicalized on one side that would not
// canonicalize identically on the other would silently break the join.

// CanonicalizeEmail trims and lowercases. ok=false when the result is not a
// plausible address (must contain a non-leading @). Deliberately NO
// provider-specific rewriting (gmail dot/plus stripping) — plus-tag variants
// are a future dedup-candidate signal, never identity equality.
func CanonicalizeEmail(raw string) (string, bool) {
	email := strings.ToLower(strings.TrimSpace(raw))
	at := strings.Index(email, "@")
	if at <= 0 || at == len(email)-1 {
		return "", false
	}
	return email, true
}

// CanonicalizePhone strips separator noise and requires a leading + followed
// by 5–15 digits (E.164 — connectors deliver country-coded numbers; adopt a
// phone library only if a channel ever delivers national formats).
func CanonicalizePhone(raw string) (string, bool) {
	cleaned := strings.Map(func(r rune) rune {
		switch r {
		case ' ', '-', '(', ')', '.', '/':
			return -1
		}
		return r
	}, strings.TrimSpace(raw))
	if !strings.HasPrefix(cleaned, "+") {
		return "", false
	}
	digits := cleaned[1:]
	if len(digits) < 5 || len(digits) > 15 {
		return "", false
	}
	for _, r := range digits {
		if r < '0' || r > '9' {
			return "", false
		}
	}
	return cleaned, true
}

// linkedinPublicIdentifier is the shape of a LinkedIn public identifier (the
// `/in/<slug>` segment of a profile URL) once lowercased: 3–100 chars of
// letters, digits and hyphens, never starting or ending with a hyphen.
var linkedinPublicIdentifier = regexp.MustCompile(`^[a-z0-9][a-z0-9-]{1,98}[a-z0-9]$`)

// linkedinProviderIdPrefixes mark a LinkedIn profile id the way Unipile
// documents it: "a provider_id starting with ACo for classic, ACw for Sales
// Navigator, AE for Recruiter". The id is case-sensitive URL-safe base64.
var linkedinProviderIdPrefixes = []string{"ACo", "ACw", "AE"}

// linkedinProviderIdShape is the character set of a profile id — no dots,
// slashes or spaces, so a URL, an email or a sentence never passes.
var linkedinProviderIdShape = regexp.MustCompile(`^[A-Za-z0-9_-]+$`)

// IsLinkedinProviderId reports whether value has the shape of a LinkedIn
// profile id (the `linkedin_id` kind) as opposed to a public identifier — the
// ONE discriminator for free-text LinkedIn identifiers (a typed target, a
// webhook sender, a record value). The prefix match is case-sensitive, so a
// public identifier typed with exactly that casing ("AErial-photo") is refused
// as ambiguous; public identifiers are case-insensitive on LinkedIn, so the
// lowercase form always works.
func IsLinkedinProviderId(value string) bool {
	value = strings.TrimSpace(value)
	if !linkedinProviderIdShape.MatchString(value) {
		return false
	}
	for _, prefix := range linkedinProviderIdPrefixes {
		if strings.HasPrefix(value, prefix) && len(value) > len(prefix) {
			return true
		}
	}
	return false
}

// CanonicalizeLinkedin accepts a profile URL (https://www.linkedin.com/in/Ada-Lovelace/?trk=x,
// linkedin.com/in/ada-lovelace, /in/ada-lovelace) or a bare public identifier
// and returns the lowercased public identifier. ok=false when no slug survives
// (company pages, empty input) or when the value — bare or as the `/in/`
// segment of a URL, the way LinkedIn itself renders slug-less profiles — is a
// profile id: those belong to the `linkedin_id` kind and must never be
// lowercased into a slug that matches nothing.
func CanonicalizeLinkedin(raw string) (string, bool) {
	segment, ok := linkedinProfileSegment(strings.TrimSpace(raw))
	if !ok || segment == "" || IsLinkedinProviderId(segment) {
		return "", false
	}
	lower := strings.ToLower(segment)
	if !linkedinPublicIdentifier.MatchString(lower) {
		return "", false
	}
	return lower, true
}

// linkedinProfileSegment strips scheme, LinkedIn host, `in/`, query/fragment
// and surrounding slashes from a profile reference while PRESERVING the
// segment's casing (the provider-id check is case-sensitive). ok=false for a
// bare LinkedIn host. The walk runs on a lowercased copy and slices the
// original by the same offsets; a value whose lowercase form changes byte
// length (non-ASCII casing) is handled on the lowercase copy alone — no
// profile id contains such characters.
func linkedinProfileSegment(value string) (string, bool) {
	lower := strings.ToLower(value)
	if len(lower) != len(value) {
		value = lower
	}
	for _, scheme := range []string{"https://", "http://"} {
		if strings.HasPrefix(lower, scheme) {
			value, lower = value[len(scheme):], lower[len(scheme):]
			break
		}
	}
	if slash := strings.Index(lower, "/"); slash >= 0 {
		if isLinkedinHost(lower[:slash]) {
			value, lower = value[slash+1:], lower[slash+1:]
		}
	} else if isLinkedinHost(lower) {
		return "", false
	}
	if cut := strings.IndexAny(lower, "?#"); cut >= 0 {
		value, lower = value[:cut], lower[:cut]
	}
	value, lower = strings.Trim(value, "/"), strings.Trim(lower, "/")
	if strings.HasPrefix(lower, "in/") {
		value = value[len("in/"):]
	}
	return strings.Trim(value, "/"), true
}

func isLinkedinHost(host string) bool {
	return host == "linkedin.com" || strings.HasSuffix(host, ".linkedin.com")
}

// RecordContactAddressKinds are the address kinds a `contact-address` record
// attribute may declare — the scope-less kinds a human can type. Provider-id
// kinds (linkedin_id, slack, messenger, telegram, …) are minted by connectors
// only.
var RecordContactAddressKinds = []ContactAddressKind{
	ContactAddressKindEmail,
	ContactAddressKindPhone,
	ContactAddressKindLinkedinPublicIdentifier,
}

// IsRecordContactAddressKind reports whether kind may back a record attribute.
func IsRecordContactAddressKind(kind ContactAddressKind) bool {
	return slices.Contains(RecordContactAddressKinds, kind)
}

// CanonicalizeContactAddress is the record-attribute dispatcher: the ONE rule
// data-service (attribute validation) and conversation-service (resolution)
// both apply. ok=false for an unsupported kind or a value that does not
// canonicalize for its kind.
func CanonicalizeContactAddress(kind ContactAddressKind, raw string) (string, bool) {
	switch kind {
	case ContactAddressKindEmail:
		return CanonicalizeEmail(raw)
	case ContactAddressKindPhone:
		return CanonicalizePhone(raw)
	case ContactAddressKindLinkedinPublicIdentifier:
		return CanonicalizeLinkedin(raw)
	default:
		return "", false
	}
}

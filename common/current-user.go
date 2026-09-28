package common

// CurrentUser is the sentinel that stands for the signed-in user AS A VALUE —
// "records owned by whoever is looking". Wherever a string value is compared
// against stored data (a record filter's value, a page condition, a
// component prop, a calendar binding, a string literal in the records SQL
// dialect) this token is swapped for the caller's user id per request instead
// of being stored or compared literally.
//
// The `$` prefix is what makes the rule universal: no user id, slug, email or
// ordinary text starts with `$`, so the token is substituted on ANY field
// without consulting the attribute type, and the bare Postgres keyword
// `current_user` (the database role) can never be mistaken for it. The
// write-time default value keeps its own object shape
// (`{"type": "current_user"}`, meta.DefaultValueCurrentUser) — that slot is
// typed already.
const CurrentUser = "$current_user"

// IsCurrentUser reports whether a value is the sentinel. Exact match — no
// trimming or case folding.
func IsCurrentUser(value string) bool {
	return value == CurrentUser
}

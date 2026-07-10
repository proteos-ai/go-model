package common

// PlatformUserId is the sentinel id stored in created_by/updated_by for writes
// that did not originate from a real user — bootstrap (the first org, admin
// role, initial user, role permissions) and unauthenticated/system fallbacks.
// Audit fields are otherwise a real user id (a UUID string); this sentinel is
// deliberately not a UUID so it can never collide with a real user.
const PlatformUserId = "platform"

// SystemUserRef is the UserRef stamped on system/bootstrap writes that have no
// real user behind them. The chip resolves the "platform" id to a "Platform"
// label regardless of Type.
func SystemUserRef() UserRef {
	return UserRef{Type: UserTypePerson, Id: PlatformUserId}
}

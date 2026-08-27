package common

const (
	ContextTagSource     = "source"
	ContextTagToken      = "token"
	ContextTagOrgId      = "orgId"
	ContextTagAuthMethod = "authMethod"
	// ContextTagPlatformRoles carries the token's platform roles, mirrored off
	// the claims so a reader needs no second lookup. Only populated on the JWT
	// path — see auth.PlatformAdminFromToken.
	ContextTagPlatformRoles = "platformRoles"
)

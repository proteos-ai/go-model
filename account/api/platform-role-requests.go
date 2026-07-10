package accountapi

// GrantPlatformRoleRequest is the body for POST /v1/platform-roles — grant a
// user a platform-level (org-independent) role, e.g. "admin".
type GrantPlatformRoleRequest struct {
	UserId string `json:"user_id" form:"user_id" validate:"required"`
	Role   string `json:"role" form:"role"`
}

package accountmodel

import (
	"time"

	"go.proteos.ai/model/common"
)

// PlatformRole is an org-independent role grant to a user (e.g. "admin" = a
// platform superuser able to act across all orgs). It is the source of truth for
// the OpenFGA tuple `user:<id> <role> platform:proteos` and the Auth0
// `proteos_platform_roles` claim, which are reconciled from this table.
type PlatformRole struct {
	UserId    string         `json:"user_id" sortable:""`
	Role      string         `json:"role" sortable:""`
	CreatedAt time.Time      `json:"created_at" sortable:""`
	CreatedBy common.UserRef `json:"created_by" sortable:""`
}

// PlatformRoleAdmin is the only platform role today (full cross-org bypass).
// Additional platform roles can be added as the FGA model gains capability
// relations on the `platform` type.
const PlatformRoleAdmin = "admin"

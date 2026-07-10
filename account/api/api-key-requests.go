package accountapi

import (
	"time"

	"go.proteos.ai/model/account"
)

type CreateApiKeyRequest struct {
	Name      string     `json:"name" form:"name" validate:"required"`
	OrgId     *string    `json:"org_id,omitempty" form:"org_id,omitempty"`
	ExpiresAt *time.Time `json:"expires_at,omitempty" form:"expires_at,omitempty"`
}

// CreateApiKeyResponse carries the full plaintext token exactly once — at
// creation. It is never stored or shown again.
type CreateApiKeyResponse struct {
	accountmodel.ApiKey
	Token string `json:"token"`
}

type GetManyApiKeysResponse struct {
	Data []accountmodel.ApiKey `json:"data"`
}

type VerifyApiKeyRequest struct {
	Token string `json:"token" validate:"required"`
}

// VerifyApiKeyResponse is the resolved owner identity for a valid key.
// UserType is "person" today; future non-person user types (service/api
// users) will surface here so the auth middleware can plant the right UserRef.
type VerifyApiKeyResponse struct {
	UserId     string `json:"user_id"`
	UserType   string `json:"user_type"`
	Email      string `json:"email"`
	GivenName  string `json:"given_name"`
	FamilyName string `json:"family_name"`
	OrgId      string `json:"org_id"`
}

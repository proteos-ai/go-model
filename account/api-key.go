package accountmodel

import (
	"time"

	"go.proteos.ai/model/common"
)

// ApiKey is the public shape of a user-scoped API key. The token hash is
// deliberately absent: it must never be selectable into (or serialized from)
// this model — the plaintext token exists only in the create response.
type ApiKey struct {
	Id          string         `json:"id" sortable:""`
	UserId      string         `json:"user_id" sortable:""`
	OrgId       string         `json:"org_id" sortable:""`
	Name        string         `json:"name" sortable:""`
	TokenPrefix string         `json:"token_prefix"`
	TokenSuffix string         `json:"token_suffix"`
	ExpiresAt   *time.Time     `json:"expires_at"`
	LastUsedAt  *time.Time     `json:"last_used_at"`
	CreatedAt   time.Time      `json:"created_at" sortable:""`
	CreatedBy   common.UserRef `json:"created_by"`
	UpdatedAt   time.Time      `json:"updated_at"`
	UpdatedBy   common.UserRef `json:"updated_by"`
}

package conversationapi

import (
	"go.proteos.ai/model/common"
	conversationmodel "go.proteos.ai/model/conversation"
)

// CreateConnectionRequest creates a connection in status=pending. The service
// validates ConnectorKey against the runtime connector registry (not the enum —
// a key may exist while the connector isn't configured in this environment) and
// stamps the connector's channel. Credentials are never accepted here: they are
// written by the connector's install flow (OAuth callback) via the credential
// store.
type CreateConnectionRequest struct {
	ConnectorKey conversationmodel.ConnectorKey    `json:"connector_key" validate:"required"`
	Name         string                            `json:"name" validate:"required"`
	Scope        conversationmodel.ConnectionScope `json:"scope" validate:"required"`
	Settings     map[string]any                    `json:"settings"`
}

type UpdateConnectionRequest struct {
	Name     *string         `json:"name,omitempty"`
	Settings *map[string]any `json:"settings,omitempty"`
}

// SyncConnectionRequest triggers a historical backfill for an email connection
// (POST /connections/:id/sync): fetch provider messages received within Range
// and drive each through the normal ingest funnel, so conversation filters,
// dedupe, and agent-listener dispatch apply exactly as for live mail. Progress
// lands in Connection.Settings under the sync_* keys; poll GET /connections/:id
// until sync_status leaves in_progress.
type SyncConnectionRequest struct {
	Range conversationmodel.ConnectionSyncRange `json:"range" validate:"required,oneof=30d 90d 365d all"`
}

// InstallConnectionRequest is the OPTIONAL body of POST
// /connections/:id/install. Input carries user-supplied install credentials
// for direct-install connectors (twilio-phone, aircall — connector-documented
// snake_case keys); empty for OAuth/hosted connectors, whose install is a
// browser flow. Values are secrets: never logged, never echoed back.
type InstallConnectionRequest struct {
	Input map[string]string `json:"input,omitempty"`
}

// InstallConnectionResponse is returned by POST /connections/:id/install.
// OAuth/hosted connectors return AuthorizationUrl (the browser opens it in a
// popup; completion lands on the connector's pre-auth callback route).
// Direct-install connectors complete server-side in this request and return
// IsCompleted=true with no URL — the client just refetches the connection.
type InstallConnectionResponse struct {
	AuthorizationUrl string `json:"authorization_url,omitempty"`
	IsCompleted      bool   `json:"is_completed,omitempty"`
	// Setup is a direct-install connector's one-time setup handout (callback
	// URLs to paste provider-side; a minted token shown exactly once — only
	// its hash is stored). Present only on the install response; never
	// retrievable again.
	Setup map[string]string `json:"setup,omitempty"`
}

// DeleteConnectionQuery carries the delete's one escape hatch. A delete first
// makes the connector release its provider-side registration (Recall calendar +
// the OAuth grant behind it) and 502s with connector_uninstall_failed when that
// fails, because the row is the only handle on that state. IsForced deletes the
// row anyway, accepting the provider-side leftovers as a manual cleanup — the
// is_self_domain_acknowledged precedent: a footgun you opt into explicitly.
type DeleteConnectionQuery struct {
	IsForced bool `json:"is_forced" form:"is_forced"`
}

type GetManyConnectionsQuery struct {
	ConnectorKey *string `json:"connector_key" form:"connector_key" db:"connector_key"`
	Channel      *string `json:"channel" form:"channel" db:"channel"`
	Scope        *string `json:"scope" form:"scope" db:"scope"`
	// OwnerId narrows to the user-scoped connections one person owns. Owner is
	// a common.UserRef JSONB, so the filter reads its id member — the same
	// ->>'id' shape the tone-synthesis queries use. Org-scoped connections have
	// no owner and never match.
	OwnerId *string `json:"owner_id" form:"owner_id" db:"owner->>'id'"`
	Status  *string `json:"status" form:"status" db:"status"`
	common.Pagination
	common.Sorting
}

type GetManyConnectionsResponse struct {
	Meta common.ResponseMeta            `json:"meta"`
	Data []conversationmodel.Connection `json:"data"`
}

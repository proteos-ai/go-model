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

// InstallConnectionResponse is returned by POST /connections/:id/install: the
// browser opens AuthorizationUrl in a popup; install completion lands on the
// connector's pre-auth oauth-callback ingest route.
type InstallConnectionResponse struct {
	AuthorizationUrl string `json:"authorization_url"`
}

type GetManyConnectionsQuery struct {
	ConnectorKey *string `json:"connector_key" form:"connector_key" db:"connector_key"`
	Channel      *string `json:"channel" form:"channel" db:"channel"`
	Scope        *string `json:"scope" form:"scope" db:"scope"`
	Status       *string `json:"status" form:"status" db:"status"`
	common.Pagination
	common.Sorting
}

type GetManyConnectionsResponse struct {
	Meta common.ResponseMeta            `json:"meta"`
	Data []conversationmodel.Connection `json:"data"`
}

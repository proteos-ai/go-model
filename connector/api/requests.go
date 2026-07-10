package connectorapi

import (
	"time"

	"go.proteos.ai/model/common"
	connectormodel "go.proteos.ai/model/connector"
	metamodel "go.proteos.ai/model/meta"
)

// CreateConnectionRequest creates a connection. OAuth-kind connectors start
// in status=pending and complete via the install flow; static-credential
// kinds (api_key, basic, bot_token) may supply Credentials write-only here —
// validated against the manifest's credential_kind, vault-encrypted, and the
// connection goes straight to active. Credentials are never read back
// unredacted.
type CreateConnectionRequest struct {
	ConnectorKey string                         `json:"connector_key" validate:"required"`
	DisplayName  string                         `json:"display_name" validate:"required"`
	Scope        connectormodel.ConnectionScope `json:"scope" validate:"required"`
	Settings     map[string]any                 `json:"settings"`
	Credentials  *WriteCredentialsRequest       `json:"credentials,omitempty"`
}

// WriteCredentialsRequest is the write-only credential payload for
// static-credential kinds. Exactly the fields matching the manifest's
// credential_kind must be set.
type WriteCredentialsRequest struct {
	Kind     connectormodel.CredentialKind `json:"kind" validate:"required"`
	ApiKey   string                        `json:"api_key,omitempty"`
	Username string                        `json:"username,omitempty"`
	Password string                        `json:"password,omitempty"`
	BotToken string                        `json:"bot_token,omitempty"`
}

type UpdateConnectionRequest struct {
	DisplayName *string                  `json:"display_name,omitempty"`
	Settings    *map[string]any          `json:"settings,omitempty"`
	Credentials *WriteCredentialsRequest `json:"credentials,omitempty"`
}

// InstallConnectionResponse is returned by POST /connections/:id/install: the
// browser opens AuthorizationUrl in a popup; completion lands on the broker's
// single per-environment GET /connectors/v1/oauth/callback.
type InstallConnectionResponse struct {
	AuthorizationUrl string `json:"authorization_url"`
}

type GetManyConnectionsQuery struct {
	ConnectorKey *string `json:"connector_key" form:"connector_key" db:"connector_key"`
	Scope        *string `json:"scope" form:"scope" db:"scope"`
	Status       *string `json:"status" form:"status" db:"status"`
	common.Pagination
	common.Sorting
}

type GetManyConnectionsResponse struct {
	Meta common.ResponseMeta         `json:"meta"`
	Data []connectormodel.Connection `json:"data"`
}

// UpsertConnectorRequest is the custom-manifest deploy target (PUT
// /connectors/v1/connectors/:key), called by the module deployer. Origin is
// forced to custom server-side; keys colliding with a pre-built manifest are
// rejected. The Methods list may be empty at manifest-deploy time — custom
// method rows land on function-service afterwards and are resolved lazily at
// catalog time.
type UpsertConnectorRequest struct {
	Title          string                             `json:"title" validate:"required"`
	Description    string                             `json:"description,omitempty"`
	Icon           string                             `json:"icon,omitempty"`
	CredentialKind connectormodel.CredentialKind      `json:"credential_kind" validate:"required"`
	OAuth          *connectormodel.OAuthConfig        `json:"oauth,omitempty"`
	ConfigSchema   []metamodel.Attribute              `json:"config_schema,omitempty"`
	Methods        []connectormodel.MethodDeclaration `json:"methods,omitempty"`
	ModuleSlug     string                             `json:"module_slug,omitempty"`
}

type GetManyConnectorsQuery struct {
	Status *string `json:"status" form:"status" db:"status"`
	common.Pagination
	common.Sorting
}

type GetManyConnectorsResponse struct {
	Meta common.ResponseMeta                `json:"meta"`
	Data []connectormodel.ConnectorManifest `json:"data"`
}

// ConnectionTokenResponse is returned by POST /connections/:id/token: the
// usable credential material for the connection's kind — and ONLY usable
// material. refresh_token and OAuth app client secrets never appear here.
type ConnectionTokenResponse struct {
	Kind        connectormodel.CredentialKind `json:"kind"`
	AccessToken string                        `json:"access_token,omitempty"`
	ExpiresAt   *time.Time                    `json:"expires_at,omitempty"`
	ApiKey      string                        `json:"api_key,omitempty"`
	Username    string                        `json:"username,omitempty"`
	Password    string                        `json:"password,omitempty"`
	BotToken    string                        `json:"bot_token,omitempty"`
}

// InvokeMethodResponse wraps a method invocation result, mirroring
// function-service's InvokeActionResponse envelope.
type InvokeMethodResponse struct {
	Result any `json:"result"`
}

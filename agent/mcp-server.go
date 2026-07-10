package agentmodel

import (
	"time"

	"go.proteos.ai/model/common"
)

// McpServer is a registered MCP server an org's agents can call — the source of
// kind=mcp tools (McpBinding.server_key points here). Registering the connection
// once keeps the url + auth/secrets off every individual Tool, and lets you
// expose many tools from one server. This is the consumer-side registry (servers
// our agents call), distinct from mcp-service (the servers Proteos exposes).
// Keyed by (org_id, key).
type McpServer struct {
	OrgId      string         `json:"org_id"`
	Key        string         `json:"key" sortable:""`
	Name       string         `json:"name" sortable:""`
	ModuleSlug string         `json:"module_slug" sortable:""`
	Url        string         `json:"url"`
	Auth       McpServerAuth  `json:"auth"`
	CreatedAt  time.Time      `json:"created_at" sortable:""`
	CreatedBy  common.UserRef `json:"created_by"`
	UpdatedAt  time.Time      `json:"updated_at" sortable:""`
	UpdatedBy  common.UserRef `json:"updated_by"`
}

// McpServerAuth is the auth config for reaching an MCP server. IsSecret flags the
// bearer token as secret-managed; it is still returned on read, so callers with
// access see the real value. For type=oauth the bearer is not stored here at all —
// access/refresh tokens live in the adapter's Redis store; only the durable client
// identity + discovered endpoints persist (in OAuth). (header variant deferred.)
type McpServerAuth struct {
	Type     string          `json:"type"`            // none | bearer | oauth
	IsSecret bool            `json:"is_secret"`       // bearer token is secret-managed
	Token    string          `json:"token,omitempty"` // bearer token (type==bearer)
	OAuth    *McpServerOAuth `json:"oauth,omitempty"` // durable oauth client config (type==oauth)
}

// McpServerOAuth is the durable OAuth client configuration for an MCP server.
// Tokens are NOT stored here (they live in the adapter's Redis store); only the
// client identity + the discovery hints needed to refresh and re-authorize.
// Either DCR fills ClientId/ClientSecret on first connect (RegistrationMode="dcr"),
// or the user supplies them when the authorization server has no registration
// endpoint (RegistrationMode="manual").
type McpServerOAuth struct {
	ClientId         string   `json:"client_id,omitempty"`
	ClientSecret     string   `json:"client_secret,omitempty"`
	RegistrationMode string   `json:"registration_mode,omitempty"` // dcr | manual
	Scopes           []string `json:"scopes,omitempty"`

	// Discovered / cached authorization-server endpoints (populated on connect).
	Issuer                string `json:"issuer,omitempty"`
	AuthorizationEndpoint string `json:"authorization_endpoint,omitempty"`
	TokenEndpoint         string `json:"token_endpoint,omitempty"`
	RegistrationEndpoint  string `json:"registration_endpoint,omitempty"`
	Resource              string `json:"resource,omitempty"` // RFC 8707 canonical resource indicator
}

// McpOAuthClientCreds is the small value the McpClient port returns from
// BeginAuthorization when Dynamic Client Registration just minted a client, so the
// domain service can persist it onto the McpServer's OAuth config.
type McpOAuthClientCreds struct {
	ClientId     string `json:"client_id"`
	ClientSecret string `json:"client_secret,omitempty"`
}

// McpToolSummary is one entry of an MCP server's tools/list — the shape returned by
// the /mcp-servers/:key/tools subresource.
type McpToolSummary struct {
	Name        string `json:"name"`
	Title       string `json:"title,omitempty"`
	Description string `json:"description,omitempty"`
}

// McpConnectionState is the derived connection state of an oauth MCP server.
type McpConnectionState string

const (
	McpConnectionDisconnected McpConnectionState = "disconnected" // no stored token
	McpConnectionConnected    McpConnectionState = "connected"    // token present (refreshable)
	McpConnectionExpired      McpConnectionState = "expired"      // access expired, no usable refresh token
	McpConnectionError        McpConnectionState = "error"        // last refresh / connect failed
)

// McpConnectionStatus reports an MCP server's connection state. State is derived
// from token presence + expiry; LastConnectedAt / Error are carried alongside.
type McpConnectionStatus struct {
	Key             string             `json:"key"`
	State           McpConnectionState `json:"state"`
	LastConnectedAt *time.Time         `json:"last_connected_at,omitempty"`
	Error           string             `json:"error,omitempty"`
}

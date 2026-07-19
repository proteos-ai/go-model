package connectormodel

import (
	"time"

	"go.proteos.ai/model/common"
	metamodel "go.proteos.ai/model/meta"
)

// ConnectorManifest describes one connector TYPE (google-calendar, a
// customer's vertex-tax, …): how it authenticates, what a connection of it
// needs configured, and which methods it exposes. Pre-built manifests are
// seeded at boot from the compiled-in Go registry (OrgId "" = platform-global,
// the workflow_node_types convention); custom manifests are upserted per org
// by module deploy.
//
// Ubiquitous naming: a CONNECTOR is the integration type (manifest/adapter),
// a CONNECTION is one configured instance, a METHOD is one callable operation.
// Same words in DB/Go/wire/SDK/UI.
type ConnectorManifest struct {
	Key            string         `json:"key" sortable:""`
	OrgId          string         `json:"org_id" sortable:""`
	Title          string         `json:"title" sortable:""`
	Description    string         `json:"description,omitempty"`
	Icon           string         `json:"icon,omitempty"`
	CredentialKind CredentialKind `json:"credential_kind"`
	// OAuth carries the provider config the broker needs; required when
	// CredentialKind is oauth, nil otherwise.
	OAuth *OAuthConfig `json:"oauth,omitempty"`
	// ConfigSchema declares the user-supplied part of Connection.Settings in
	// the platform's attribute language, so the UI can render a form and the
	// service can validate.
	ConfigSchema []metamodel.Attribute `json:"config_schema,omitempty"`
	Methods      []MethodDeclaration   `json:"methods"`
	// OAuthRedirectUri is COMPUTED on API reads (never stored): the broker's
	// single per-environment callback the provider app must whitelist. Shown
	// in the connect wizard so operators can register it without digging
	// through deployment config.
	OAuthRedirectUri string          `json:"oauth_redirect_uri,omitempty"`
	Origin           ConnectorOrigin `json:"origin" sortable:""`
	// ModuleSlug names the deploying module when Origin is custom; empty for
	// pre-built connectors.
	ModuleSlug string          `json:"module_slug,omitempty"`
	Status     ConnectorStatus `json:"status" sortable:""`
	CreatedAt  time.Time       `json:"created_at" sortable:""`
	CreatedBy  common.UserRef  `json:"created_by"`
	UpdatedAt  time.Time       `json:"updated_at" sortable:""`
	UpdatedBy  common.UserRef  `json:"updated_by"`
}

// ConnectorOrigin distinguishes compiled-in Go connectors from module-deployed
// wasm connectors. It decides which method executor runs an invocation.
type ConnectorOrigin string

const (
	ConnectorOriginPreBuilt ConnectorOrigin = "pre_built"
	ConnectorOriginCustom   ConnectorOrigin = "custom"
)

// ConnectorStatus: pre-built manifests flip to inactive when the compiled-in
// registry no longer ships them (descriptorsync semantics); custom manifests
// flip on module deactivation.
type ConnectorStatus string

const (
	ConnectorStatusActive   ConnectorStatus = "active"
	ConnectorStatusInactive ConnectorStatus = "inactive"
)

// OAuthConfig is the per-provider OAuth wiring the broker executes. The OAuth
// APP credentials (client id/secret) are deliberately NOT stored here: they
// are metadata-service variables, referenced by key, resolved per org at
// install/refresh time — app credentials are variables, tokens are
// connections (Connector Platform D2).
type OAuthConfig struct {
	AuthUrl  string   `json:"auth_url"`
	TokenUrl string   `json:"token_url"`
	Scopes   []string `json:"scopes"`
	// AuthParams are extra authorization-URL query params some providers need
	// (access_type=offline, prompt=consent, tenant hints, …).
	AuthParams    map[string]string `json:"auth_params,omitempty"`
	IsPkceEnabled bool              `json:"is_pkce_enabled,omitempty"`
	// ClientIdVariableKey / ClientSecretVariableKey name the metadata-service
	// variables holding the OAuth app credentials for the org.
	ClientIdVariableKey     string `json:"client_id_variable_key"`
	ClientSecretVariableKey string `json:"client_secret_variable_key"`
	// IdentityProbe derives external_account_id + display_name after the code
	// exchange; nil means the connection keeps an empty external identity.
	IdentityProbe *IdentityProbe `json:"identity_probe,omitempty"`
}

// IdentityProbe describes where the broker reads the remote account identity
// from after a successful code exchange: claims of the returned id_token, or
// a GET on a provider endpoint with the fresh access token.
type IdentityProbe struct {
	Source IdentityProbeSource `json:"source"`
	// Url is the endpoint to GET when Source is endpoint; unused for id_token.
	Url string `json:"url,omitempty"`
	// AccountIdPath / DisplayNamePath are dot-paths into the claims/response
	// (e.g. "email", "user.real_name").
	AccountIdPath   string `json:"account_id_path"`
	DisplayNamePath string `json:"display_name_path,omitempty"`
}

// IdentityProbeSource enumerates where identity claims come from.
type IdentityProbeSource string

const (
	IdentityProbeSourceIdToken  IdentityProbeSource = "id_token"
	IdentityProbeSourceEndpoint IdentityProbeSource = "endpoint"
)

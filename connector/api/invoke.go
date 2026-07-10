package connectorapi

import (
	"encoding/json"
	"time"
)

// The connector-service ⇄ function-service contract for CUSTOM connector
// methods. connector-service resolves the connection (vault decrypt +
// broker-refreshed access token) and forwards the invocation to
// function-service's POST /functions/v1/connector-methods/:connectorKey/:method/invoke,
// which runs the module-deployed wasm with the connection context injected.
//
// Transport headers on that call:
//   - Authorization: the ORIGINAL caller's bearer, forwarded verbatim, so the
//     guest's host-fn calls authorize as the invoking user.
//   - X-Internal-Token: proves the caller is connector-service — ordinary
//     users must not reach this route and spoof the connection block.
//   - X-Invocation-Depth: cycle guard; rejected above MaxInvocationDepth.
const (
	InternalTokenHeader   = "X-Internal-Token"
	InvocationDepthHeader = "X-Invocation-Depth"
	// MaxInvocationDepth bounds connector-method → host fn → connector-method
	// re-entry. Depth 1 = a user-initiated invocation; 2 = one nested hop.
	MaxInvocationDepth = 2
)

// InvokeConnectorMethodRequest is the request body function-service receives.
// Method is the manifest method KEY (what the guest registered via
// fn.RegisterConnectorMethod) — the URL carries the action slug, which
// addresses the wasm, not the handler. The access token travels in this
// intra-cluster envelope exactly once and is never persisted or logged by
// either side.
type InvokeConnectorMethodRequest struct {
	Method     string            `json:"method"`
	Connection ConnectionContext `json:"connection"`
	Params     json.RawMessage   `json:"params"`
}

// ConnectionContext is the resolved-connection slice a custom method may see:
// usable token material and settings — structurally no refresh_token, no
// client secret.
type ConnectionContext struct {
	Id                string         `json:"id"`
	ConnectorKey      string         `json:"connector_key"`
	Scope             string         `json:"scope"`
	ExternalAccountId string         `json:"external_account_id,omitempty"`
	Settings          map[string]any `json:"settings,omitempty"`
	AccessToken       string         `json:"access_token,omitempty"`
	TokenExpiresAt    *time.Time     `json:"token_expires_at,omitempty"`
}

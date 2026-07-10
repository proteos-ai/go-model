package connectormodel

import (
	"encoding/json"
	"time"

	"go.proteos.ai/model/common"
)

// Connection is one configured integration instance for an org (or a single
// user): a Google Calendar grant installed via OAuth, an API-key integration,
// etc. ConnectorKey names WHICH connector owns it — an open string resolved
// against the connector manifest registry at runtime, since user-defined
// connectors register keys the platform has never seen at compile time.
type Connection struct {
	Id           string `json:"id"`
	OrgId        string `json:"org_id"`
	ConnectorKey string `json:"connector_key" sortable:""`
	DisplayName  string `json:"display_name" sortable:""`
	// Scope: org-wide integration vs a single user's grant. Owner is a
	// common.UserRef ({type,id}) set only when scope=user (nil for org-wide),
	// consistent with the platform's user-reference convention.
	Scope ConnectionScope `json:"scope"`
	Owner *common.UserRef `json:"owner,omitempty"`
	// ExternalAccountId is the provider-side identity of the grant (the
	// connected Google account's email, a Slack team_id, …), stamped by the
	// OAuth broker's identity probe. UNIQUE per (org_id, connector_key) where
	// non-empty — inbound webhooks recover the connection from it, and a second
	// install of the same remote account updates rather than duplicates.
	ExternalAccountId string `json:"external_account_id" sortable:""`
	// Credentials is the connector-owned secret material, a kind-discriminated
	// union (see connection-credentials.go) — nil for a pending connection.
	// At rest the whole envelope is vault-encrypted into a single TEXT column;
	// encryption/decryption happens only in the repository adapter, and
	// logic.RedactConnection masks secret fields on every API read.
	Credentials ConnectionCredentials `json:"credentials,omitempty"`
	// Settings is connector-operational state (JSONB): machine-written sync
	// cursors, watch-channel bookkeeping, plus any user-supplied configuration
	// the manifest's config_schema declares.
	Settings map[string]any   `json:"settings"`
	Status   ConnectionStatus `json:"status" sortable:""`
	// StatusDetail carries the last error (token refresh failure, revoked
	// grant) for the status badge; empty while healthy.
	StatusDetail string         `json:"status_detail,omitempty"`
	CreatedAt    time.Time      `json:"created_at" sortable:""`
	CreatedBy    common.UserRef `json:"created_by"`
	UpdatedAt    time.Time      `json:"updated_at" sortable:""`
	UpdatedBy    common.UserRef `json:"updated_by"`
}

// MarshalJSON / UnmarshalJSON round-trip Credentials through the same
// {kind, data} envelope used at rest (MarshalConnectionCredentials /
// DecodeConnectionCredentials). Without this, encoding/json marshals the
// interface fine (it just writes the concrete value) but CANNOT unmarshal
// into it — there is no way to know which concrete type to build from
// {"api_key":"********"} alone. That breaks every Go API consumer (any
// generic json.Unmarshal into Connection, e.g. the Go SDK) even though the
// wire JSON is perfectly valid; TS/JS callers never notice since they don't
// type-check the field. The alias trick overrides just this one field so the
// rest of Connection's fields keep their normal (un)marshaling.
func (connection Connection) MarshalJSON() ([]byte, error) {
	type alias Connection
	credentials, err := MarshalConnectionCredentials(connection.Credentials)
	if err != nil {
		return nil, err
	}
	return json.Marshal(struct {
		Credentials json.RawMessage `json:"credentials,omitempty"`
		alias
	}{Credentials: credentials, alias: alias(connection)})
}

func (connection *Connection) UnmarshalJSON(data []byte) error {
	type alias Connection
	wire := struct {
		Credentials json.RawMessage `json:"credentials,omitempty"`
		*alias
	}{alias: (*alias)(connection)}
	if err := json.Unmarshal(data, &wire); err != nil {
		return err
	}
	credentials, err := DecodeConnectionCredentials(wire.Credentials)
	if err != nil {
		return err
	}
	connection.Credentials = credentials
	return nil
}

// ConnectionScope distinguishes org-wide integrations from per-user grants.
type ConnectionScope string

const (
	ConnectionScopeOrg  ConnectionScope = "org"
	ConnectionScopeUser ConnectionScope = "user"
)

// ConnectionStatus is the install lifecycle. pending → active on a completed
// install; error on a failing refresh/sync (recoverable); revoked when the
// remote side invalidated the grant (terminal until re-install).
type ConnectionStatus string

const (
	ConnectionStatusPending ConnectionStatus = "pending"
	ConnectionStatusActive  ConnectionStatus = "active"
	ConnectionStatusError   ConnectionStatus = "error"
	ConnectionStatusRevoked ConnectionStatus = "revoked"
)

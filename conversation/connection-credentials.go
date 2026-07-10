package conversationmodel

import (
	"encoding/json"
	"fmt"
	"time"
)

// ConnectionCredentials is the kind-discriminated secret payload of a Connection
// (JSONB), mirroring the agentmodel.ToolBinding tagged-union pattern. Variants
// are keyed by AUTH SHAPE, not by connector — gmail, outlook and every future
// Google/Microsoft OAuth connector share OAuthCredentials, so adding such a
// connector needs no change here. nil = a pending connection with nothing
// installed yet.
//
// Redact returns a copy with secret material masked (presence preserved) for
// API reads — each variant knows which of its own fields are secret, so central
// redaction logic never has to. Kind is the stored discriminator.
type ConnectionCredentials interface {
	isConnectionCredentials()
	Kind() CredentialKind
	Redact() ConnectionCredentials
}

// CredentialKind discriminates the credential variants by authentication
// mechanism. New connectors reuse an existing kind whenever their auth shape
// matches; a new kind is added only for a genuinely new mechanism.
type CredentialKind string

const (
	// CredentialKindBotToken: a long-lived bot token + bot user id (Slack; other
	// bot-token chat integrations later).
	CredentialKindBotToken CredentialKind = "bot_token"
	// CredentialKindOAuth: an OAuth offline grant — long-lived refresh token plus
	// a cached access token (Gmail; Outlook and other Google/Microsoft connectors
	// reuse this).
	CredentialKindOAuth CredentialKind = "oauth"
	// CredentialKindHostedAccount: an account held by an aggregator under the
	// platform's own tenancy (Unipile). The platform-level API key lives in env
	// config, so the connection stores only the aggregator-side account identity —
	// nothing secret.
	CredentialKindHostedAccount CredentialKind = "hosted_account"
)

// redactedPlaceholder is what a set secret field reads as on the API — presence
// is visible, the value never leaves the service.
const redactedPlaceholder = "********"

// BotTokenCredentials backs bot-token installs (Slack). bot_user_id is not a
// secret — the UI needs it (mention triggers show it), so it survives redaction.
type BotTokenCredentials struct {
	BotToken      string `json:"bot_token,omitempty"`
	BotUserId     string `json:"bot_user_id,omitempty"`
	GrantedScopes string `json:"granted_scopes,omitempty"`
}

func (BotTokenCredentials) isConnectionCredentials() {}
func (BotTokenCredentials) Kind() CredentialKind     { return CredentialKindBotToken }

func (creds BotTokenCredentials) Redact() ConnectionCredentials {
	if creds.BotToken != "" {
		creds.BotToken = redactedPlaceholder
	}
	return creds
}

// OAuthCredentials backs OAuth-offline connectors (Gmail; Outlook/etc later).
type OAuthCredentials struct {
	RefreshToken   string     `json:"refresh_token,omitempty"`
	AccessToken    string     `json:"access_token,omitempty"`
	TokenExpiresAt *time.Time `json:"token_expires_at,omitempty"`
	// GrantedScopes records what the external side actually granted (diagnostics).
	GrantedScopes string `json:"granted_scopes,omitempty"`
}

func (OAuthCredentials) isConnectionCredentials() {}
func (OAuthCredentials) Kind() CredentialKind     { return CredentialKindOAuth }

func (creds OAuthCredentials) Redact() ConnectionCredentials {
	if creds.RefreshToken != "" {
		creds.RefreshToken = redactedPlaceholder
	}
	if creds.AccessToken != "" {
		creds.AccessToken = redactedPlaceholder
	}
	return creds
}

// HostedAccountCredentials backs aggregator-hosted connectors (the unipile-*
// family). AccountId is the aggregator-side account handle all API calls key
// on; OwnerProviderId is the connected user's provider-side identity used to
// classify message direction (a webhook sender matching it = outbound). It can
// arrive empty from the hosted-auth callback and is self-healed from webhook
// payloads. Neither field is secret, so Redact is the identity.
type HostedAccountCredentials struct {
	AccountId       string `json:"account_id,omitempty"`
	OwnerProviderId string `json:"owner_provider_id,omitempty"`
}

func (HostedAccountCredentials) isConnectionCredentials() {}
func (HostedAccountCredentials) Kind() CredentialKind     { return CredentialKindHostedAccount }

func (creds HostedAccountCredentials) Redact() ConnectionCredentials {
	return creds
}

// credentialEnvelope is the stored (JSONB) form: a self-describing {kind, data}
// wrapper so the decode switch keys on the auth mechanism, decoupled from
// connector_key. Only the persistence + API-output paths use this; the wire-in
// path never carries credentials (install writes them server-side).
type credentialEnvelope struct {
	Kind CredentialKind  `json:"kind"`
	Data json.RawMessage `json:"data"`
}

// MarshalConnectionCredentials encodes a variant into the stored envelope. nil
// (a pending connection) marshals to nil so the column keeps its '{}' default.
//
// The result is typed json.RawMessage, NOT a bare []byte, on purpose: bun binds
// a []byte parameter as bytea, and Postgres then rejects the bytea→jsonb cast
// (SQLSTATE 22P02). json.RawMessage takes bun's jsonRawMessageType path, which
// emits the bytes as a string literal that casts to jsonb cleanly.
func MarshalConnectionCredentials(credentials ConnectionCredentials) (json.RawMessage, error) {
	if credentials == nil {
		return nil, nil
	}
	data, err := json.Marshal(credentials)
	if err != nil {
		return nil, err
	}
	return json.Marshal(credentialEnvelope{Kind: credentials.Kind(), Data: data})
}

// DecodeConnectionCredentials decodes the stored envelope back into the typed
// variant. Empty / '{}' / null (a pending connection) → nil. New auth
// mechanisms add a case here; connectors reusing an existing kind do not.
func DecodeConnectionCredentials(raw []byte) (ConnectionCredentials, error) {
	trimmed := string(raw)
	if len(raw) == 0 || trimmed == "null" || trimmed == "{}" {
		return nil, nil
	}
	var envelope credentialEnvelope
	if err := json.Unmarshal(raw, &envelope); err != nil {
		return nil, err
	}
	if envelope.Kind == "" {
		return nil, nil
	}
	switch envelope.Kind {
	case CredentialKindBotToken:
		var credentials BotTokenCredentials
		if err := json.Unmarshal(envelope.Data, &credentials); err != nil {
			return nil, err
		}
		return credentials, nil
	case CredentialKindOAuth:
		var credentials OAuthCredentials
		if err := json.Unmarshal(envelope.Data, &credentials); err != nil {
			return nil, err
		}
		return credentials, nil
	case CredentialKindHostedAccount:
		var credentials HostedAccountCredentials
		if err := json.Unmarshal(envelope.Data, &credentials); err != nil {
			return nil, err
		}
		return credentials, nil
	default:
		return nil, fmt.Errorf("unknown credential kind %q", envelope.Kind)
	}
}

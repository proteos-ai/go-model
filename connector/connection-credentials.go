package connectormodel

import (
	"encoding/json"
	"fmt"
	"time"
)

// ConnectionCredentials is the kind-discriminated secret payload of a
// Connection, ported from conversationmodel's proven pattern. Variants are
// keyed by AUTH SHAPE, not by connector — every OAuth connector shares
// OAuthCredentials, so adding a connector needs no change here. nil = a
// pending connection with nothing installed yet.
//
// Redact returns a copy with secret material masked (presence preserved) for
// API reads — each variant knows which of its own fields are secret, so
// central redaction logic never has to. Kind is the stored discriminator.
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
	// CredentialKindOAuth: an OAuth offline grant — long-lived refresh token
	// plus a cached access token, refreshed exclusively by the broker.
	CredentialKindOAuth CredentialKind = "oauth"
	// CredentialKindApiKey: a single static API key.
	CredentialKindApiKey CredentialKind = "api_key"
	// CredentialKindBasic: HTTP basic auth username + password.
	CredentialKindBasic CredentialKind = "basic"
	// CredentialKindBotToken: a long-lived bot token + bot user id (Slack-style
	// integrations, once conversation connectors migrate here).
	CredentialKindBotToken CredentialKind = "bot_token"
)

// redactedPlaceholder is what a set secret field reads as on the API —
// presence is visible, the value never leaves the service.
const redactedPlaceholder = "********"

// OAuthCredentials backs OAuth-offline connectors. The broker owns refresh;
// consumers only ever see the access token (via the token endpoint or the
// injected method context), never RefreshToken.
type OAuthCredentials struct {
	RefreshToken   string     `json:"refresh_token,omitempty"`
	AccessToken    string     `json:"access_token,omitempty"`
	TokenExpiresAt *time.Time `json:"token_expires_at,omitempty"`
	// GrantedScopes records what the external side actually granted (diagnostics).
	GrantedScopes string `json:"granted_scopes,omitempty"`
}

func (OAuthCredentials) isConnectionCredentials() {}
func (OAuthCredentials) Kind() CredentialKind     { return CredentialKindOAuth }

func (credentials OAuthCredentials) Redact() ConnectionCredentials {
	if credentials.RefreshToken != "" {
		credentials.RefreshToken = redactedPlaceholder
	}
	if credentials.AccessToken != "" {
		credentials.AccessToken = redactedPlaceholder
	}
	return credentials
}

// ApiKeyCredentials backs static-API-key connectors.
type ApiKeyCredentials struct {
	ApiKey string `json:"api_key,omitempty"`
}

func (ApiKeyCredentials) isConnectionCredentials() {}
func (ApiKeyCredentials) Kind() CredentialKind     { return CredentialKindApiKey }

func (credentials ApiKeyCredentials) Redact() ConnectionCredentials {
	if credentials.ApiKey != "" {
		credentials.ApiKey = redactedPlaceholder
	}
	return credentials
}

// BasicCredentials backs HTTP-basic-auth connectors. The username is not a
// secret — the UI shows it so users can tell accounts apart.
type BasicCredentials struct {
	Username string `json:"username,omitempty"`
	Password string `json:"password,omitempty"`
}

func (BasicCredentials) isConnectionCredentials() {}
func (BasicCredentials) Kind() CredentialKind     { return CredentialKindBasic }

func (credentials BasicCredentials) Redact() ConnectionCredentials {
	if credentials.Password != "" {
		credentials.Password = redactedPlaceholder
	}
	return credentials
}

// BotTokenCredentials backs bot-token installs. bot_user_id is not a secret —
// the UI needs it, so it survives redaction.
type BotTokenCredentials struct {
	BotToken      string `json:"bot_token,omitempty"`
	BotUserId     string `json:"bot_user_id,omitempty"`
	GrantedScopes string `json:"granted_scopes,omitempty"`
}

func (BotTokenCredentials) isConnectionCredentials() {}
func (BotTokenCredentials) Kind() CredentialKind     { return CredentialKindBotToken }

func (credentials BotTokenCredentials) Redact() ConnectionCredentials {
	if credentials.BotToken != "" {
		credentials.BotToken = redactedPlaceholder
	}
	return credentials
}

// credentialEnvelope is the serialized form: a self-describing {kind, data}
// wrapper so the decode switch keys on the auth mechanism, decoupled from
// connector_key. connector-service vault-encrypts this envelope's JSON into a
// TEXT column; the wire-in path only carries credentials on create/update of
// non-OAuth kinds (write-only, never read back).
type credentialEnvelope struct {
	Kind CredentialKind  `json:"kind"`
	Data json.RawMessage `json:"data"`
}

// MarshalConnectionCredentials encodes a variant into the envelope. nil (a
// pending connection) marshals to nil.
//
// The result is typed json.RawMessage, NOT a bare []byte, on purpose: bun
// binds a []byte parameter as bytea, and Postgres then rejects the
// bytea→jsonb cast (SQLSTATE 22P02). json.RawMessage takes bun's
// jsonRawMessageType path, which emits the bytes as a string literal.
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

// DecodeConnectionCredentials decodes the envelope back into the typed
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
	case CredentialKindOAuth:
		var credentials OAuthCredentials
		if err := json.Unmarshal(envelope.Data, &credentials); err != nil {
			return nil, err
		}
		return credentials, nil
	case CredentialKindApiKey:
		var credentials ApiKeyCredentials
		if err := json.Unmarshal(envelope.Data, &credentials); err != nil {
			return nil, err
		}
		return credentials, nil
	case CredentialKindBasic:
		var credentials BasicCredentials
		if err := json.Unmarshal(envelope.Data, &credentials); err != nil {
			return nil, err
		}
		return credentials, nil
	case CredentialKindBotToken:
		var credentials BotTokenCredentials
		if err := json.Unmarshal(envelope.Data, &credentials); err != nil {
			return nil, err
		}
		return credentials, nil
	default:
		return nil, fmt.Errorf("unknown credential kind %q", envelope.Kind)
	}
}

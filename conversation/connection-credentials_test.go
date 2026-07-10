package conversationmodel

import (
	"testing"
	"time"
)

func TestConnectionCredentialsRoundTrip(t *testing.T) {
	expiry := time.Date(2026, 7, 2, 12, 0, 0, 0, time.UTC)
	cases := []struct {
		name  string
		creds ConnectionCredentials
		kind  CredentialKind
	}{
		{"bot token", BotTokenCredentials{BotToken: "xoxb", BotUserId: "UBOT", GrantedScopes: "chat:write"}, CredentialKindBotToken},
		{"oauth", OAuthCredentials{RefreshToken: "1//r", AccessToken: "ya29", TokenExpiresAt: &expiry, GrantedScopes: "gmail.readonly"}, CredentialKindOAuth},
		{"hosted account", HostedAccountCredentials{AccountId: "acc_9f2", OwnerProviderId: "4915112345678@s.whatsapp.net"}, CredentialKindHostedAccount},
	}
	for _, testCase := range cases {
		t.Run(testCase.name, func(t *testing.T) {
			raw, err := MarshalConnectionCredentials(testCase.creds)
			if err != nil {
				t.Fatalf("marshal: %v", err)
			}
			decoded, err := DecodeConnectionCredentials(raw)
			if err != nil {
				t.Fatalf("decode: %v", err)
			}
			if decoded.Kind() != testCase.kind {
				t.Fatalf("kind = %q, want %q", decoded.Kind(), testCase.kind)
			}
			// The decoded value must be the same concrete variant with intact fields.
			switch want := testCase.creds.(type) {
			case BotTokenCredentials:
				got, ok := decoded.(BotTokenCredentials)
				if !ok || got != want {
					t.Fatalf("bot-token round trip lost data: got %#v want %#v", decoded, want)
				}
			case OAuthCredentials:
				got, ok := decoded.(OAuthCredentials)
				if !ok || got.RefreshToken != want.RefreshToken || got.AccessToken != want.AccessToken ||
					got.TokenExpiresAt == nil || !got.TokenExpiresAt.Equal(*want.TokenExpiresAt) {
					t.Fatalf("oauth round trip lost data: got %#v want %#v", decoded, want)
				}
			case HostedAccountCredentials:
				got, ok := decoded.(HostedAccountCredentials)
				if !ok || got != want {
					t.Fatalf("hosted-account round trip lost data: got %#v want %#v", decoded, want)
				}
				// Nothing in a hosted account is secret — redaction is the identity.
				if got.Redact() != ConnectionCredentials(want) {
					t.Fatalf("hosted-account redaction must be the identity")
				}
			}
		})
	}
}

func TestConnectionCredentialsNilAndEmpty(t *testing.T) {
	// nil marshals to nil bytes (column keeps its '{}' default).
	raw, err := MarshalConnectionCredentials(nil)
	if err != nil || raw != nil {
		t.Fatalf("nil credentials must marshal to nil: %v %v", raw, err)
	}
	// The empty/pending forms all decode back to nil.
	for _, empty := range []string{"", "null", "{}"} {
		decoded, err := DecodeConnectionCredentials([]byte(empty))
		if err != nil {
			t.Fatalf("decode %q: %v", empty, err)
		}
		if decoded != nil {
			t.Fatalf("decode %q must be nil, got %#v", empty, decoded)
		}
	}
}

func TestDecodeConnectionCredentialsUnknownKind(t *testing.T) {
	if _, err := DecodeConnectionCredentials([]byte(`{"kind":"carrier_pigeon","data":{}}`)); err == nil {
		t.Fatalf("an unknown credential kind must error, not silently drop")
	}
}

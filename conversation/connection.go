package conversationmodel

import (
	"time"

	"go.proteos.ai/model/common"
)

// Connection is one configured integration instance for an org (or a single
// user): a Slack workspace installed via OAuth, a user's Gmail grant, etc.
// ConnectorKey names WHICH integration owns it (the routing key for Send and
// Install); Channel is the medium it carries, denormalized from the connector at
// create time so conversations can be filtered without a registry lookup.
// adhoc/meeting conversations need no connection at all, and two connections on
// the same channel (gmail + outlook → email) coexist.
type Connection struct {
	Id           string       `json:"id"`
	OrgId        string       `json:"org_id"`
	ConnectorKey ConnectorKey `json:"connector_key" sortable:""`
	Channel      Channel      `json:"channel" sortable:""`
	Name         string       `json:"name" sortable:""`
	// Scope: org-wide integration vs a single user's grant. Owner is a
	// common.UserRef ({type,id}) set only when scope=user (nil for org-wide),
	// consistent with the platform's user-reference convention (created_by,
	// the `user` attribute) rather than a bare id string.
	Scope ConnectionScope `json:"scope"`
	Owner *common.UserRef `json:"owner,omitempty"`
	// ExternalAccountId is the integration-side tenant identity (Slack team_id,
	// the connected Gmail address, …). UNIQUE per connector_key where non-empty —
	// inbound webhooks recover the org + connection from it.
	ExternalAccountId string `json:"external_account_id"`
	// Credentials is the connector-owned secret material (JSONB), a
	// kind-discriminated union (see connection-credentials.go) — nil for a
	// pending connection. Stored plaintext per current repo convention
	// (encryption-at-rest is a platform-wide follow-up); redacted by
	// logic.RedactConnection on every API read.
	Credentials ConnectionCredentials `json:"credentials,omitempty"`
	// Settings is connector-specific non-secret configuration (JSONB), e.g. the
	// Gmail label/query filter and sync-start date, watch expiry, last historyId.
	Settings map[string]any   `json:"settings"`
	Status   ConnectionStatus `json:"status" sortable:""`
	// SupportsReactions + Reactions are COMPUTED on read (never stored): whether
	// the connection's connector implements the reaction capability, and its
	// declarative descriptor (vocabulary kind, allowed set, cardinality) so the
	// UI adapts per connection. Reactions is nil when unsupported.
	SupportsReactions bool                `json:"supports_reactions"`
	Reactions         *ReactionCapability `json:"reactions,omitempty"`
	// Provider is COMPUTED on read like SupportsReactions: who operates the
	// integration mechanics (native | unipile), taken from the connector. Empty
	// when the connector is not registered in this environment.
	Provider  ConnectorProvider `json:"provider,omitempty"`
	CreatedAt time.Time         `json:"created_at" sortable:""`
	CreatedBy common.UserRef    `json:"created_by"`
	UpdatedAt time.Time         `json:"updated_at" sortable:""`
	UpdatedBy common.UserRef    `json:"updated_by"`
}

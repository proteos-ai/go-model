package conversationmodel

import (
	"time"

	"go.proteos.ai/model/common"
)

// Participant is a PERSON known on a connection — a Slack workspace member, an
// email correspondent. It is the person half of the old recipient directory
// (the venue half is Room): rows are populated by a connector directory sweep
// (Source=sync, e.g. Slack users.list) or derived from ingested messages
// (Source=ingest — the sole, scope-free source for Gmail correspondents), and
// are unique per (org, connection, external_id).
//
// A participant is what a message Sender and a reaction resolve against: the
// directory row is snapshotted inline as a ParticipantRef at ingest. User is
// the resolved platform identity (matched by email at sync/ingest time,
// best-effort); nil until a match is found.
type Participant struct {
	Id           string `json:"id"`
	OrgId        string `json:"org_id"`
	ConnectionId string `json:"connection_id"`
	// ExternalId is the connector-side identity used to address the person: a
	// Slack user id, or an email address.
	ExternalId string `json:"external_id" sortable:""`
	Name       string `json:"name" sortable:""`
	// Email is the person's address on the provider side (Slack profile email,
	// the correspondent address itself for email). It backs the platform-user
	// resolution and may be empty when the provider doesn't expose one.
	Email string `json:"email,omitempty"`
	// PlatformUser is the resolved platform identity (email match against the
	// account directory); nil when unresolved.
	PlatformUser *common.UserRef   `json:"platform_user,omitempty"`
	Metadata     map[string]any    `json:"metadata"`
	Source       ParticipantSource `json:"source"`
	// LastSeenAt is refreshed on every sync/ingest touch; sync pruning deletes
	// Source=sync rows not touched by the latest sweep.
	LastSeenAt time.Time      `json:"last_seen_at" sortable:""`
	CreatedAt  time.Time      `json:"created_at" sortable:""`
	CreatedBy  common.UserRef `json:"created_by"`
	UpdatedAt  time.Time      `json:"updated_at" sortable:""`
	UpdatedBy  common.UserRef `json:"updated_by"`
}

// Ref projects the directory row onto the inline snapshot shape carried by
// messages and reactions.
func (participant Participant) Ref() ParticipantRef {
	return ParticipantRef{
		ExternalId:   participant.ExternalId,
		Name:         participant.Name,
		Email:        participant.Email,
		PlatformUser: participant.PlatformUser,
	}
}

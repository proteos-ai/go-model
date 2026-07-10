package conversationmodel

import (
	"time"

	"go.proteos.ai/model/common"
)

// Room is an addressable VENUE on a connection — a Slack channel, a WhatsApp
// group, a meeting room later. It is the venue half of the old recipient
// directory (the person half is Participant): rows come from a connector
// directory sweep (Source=sync, e.g. Slack conversations.list) or are minted
// at ingest for a venue the sweep hasn't seen — or can't see: WhatsApp groups
// have no sweepable directory — (Source=ingest), and are unique per (org,
// connection, external_id). A room is a send target only — it is never a
// message sender and never resolves to a platform user.
type Room struct {
	Id           string `json:"id"`
	OrgId        string `json:"org_id"`
	ConnectionId string `json:"connection_id"`
	// ExternalId is the connector-side identity used to post into the room: a
	// Slack channel id.
	ExternalId string         `json:"external_id" sortable:""`
	Name       string         `json:"name" sortable:""`
	Metadata   map[string]any `json:"metadata"`
	Source     RoomSource     `json:"source"`
	// LastSeenAt is refreshed on every sync/ingest touch; sync pruning deletes
	// Source=sync rows not touched by the latest sweep.
	LastSeenAt time.Time      `json:"last_seen_at" sortable:""`
	CreatedAt  time.Time      `json:"created_at" sortable:""`
	CreatedBy  common.UserRef `json:"created_by"`
	UpdatedAt  time.Time      `json:"updated_at" sortable:""`
	UpdatedBy  common.UserRef `json:"updated_by"`
}

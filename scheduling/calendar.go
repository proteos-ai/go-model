package schedulingmodel

import (
	"time"

	"go.proteos.ai/model/common"
)

// CalendarSync is the mirror state of one calendar: the change cursor, the
// backfilled window, and the last outcome.
type CalendarSync struct {
	Cursor       string             `json:"cursor,omitempty"`
	BackfillFrom *time.Time         `json:"backfill_from,omitempty"`
	BackfillTo   *time.Time         `json:"backfill_to,omitempty"`
	LastSyncedAt *time.Time         `json:"last_synced_at,omitempty"`
	Status       CalendarSyncStatus `json:"status"`
	Error        string             `json:"error,omitempty"`
}

// CalendarWatch is the open push subscription for one calendar. Google fills
// ChannelId + ResourceId, Microsoft fills SubscriptionId. Only the HASH of the
// secret the provider echoes is stored — the plaintext exists once, at
// registration.
type CalendarWatch struct {
	ChannelId       string     `json:"channel_id,omitempty"`
	ResourceId      string     `json:"resource_id,omitempty"`
	SubscriptionId  string     `json:"subscription_id,omitempty"`
	ClientStateHash string     `json:"client_state_hash,omitempty"`
	ExpiresAt       *time.Time `json:"expires_at,omitempty"`
}

// Calendar is one provider calendar under a connection, plus the platform's
// roles for it: IsSynced (mirrored + watched), IsAvailabilitySource (its busy
// time counts against the owner's slots), IsTarget (the ONE calendar per user
// the platform writes exclusive events to), and Visibility (what colleagues
// see of its events).
type Calendar struct {
	Id                   string `json:"id"`
	OrgId                string `json:"org_id"`
	CalendarConnectionId string `json:"calendar_connection_id"`
	// Owner is the platform user whose account the calendar belongs to.
	Owner                     common.UserRef     `json:"owner"`
	ProviderCalendarId        string             `json:"provider_calendar_id"`
	Name                      string             `json:"name" sortable:""`
	Color                     string             `json:"color,omitempty"`
	Timezone                  string             `json:"timezone,omitempty"`
	IsPrimary                 bool               `json:"is_primary"`
	IsReadonly                bool               `json:"is_readonly"`
	HasIntegratedConferencing bool               `json:"has_integrated_conferencing"`
	IsSynced                  bool               `json:"is_synced" sortable:""`
	IsAvailabilitySource      bool               `json:"is_availability_source"`
	IsTarget                  bool               `json:"is_target"`
	Visibility                CalendarVisibility `json:"visibility" sortable:""`
	Sync                      CalendarSync       `json:"sync"`
	Watch                     CalendarWatch      `json:"watch"`
	CreatedAt                 time.Time          `json:"created_at" sortable:""`
	CreatedBy                 common.UserRef     `json:"created_by"`
	UpdatedAt                 time.Time          `json:"updated_at" sortable:""`
	UpdatedBy                 common.UserRef     `json:"updated_by"`
}

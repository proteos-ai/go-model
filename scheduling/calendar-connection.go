package schedulingmodel

import (
	"time"

	"go.proteos.ai/model/common"
)

// CalendarConnection is our binding to one connector-service grant (a
// `google-calendar` / `microsoft-calendar` connection of scope user): which
// provider, which account, whether the grant is still good. A user may hold
// several (work Google + personal Gmail + Outlook); the same provider account
// cannot be bound twice. Calendars hang off it; the scheduling roles
// (target, availability source, synced, visibility) live on the calendars.
type CalendarConnection struct {
	Id    string `json:"id"`
	OrgId string `json:"org_id"`
	// Owner is the platform user who holds the grant (connector-service's
	// Connection.Owner, mirrored here).
	Owner        common.UserRef   `json:"owner"`
	ConnectionId string           `json:"connection_id"`
	Provider     CalendarProvider `json:"provider" sortable:""`
	// ExternalAccountId is the provider account (the e-mail address / UPN),
	// copied from the grant so the mirror can tell the owner's own attendee row.
	ExternalAccountId string                   `json:"external_account_id" sortable:""`
	Status            CalendarConnectionStatus `json:"status" sortable:""`
	StatusDetail      string                   `json:"status_detail,omitempty"`
	LastSyncedAt      *time.Time               `json:"last_synced_at,omitempty" sortable:""`
	CreatedAt         time.Time                `json:"created_at" sortable:""`
	CreatedBy         common.UserRef           `json:"created_by"`
	UpdatedAt         time.Time                `json:"updated_at" sortable:""`
	UpdatedBy         common.UserRef           `json:"updated_by"`
}

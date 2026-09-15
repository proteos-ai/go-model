// Package schedulingapi holds the request / query / response DTOs of
// scheduling-service's HTTP API (validation tags live here, never on the
// entities).
package schedulingapi

import (
	"go.proteos.ai/model/common"
	schedulingmodel "go.proteos.ai/model/scheduling"
)

// CreateCalendarConnectionRequest binds a connector-service grant (a
// `google-calendar` / `microsoft-calendar` connection of scope user owned by
// the caller) to the calendar mirror. The service resolves the grant, lists
// the account's calendars and creates their rows with default roles.
type CreateCalendarConnectionRequest struct {
	ConnectionId string `json:"connection_id" validate:"required,max=64"`
}

// UpdateCalendarConnectionRequest re-activates a connection an admin repaired
// (e.g. after the grant was reinstalled). Only `active` can be requested; the
// other statuses are set by the sync path.
type UpdateCalendarConnectionRequest struct {
	Status *schedulingmodel.CalendarConnectionStatus `json:"status,omitempty" validate:"omitempty,oneof=active"`
}

type GetManyCalendarConnectionsQuery struct {
	Provider *string `json:"provider" form:"provider" validate:"omitempty,oneof=google microsoft"`
	Status   *string `json:"status" form:"status" validate:"omitempty,oneof=active error revoked"`
	common.Pagination
	common.Sorting
}

type GetManyCalendarConnectionsResponse struct {
	Meta common.ResponseMeta                  `json:"meta"`
	Data []schedulingmodel.CalendarConnection `json:"data"`
}

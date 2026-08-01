package eventapi

import (
	eventmodel "go.proteos.ai/model/events"
)

// GetCatalogResponse wraps the event catalog: every topic the caller may
// subscribe to, with the event types it carries.
type GetCatalogResponse struct {
	Data []eventmodel.CatalogTopic `json:"data"`
}

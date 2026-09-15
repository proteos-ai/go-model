package dataapi

import (
	"go.proteos.ai/model/common"
	datamodel "go.proteos.ai/model/data"
)

// GetManyRecordDuplicatesQuery lists the duplicate pairs of one entity
// (the entity rides the path — it is the permission target). RecordId narrows
// to the pairs a record takes part in, on either side; Status defaults to
// open — the pairs a record page shows.
type GetManyRecordDuplicatesQuery struct {
	RecordId *string                          `json:"record_id" form:"record_id"`
	Status   *datamodel.RecordDuplicateStatus `json:"status" form:"status"`
	common.Pagination
	common.Sorting
}

type GetManyRecordDuplicatesResponse struct {
	Meta common.ResponseMeta         `json:"meta"`
	Data []datamodel.RecordDuplicate `json:"data"`
}

// PublishContactObservationsResponse is the body of
// POST /data/v1/records/:entitySlug/contact-observations — ONE page of the
// entity's records replayed as record_contact_observation.created events (the
// backfill for a retyped or newly added contact-address attribute). The caller
// pages until page >= pages_total.
type PublishContactObservationsResponse struct {
	EntitySlug string `json:"entity_slug"`
	Page       int    `json:"page"`
	PageSize   int    `json:"page_size"`
	PagesTotal int    `json:"pages_total"`
	ItemsTotal int    `json:"items_total"`
	Published  int    `json:"published"`
	// Skipped counts rows with no address value — nothing to observe.
	Skipped int `json:"skipped"`
}

package discoveryapi

import (
	"time"

	"go.proteos.ai/model/common"
	discoverymodel "go.proteos.ai/model/discovery"
)

// CreateDiscoveryEvidenceRequest creates one piece of evidence
// (POST /discovery/v1/evidence). The client has already uploaded every file
// to storage-service under its own token; the request carries file ids. Shape
// rules per kind (one frame at 0 for a screenshot, a recording file for a
// recording, notes for feedback) are enforced by the domain, not by tags.
type CreateDiscoveryEvidenceRequest struct {
	Kind     discoverymodel.DiscoveryEvidenceKind     `json:"kind" validate:"required,oneof=screenshot recording feedback"`
	Category discoverymodel.DiscoveryEvidenceCategory `json:"category" validate:"required,oneof=bug idea question praise other"`
	Title    string                                   `json:"title" validate:"max=200"`
	Notes    string                                   `json:"notes" validate:"max=20000"`
	Source   discoverymodel.DiscoveryEvidenceSource   `json:"source" validate:"required,oneof=web extension"`
	Context  discoverymodel.DiscoveryEvidenceContext  `json:"context"`

	RecordingFileId     string `json:"recording_file_id" validate:"max=64"`
	RecordingDurationMs int    `json:"recording_duration_ms" validate:"min=0"`
	RecordingMimeType   string `json:"recording_mime_type" validate:"max=64"`
	RecordingSizeBytes  int64  `json:"recording_size_bytes" validate:"min=0"`
	HasAudio            bool   `json:"has_audio"`

	Links  []CreateDiscoveryEvidenceLinkRequest  `json:"links" validate:"max=10,dive"`
	Frames []CreateDiscoveryEvidenceFrameRequest `json:"frames" validate:"max=200,dive"`
}

type CreateDiscoveryEvidenceLinkRequest struct {
	Url   string `json:"url" validate:"required,max=2048"`
	Title string `json:"title" validate:"max=200"`
}

// CreateDiscoveryEvidenceFrameRequest is one annotated still. Position is
// assigned by the server from the array order.
type CreateDiscoveryEvidenceFrameRequest struct {
	AtMs                 int                                    `json:"at_ms" validate:"min=0"`
	ImageFileId          string                                 `json:"image_file_id" validate:"required,max=64"`
	AnnotatedImageFileId string                                 `json:"annotated_image_file_id" validate:"max=64"`
	Comment              string                                 `json:"comment" validate:"max=5000"`
	Strokes              []CreateDiscoveryEvidenceStrokeRequest `json:"strokes" validate:"max=500,dive"`
}

// CreateDiscoveryEvidenceStrokeRequest carries the client's stroke id when it
// has one (the review UI keys on it); the server mints one otherwise.
type CreateDiscoveryEvidenceStrokeRequest struct {
	Id          string                    `json:"id" validate:"max=64"`
	StartedAtMs int                       `json:"started_at_ms" validate:"min=0"`
	EndedAtMs   int                       `json:"ended_at_ms" validate:"min=0"`
	Color       string                    `json:"color" validate:"required,max=32"`
	Width       float64                   `json:"width" validate:"gt=0,lte=200"`
	Points      []discoverymodel.InkPoint `json:"points" validate:"required,min=1,max=5000"`
}

// UpdateDiscoveryEvidenceRequest is the submitter's partial edit: the words,
// never the capture. Frames patch comments by frame id; a frame id the record
// does not have is a 400.
type UpdateDiscoveryEvidenceRequest struct {
	Title    *string                                   `json:"title,omitempty" validate:"omitempty,max=200"`
	Notes    *string                                   `json:"notes,omitempty" validate:"omitempty,max=20000"`
	Category *discoverymodel.DiscoveryEvidenceCategory `json:"category,omitempty" validate:"omitempty,oneof=bug idea question praise other"`
	Links    *[]CreateDiscoveryEvidenceLinkRequest     `json:"links,omitempty" validate:"omitempty,max=10,dive"`
	Frames   *[]UpdateDiscoveryEvidenceFrameRequest    `json:"frames,omitempty" validate:"omitempty,max=200,dive"`
}

type UpdateDiscoveryEvidenceFrameRequest struct {
	Id      string `json:"id" validate:"required,max=64"`
	Comment string `json:"comment" validate:"max=5000"`
}

// UpdateDiscoveryEvidenceStatusRequest is the platform-admin review verb
// (PATCH /platform/evidence/:id/status).
type UpdateDiscoveryEvidenceStatusRequest struct {
	Status discoverymodel.DiscoveryEvidenceStatus `json:"status" validate:"required,oneof=new reviewed archived"`
}

// GetManyDiscoveryEvidenceQuery lists evidence in the caller's org (own rows
// under a scoped grant, every row under the unscoped read).
type GetManyDiscoveryEvidenceQuery struct {
	Kind     *string `json:"kind" form:"kind" validate:"omitempty,oneof=screenshot recording feedback"`
	Category *string `json:"category" form:"category" validate:"omitempty,oneof=bug idea question praise other"`
	Status   *string `json:"status" form:"status" validate:"omitempty,oneof=new reviewed archived"`
	Source   *string `json:"source" form:"source" validate:"omitempty,oneof=web extension"`
	common.Pagination
	common.Sorting
}

// GetManyPlatformDiscoveryEvidenceQuery is the cross-org listing for platform
// admins: the org filter plus a created window.
type GetManyPlatformDiscoveryEvidenceQuery struct {
	GetManyDiscoveryEvidenceQuery
	OrgId       *string    `json:"org_id" form:"org_id"`
	CreatedFrom *time.Time `json:"created_from" form:"created_from" time_format:"2006-01-02T15:04:05Z07:00"`
	CreatedTo   *time.Time `json:"created_to" form:"created_to" time_format:"2006-01-02T15:04:05Z07:00"`
}

type GetManyDiscoveryEvidenceResponse struct {
	Meta common.ResponseMeta                `json:"meta"`
	Data []discoverymodel.DiscoveryEvidence `json:"data"`
}

type GetDiscoveryEvidenceOrgSummaryResponse struct {
	Data []discoverymodel.DiscoveryEvidenceOrgSummary `json:"data"`
}

// DiscoveryEvidenceDownloadUrlResponse is a presigned storage-service URL the
// browser fetches directly; the service never streams bytes.
type DiscoveryEvidenceDownloadUrlResponse struct {
	Url       string    `json:"url"`
	ExpiresAt time.Time `json:"expires_at"`
}

package discoverymodel

import (
	"time"

	"go.proteos.ai/model/common"
)

// DiscoveryEvidence is one thing a user sent Proteos: a screenshot (exactly one
// frame at 0 ms), a screen recording (a file plus the frames where ink was
// drawn), or written feedback (neither). Owned by its creator, stored in the
// submitter's org, read across orgs by platform admins.
//
// Frames are the bulky part: lists omit them, the detail carries them, and the
// event payload drops them.
type DiscoveryEvidence struct {
	Id       string                    `json:"id"`
	OrgId    string                    `json:"org_id" sortable:""`
	Kind     DiscoveryEvidenceKind     `json:"kind" sortable:""`
	Category DiscoveryEvidenceCategory `json:"category" sortable:""`
	Title    string                    `json:"title"`
	Notes    string                    `json:"notes"`
	Status   DiscoveryEvidenceStatus   `json:"status" sortable:""`
	Source   DiscoveryEvidenceSource   `json:"source" sortable:""`
	Context  DiscoveryEvidenceContext  `json:"context"`
	// Recording facts; empty / zero on a screenshot or feedback. Duration is the
	// client's wall clock and is the authoritative timeline length (a webm
	// written by MediaRecorder carries none).
	RecordingFileId     string `json:"recording_file_id"`
	RecordingDurationMs int    `json:"recording_duration_ms"`
	RecordingMimeType   string `json:"recording_mime_type"`
	RecordingSizeBytes  int64  `json:"recording_size_bytes"`
	HasAudio            bool   `json:"has_audio"`
	// FrameCount is kept on the row so a list (which omits frames) can show it.
	FrameCount int                      `json:"frame_count"`
	Links      []DiscoveryEvidenceLink  `json:"links"`
	Frames     []DiscoveryEvidenceFrame `json:"frames,omitempty"`
	ReviewedAt *time.Time               `json:"reviewed_at,omitempty"`
	ReviewedBy *common.UserRef          `json:"reviewed_by,omitempty"`
	CreatedAt  time.Time                `json:"created_at" sortable:""`
	CreatedBy  common.UserRef           `json:"created_by"`
	UpdatedAt  time.Time                `json:"updated_at" sortable:""`
	UpdatedBy  common.UserRef           `json:"updated_by"`
}

// DiscoveryEvidenceContext is what the client knew when it captured: where the
// user was and on what. posthog_session_id lets Proteos jump to the session
// replay when the project records one (web only).
type DiscoveryEvidenceContext struct {
	Url              string  `json:"url"`
	Path             string  `json:"path"`
	PageTitle        string  `json:"page_title"`
	AppVersion       string  `json:"app_version"`
	ExtensionVersion string  `json:"extension_version"`
	UserAgent        string  `json:"user_agent"`
	ViewportWidth    int     `json:"viewport_width"`
	ViewportHeight   int     `json:"viewport_height"`
	DevicePixelRatio float64 `json:"device_pixel_ratio"`
	Locale           string  `json:"locale"`
	Timezone         string  `json:"timezone"`
	PosthogSessionId string  `json:"posthog_session_id"`
}

// DiscoveryEvidenceLink is a URL the user attached — a Loom share link, a
// document, a ticket.
type DiscoveryEvidenceLink struct {
	Url   string `json:"url"`
	Title string `json:"title"`
}

// DiscoveryEvidenceFrame is one annotated still: the image at at_ms (0 for a
// screenshot), the ink drawn on it as vector strokes, and the comment the
// user added at review time. annotated_image_file_id is the composite with
// the ink burned in, for thumbnails and for agents that read pixels.
type DiscoveryEvidenceFrame struct {
	Id                   string      `json:"id"`
	Position             int         `json:"position"`
	AtMs                 int         `json:"at_ms"`
	ImageFileId          string      `json:"image_file_id"`
	AnnotatedImageFileId string      `json:"annotated_image_file_id"`
	Comment              string      `json:"comment"`
	Strokes              []InkStroke `json:"strokes"`
}

// InkStroke is one pen-down → pen-up gesture. Points are normalized to the
// capture surface (0..1 on both axes); times are ms since recording start.
type InkStroke struct {
	Id          string     `json:"id"`
	StartedAtMs int        `json:"started_at_ms"`
	EndedAtMs   int        `json:"ended_at_ms"`
	Color       string     `json:"color"`
	Width       float64    `json:"width"`
	Points      []InkPoint `json:"points"`
}

type InkPoint struct {
	X float64 `json:"x"`
	Y float64 `json:"y"`
	T int     `json:"t"`
}

// DiscoveryEvidenceOrgSummary is one row of the platform-admin overview: how
// much evidence an org has sent, how much is still new, and when the latest
// arrived.
type DiscoveryEvidenceOrgSummary struct {
	OrgId    string     `json:"org_id"`
	Count    int        `json:"count"`
	NewCount int        `json:"new_count"`
	NewestAt *time.Time `json:"newest_at,omitempty"`
}

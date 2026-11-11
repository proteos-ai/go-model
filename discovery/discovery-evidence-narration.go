package discoverymodel

import (
	"time"

	"go.proteos.ai/model/common"
)

// DiscoveryEvidenceNarration is the job that turns what the sender said while
// recording into drafts: a summary for the feedback text and one comment per
// annotated frame. It starts when the recording stops, from an audio-only
// file, so the drafts can land before the sender presses Send. Evidence sent
// while the job still runs names it (narration_id) and receives the drafts in
// whatever is still empty.
//
// The transcript is input for the model calls only: it never travels on the
// wire.
type DiscoveryEvidenceNarration struct {
	Id            string                                   `json:"id"`
	OrgId         string                                   `json:"org_id"`
	Status        DiscoveryEvidenceNarrationStatus         `json:"status"`
	Summary       string                                   `json:"summary"`
	FrameComments []DiscoveryEvidenceNarrationFrameComment `json:"frame_comments"`
	CreatedAt     time.Time                                `json:"created_at"`
	CreatedBy     common.UserRef                           `json:"created_by"`
	UpdatedAt     time.Time                                `json:"updated_at"`

	AudioFileId     string                            `json:"-"`
	Frames          []DiscoveryEvidenceNarrationFrame `json:"-"`
	TranscriptionId string                            `json:"-"`
	Transcript      string                            `json:"-"`
	// EvidenceId is set once evidence was sent with this narration.
	EvidenceId string `json:"-"`
}

// DiscoveryEvidenceNarrationFrame is a frame as the client knew it at Stop:
// its own key (the evidence does not exist yet) and where it sits in the
// recording.
type DiscoveryEvidenceNarrationFrame struct {
	Key  string `json:"key"`
	AtMs int    `json:"at_ms"`
}

// DiscoveryEvidenceNarrationFrameComment is the draft for one frame, with the
// stretch of narration it was taken from.
type DiscoveryEvidenceNarrationFrameComment struct {
	Key     string `json:"key"`
	Comment string `json:"comment"`
	FromMs  int    `json:"from_ms"`
	ToMs    int    `json:"to_ms"`
}

// IsTerminal reports whether the job is done, one way or the other.
func (narration DiscoveryEvidenceNarration) IsTerminal() bool {
	return narration.Status == DiscoveryEvidenceNarrationStatusCompleted || narration.Status == DiscoveryEvidenceNarrationStatusFailed
}

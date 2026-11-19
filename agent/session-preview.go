package agentmodel

import "time"

// SessionPreviewStatus says what a session's preview shows.
type SessionPreviewStatus string

const (
	// SessionPreviewServing: the preview is served at SessionPreview.Url.
	SessionPreviewServing SessionPreviewStatus = "serving"
	// SessionPreviewNoModule: the session's sandbox runs and holds no module
	// to preview yet.
	SessionPreviewNoModule SessionPreviewStatus = "no_module"
	// SessionPreviewNotStarted: the session has no sandbox yet. Its first
	// message creates one.
	SessionPreviewNotStarted SessionPreviewStatus = "not_started"
)

// SessionPreview is the live preview of the module an agent works on in its
// session's sandbox (see RuntimeSandbox), as a page a browser shows in a frame.
type SessionPreview struct {
	Status SessionPreviewStatus `json:"status"`
	// Url is the preview's address, set while Status is serving. The address
	// is the whole credential: whoever has it sees the preview until
	// ExpiresAt, so it is shown to the session's owner only and not stored.
	Url *string `json:"url"`
	// ExpiresAt is when Url stops working. Opening the preview again gives a
	// new one.
	ExpiresAt *time.Time `json:"expires_at"`
}

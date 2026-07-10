package storagemodel

import "time"

type FileAccessTokenKind string

const (
	FileAccessTokenKindDownload FileAccessTokenKind = "download"
	FileAccessTokenKindUpload   FileAccessTokenKind = "upload"
)

// FileAccessToken is a short-lived credential for public download or upload of
// one file. By default it is single-use — the row's consumed_at is set on the
// first redeem (the CAS race gate). A download token minted with AllowsMultiUse
// is NOT consumed and stays valid until it expires, so an external consumer can
// hit the URL more than once (a HEAD/probe then GET, a redirect-follow, a
// retry). The token value itself is the sole credential at redeem time — the
// row is never returned over the API; only the built URL and expires_at leave
// the service.
type FileAccessToken struct {
	Token          string              `json:"token"`
	FileId         string              `json:"file_id"`
	Kind           FileAccessTokenKind `json:"kind"`
	OrgId          *string             `json:"org_id,omitempty"`
	CreatedBy      *string             `json:"created_by,omitempty"`
	AllowsMultiUse bool                `json:"allows_multi_use"`
	ExpiresAt      time.Time           `json:"expires_at"`
	ConsumedAt     *time.Time          `json:"consumed_at,omitempty"`
	CreatedAt      time.Time           `json:"created_at"`
}

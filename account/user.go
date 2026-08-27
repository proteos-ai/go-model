package accountmodel

import (
	"go.proteos.ai/model/common"
	"time"
)

type User struct {
	Id           string `json:"id" sortable:""`
	Email        string `json:"email" sortable:""`
	GivenName    string `json:"given_name" sortable:""`
	FamilyName   string `json:"family_name" sortable:""`
	ExternalId   string `json:"external_id" sortable:""`
	DefaultOrgId string `json:"default_org_id" sortable:""`
	// ProfileImage references the user's photo in the storage-service. Nil =
	// no photo (clients render initials). Not sortable — it is a JSONB object.
	ProfileImage *common.FileRef `json:"profile_image,omitempty"`
	CreatedAt    time.Time       `json:"created_at" sortable:""`
	CreatedBy    common.UserRef  `json:"created_by" sortable:""`
	UpdatedAt    time.Time       `json:"updated_at" sortable:""`
	UpdatedBy    common.UserRef  `json:"updated_by" sortable:""`
}

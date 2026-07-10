package storagemodel

import "time"

type File struct {
	Id             string    `json:"id" example:"0027.8f5b88e0-b875-4aa7-a73e-08995671047e"`
	OrgId          *string   `json:"org_id,omitempty" sortable:""`
	Name           string    `json:"name" example:"pdr-max-mustermann.pdf"`
	ContentType    string    `json:"content_type" example:"application/pdf"`
	CurrentVersion *Version  `json:"current_version"`
	IsDeleted      bool      `json:"is_deleted" example:"false"`
	IsPersisted    bool      `json:"is_persisted" example:"false"`
	IsLocked       bool      `json:"is_locked" example:"false"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

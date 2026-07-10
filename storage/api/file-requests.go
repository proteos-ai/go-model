package storageapi

import (
	"go.proteos.ai/model/common"
	"go.proteos.ai/model/storage"
)

type CreateFileRequest struct {
	Name        string `json:"name" form:"name" validate:"required"`
	ContentType string `json:"content_type" form:"content_type"`
}

// MintDownloadUrlRequest is the optional body of POST
// /files/:id/generate-download-url. allows_multi_use mints a token that
// survives repeated reads (and a longer TTL) for external consumers that probe
// the URL before downloading; omitted/false yields the default single-use URL.
type MintDownloadUrlRequest struct {
	AllowsMultiUse bool `json:"allows_multi_use"`
}

type UpdateFileRequest struct {
	Name           *string               `json:"name,omitempty" form:"name"`
	IsPersisted    *bool                 `json:"-,omitempty" form:"is_persisted"`
	IsLocked       *bool                 `json:"is_locked,omitempty" form:"is_locked"`
	ContentType    *string               `json:"-,omitempty"`
	CurrentVersion *storagemodel.Version `json:"-"`
}

type GetManyFilesQuery struct {
	Id          string `form:"id,omitempty" db:"id"`
	Name        string `form:"name,omitempty" db:"name"`
	Type        string `form:"type,omitempty" db:"type"`
	IsDirectory bool   `form:"is_directory,omitempty" db:"is_directory"`
	common.Pagination
	FileSorting
}

type FileSorting struct {
	Sort   common.SortDirection `json:"sort_direction" form:"sort_direction" validate:"omitempty,oneof=asc desc"`
	SortBy string               `json:"sort_by" form:"sort_by" validate:"is_file_sortable_attribute"`
}

func (self FileSorting) GetSort() common.SortDirection {
	return self.Sort
}

func (self FileSorting) GetSortBy() string {
	return self.SortBy
}

func (self FileSorting) GetSortableAttributes() map[string]string {
	return GetSortableFileAttributes()
}

func GetSortableFileAttributes() map[string]string {
	return map[string]string{
		"id":               "id",
		"name":             "name",
		"created_at":       "created_at",
		"updated_at":       "updated_at",
		"last_accessed_at": "last_accessed_at",
	}
}

package storageapi

import "go.proteos.ai/model/common"

type CreateVersionRequest struct {
	FileId      string `json:"file_id"`
	SizeInBytes int64  `json:"-,omitempty"`
}

type GetManyVersionsQuery struct {
	common.Pagination
	VersionSorting
}

type VersionSorting struct {
	Sort   common.SortDirection `json:"sort_direction" form:"sort_direction" validate:"omitempty,oneof=asc desc"`
	SortBy string               `json:"sort_by" form:"sort_by" validate:"is_version_sortable_attribute"`
}

func (self VersionSorting) GetSort() common.SortDirection {
	return self.Sort
}

func (self VersionSorting) GetSortBy() string {
	return self.SortBy
}

func (self VersionSorting) GetSortableAttributes() map[string]string {
	return GetSortableVersionAttributes()
}

func GetSortableVersionAttributes() map[string]string {
	return map[string]string{
		"number":           "number",
		"created_at":       "created_at",
		"last_accessed_at": "last_accessed_at",
	}
}

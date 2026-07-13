package metaapi

import (
	"go.proteos.ai/model/common"
	"go.proteos.ai/model/meta"
)

type CreateComponentRequest struct {
	Slug         string         `json:"slug" validate:"required"`
	Name         string         `json:"name"`
	ModuleSlug   string         `json:"module_slug"`
	Description  string         `json:"description"`
	BundleFileId string         `json:"bundle_file_id"`
	SourceFileId string         `json:"source_file_id"`
	PropsSchema  map[string]any `json:"props_schema"`
	// IsPublic opts the compiled bundle into unauthenticated serving (see
	// metamodel.Component.IsPublic). Manifest-driven: a deploy without the
	// field sets false.
	IsPublic bool `json:"is_public"`
}

type UpdateComponentRequest struct {
	Name         *string        `json:"name,omitempty"`
	Description  *string        `json:"description,omitempty"`
	BundleFileId *string        `json:"bundle_file_id,omitempty"`
	SourceFileId *string        `json:"source_file_id,omitempty"`
	PropsSchema  map[string]any `json:"props_schema,omitempty"`
	IsPublic     *bool          `json:"is_public,omitempty"`
}

type GetManyComponentsQuery struct {
	Slug       *string `json:"slug" db:"slug"`
	Name       *string `json:"name" db:"name"`
	ModuleSlug *string `json:"module_slug" db:"module_slug"`
	CreatedBy  *string `json:"created_by" db:"created_by->>'id'"`
	UpdatedBy  *string `json:"updated_by" db:"updated_by->>'id'"`
	common.Pagination
	common.Sorting
}

type GetManyComponentsResponse struct {
	Meta common.ResponseMeta   `json:"meta"`
	Data []metamodel.Component `json:"data"`
}

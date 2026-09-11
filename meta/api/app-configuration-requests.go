package metaapi

import (
	"go.proteos.ai/model/common"
	"go.proteos.ai/model/meta"
)

type CreateAppConfigurationRequest struct {
	Slug            string             `json:"slug" validate:"required"`
	ModuleSlug      string             `json:"module_slug"`
	AppSlug         string             `json:"app_slug" validate:"required"`
	ProfileSlug     string             `json:"profile_slug"`
	Home            *metamodel.AppHome `json:"home,omitempty"`
	MenuSlug        string             `json:"menu_slug,omitempty"`
	DefaultAgentKey string             `json:"default_agent_key,omitempty"`
	AgentKeys       []string           `json:"agent_keys,omitempty"`
	RecordPages     map[string]string  `json:"record_pages,omitempty"`
}

// UpdateAppConfigurationRequest is a partial update. app_slug and profile_slug
// are the row's identity (unique together) and are not updatable — delete and
// recreate to rebind. Home / AgentKeys / RecordPages are tri-state via
// common.Optional so a caller can clear them (null) as well as set them.
type UpdateAppConfigurationRequest struct {
	ModuleSlug      *string                            `json:"module_slug,omitempty"`
	Home            common.Optional[metamodel.AppHome] `json:"home" bun:"-"`
	MenuSlug        *string                            `json:"menu_slug,omitempty"`
	DefaultAgentKey *string                            `json:"default_agent_key,omitempty"`
	AgentKeys       *[]string                          `json:"agent_keys,omitempty"`
	RecordPages     *map[string]string                 `json:"record_pages,omitempty"`
}

type GetManyAppConfigurationsQuery struct {
	Slug        *string `json:"slug" db:"slug"`
	ModuleSlug  *string `json:"module_slug" db:"module_slug"`
	AppSlug     *string `json:"app_slug" db:"app_slug"`
	ProfileSlug *string `json:"profile_slug" db:"profile_slug"`
	// IsDefault selects the default rows (profile_slug = '') when true, the profile
	// overrides when false. A dedicated flag because an empty profile_slug
	// cannot travel as a query parameter (every SDK drops empty strings). No
	// db tag: the repository applies it explicitly.
	IsDefault *bool   `json:"is_default"`
	MenuSlug  *string `json:"menu_slug" db:"menu_slug"`
	CreatedBy *string `json:"created_by" db:"created_by->>'id'"`
	UpdatedBy *string `json:"updated_by" db:"updated_by->>'id'"`
	common.Pagination
	common.Sorting
}

type GetManyAppConfigurationsResponse struct {
	Meta common.ResponseMeta          `json:"meta"`
	Data []metamodel.AppConfiguration `json:"data"`
}

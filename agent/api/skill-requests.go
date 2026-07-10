package agentapi

import (
	agentmodel "go.proteos.ai/model/agent"
	"go.proteos.ai/model/common"
)

// Skills are created/updated via the multipart deploy endpoint (a bundle upload),
// not a JSON body — so there is no Create/Update request struct here, only the
// read DTOs.

type GetManySkillsQuery struct {
	Key          *string `json:"key" form:"key" db:"key"`
	Name         *string `json:"name" form:"name" db:"name"`
	ModuleSlug   *string `json:"module_slug" form:"module_slug" db:"module_slug"`
	NameContains *string `json:"name[contains]" form:"name[contains]" db:"name" op:"contains"`
	common.Pagination
	common.Sorting
}

type GetManySkillsResponse struct {
	Meta common.ResponseMeta `json:"meta"`
	Data []agentmodel.Skill  `json:"data"`
}

type GetManySkillVersionsResponse struct {
	Data []agentmodel.SkillVersion `json:"data"`
}

package agentapi

import (
	agentmodel "go.proteos.ai/model/agent"
	"go.proteos.ai/model/common"
)

type CreateAgentRequest struct {
	Key          string                 `json:"key" validate:"required"`
	Name         string                 `json:"name" validate:"required"`
	ModuleSlug   string                 `json:"module_slug"`
	Description  string                 `json:"description"`
	SystemPrompt string                 `json:"system_prompt"`
	ModelConfig  agentmodel.ModelConfig `json:"model_config"`
	Skills       []string               `json:"skills"`
	Tools        []string               `json:"tools"`
	Subagents    []string               `json:"subagents"`
	McpServers   []string               `json:"mcp_servers"`
	IsOrgDefault bool                   `json:"is_org_default"`
}

type UpdateAgentRequest struct {
	Name         *string                 `json:"name,omitempty"`
	ModuleSlug   *string                 `json:"module_slug,omitempty"`
	Description  *string                 `json:"description,omitempty"`
	SystemPrompt *string                 `json:"system_prompt,omitempty"`
	ModelConfig  *agentmodel.ModelConfig `json:"model_config,omitempty"`
	Skills       *[]string               `json:"skills,omitempty"`
	Tools        *[]string               `json:"tools,omitempty"`
	Subagents    *[]string               `json:"subagents,omitempty"`
	McpServers   *[]string               `json:"mcp_servers,omitempty"`
	IsOrgDefault *bool                   `json:"is_org_default,omitempty"`
}

type GetManyAgentsQuery struct {
	Key          *string `json:"key" form:"key" db:"key"`
	Name         *string `json:"name" form:"name" db:"name"`
	ModuleSlug   *string `json:"module_slug" form:"module_slug" db:"module_slug"`
	NameContains *string `json:"name[contains]" form:"name[contains]" db:"name" op:"contains"`
	common.Pagination
	common.Sorting
}

type GetManyAgentsResponse struct {
	Meta common.ResponseMeta `json:"meta"`
	Data []agentmodel.Agent  `json:"data"`
}

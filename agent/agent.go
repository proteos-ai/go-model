package agentmodel

import (
	"time"

	"go.proteos.ai/model/common"
)

// Agent is a configured persona: a system prompt + model config + the skills,
// tools, sub-agents and MCP servers it can use. All references are by key —
// SystemPrompt is a Prompt key, Skills/Tools/Subagents/McpServers are
// Skill/Tool/Agent/McpServer keys — and keys are immutable so references survive
// renames. Sub-agents are surfaced to the model as agents-as-tools. McpServers
// attach the whole server: every tool the server exposes becomes available to
// the agent (vs a kind=mcp Tool, which binds a single tool). Keyed by (org_id, key).
type Agent struct {
	OrgId        string      `json:"org_id"`
	Key          string      `json:"key" sortable:""`
	Name         string      `json:"name" sortable:""`
	ModuleSlug   string      `json:"module_slug" sortable:""`
	Description  string      `json:"description"`
	SystemPrompt string      `json:"system_prompt"`
	ModelConfig  ModelConfig `json:"model_config"`
	Skills       []string    `json:"skills"`
	Tools        []string    `json:"tools"`
	Subagents    []string    `json:"subagents"`
	McpServers   []string    `json:"mcp_servers"`
	// IsOrgDefault marks the single agent surfaced by default for the org (e.g. the
	// Ask Proteos assistant). At most one agent per org may carry it.
	IsOrgDefault bool           `json:"is_org_default"`
	Version      int            `json:"version"`
	CreatedAt    time.Time      `json:"created_at" sortable:""`
	CreatedBy    common.UserRef `json:"created_by"`
	UpdatedAt    time.Time      `json:"updated_at" sortable:""`
	UpdatedBy    common.UserRef `json:"updated_by"`
}

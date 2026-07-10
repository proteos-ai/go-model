package agentapi

import (
	agentmodel "go.proteos.ai/model/agent"
	"go.proteos.ai/model/common"
)

type CreateMcpServerRequest struct {
	Key        string                   `json:"key" validate:"required"`
	Name       string                   `json:"name" validate:"required"`
	ModuleSlug string                   `json:"module_slug"`
	Url        string                   `json:"url" validate:"required"`
	Auth       agentmodel.McpServerAuth `json:"auth"`
}

type UpdateMcpServerRequest struct {
	Name       *string                   `json:"name,omitempty"`
	ModuleSlug *string                   `json:"module_slug,omitempty"`
	Url        *string                   `json:"url,omitempty"`
	Auth       *agentmodel.McpServerAuth `json:"auth,omitempty"`
}

type GetManyMcpServersQuery struct {
	Key          *string `json:"key" form:"key" db:"key"`
	Name         *string `json:"name" form:"name" db:"name"`
	ModuleSlug   *string `json:"module_slug" form:"module_slug" db:"module_slug"`
	NameContains *string `json:"name[contains]" form:"name[contains]" db:"name" op:"contains"`
	common.Pagination
	common.Sorting
}

type GetManyMcpServersResponse struct {
	Meta common.ResponseMeta    `json:"meta"`
	Data []agentmodel.McpServer `json:"data"`
}

// StartMcpOAuthResponse is returned by the authenticated connect-start endpoint:
// the URL to send the browser to, plus the opaque single-use state.
type StartMcpOAuthResponse struct {
	AuthorizationUrl string `json:"authorization_url"`
	State            string `json:"state"`
}

// McpOAuthCallbackQuery binds the authorization server's redirect query params on
// the UNAUTHENTICATED callback. Org/user/server identity is recovered from the
// single-use flow state keyed by State — never trusted from these params.
type McpOAuthCallbackQuery struct {
	Code             string `form:"code"`
	State            string `form:"state"`
	Error            string `form:"error"`
	ErrorDescription string `form:"error_description"`
}

// ListMcpServerToolsResponse is the /mcp-servers/:key/tools subresource payload.
type ListMcpServerToolsResponse struct {
	Data []agentmodel.McpToolSummary `json:"data"`
}

// GetMcpConnectionStatusResponse wraps an MCP server's derived connection status.
type GetMcpConnectionStatusResponse struct {
	Data agentmodel.McpConnectionStatus `json:"data"`
}

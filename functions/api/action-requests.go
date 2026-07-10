package functionsapi

import (
	"encoding/json"

	"go.proteos.ai/model/common"
	functionsmodel "go.proteos.ai/model/functions"
	metamodel "go.proteos.ai/model/meta"
)

// DeployActionRequest is the JSON metadata side of the multipart body
// for `POST /api/v1/actions` and `PUT /api/v1/actions/:slug` (upsert).
// `EntitySlug` is required when `Scope == ActionScopeEntity`; validation
// is enforced server-side (CHECK constraint on the `actions` table +
// service-layer guard). The wasm bytes come from the multipart `wasm`
// part; the entry-point Go file is by convention always `./main.go`
// next to the manifest, so no source path is carried on the wire.
type DeployActionRequest struct {
	Slug       string                     `json:"slug" validate:"required"`
	ModuleSlug string                     `json:"module_slug" validate:"required"`
	Scope      functionsmodel.ActionScope `json:"scope" validate:"required"`
	EntitySlug *string                    `json:"entity,omitempty"`
	Name       string                     `json:"name" validate:"required"`
	// IsPublic opts a global action into public (unauthenticated) dispatch.
	// Defaults false; only honoured for Scope == ActionScopeGlobal.
	IsPublic bool                  `json:"is_public,omitempty"`
	Params   []metamodel.Attribute `json:"params"`
	Returns  []metamodel.Attribute `json:"returns"`
}

// PatchActionRequest is the body for `PATCH /api/v1/actions/:slug` —
// partial update of non-lifecycle fields. Activate / deactivate have
// dedicated sub-resources.
type PatchActionRequest struct {
	IsActive *bool `json:"is_active,omitempty"`
	IsPublic *bool `json:"is_public,omitempty"`
}

type GetManyActionsQuery struct {
	Slug       *string                     `json:"slug,omitempty" db:"slug"`
	ModuleSlug *string                     `json:"module_slug,omitempty" db:"module_slug"`
	EntitySlug *string                     `json:"entity,omitempty" db:"entity_slug"`
	Scope      *functionsmodel.ActionScope `json:"scope,omitempty" db:"scope"`
	IsActive   *bool                       `json:"is_active,omitempty" db:"is_active"`
	IsPublic   *bool                       `json:"is_public,omitempty" db:"is_public"`
	common.Pagination
	common.Sorting
}

type GetManyActionsResponse struct {
	Meta common.ResponseMeta     `json:"meta"`
	Data []functionsmodel.Action `json:"data"`
}

// ActionSummary is the lightweight projection returned by
// `GET /api/v1/entities/:entity/actions` (ListForEntity discovery).
type ActionSummary struct {
	Slug        string                     `json:"slug"`
	Name        string                     `json:"name"`
	Scope       functionsmodel.ActionScope `json:"scope"`
	EntitySlug  *string                    `json:"entity,omitempty"`
	Description string                     `json:"description,omitempty"`
	Params      []metamodel.Attribute      `json:"params"`
	Returns     []metamodel.Attribute      `json:"returns"`
}

// GetActionLogsQuery — query string for `GET /api/v1/actions/:slug/logs`.
// Same semantics as GetHookLogsQuery; see that type for `Since`, `Level`,
// and `Follow` behaviour.
type GetActionLogsQuery struct {
	Follow bool   `json:"follow,omitempty"`
	Since  string `json:"since,omitempty"`
	Level  string `json:"level,omitempty"`
}

// ---- Dispatch wire shapes ----
//
// `POST /api/v1/actions/:slug/invoke` and
// `POST /api/v1/entities/:entity/records/:recordId/actions/:slug/invoke`
// take the action's params directly as the request body (an arbitrary
// JSON object shaped by the action's `params` schema) and return the
// action's result as the response body (shaped by `returns`).
//
// No additional envelope wraps either side — the slug, entity, and
// recordId come from the URL path.

type InvokeActionResponse struct {
	Result json.RawMessage `json:"result"`
}

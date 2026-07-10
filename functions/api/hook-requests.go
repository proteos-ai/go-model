package functionsapi

import (
	"encoding/json"

	"go.proteos.ai/model/common"
	functionsmodel "go.proteos.ai/model/functions"
)

// DeployHookRequest is the JSON metadata side of the multipart body for
// `POST /api/v1/hooks` and `PUT /api/v1/hooks/:slug` (upsert). The wasm
// bytes come from the multipart `wasm` part; the entry-point Go file is
// by convention always `./main.go` next to the manifest, so no source
// path is carried on the wire.
type DeployHookRequest struct {
	Slug       string                   `json:"slug" validate:"required"`
	ModuleSlug string                   `json:"module_slug" validate:"required"`
	EntitySlug string                   `json:"entity" validate:"required"`
	Event      functionsmodel.HookEvent `json:"event" validate:"required"`
}

// PatchHookRequest is the body for `PATCH /api/v1/hooks/:slug` — partial
// update of non-lifecycle fields. Lifecycle transitions (activate /
// deactivate) have dedicated sub-resources.
type PatchHookRequest struct {
	IsActive *bool `json:"is_active,omitempty"`
}

type GetManyHooksQuery struct {
	Slug       *string                   `json:"slug,omitempty" db:"slug"`
	ModuleSlug *string                   `json:"module_slug,omitempty" db:"module_slug"`
	EntitySlug *string                   `json:"entity,omitempty" db:"entity_slug"`
	Event      *functionsmodel.HookEvent `json:"event,omitempty" db:"event"`
	IsActive   *bool                     `json:"is_active,omitempty" db:"is_active"`
	common.Pagination
	common.Sorting
}

type GetManyHooksResponse struct {
	Meta common.ResponseMeta   `json:"meta"`
	Data []functionsmodel.Hook `json:"data"`
}

// GetHookLogsQuery — query string for `GET /api/v1/hooks/:slug/logs`.
//
// `Since` is a Go duration string ("10m", "1h"); the service interprets
// it as "entries newer than now-Since". `Level` filters to entries at or
// above the threshold (`debug < info < warn < error`). `Follow` switches
// the response to chunked `application/x-ndjson`, one entry per line.
type GetHookLogsQuery struct {
	Follow bool   `json:"follow,omitempty"`
	Since  string `json:"since,omitempty"`
	Level  string `json:"level,omitempty"`
}

// ---- Dispatch envelopes — wire shapes for ----
//   POST /api/v1/entities/:entity/hooks/on-{before,after}-{create,update,delete}
//
// Field names are pinned: data-service serialises against these exact
// JSON keys (`record` / `currentRecord` / `previousRecord`).

type OnBeforeCreateRequest struct {
	Record json.RawMessage `json:"record"`
}

type OnBeforeUpdateRequest struct {
	Record        json.RawMessage `json:"record"`
	CurrentRecord json.RawMessage `json:"current_record"`
}

type OnBeforeDeleteRequest struct {
	Record json.RawMessage `json:"record"`
}

type OnAfterCreateRequest struct {
	Record json.RawMessage `json:"record"`
}

type OnAfterUpdateRequest struct {
	Record         json.RawMessage `json:"record"`
	PreviousRecord json.RawMessage `json:"previous_record"`
}

type OnAfterDeleteRequest struct {
	Record json.RawMessage `json:"record"`
}

// OnBeforeResponse is the body of every successful before-hook response.
// The hook chain may have mutated the record; persistence uses the
// returned bytes.
type OnBeforeResponse struct {
	Record json.RawMessage `json:"record"`
}

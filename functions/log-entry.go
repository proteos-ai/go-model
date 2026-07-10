package functionsmodel

import "time"

// LogEntry is one line returned by the log endpoints:
//
//	GET /api/v1/logs                  (org-wide)
//	GET /api/v1/hooks/:slug/logs      (per-hook filter over the same stream)
//	GET /api/v1/actions/:slug/logs    (per-action filter over the same stream)
//
// `ResourceType` + `ResourceSlug` let the org-wide list disambiguate hook
// vs action emissions without out-of-band context. `HookSlug` stays for
// back-compat with the pre-LUM-68 wire and is populated from
// `ResourceSlug` whenever `ResourceType == "hook"`.
type LogEntry struct {
	Timestamp    time.Time      `json:"timestamp"`
	Level        string         `json:"level"`
	Message      string         `json:"message"`
	ResourceType string         `json:"resource_type,omitempty"`
	ResourceSlug string         `json:"resource_slug,omitempty"`
	InvocationId string         `json:"invocation_id,omitempty"`
	HookSlug     string         `json:"hook_slug,omitempty"`
	Fields       map[string]any `json:"fields,omitempty"`
}

// The values `LogEntry.ResourceType` takes. Centralised here so the host log
// handler, dispatchers, and read paths all agree on the literals.
const (
	LogResourceTypeHook            = "hook"
	LogResourceTypeAction          = "action"
	LogResourceTypeConnectorMethod = "connector_method"
)

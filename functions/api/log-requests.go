package functionsapi

// GetLogsQuery — query string for `GET /api/v1/logs` (org-wide).
//
// `Hooks` and `Actions` are optional repeatable filters; when both are
// empty the endpoint returns every entry for the auth'd org. When one or
// both are populated the result is a union: an entry passes iff its
// `resourceType` matches one of the populated sides and its
// `resourceSlug` is in the corresponding slug set.
//
// `Since`, `Level`, and `Follow` share the per-resource semantics — see
// `GetHookLogsQuery`.
type GetLogsQuery struct {
	Hooks   []string `json:"hook,omitempty"`
	Actions []string `json:"action,omitempty"`
	Follow  bool     `json:"follow,omitempty"`
	Since   string   `json:"since,omitempty"`
	Level   string   `json:"level,omitempty"`
}

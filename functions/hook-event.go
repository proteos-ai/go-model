package functionsmodel

// HookEvent identifies which record-lifecycle event a hook responds to.
// Hooks fire on data-service Create / Update / Delete; the six events
// split before-commit (mutating, abort-on-error) from after-commit
// (best-effort, log-skip).
type HookEvent string

const (
	HookEventBeforeCreate HookEvent = "before_create"
	HookEventBeforeUpdate HookEvent = "before_update"
	HookEventBeforeDelete HookEvent = "before_delete"
	HookEventAfterCreate  HookEvent = "after_create"
	HookEventAfterUpdate  HookEvent = "after_update"
	HookEventAfterDelete  HookEvent = "after_delete"
)

var HookEvents = []HookEvent{
	HookEventBeforeCreate,
	HookEventBeforeUpdate,
	HookEventBeforeDelete,
	HookEventAfterCreate,
	HookEventAfterUpdate,
	HookEventAfterDelete,
}

func (HookEvent) Enum() []interface{} {
	out := make([]interface{}, len(HookEvents))
	for i, e := range HookEvents {
		out[i] = e
	}
	return out
}

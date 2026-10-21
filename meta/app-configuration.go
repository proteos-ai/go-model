package metamodel

import (
	"time"

	"go.proteos.ai/model/common"
)

// AppHomeType says what an app's home reference names — the same grammar as
// MenuItem.Type/Reference, restricted to the two things a home can open.
type AppHomeType string

const (
	AppHomeTypeList AppHomeType = "list"
	AppHomeTypePage AppHomeType = "page"
)

// SupportedAppHomeTypes is the accepted set.
var SupportedAppHomeTypes = []AppHomeType{AppHomeTypeList, AppHomeTypePage}

func (AppHomeType) Enum() []interface{} {
	out := make([]interface{}, len(SupportedAppHomeTypes))
	for i, homeType := range SupportedAppHomeTypes {
		out[i] = homeType
	}
	return out
}

// AppHome is what an app opens on: a list (records table) or a platform page.
type AppHome struct {
	Type      AppHomeType `json:"type"`
	Reference string      `json:"reference"`
}

// RecordPageBinding names the record pages one entity opens in, per mode
// (PageMode). Every field is optional; the web falls back create → edit →
// view → the entity's first record page. Strict object on the wire (no string
// shorthand).
type RecordPageBinding struct {
	View   string `json:"view,omitempty"`
	Create string `json:"create,omitempty"`
	Edit   string `json:"edit,omitempty"`
}

// IsEmpty reports whether no mode is bound.
func (binding RecordPageBinding) IsEmpty() bool {
	return binding.View == "" && binding.Create == "" && binding.Edit == ""
}

// ForMode returns the page slug bound to one mode ("" when unbound).
func (binding RecordPageBinding) ForMode(mode PageMode) string {
	switch mode {
	case PageModeView:
		return binding.View
	case PageModeCreate:
		return binding.Create
	case PageModeEdit:
		return binding.Edit
	}
	return ""
}

// TriggerPageBinding names the trigger pages one action / workflow opens, per
// dialog phase (TriggerPhase). Every field is optional; an unbound phase falls
// back to the invocable's first trigger page, then to the dialog's generated
// layout. Strict object on the wire.
type TriggerPageBinding struct {
	Input   string `json:"input,omitempty"`
	Running string `json:"running,omitempty"`
	Result  string `json:"result,omitempty"`
}

// IsEmpty reports whether no phase is bound.
func (binding TriggerPageBinding) IsEmpty() bool {
	return binding.Input == "" && binding.Running == "" && binding.Result == ""
}

// ForPhase returns the page slug bound to one phase ("" when unbound).
func (binding TriggerPageBinding) ForPhase(phase TriggerPhase) string {
	switch phase {
	case TriggerPhaseInput:
		return binding.Input
	case TriggerPhaseRunning:
		return binding.Running
	case TriggerPhaseResult:
		return binding.Result
	}
	return ""
}

// AppConfiguration is a TYPED binding row: "how app X presents itself" —
// either to everyone (ProfileSlug empty = the app's default configuration) or to
// one profile (ProfileSlug set = an override). One row per (app, profile);
// the web merges override ⊕ default field-wise and falls through to structural
// defaults (menu is_default, first menu leaf, org-default agent, first page)
// for anything still unset. A standalone resource — not embedded in App — so
// a module redeploy never clobbers a customer's own profile override, and so
// module A can configure module B's app.
//
// Defaults keyed by profile alone (default app, app visibility) do NOT live
// here; they are attributes of the account-service Profile.
type AppConfiguration struct {
	Slug       string `json:"slug" sortable:""`
	OrgId      string `json:"org_id" sortable:""`
	ModuleSlug string `json:"module_slug" sortable:""`
	AppSlug    string `json:"app_slug" sortable:""`
	// ProfileSlug binds the row to one account-service profile; empty = default.
	// A cross-service reference the server does not validate (fail-open: an
	// unknown profile simply never matches anyone).
	ProfileSlug string `json:"profile_slug" sortable:""`
	// Home is what the app opens on. Nil = the first list/page leaf of the
	// resolved menu.
	Home *AppHome `json:"home,omitempty"`
	// MenuSlug names the menu to show. Empty = the app's is_default menu.
	MenuSlug string `json:"menu_slug,omitempty" sortable:""`
	// DefaultAgentKey is the agent Ask Proteos preselects inside this app.
	// Empty = the org-default agent. Must be one of AgentKeys when that is set.
	DefaultAgentKey string `json:"default_agent_key,omitempty" sortable:""`
	// AgentKeys restricts the agents offered inside this app. Empty = every
	// org agent. Agent keys are agent-service references, unvalidated here.
	AgentKeys []string `json:"agent_keys,omitempty"`
	// RecordPages maps entity_slug → the record pages a record of that entity
	// opens in from this app, per mode (view / create / edit). Each must be a
	// record page over that entity. Modes absent fall back create → edit →
	// view → the entity's first page.
	RecordPages map[string]RecordPageBinding `json:"record_pages,omitempty"`
	// ActionPages maps action slug → the trigger pages its invoke dialog
	// renders, per phase (input / running / result); WorkflowPages the same
	// for workflow keys. Each must be a trigger page bound to exactly that
	// kind + reference. Phases absent fall back to the invocable's first
	// trigger page, then to the dialog's generated layout.
	ActionPages   map[string]TriggerPageBinding `json:"action_pages,omitempty"`
	WorkflowPages map[string]TriggerPageBinding `json:"workflow_pages,omitempty"`
	// SearchableEntities lists the entity slugs whose records the command
	// palette searches (by record_title) while inside this app. Empty = no
	// record search. Each must name an existing entity (validated at save);
	// only entities with a simple title_template return results at runtime.
	SearchableEntities []string       `json:"searchable_entities,omitempty"`
	CreatedAt          time.Time      `json:"created_at" sortable:""`
	CreatedBy          common.UserRef `json:"created_by" sortable:""`
	UpdatedAt          time.Time      `json:"updated_at" sortable:""`
	UpdatedBy          common.UserRef `json:"updated_by" sortable:""`
}

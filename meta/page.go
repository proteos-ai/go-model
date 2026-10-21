package metamodel

import (
	"go.proteos.ai/model/common"
	"time"
)

// PageType encodes what a page binds to and how it is served — both the
// chrome (app shell vs bare) and the auth posture follow from it:
//
//   - record:   rendered against a single record of EntitySlug; app chrome;
//     authenticated.
//   - platform: standalone, no record context (e.g. a dashboard launched from
//     a menu item); app chrome; authenticated.
//   - kiosk:    standalone, NO app chrome (bare page at /k/…); authenticated.
//   - public:   standalone, NO app chrome (bare page at /p/…);
//     UNAUTHENTICATED — the page layout is world-readable and its components
//     may only call is_public global actions. Never bind org secrets into a
//     public page's layout or props.
//   - trigger:  what the person sees when they trigger an action or a
//     workflow — rendered inside the invoke dialog, authenticated. Bound by
//     Kind + Reference to ONE invocable; its field elements bind to that
//     invocable's params and returns instead of an entity's attributes. Which
//     dialog phase (input / running / result) a trigger page renders in is
//     the app configuration's binding, like which record mode a record page
//     opens in.
//
// Empty normalizes to PageTypeRecord for backward compatibility.
type PageType string

const (
	PageTypeRecord   PageType = "record"
	PageTypePlatform PageType = "platform"
	PageTypeKiosk    PageType = "kiosk"
	PageTypePublic   PageType = "public"
	PageTypeTrigger  PageType = "trigger"
)

// SupportedPageTypes is the set of types this version accepts. (The formerly
// reserved "external" placeholder was retired in favour of kiosk + public.)
var SupportedPageTypes = []PageType{PageTypeRecord, PageTypePlatform, PageTypeKiosk, PageTypePublic, PageTypeTrigger}

// IsStandalone reports whether the page renders without a record context —
// everything except record pages. Standalone pages forbid entity_slug and
// record-bound layout elements (related_list, related_record, …).
func (pageType PageType) IsStandalone() bool {
	return pageType.Normalized() != PageTypeRecord
}

// AllowsFields reports whether `field` elements have something to bind to: a
// record's attributes (record pages) or a trigger's params / returns (trigger
// pages). Orthogonal to IsStandalone — a trigger page is standalone AND
// allows fields.
func (pageType PageType) AllowsFields() bool {
	normalized := pageType.Normalized()
	return normalized == PageTypeRecord || normalized == PageTypeTrigger
}

// Normalized returns the type with the empty value defaulting to
// PageTypeRecord, so existing record pages (persisted before `type` existed)
// keep their behavior.
func (pageType PageType) Normalized() PageType {
	if pageType == "" {
		return PageTypeRecord
	}
	return pageType
}

// PageActionKind and PageAction live in the common model package so
// conversation-service (conversation types carry post-call actions) shares
// the one contract — meta already imports conversation, so the type cannot
// live here. The aliases keep every metamodel reference compiling.
type PageActionKind = common.PageActionKind

const (
	PageActionKindAction   = common.PageActionKindAction
	PageActionKindWorkflow = common.PageActionKindWorkflow
	PageActionKindAgent    = common.PageActionKindAgent
	PageActionKindIntent   = common.PageActionKindIntent
)

// PageAction is one toolbar button on a page or list — see common.PageAction.
type PageAction = common.PageAction

// PageMode is how the viewer interacts with a record page — the `mode` the
// Liquid scope and the component protocol carry. A record page serves every
// mode; which page a mode opens in is the app configuration's RecordPageBinding.
//
//   - view:   a read-only record page.
//   - create: the form for a record that is not saved yet.
//   - edit:   the form over a saved record.
type PageMode string

const (
	PageModeView   PageMode = "view"
	PageModeCreate PageMode = "create"
	PageModeEdit   PageMode = "edit"
)

// RecordPageModes is the set a record page is bound for.
var RecordPageModes = []PageMode{PageModeView, PageModeCreate, PageModeEdit}

// TriggerPhase is where an invocation is in the dialog's lifecycle — the
// `phase` the dialog's Liquid scope and component protocol carry, and the key
// of a TriggerPageBinding. Orthogonal to PageMode: the input phase renders in
// create mode (an editable form over an unsaved seed), running and result in
// view mode. A trigger page serves every phase; which page a phase renders is
// the app configuration's binding.
//
//   - input:   the form (and the confirm step, rendered read-only) over the
//     params.
//   - running: while the invoke call runs; for a workflow, the "started" page
//     after the run was accepted. The params, read-only; the scope also
//     carries `execution` for workflows.
//   - result:  once the action returned (the params plus its returns) or the
//     workflow run reached a terminal status (`execution` in scope).
type TriggerPhase string

const (
	TriggerPhaseInput   TriggerPhase = "input"
	TriggerPhaseRunning TriggerPhase = "running"
	TriggerPhaseResult  TriggerPhase = "result"
)

// SupportedTriggerPhases is the accepted set.
var SupportedTriggerPhases = []TriggerPhase{TriggerPhaseInput, TriggerPhaseRunning, TriggerPhaseResult}

// TriggerSchema is what a trigger page's fields may bind to: the invocable's
// input params and, for actions, its returns (empty for workflows).
type TriggerSchema struct {
	Params  []Attribute `json:"params"`
	Returns []Attribute `json:"returns"`
}

// Attributes is the bindable set of a trigger page: the params, then every
// return whose name no param already takes (a collision resolves to the
// param — the one the person fills in; the result phase still shows the
// returned value under that name).
func (schema TriggerSchema) Attributes() []Attribute {
	taken := make(map[string]struct{}, len(schema.Params))
	out := make([]Attribute, 0, len(schema.Params)+len(schema.Returns))
	for _, attribute := range schema.Params {
		taken[attribute.Name] = struct{}{}
		out = append(out, attribute)
	}
	for _, attribute := range schema.Returns {
		if _, exists := taken[attribute.Name]; exists {
			continue
		}
		out = append(out, attribute)
	}
	return out
}

type Page struct {
	Slug       string   `json:"slug" sortable:""`
	OrgId      string   `json:"org_id" sortable:""`
	ModuleSlug string   `json:"module_slug" sortable:""`
	Type       PageType `json:"type" sortable:""`
	EntitySlug string   `json:"entity_slug,omitempty" sortable:""`
	// Kind / Reference bind a trigger page to its invocable (an action slug
	// for kind=action, a workflow key for kind=workflow — the same kinds a
	// PageAction button names; agent and intent buttons have no params and
	// take no page). Set on type=trigger only; immutable after create, like
	// EntitySlug.
	Kind      PageActionKind `json:"kind,omitempty" sortable:""`
	Reference string         `json:"reference,omitempty" sortable:""`
	Name      string         `json:"name" sortable:""`
	Actions   []PageAction   `json:"actions"`
	Layout    PageLayout     `json:"layout"`
	CreatedAt time.Time      `json:"created_at" sortable:""`
	CreatedBy common.UserRef `json:"created_by" sortable:""`
	UpdatedAt time.Time      `json:"updated_at" sortable:""`
	UpdatedBy common.UserRef `json:"updated_by" sortable:""`
}

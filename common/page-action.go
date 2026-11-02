package common

import (
	"errors"
	"fmt"
	"strings"

	"golang.org/x/exp/slices"
)

// Icon names an icon (a Lucide PascalCase name, or one of the legacy
// kebab-case menu icons). It lives here, not in metamodel, because PageAction
// is shared by metadata-service (pages, lists) and conversation-service
// (conversation types) — and the meta package already imports conversation.
type Icon string

// PageActionKind selects what a PageAction button invokes. Empty normalizes to
// PageActionKindAction so buttons persisted before `kind` existed keep working.
type PageActionKind string

const (
	// PageActionKindAction invokes a function-service Action (by slug).
	PageActionKindAction PageActionKind = "action"
	// PageActionKindWorkflow starts a manual run of a workflow (by key).
	PageActionKindWorkflow PageActionKind = "workflow"
	// PageActionKindAgent opens an Ava chat with the agent `Agent`: `Prompt`
	// becomes the visible first message, `Context` is handed to the model
	// alongside it and never shown.
	PageActionKindAgent PageActionKind = "agent"
	// PageActionKindIntent prepares a record proposal (one model call over
	// `Prompt` + `Context`), shows it as an editable review and saves only
	// what the person confirms — the command palette's create/update intents.
	PageActionKindIntent PageActionKind = "intent"
	// PageActionKindLink opens `Url` in a new browser tab. The template
	// renders synchronously on the press — no model call, nothing awaited — so
	// the open stays inside the click and is never popup-blocked.
	PageActionKindLink PageActionKind = "link"
)

// Normalized returns the kind with the empty value defaulting to action.
func (kind PageActionKind) Normalized() PageActionKind {
	if kind == "" {
		return PageActionKindAction
	}
	return kind
}

// PageActionIntent is what a kind=intent button prepares. The values are the
// command palette's intent names.
type PageActionIntent string

const (
	// PageActionIntentCreateRecord prepares a NEW record of `Entity`.
	PageActionIntentCreateRecord PageActionIntent = "create_record"
	// PageActionIntentUpdateRecord prepares changes to one existing record of
	// `Entity` — `Record` when set, else the surface's linked record of that
	// entity (picked by the person when there are several).
	PageActionIntentUpdateRecord PageActionIntent = "update_record"
	// PageActionIntentUpsertRecord lets the model choose: change one of the
	// records `Match` finds, or prepare a NEW record of `Entity` when none of
	// them fits. The rules for choosing live in `Prompt`.
	PageActionIntentUpsertRecord PageActionIntent = "upsert_record"
	// PageActionIntentInvokeAction opens the invoke dialog of `Action` with
	// its params filled by the model from `Prompt` + `Context`; `Params`
	// (Liquid) still prefill and lock what the surface already knows. The
	// person always reviews before it runs.
	PageActionIntentInvokeAction PageActionIntent = "invoke_action"
	// PageActionIntentInvokeWorkflow is invoke_action for a manual run of
	// `Workflow`, its Manual Trigger inputs filled by the model; `Inputs`
	// prefill and lock.
	PageActionIntentInvokeWorkflow PageActionIntent = "invoke_workflow"
)

// PageActionIntents is the closed set a kind=intent button may name.
var PageActionIntents = []PageActionIntent{
	PageActionIntentCreateRecord,
	PageActionIntentUpdateRecord,
	PageActionIntentUpsertRecord,
	PageActionIntentInvokeAction,
	PageActionIntentInvokeWorkflow,
}

// isRecordIntent reports whether the intent prepares a record of `Entity`.
func (intent PageActionIntent) isRecordIntent() bool {
	return intent == PageActionIntentCreateRecord ||
		intent == PageActionIntentUpdateRecord ||
		intent == PageActionIntentUpsertRecord
}

const pageActionIntentList = "create_record, update_record, upsert_record, invoke_action or invoke_workflow"

// PageActionMatchMaxLimit caps how many candidates an upsert_record button
// hands the model.
const PageActionMatchMaxLimit = 25

// PageActionMatch is the list query an upsert_record button finds its
// candidate records with — a page of `Entity` records, fetched when the
// button is pressed. Filter has the list endpoint's `_filter` shape; every
// element's Value is a Liquid template rendered against the surface's scope,
// and a value that renders empty yields no candidates at all (never a wider
// query). Sort is the list endpoint's `field:asc|desc`; Limit defaults to 10.
type PageActionMatch struct {
	Filter *FilterGroup `json:"filter,omitempty"`
	Sort   string       `json:"sort,omitempty"`
	Limit  int          `json:"limit,omitempty"`
}

// PageAction is one action button — a page's or a list's toolbar button, or a
// conversation type's action (the strip in the call or meeting surface's
// footer, during and after the conversation). Which target it invokes
// follows from Kind:
//
//   - action   → Action (slug) + Params
//   - workflow → Workflow (key) + Inputs
//   - agent    → Agent (key) + Prompt + Context
//   - intent   → Intent + Prompt + Context (+ FieldScope), and by intent:
//     create/update/upsert_record → Entity (+ Record | Match);
//     invoke_action → Action (+ Params); invoke_workflow → Workflow (+ Inputs)
//   - link     → Url
//
// Params / Inputs map target field names to Liquid templates; Prompt, Context,
// Entity, Record and Url are Liquid templates too. Everything renders client-side
// against the surface's scope (a page: {record, entity, params, current_user};
// a conversation: {conversation, contact, records, records_by_entity, tags,
// tags_by_definition, brief, current_user}, where `| summarize: "<lens>"` on
// the conversation or a tag writes a focused summary of what was said when the
// button is pressed — except in Url, which renders synchronously and leaves
// `summarize` empty). Prompt is literal message TEXT, never a Prompt resource
// key. VisibleWhen hides the button unless it matches the surface's subject
// (the record on a page, the whole scope on a conversation — the stage is
// `conversation.status eq active` / `in ended|archived`, a detected tag
// `tags_by_definition.<key> not_empty`). SkipConfirmation fires an
// action/workflow immediately when every required field resolved.
type PageAction struct {
	Label    string           `json:"label"`
	Icon     Icon             `json:"icon"`
	Kind     PageActionKind   `json:"kind,omitempty"`
	Action   string           `json:"action,omitempty"`
	Workflow string           `json:"workflow,omitempty"`
	Agent    string           `json:"agent,omitempty"`
	Intent   PageActionIntent `json:"intent,omitempty"`
	Entity   string           `json:"entity,omitempty"`
	Record   string           `json:"record,omitempty"`
	Match    *PageActionMatch `json:"match,omitempty"`
	// FieldScope (intent only) lists the attribute names the proposal may
	// fill — plain names, not templates. Set, it replaces the app's bound
	// create / edit page as the proposal's scope; required attributes
	// without a default are always added so a create can be saved.
	FieldScope       []string          `json:"field_scope,omitempty"`
	Url              string            `json:"url,omitempty"`
	Params           map[string]string `json:"params,omitempty"`
	Inputs           map[string]string `json:"inputs,omitempty"`
	Prompt           string            `json:"prompt,omitempty"`
	Context          string            `json:"context,omitempty"`
	SkipConfirmation bool              `json:"skip_confirmation,omitempty"`
	VisibleWhen      *FilterGroup      `json:"visible_when,omitempty"`
	// ExtensionKey is the list-extension stamp for a toolbar action a list
	// extension appended to its host list — same contract as
	// Column.ExtensionKey. Never set on a page's own actions.
	ExtensionKey string `json:"extension_key,omitempty"`
}

// PageActionValidationError is one broken rule of a PageAction list. Code is
// the stable wire code each service reports it under (as a 400).
type PageActionValidationError struct {
	Index   int
	Code    string
	Message string
}

func (validationError *PageActionValidationError) Error() string {
	return fmt.Sprintf("actions[%d]: %s", validationError.Index, validationError.Message)
}

// ValidatePageActions enforces the shape of an action button set: each button
// names exactly the target its kind requires and carries only the fields that
// belong to that kind. Pure — referenced actions, workflows, agents and
// entities are resolved by the client at click time, never here; templates
// are not parsed. Returns the first broken rule as a *PageActionValidationError.
func ValidatePageActions(actions []PageAction) error {
	for index, action := range actions {
		if err := validatePageAction(action); err != nil {
			return &PageActionValidationError{Index: index, Code: err.code, Message: err.message}
		}
	}
	return nil
}

type pageActionRuleError struct {
	code    string
	message string
}

func pageActionRule(code, format string, args ...any) *pageActionRuleError {
	return &pageActionRuleError{code: code, message: fmt.Sprintf(format, args...)}
}

const (
	pageActionCodeLabelRequired      = "page_action_label_required"
	pageActionCodeTargetRequired     = "page_action_target_required"
	pageActionCodeTargetConflict     = "page_action_target_conflict"
	pageActionCodePrefillMismatch    = "page_action_prefill_mismatch"
	pageActionCodePromptMismatch     = "page_action_prompt_mismatch"
	pageActionCodeUnsupportedKind    = "unsupported_page_action_kind"
	pageActionCodeUnsupportedIntent  = "unsupported_page_action_intent"
	pageActionCodeVisibleWhenInvalid = "page_action_visible_when_invalid"
	pageActionCodeMatchInvalid       = "page_action_match_invalid"
	pageActionCodeFieldScopeInvalid  = "page_action_field_scope_invalid"
)

func validatePageAction(action PageAction) *pageActionRuleError {
	if strings.TrimSpace(action.Label) == "" {
		return pageActionRule(pageActionCodeLabelRequired, "label is required")
	}
	kind := action.Kind.Normalized()
	if err := validatePageActionTarget(kind, action); err != nil {
		return err
	}
	if err := validatePageActionFields(kind, action); err != nil {
		return err
	}
	if err := validatePageActionMatch(action.Match); err != nil {
		return err
	}
	if err := validatePageActionFieldScope(action.FieldScope); err != nil {
		return err
	}
	if action.VisibleWhen != nil {
		if problems := ValidateFilterGroup("visible_when", action.VisibleWhen); len(problems) > 0 {
			first := problems[0]
			return pageActionRule(pageActionCodeVisibleWhenInvalid, "%s: %s", first.Path, first.Message)
		}
	}
	return nil
}

// validatePageActionTarget checks the one target field the kind requires.
func validatePageActionTarget(kind PageActionKind, action PageAction) *pageActionRuleError {
	switch kind {
	case PageActionKindAction:
		if action.Action == "" {
			return pageActionRule(pageActionCodeTargetRequired, "an action button requires `action` (the action slug)")
		}
	case PageActionKindWorkflow:
		if action.Workflow == "" {
			return pageActionRule(pageActionCodeTargetRequired, "a workflow button requires `workflow` (the workflow key)")
		}
	case PageActionKindAgent:
		if action.Agent == "" {
			return pageActionRule(pageActionCodeTargetRequired, "an agent button requires `agent` (the agent key)")
		}
	case PageActionKindIntent:
		if action.Intent == "" {
			return pageActionRule(pageActionCodeTargetRequired, "an intent button requires `intent` (%s)", pageActionIntentList)
		}
		if !isPageActionIntent(action.Intent) {
			return pageActionRule(pageActionCodeUnsupportedIntent,
				"intent %q is not supported (expected %s)", action.Intent, pageActionIntentList)
		}
		switch {
		case action.Intent == PageActionIntentInvokeAction && action.Action == "":
			return pageActionRule(pageActionCodeTargetRequired, "an invoke_action intent requires `action` (the action slug)")
		case action.Intent == PageActionIntentInvokeWorkflow && action.Workflow == "":
			return pageActionRule(pageActionCodeTargetRequired, "an invoke_workflow intent requires `workflow` (the workflow key)")
		case action.Intent.isRecordIntent() && strings.TrimSpace(action.Entity) == "":
			return pageActionRule(pageActionCodeTargetRequired, "a %s intent requires `entity` (the entity slug)", action.Intent)
		}
	case PageActionKindLink:
		if strings.TrimSpace(action.Url) == "" {
			return pageActionRule(pageActionCodeTargetRequired, "a link button requires `url` (the address to open)")
		}
	default:
		return pageActionRule(pageActionCodeUnsupportedKind,
			"kind %q is not supported (expected action, workflow, agent, intent or link)", action.Kind)
	}
	return nil
}

// pageActionOwnedFields is what a button of this kind (and intent) may set
// beyond label, icon, prompt/context and visible_when.
func pageActionOwnedFields(kind PageActionKind, intent PageActionIntent) []string {
	switch kind {
	case PageActionKindAction:
		return []string{"action", "params", "skip_confirmation"}
	case PageActionKindWorkflow:
		return []string{"workflow", "inputs", "skip_confirmation"}
	case PageActionKindAgent:
		return []string{"agent"}
	case PageActionKindLink:
		return []string{"url"}
	case PageActionKindIntent:
		switch intent {
		case PageActionIntentUpdateRecord:
			return []string{"intent", "entity", "record", "field_scope"}
		case PageActionIntentUpsertRecord:
			return []string{"intent", "entity", "match", "field_scope"}
		case PageActionIntentInvokeAction:
			return []string{"intent", "action", "params", "field_scope"}
		case PageActionIntentInvokeWorkflow:
			return []string{"intent", "workflow", "inputs", "field_scope"}
		default:
			return []string{"intent", "entity", "field_scope"}
		}
	}
	return nil
}

// validatePageActionFields rejects every field the button's kind (and, for
// an intent, its intent) does not carry.
func validatePageActionFields(kind PageActionKind, action PageAction) *pageActionRuleError {
	owned := pageActionOwnedFields(kind, action.Intent)
	fields := []struct {
		field string
		isSet bool
		code  string
	}{
		{"action", action.Action != "", pageActionCodeTargetConflict},
		{"workflow", action.Workflow != "", pageActionCodeTargetConflict},
		{"agent", action.Agent != "", pageActionCodeTargetConflict},
		{"intent", action.Intent != "", pageActionCodeTargetConflict},
		{"entity", action.Entity != "", pageActionCodeTargetConflict},
		{"record", action.Record != "", pageActionCodeTargetConflict},
		{"match", action.Match != nil, pageActionCodeTargetConflict},
		{"field_scope", len(action.FieldScope) > 0, pageActionCodeTargetConflict},
		{"url", action.Url != "", pageActionCodeTargetConflict},
		{"params", len(action.Params) > 0, pageActionCodePrefillMismatch},
		{"inputs", len(action.Inputs) > 0, pageActionCodePrefillMismatch},
		{"skip_confirmation", action.SkipConfirmation, pageActionCodeTargetConflict},
	}
	for _, field := range fields {
		if !field.isSet || slices.Contains(owned, field.field) {
			continue
		}
		if kind == PageActionKindIntent {
			return pageActionRule(field.code, "a %s intent must not set `%s`", action.Intent, field.field)
		}
		return pageActionRule(field.code, "a %s button must not set `%s`", kind, field.field)
	}
	carriesPrompt := kind == PageActionKindAgent || kind == PageActionKindIntent
	if !carriesPrompt && (action.Prompt != "" || action.Context != "") {
		return pageActionRule(pageActionCodePromptMismatch, "a %s button must not set `prompt` or `context` (agent and intent buttons only)", kind)
	}
	return nil
}

// validatePageActionMatch checks an upsert_record button's candidate query.
func validatePageActionMatch(match *PageActionMatch) *pageActionRuleError {
	if match == nil {
		return nil
	}
	if match.Limit < 0 || match.Limit > PageActionMatchMaxLimit {
		return pageActionRule(pageActionCodeMatchInvalid, "match.limit must be between 1 and %d", PageActionMatchMaxLimit)
	}
	if match.Filter != nil {
		if problems := ValidateFilterGroup("match.filter", match.Filter); len(problems) > 0 {
			first := problems[0]
			return pageActionRule(pageActionCodeMatchInvalid, "%s: %s", first.Path, first.Message)
		}
	}
	return nil
}

// validatePageActionFieldScope checks an intent button's field list: plain,
// distinct attribute names. Whether they exist is the client's to resolve —
// the entity is a template.
func validatePageActionFieldScope(fieldScope []string) *pageActionRuleError {
	seen := make(map[string]bool, len(fieldScope))
	for index, name := range fieldScope {
		trimmed := strings.TrimSpace(name)
		if trimmed == "" {
			return pageActionRule(pageActionCodeFieldScopeInvalid, "field_scope[%d] is empty", index)
		}
		if seen[trimmed] {
			return pageActionRule(pageActionCodeFieldScopeInvalid, "field_scope[%d]: %q is listed twice", index, trimmed)
		}
		seen[trimmed] = true
	}
	return nil
}

func isPageActionIntent(intent PageActionIntent) bool {
	return slices.Contains(PageActionIntents, intent)
}

// AsPageActionValidationError unwraps err into the typed rule violation.
func AsPageActionValidationError(err error) (*PageActionValidationError, bool) {
	var validationError *PageActionValidationError
	if errors.As(err, &validationError) {
		return validationError, true
	}
	return nil, false
}

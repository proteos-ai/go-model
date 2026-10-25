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
)

// PageActionIntents is the closed set a kind=intent button may name.
var PageActionIntents = []PageActionIntent{
	PageActionIntentCreateRecord,
	PageActionIntentUpdateRecord,
}

// PageAction is one action button — a page's or a list's toolbar button, or a
// conversation type's next step (a live call's quick action, the wrap-up's
// "Next steps"). Which target it invokes follows from Kind:
//
//   - action   → Action (slug) + Params
//   - workflow → Workflow (key) + Inputs
//   - agent    → Agent (key) + Prompt + Context
//   - intent   → Intent + Entity (+ Record) + Prompt + Context
//
// Params / Inputs map target field names to Liquid templates; Prompt, Context,
// Entity and Record are Liquid templates too. Everything renders client-side
// against the surface's scope (a page: {record, entity, params, current_user};
// a conversation: {conversation, contact, records, records_by_entity, tags,
// tags_by_definition, brief, current_user}, where `| summarize: "<lens>"` on
// the conversation or a tag writes a focused summary of what was said when the
// button is pressed). Prompt is literal message TEXT, never a Prompt resource
// key. VisibleWhen hides the button unless it matches the surface's subject
// (the record on a page, the whole scope on a conversation — the stage is
// `conversation.status eq active` / `in ended|archived`, a detected tag
// `tags_by_definition.<key> not_empty`). SkipConfirmation fires an
// action/workflow immediately when every required field resolved.
type PageAction struct {
	Label            string            `json:"label"`
	Icon             Icon              `json:"icon"`
	Kind             PageActionKind    `json:"kind,omitempty"`
	Action           string            `json:"action,omitempty"`
	Workflow         string            `json:"workflow,omitempty"`
	Agent            string            `json:"agent,omitempty"`
	Intent           PageActionIntent  `json:"intent,omitempty"`
	Entity           string            `json:"entity,omitempty"`
	Record           string            `json:"record,omitempty"`
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
			return pageActionRule(pageActionCodeTargetRequired, "an intent button requires `intent` (create_record or update_record)")
		}
		if !isPageActionIntent(action.Intent) {
			return pageActionRule(pageActionCodeUnsupportedIntent,
				"intent %q is not supported (expected create_record or update_record)", action.Intent)
		}
		if strings.TrimSpace(action.Entity) == "" {
			return pageActionRule(pageActionCodeTargetRequired, "an intent button requires `entity` (the entity slug)")
		}
	default:
		return pageActionRule(pageActionCodeUnsupportedKind,
			"kind %q is not supported (expected action, workflow, agent or intent)", action.Kind)
	}
	return nil
}

// validatePageActionFields rejects every field that belongs to another kind.
func validatePageActionFields(kind PageActionKind, action PageAction) *pageActionRuleError {
	targets := []struct {
		field   string
		isSet   bool
		ownedBy PageActionKind
	}{
		{"action", action.Action != "", PageActionKindAction},
		{"workflow", action.Workflow != "", PageActionKindWorkflow},
		{"agent", action.Agent != "", PageActionKindAgent},
		{"intent", action.Intent != "", PageActionKindIntent},
		{"entity", action.Entity != "", PageActionKindIntent},
		{"record", action.Record != "", PageActionKindIntent},
	}
	for _, target := range targets {
		if target.isSet && target.ownedBy != kind {
			return pageActionRule(pageActionCodeTargetConflict, "a %s button must not set `%s`", kind, target.field)
		}
	}
	if kind == PageActionKindIntent && action.Intent == PageActionIntentCreateRecord && action.Record != "" {
		return pageActionRule(pageActionCodeTargetConflict, "a create_record intent must not set `record`")
	}
	if len(action.Params) > 0 && kind != PageActionKindAction {
		return pageActionRule(pageActionCodePrefillMismatch, "a %s button must not set `params` (action buttons only)", kind)
	}
	if len(action.Inputs) > 0 && kind != PageActionKindWorkflow {
		return pageActionRule(pageActionCodePrefillMismatch, "a %s button must not set `inputs` (workflow buttons only)", kind)
	}
	carriesPrompt := kind == PageActionKindAgent || kind == PageActionKindIntent
	if !carriesPrompt && (action.Prompt != "" || action.Context != "") {
		return pageActionRule(pageActionCodePromptMismatch, "a %s button must not set `prompt` or `context` (agent and intent buttons only)", kind)
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

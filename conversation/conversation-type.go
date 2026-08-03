package conversationmodel

import (
	"time"

	"go.proteos.ai/model/common"
)

// ConversationTypeConfig is the optional per-type behavior configuration.
// Stored as JSONB so later knobs (model overrides, auto-actions, …) need no
// DDL.
type ConversationTypeConfig struct {
	// SummaryPromptKey references an agent-service prompt by key. When a
	// conversation is classified as this type, the prompt's CURRENT version
	// body replaces the built-in summary system prompt (the instruction —
	// metadata, glossary, transcript — is unchanged). Empty ⇒ built-in prompt.
	SummaryPromptKey string `json:"summary_prompt_key,omitempty"`
}

// ConversationType is a per-org conversation taxonomy entry. Before a meeting
// summary is generated, a cheap classifier model reads the transcript head
// plus every type's Definition and picks the matching Key (or none); the
// result lands on Conversation.TypeKey and may swap the summary system prompt
// via Config.SummaryPromptKey. Keyed (org_id, key) and module-deployable
// (conversation-types/<key>.json), like agent-service prompts.
type ConversationType struct {
	OrgId string `json:"org_id"`
	// Key is the immutable identity within the org — the value the classifier
	// answers with and the value stored on Conversation.TypeKey.
	Key  string `json:"key" sortable:""`
	Name string `json:"name,omitempty"`
	// Definition tells the classifier when a conversation IS this type — it is
	// injected verbatim into the classification prompt, so it should describe
	// the conversation's content/participants/purpose, not internal process.
	Definition string                 `json:"definition"`
	Config     ConversationTypeConfig `json:"config"`
	// ModuleSlug attributes the type to the module that deployed it; empty =
	// not module-owned.
	ModuleSlug string         `json:"module_slug,omitempty"`
	CreatedAt  time.Time      `json:"created_at" sortable:""`
	CreatedBy  common.UserRef `json:"created_by"`
	UpdatedAt  time.Time      `json:"updated_at" sortable:""`
	UpdatedBy  common.UserRef `json:"updated_by"`
}

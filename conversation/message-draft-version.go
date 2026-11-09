package conversationmodel

import (
	"time"

	"go.proteos.ai/model/common"
)

// MessageDraftVersionKind says how a version of a draft came to be.
type MessageDraftVersionKind string

const (
	// MessageDraftVersionKindDrafted is version 1: what was proposed.
	MessageDraftVersionKindDrafted MessageDraftVersionKind = "drafted"
	// MessageDraftVersionKindEdited is a person's edit. Saves by the same person
	// update one version, so typing makes one version, not hundreds.
	MessageDraftVersionKindEdited MessageDraftVersionKind = "edited"
	// MessageDraftVersionKindRedrafted is one model rewrite.
	MessageDraftVersionKindRedrafted MessageDraftVersionKind = "redrafted"
	// MessageDraftVersionKindRestored copies an earlier version forward.
	MessageDraftVersionKindRestored MessageDraftVersionKind = "restored"
)

// MessageDraftVersion is one step in a draft message's history.
type MessageDraftVersion struct {
	Id        string                  `json:"id"`
	OrgId     string                  `json:"org_id"`
	MessageId string                  `json:"message_id"`
	Version   int                     `json:"version"`
	Kind      MessageDraftVersionKind `json:"kind"`
	// Subject is the effective subject at this version, wherever it is stored.
	Subject       string             `json:"subject"`
	Recipients    []MessageRecipient `json:"recipients"`
	Content       []ContentBlock     `json:"content"`
	ToneProfileId string             `json:"tone_profile_id,omitempty"`
	// Instruction is the reviewer's redraft prompt (kind=redrafted).
	Instruction string `json:"instruction,omitempty"`
	// SourceVersion is the version a restore copied from (kind=restored).
	SourceVersion int            `json:"source_version,omitempty"`
	CreatedAt     time.Time      `json:"created_at"`
	CreatedBy     common.UserRef `json:"created_by"`
	UpdatedAt     time.Time      `json:"updated_at"`
	UpdatedBy     common.UserRef `json:"updated_by"`
}

// MessageDraftChange is one field that differs between what was proposed and
// the latest version. Field is body | subject | to | cc | bcc.
type MessageDraftChange struct {
	Field string `json:"field"`
	From  any    `json:"from"`
	To    any    `json:"to"`
}

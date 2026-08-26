package conversationmodel

import (
	"time"

	"go.proteos.ai/model/common"
)

// ToneProfileSetup is the per-user opt-in for tone-of-voice synthesis: only
// set-up users are swept. The row also carries the whole-user synthesis claim
// (Status/StartedAt) — one run at a time per user, whether triggered by the
// sweep or the manual synthesize action.
type ToneProfileSetup struct {
	Id    string `json:"id"`
	OrgId string `json:"org_id"`
	// OwnedBy is the profiled platform user.
	OwnedBy           common.UserRef         `json:"owned_by"`
	Status            ToneProfileSetupStatus `json:"status" sortable:""`
	StartedAt         *time.Time             `json:"started_at,omitempty"`
	LastSynthesizedAt *time.Time             `json:"last_synthesized_at,omitempty" sortable:""`
	CreatedAt         time.Time              `json:"created_at" sortable:""`
	CreatedBy         common.UserRef         `json:"created_by"`
	UpdatedAt         time.Time              `json:"updated_at" sortable:""`
	UpdatedBy         common.UserRef         `json:"updated_by"`
}

// ToneProfile is one generated tone-of-voice instruction row. Rows form a
// specificity hierarchy per profiled user — user aggregate (the constant
// voice), per-channel base, per contact-group, per individual contact — and
// the MOST SPECIFIC matching row wins at read time. Every row is
// self-contained: a drafting consumer injects exactly one row's Instructions
// verbatim and never concatenates tiers (override cascades degrade model
// instruction-following).
type ToneProfile struct {
	Id    string `json:"id"`
	OrgId string `json:"org_id"`
	// OwnedBy is the profiled platform user.
	OwnedBy common.UserRef `json:"owned_by"`
	// Channel scopes the row to one medium; empty = the cross-channel user
	// aggregate (the voice proper).
	Channel Channel `json:"channel,omitempty" sortable:""`
	// ContactGroupKey scopes the row to one contact group on the channel;
	// empty = the (user, channel) base row.
	ContactGroupKey string `json:"contact_group_key,omitempty"`
	// ContactId scopes the row to one individual contact within the group;
	// only written when the contact had enough corpus samples.
	ContactId string `json:"contact_id,omitempty"`
	// Scope is derived from which scope fields are set — never stored.
	Scope ToneProfileScope `json:"scope"`
	// Instructions is the COMPLETE markdown instruction set for this tier,
	// served verbatim to a drafting model.
	Instructions string `json:"instructions"`
	// Differences is the short delta vs the tier above — what a human scans in
	// the UI, and the guard that keeps a tier from being a near-copy of its
	// parent (no meaningful deviation ⇒ no row). Empty on root tiers.
	Differences string `json:"differences,omitempty"`
	// SampleCount is the largest corpus (message count) any synthesis run has
	// used for this row; the (user, channel) base row's count drives the
	// mature refresh cadence (mature = the corpus cap is saturated).
	SampleCount       int            `json:"sample_count"`
	LastMessageAt     *time.Time     `json:"last_message_at,omitempty"`
	LastSynthesizedAt *time.Time     `json:"last_synthesized_at,omitempty" sortable:""`
	CreatedAt         time.Time      `json:"created_at" sortable:""`
	CreatedBy         common.UserRef `json:"created_by"`
	UpdatedAt         time.Time      `json:"updated_at" sortable:""`
	UpdatedBy         common.UserRef `json:"updated_by"`
}

// ResolveScope derives the tier from which scope fields are set.
func (profile ToneProfile) ResolveScope() ToneProfileScope {
	switch {
	case profile.ContactId != "":
		return ToneProfileScopeContact
	case profile.ContactGroupKey != "":
		return ToneProfileScopeGroup
	case profile.Channel != "":
		return ToneProfileScopeChannel
	default:
		return ToneProfileScopeUser
	}
}

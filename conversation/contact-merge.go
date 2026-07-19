package conversationmodel

import (
	"time"

	"go.proteos.ai/model/common"
)

// ContactMerge is the append-only audit record of one executed contact merge:
// which contact absorbed which, how the merge fired, on what evidence, and a
// full pre-merge snapshot of the losing side so support can reconstruct (and a
// future split/unmerge can consume) the prior state. Rows are never updated.
type ContactMerge struct {
	Id              string `json:"id"`
	OrgId           string `json:"org_id"`
	WinnerContactId string `json:"winner_contact_id"`
	LoserContactId  string `json:"loser_contact_id"`
	// InitiatedBy says HOW the merge fired (deterministic logic, LLM
	// adjudication, human decision) — a separate axis from CreatedBy (WHO wrote
	// the row; system merges carry the system user ref).
	InitiatedBy MergeInitiator `json:"initiated_by"`
	// ProposalId links the driving merge proposal when the merge came through
	// the review queue; empty for direct/deterministic merges.
	ProposalId string `json:"proposal_id,omitempty"`
	// Evidence is why the two were judged the same person — the co-carried
	// observation for deterministic merges, signals/verdict for adjudicated
	// ones.
	Evidence map[string]any `json:"evidence,omitempty"`
	// PreMergeSnapshot is the loser contact row plus its address rows as they
	// stood before the merge — the undo material.
	PreMergeSnapshot map[string]any `json:"pre_merge_snapshot"`
	CreatedAt        time.Time      `json:"created_at"`
	CreatedBy        common.UserRef `json:"created_by"`
}

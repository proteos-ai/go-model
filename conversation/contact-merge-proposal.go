package conversationmodel

import (
	"time"

	"go.proteos.ai/model/common"
)

// ContactMergeProposal is one candidate duplicate pair awaiting review: two
// contacts that some evidence says are the same person, but that safety rules
// (non-thin loser) or confidence kept from auto-merging. The pair is stored
// canonically ordered (ContactAId = smaller id) so the partial unique index
// admits at most one open proposal per pair.
type ContactMergeProposal struct {
	Id    string `json:"id"`
	OrgId string `json:"org_id"`
	// Canonical pair ordering: ContactAId < ContactBId (by string compare).
	ContactAId string `json:"contact_a_id"`
	ContactBId string `json:"contact_b_id"`
	// SuggestedWinnerId is logic.SelectWinner's pick — the UI default when a
	// human approves. Empty when no preference.
	SuggestedWinnerId string `json:"suggested_winner_id,omitempty"`
	// Confidence is 1.0 for deterministic evidence (a co-carried observation
	// proves the pair); an adjudication score otherwise.
	Confidence *float64 `json:"confidence,omitempty"`
	// Evidence is the machine-readable why: {kind: deterministic|blocked,
	// observation, signals}.
	Evidence map[string]any `json:"evidence,omitempty"`
	// Verdict is the LLM adjudication memory ({same_person, confidence,
	// reasoning, model, adjudicated_at}); nil until the dedup runner exists.
	Verdict    map[string]any      `json:"verdict,omitempty"`
	Status     MergeProposalStatus `json:"status" sortable:""`
	ResolvedBy *common.UserRef     `json:"resolved_by,omitempty"`
	ResolvedAt *time.Time          `json:"resolved_at,omitempty"`
	CreatedAt  time.Time           `json:"created_at" sortable:""`
	CreatedBy  common.UserRef      `json:"created_by"`
	UpdatedAt  time.Time           `json:"updated_at" sortable:""`
	UpdatedBy  common.UserRef      `json:"updated_by"`
}

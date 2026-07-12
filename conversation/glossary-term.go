package conversationmodel

import (
	"time"

	"go.proteos.ai/model/common"
)

// GlossaryTerm is a per-org custom vocabulary entry that improves transcription
// accuracy. The Term is injected into the transcription provider's keyterm
// prompting (Deepgram keyterm/keywords) at the moment a transcription starts;
// the Definition is human/agent-facing documentation only and never reaches the
// provider. Priority is optional: terms carrying a priority outrank those
// without, then ascending by number — this decides which terms win when an org
// stores more terms than the provider's per-request keyterm cap.
type GlossaryTerm struct {
	Id         string `json:"id"`
	OrgId      string `json:"org_id"`
	Term       string `json:"term" sortable:""`
	Definition string `json:"definition,omitempty"`
	// Priority ranks a term for keyterm selection: any term with a priority is
	// selected ahead of any without; within priority, lower number wins. Nil ⇒
	// unprioritized (selected only after every prioritized term).
	Priority  *int           `json:"priority,omitempty" sortable:""`
	CreatedAt time.Time      `json:"created_at" sortable:""`
	CreatedBy common.UserRef `json:"created_by"`
	UpdatedAt time.Time      `json:"updated_at" sortable:""`
	UpdatedBy common.UserRef `json:"updated_by"`
}

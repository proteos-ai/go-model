package agentmodel

import (
	"time"

	"go.proteos.ai/model/common"
)

// Skill is a versioned bundle of instructions + files in the Anthropic
// Agent-Skills shape. Key = the SKILL.md frontmatter `name` (kebab). The parent
// row caches the L1 discovery metadata (name + description) extracted from the
// frontmatter; the bundle itself (SKILL.md + files, a tar.gz in storage-service)
// is the source of truth and lives on the immutable version rows.
type Skill struct {
	OrgId          string         `json:"org_id"`
	Key            string         `json:"key" sortable:""`
	Name           string         `json:"name" sortable:""`
	ModuleSlug     string         `json:"module_slug" sortable:""`
	Description    string         `json:"description"`
	CurrentVersion *SkillVersion  `json:"current_version"`
	CreatedAt      time.Time      `json:"created_at" sortable:""`
	CreatedBy      common.UserRef `json:"created_by"`
	UpdatedAt      time.Time      `json:"updated_at" sortable:""`
	UpdatedBy      common.UserRef `json:"updated_by"`
}

// SkillVersion is an immutable snapshot pinning one storage bundle.
type SkillVersion struct {
	Id        string         `json:"id"`
	Number    uint32         `json:"number"`
	OrgId     string         `json:"org_id"`
	SkillKey  string         `json:"skill_key"`
	Bundle    SkillBundle    `json:"bundle"`
	CreatedAt time.Time      `json:"created_at"`
	CreatedBy common.UserRef `json:"created_by"`
}

// SkillBundle pins the exact immutable storage Version the SKILL.md+files tar.gz
// was uploaded as. file_id alone resolves to the File's latest version, which
// would break old-version immutability — so FileVersionId is the real pin;
// Checksum / SizeInBytes are denormalized from that storage Version for cheap
// display without a round-trip. Checksum is the canonical "<algo>:<hex>" content
// address (see common.FormatChecksum) — the same convention as hook/action/
// component artifacts, so the CLI compares all of them uniformly.
type SkillBundle struct {
	FileId        string `json:"file_id"`
	FileVersionId string `json:"file_version_id"`
	Checksum      string `json:"checksum"`
	SizeInBytes   int64  `json:"size_in_bytes"`
}

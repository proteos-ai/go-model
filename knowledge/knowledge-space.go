package knowledgemodel

import (
	"time"

	"go.proteos.ai/model/common"
)

// KnowledgeSpace is the container a node lives in — engineering, sales, people.
// A node belongs to at most one space; a node with no space is unassigned and
// visible to anyone holding knowledge-nodes:read.
//
// A space is NOT a label. Labels are many-per-node, overlapping, non-security
// tags for categorization; a space is the single place a node lives, and it is
// what instance-level access is evaluated against. Both exist and neither
// replaces the other.
//
// Slug-keyed within an org (composite PK), because spaces are deployable CONFIG
// rather than content: a module manifest declaring `space: engineering` must
// resolve to the same space in every environment, which a uuid minted at deploy
// time cannot do. That also makes the slug immutable by construction — renaming
// would mean creating a different space — so an access grant naming it can never
// be orphaned.
//
// Note what is NOT here: no owner, no team list, no visibility flag. Every
// authority statement about a space lives in the knowledge_access_grants table,
// so there is exactly one mechanism to read and one place a permission can come
// from.
type KnowledgeSpace struct {
	OrgId       string  `json:"org_id"`
	Slug        string  `json:"slug" sortable:""` // immutable; half the primary key
	Name        string  `json:"name" sortable:""`
	Description *string `json:"description,omitempty"`
	Color       *string `json:"color,omitempty"`
	Icon        *string `json:"icon,omitempty"`
	// Scope says who a space is for: SpaceScopeShared (the default — an org
	// space such as engineering or sales) or SpaceScopePersonal (one person's
	// own space, where their assistants keep memory notes). Scope is set at
	// creation and never changes.
	Scope string `json:"scope" sortable:""`
	// OwnerUserId is the person a personal space belongs to; nil on a shared
	// space. It records whose space this is — it is NOT what grants access.
	// Access still comes only from knowledge_access_grants: a personal space is
	// created holding grants for its owner and nobody else.
	OwnerUserId *string        `json:"owner_user_id,omitempty"`
	CreatedAt   time.Time      `json:"created_at" sortable:""`
	CreatedBy   common.UserRef `json:"created_by"`
	UpdatedAt   time.Time      `json:"updated_at" sortable:""`
	UpdatedBy   common.UserRef `json:"updated_by"`
}

const (
	// SpaceScopeShared is an ordinary org space. The default.
	SpaceScopeShared = "shared"
	// SpaceScopePersonal is one person's own space, created on first use by
	// the ensure-personal action rather than by POST /spaces.
	SpaceScopePersonal = "personal"
)

// PersonalSpaceSlugPrefix starts every personal space's slug: the slug is this
// prefix plus the owner's user id, so a person's space is addressable without a
// lookup and creating it twice converges on one row. POST /spaces refuses slugs
// with this prefix, so nobody can claim another person's slug ahead of them.
const PersonalSpaceSlugPrefix = "personal-"

// UnassignedSpaceSlug is the reserved filter value meaning "nodes that belong to
// no space at all". It is deliberately NOT a legal slug — slugs are kebab-case
// (`^[a-z0-9]+(?:-[a-z0-9]+)*$`), so underscores make collision with a slug an
// org actually creates impossible.
//
// It exists because an absent space filter and an "unassigned only" filter are
// different questions: no filter means every node, while this one means exactly
// the nodes no space owns. Without a sentinel the second is unaskable.
const UnassignedSpaceSlug = "__unassigned__"

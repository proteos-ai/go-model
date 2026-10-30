package datamodel

import "time"

// EntityIndexType names an index by what it does, never by the Postgres
// structure behind it.
type EntityIndexType string

const (
	// EntityIndexTypeLookup serves equality, IN, ranges and ORDER BY on its
	// attributes (a btree on the typed projection of each).
	EntityIndexTypeLookup EntityIndexType = "lookup"
	// EntityIndexTypeUnique is a lookup index that also refuses a second
	// record with the same attribute values (409 record_unique_violation).
	EntityIndexTypeUnique EntityIndexType = "unique"
	// EntityIndexTypeSearch serves the contains / starts_with / ends_with
	// filters on ONE string attribute — substring matching, not word search
	// (a trigram GIN).
	EntityIndexTypeSearch EntityIndexType = "search"
)

func (indexType EntityIndexType) IsValid() bool {
	switch indexType {
	case EntityIndexTypeLookup, EntityIndexTypeUnique, EntityIndexTypeSearch:
		return true
	}
	return false
}

// EntityIndexSource says where an index definition comes from, which decides
// who may remove it.
type EntityIndexSource string

const (
	// EntityIndexSourceDeclared indexes come from the entity definition —
	// is_unique and unique_constraints. Removed by changing the entity.
	EntityIndexSourceDeclared EntityIndexSource = "declared"
	// EntityIndexSourceAutomatic indexes the platform adds by rule: every
	// top-level relation and user attribute gets a lookup index, because
	// relation hops and the delete cascade filter on them.
	EntityIndexSourceAutomatic EntityIndexSource = "automatic"
	// EntityIndexSourceApplied indexes were added per org by an admin (or by
	// accepting an advisor recommendation). The only kind that can be deleted
	// through the index API.
	EntityIndexSourceApplied EntityIndexSource = "applied"
)

// EntityIndexState is the live state of the index on the record table.
type EntityIndexState string

const (
	EntityIndexStateReady    EntityIndexState = "ready"
	EntityIndexStateBuilding EntityIndexState = "building"
	// EntityIndexStateInvalid is a build that failed — for a unique index,
	// almost always because the table already holds duplicates. The index
	// exists but is not used and, for unique, not enforced.
	EntityIndexStateInvalid EntityIndexState = "invalid"
	// EntityIndexStateMissing is desired but not present on the table yet
	// (the reconcile has not run or failed before building it).
	EntityIndexStateMissing EntityIndexState = "missing"
)

// EntityIndex is one index on one entity's record table, as the index API
// reports it. Name/Type/Attributes/Source are the definition; the rest is live
// state read from pg_catalog when listed.
type EntityIndex struct {
	Name       string            `json:"name"`
	Type       EntityIndexType   `json:"type"`
	Attributes []string          `json:"attributes"`
	Source     EntityIndexSource `json:"source"`
	// State is filled by the status listing; zero on a bare definition.
	State     EntityIndexState `json:"state,omitempty"`
	SizeBytes int64            `json:"size_bytes,omitempty"`
	Scans     int64            `json:"scans,omitempty"`
	// CreatedAt is set for applied indexes (the row's creation time); the
	// advisor uses it to age an unused index before recommending removal.
	CreatedAt *time.Time `json:"created_at,omitempty"`
}

// IndexRecommendationAction is what the advisor suggests doing.
type IndexRecommendationAction string

const (
	IndexRecommendationAdd    IndexRecommendationAction = "add"
	IndexRecommendationRemove IndexRecommendationAction = "remove"
)

// IndexEvidence is why the advisor recommends an index: the query shape ran
// slow SlowQueries times in the window (only slow runs are recorded), those
// runs averaged AvgMs, over a table of about Rows rows. For a removal, Scans
// is the index's lifetime scan count (zero) and Rows the table size.
type IndexEvidence struct {
	SlowQueries int64   `json:"slow_queries"`
	AvgMs       float64 `json:"avg_ms"`
	Rows        int64   `json:"rows"`
	Scans       int64   `json:"scans"`
	Reason      string  `json:"reason"`
}

// IndexRecommendation is one advisor suggestion. `add` carries a complete
// index definition ready for POST .../indexes; `remove` names an applied index.
type IndexRecommendation struct {
	Action   IndexRecommendationAction `json:"action"`
	Index    EntityIndex               `json:"index"`
	Evidence IndexEvidence             `json:"evidence"`
	// Replaces, on an `add`, names the applied lookup indexes the new index
	// makes redundant (their attributes are its leading attributes). Pass it
	// through as CreateEntityIndexRequest.Replaces so the create removes them
	// in the same step — creating the wider index while they stay is refused.
	Replaces []string `json:"replaces,omitempty"`
}

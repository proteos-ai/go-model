package common

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
)

// PublicAccessOperation is one operation a resource may be exposed for on the
// UNAUTHENTICATED public surface. It is a SET, not a level: the operations are
// independent, so a resource can be public for `write` (a form-submission /
// lead-capture endpoint) WITHOUT being public for `read`.
type PublicAccessOperation string

const (
	// PublicAccessRead — records/definition/file are world-READABLE
	// unauthenticated. The only operation supported today.
	PublicAccessRead PublicAccessOperation = "read"
	// PublicAccessWrite — anonymous create/update. RESERVED — the anonymous
	// principal + hook/m2m + abuse-guardrail work is deferred; ValidatePublicAccess
	// rejects it for now.
	PublicAccessWrite PublicAccessOperation = "write"
	// PublicAccessDelete — anonymous delete. RESERVED — see PublicAccessWrite.
	PublicAccessDelete PublicAccessOperation = "delete"
)

// PublicAccess is the set of operations a resource is publicly exposed for. An
// empty (or nil) set means fully private — the default. Stored as JSONB (a JSON
// array of operation strings) via the Valuer/Scanner below, so it maps cleanly
// in bun without a struct tag and serializes to a plain array on the wire.
type PublicAccess []PublicAccessOperation

// Contains reports whether op is in the set.
func (access PublicAccess) Contains(op PublicAccessOperation) bool {
	for _, granted := range access {
		if granted == op {
			return true
		}
	}
	return false
}

// Value implements driver.Valuer — marshals to a JSON array. A nil set stores
// as `[]` (not `null`) so the column can stay NOT NULL DEFAULT '[]'.
func (access PublicAccess) Value() (driver.Value, error) {
	if access == nil {
		return []byte("[]"), nil
	}
	return json.Marshal(access)
}

// Scan implements sql.Scanner — reads the JSONB array back.
func (access *PublicAccess) Scan(src any) error {
	if src == nil {
		*access = nil
		return nil
	}
	var raw []byte
	switch value := src.(type) {
	case []byte:
		raw = value
	case string:
		raw = []byte(value)
	default:
		return fmt.Errorf("common.PublicAccess: cannot scan %T", src)
	}
	if len(raw) == 0 {
		*access = nil
		return nil
	}
	return json.Unmarshal(raw, access)
}

// ValidatePublicAccess enforces the public-access policy on a write:
//   - every operation must be a known value, and
//   - only `read` is supported today — `write`/`delete` are rejected with a
//     clear "not yet supported" error until their backends land. When they do,
//     relax the switch here (single source of truth for both entities + files).
//
// Returns nil for an empty/nil set (fully private).
func ValidatePublicAccess(access PublicAccess) error {
	for _, op := range access {
		switch op {
		case PublicAccessRead:
			// allowed
		case PublicAccessWrite, PublicAccessDelete:
			return fmt.Errorf("public %q access is not yet supported (only \"read\" is available today)", op)
		default:
			return fmt.Errorf("unknown public access operation %q (expected one of read, write, delete)", op)
		}
	}
	return nil
}

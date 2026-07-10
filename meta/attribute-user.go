package metamodel

// UserAttributeMeta holds the metadata for an attribute of type `user`.
//
// A user attribute stores a reference to a platform user — the identity owned
// by the account-service, not a record in the tenant's data schema. The
// value is the composite common.UserRef `{ type, id }`: `id` is the user id and
// `type` is the user kind (person | agent | api). It is stored as a JSONB object
// and filtered/sorted on the nested `id` (`data->'field'->>'id'`); the people
// picker resolves and renders the chip by `id`.
//
// The attribute carries no required configuration today; this struct exists for
// parity with the other AttributeMeta variants and reserves room for future
// options (e.g. multi-user selection).
type UserAttributeMeta struct {
	Description string `json:"description,omitempty"`
}

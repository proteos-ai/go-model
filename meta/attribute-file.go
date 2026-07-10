package metamodel

// FileAttributeMeta holds the metadata for an attribute of type `file`.
//
// A file attribute stores a reference to a stored file — the blob owned by the
// storage-service, not a record in the tenant's data schema. The value is the
// composite common.FileRef `{ id, name }`: `id` is the storage-service file id
// and `name` is the original filename. It is stored as a JSONB object and
// filtered on the nested `id` (`data->'field'->>'id'`); the control resolves
// size/content-type from the storage-service by `id`.
//
// The attribute carries no required configuration today; this struct exists for
// parity with the other AttributeMeta variants and reserves room for future
// options (e.g. accepted content types, max size).
type FileAttributeMeta struct {
	Description string `json:"description,omitempty"`
}

package common

// FileRef is a reference to a stored file — the blob owned by the
// storage-service, not a record in the tenant's data schema. It is the value of
// any `file` attribute, stored as a JSONB object `{ id, name }`: `id` is the
// storage-service file id and `name` is the original filename (denormalised so
// the record carries a human label without a storage round-trip).
//
// Only `id` and `name` are persisted; size and content type are resolved from
// the storage-service by `id` on read. Filtering keys on `id`
// (`data->'field'->>'id'`).
type FileRef struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

package common

// FileRef is a reference to a stored file — the blob owned by the
// storage-service, not a record in the tenant's data schema. It is the value of
// any `file` attribute, stored as a JSONB object `{ id, name, content_type? }`:
// `id` is the storage-service file id, `name` the original filename and
// `content_type` the MIME type — both denormalised so the record carries a
// human label and a renderer hint without a storage round-trip.
//
// `content_type` is optional: writers that know the MIME include it, readers
// fall back to resolving it (and size) from the storage-service by `id`. The
// storage-service stays authoritative — after an in-place replace the stored
// hint may lag until the record is next written. Filtering keys on `id`
// (`data->'field'->>'id'`).
type FileRef struct {
	Id          string `json:"id"`
	Name        string `json:"name"`
	ContentType string `json:"content_type,omitempty"`
}

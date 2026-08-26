package functionsmodel

// CompiledRef is the precompiled-machine-code triple carried by both Hook and
// Action as their `compiled_*` columns. It exists as a named type so the
// persistence port can take the three values as one argument instead of a
// positional string triple that is trivial to transpose.
//
// FileId points at a storage-service blob holding the compiled output for the
// row's wasm. Target is the "<wazero-version>_<GOOS>_<GOARCH>" the artifact is
// valid for — a runtime whose own target differs ignores the artifact and
// compiles the wasm locally, so a runtime upgrade needs no migration. Checksum
// is the canonical "<algo>:<hex>" address of the blob, verified before the
// bytes are trusted as machine code.
type CompiledRef struct {
	FileId   string `json:"compiled_file_id,omitempty"`
	Target   string `json:"compiled_target,omitempty"`
	Checksum string `json:"compiled_checksum,omitempty"`
}

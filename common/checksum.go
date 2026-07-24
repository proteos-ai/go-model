package common

import "strings"

// Content-addressing convention for deployed artifacts (wasm hooks/actions,
// component bundles, skill archives). The stored value is "<algo>:<hex>" — the
// git/OCI/SRI form — so the hash algorithm travels WITH the value instead of
// being baked into a column or field name. Switching algorithms later is then a
// value-format change that self-heals through a redeploy, never a schema rename.
//
// The value is only ever compared for equality (does the local build match
// what's deployed?), never parsed for its algorithm — the prefix exists purely
// so a future mixed-algorithm state stays unambiguous.
const ChecksumAlgoSha256 = "sha256"

// FormatChecksum renders a content address in canonical "<algo>:<hex>" form.
// An empty digest yields an empty string (meaning "unknown / not stamped").
func FormatChecksum(algo, hexDigest string) string {
	if hexDigest == "" {
		return ""
	}
	return algo + ":" + hexDigest
}

// NormalizeChecksum upgrades a bare hex digest to canonical "sha256:<hex>" form.
// A value that already carries an "<algo>:" prefix passes through unchanged, and
// empty stays empty. This is the compatibility shim for the two places raw hex
// still appears: rows deployed before checksum stamping and the CLI builder's
// dist sidecars (which fingerprint a local build as plain sha256). Comparing two
// NormalizeChecksum results is the correct artifact-equality test.
func NormalizeChecksum(s string) string {
	if s == "" || strings.Contains(s, ":") {
		return s
	}
	return ChecksumAlgoSha256 + ":" + s
}

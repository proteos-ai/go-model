package common

import (
	"fmt"
	"strings"
)

// StripBearerPrefix removes the leading `Bearer ` (case-insensitive)
// from an Authorization header value, returning the raw token. Returns
// an error when the value is missing the prefix entirely — useful for
// authentication middlewares feeding a JWT parser, which expect the
// raw token without scheme.
//
// Callers that just want to defensively normalize whatever they got
// (and don't care whether the prefix was present) should use a plain
// strings.TrimPrefix instead.
func StripBearerPrefix(value string) (string, error) {
	const prefix = "Bearer "
	if len(value) >= len(prefix) && strings.EqualFold(value[:len(prefix)], prefix) {
		return value[len(prefix):], nil
	}
	return "", fmt.Errorf("Invalid token.")
}
